export {
  Service,
  type Font,
  type Options,
  type Preview,
  type Settings,
} from "@/bindings/gabe565.com/nightscout-menu-bar/internal/settings";

export const errorMessage = (err: unknown): string =>
  typeof err === "object" && err !== null && "message" in err ? String(err.message) : String(err);
