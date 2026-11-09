import { computedAsync, refDebounced } from "@vueuse/core";
import { CancelError, type CancellablePromise, Events } from "@wailsio/runtime";
import { computed, ref, shallowRef, watch } from "vue";
import { previewReading } from "@/preview";
import { Service, errorMessage } from "@/service";
import { markSaved, payload, settings } from "@/settings";

// Connection test

export const testStatus = ref<{ kind: "" | "muted" | "ok" | "err"; text: string }>({
  kind: "",
  text: "",
});
export const testing = ref(false);

export async function testConnection() {
  testStatus.value = { kind: "muted", text: "Testing…" };
  testing.value = true;
  try {
    const reading = await Service.Test(payload.value);
    testStatus.value = { kind: "ok", text: "✓ Connected. Latest reading: " + reading };
    previewReading.value++;
  } catch (err) {
    testStatus.value = { kind: "err", text: "✗ " + errorMessage(err) };
  } finally {
    testing.value = false;
  }
}

// Nocturne

const connecting = shallowRef<CancellablePromise<void> | null>(null);
export const isConnecting = computed(() => connecting.value !== null);

// Bumped when the connection changes without the URL changing
const connectionChanged = ref(0);
export const nocturneConnected = computedAsync(async () => {
  void connectionChanged.value;
  if (connecting.value) return false;
  return Service.NocturneConnected(settings.value.url);
}, false);

const supportURL = refDebounced(
  computed(() => settings.value.url.trim()),
  600,
);
export const nocturneSupported = computedAsync(async () => {
  if (connecting.value) return true;
  const url = supportURL.value;
  return url !== "" && (await Service.NocturneSupported(url));
}, false);

const nocturneMessage = ref("");
watch(
  () => settings.value.url,
  () => (nocturneMessage.value = ""),
);
export const nocturneStatus = computed(() => {
  if (nocturneMessage.value) return nocturneMessage.value;
  if (connecting.value) return "Waiting for approval…";
  return nocturneConnected.value ? "✓ Connected" : "Not connected";
});
export const nocturneCode = ref("…");
export const approvalURL = ref("");

Events.On("nocturne:code", (event) => {
  nocturneCode.value = event.data.code;
  approvalURL.value = event.data.url;
});

export async function toggleNocturne() {
  nocturneMessage.value = "";
  if (connecting.value) {
    connecting.value.cancel();
    return;
  }
  if (nocturneConnected.value) {
    nocturneMessage.value = "Disconnecting…";
    try {
      await Service.DisconnectNocturne();
    } catch (err) {
      nocturneMessage.value = errorMessage(err);
      return;
    }
    nocturneMessage.value = "";
    connectionChanged.value++;
    return;
  }

  nocturneCode.value = "…";
  approvalURL.value = "";
  connecting.value = Service.ConnectNocturne(settings.value.url);
  try {
    await connecting.value;
    // The service saved the URL along with the connection
    markSaved("url");
    testConnection();
  } catch (err) {
    if (!(err instanceof CancelError)) nocturneMessage.value = "✗ " + errorMessage(err);
  } finally {
    connecting.value = null;
  }
}
