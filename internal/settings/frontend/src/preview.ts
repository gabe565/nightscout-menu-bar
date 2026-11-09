import { computedAsync, refDebounced } from "@vueuse/core";
import { ref } from "vue";
import { type Preview, Service } from "@/service";
import { payload, tab } from "@/settings";

// Bumped when the service has a new live reading to preview
export const previewReading = ref(0);

const debounced = refDebounced(payload, 150);

let last: Preview | null = null;
export const preview = computedAsync(async () => {
  const settings = debounced.value;
  void previewReading.value;
  if (tab.value !== "display") return last;
  return (last = await Service.Preview(settings));
}, null);
