// The service validates again on save.
import type { Settings } from "@/service";

export type Pane = "connection" | "display" | "general" | "advanced";

export const fields = {
  url: { label: "URL", pane: "connection" },
  staleThreshold: { label: "Stale after", pane: "display" },
  maxFontSize: { label: "Max size", pane: "display" },
  fetchDelay: { label: "Fetch delay", pane: "advanced" },
  fallbackInterval: { label: "Fallback", pane: "advanced" },
} as const satisfies Record<string, { label: string; pane: Pane }>;

export type Field = keyof typeof fields;
export type Errors = Partial<Record<Field, string>>;

const durationPattern = /^(0|(\d+(\.\d+)?(ns|us|µs|ms|s|m|h))+)$/;

function validateDuration(value: string): string {
  return durationPattern.test(value.trim()) ? "" : "Use a duration like 30s, 5m, or 1h30m.";
}

export function validate(settings: Settings): Errors {
  const errors: Record<Field, string> = {
    url:
      settings.url.trim() && !/^https?:\/\/\S+$/i.test(settings.url.trim())
        ? "Must start with http:// or https://."
        : "",
    staleThreshold: validateDuration(settings.staleThreshold),
    maxFontSize:
      settings.dynamicIcon && !(Number(settings.maxFontSize) > 0)
        ? "Must be a positive number."
        : "",
    fetchDelay: validateDuration(settings.fetchDelay),
    fallbackInterval: validateDuration(settings.fallbackInterval),
  };
  return Object.fromEntries(Object.entries(errors).filter(([, message]) => message));
}
