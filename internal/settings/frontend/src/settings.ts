import { computed, nextTick, ref, watch } from "vue";
import { type Font, Service, type Settings, errorMessage } from "@/service";
import { type Errors, type Field, type Pane, fields, validate } from "@/validation";

// Modules importing settings wait for them to load
const [initialOptions, initialSettings] = await Promise.all([Service.Options(), Service.Load()]);
export const options = initialOptions;
export const settings = ref(initialSettings);
export const fonts = ref<Font[]>([]);
export const tab = ref<Pane>("connection");
export const saveError = ref("");
let saveErrorField: Field | undefined;

// A cleared number input is an empty string
export const payload = computed<Settings>(() => ({
  ...settings.value,
  maxFontSize: Number(settings.value.maxFontSize),
}));

const saved = ref(JSON.stringify(payload.value));
export const isDirty = computed(() => JSON.stringify(payload.value) !== saved.value);
export const errors = computed<Errors>(() => validate(settings.value));

export function markSaved<K extends keyof Settings>(name: K) {
  const s: Settings = JSON.parse(saved.value);
  s[name] = payload.value[name];
  saved.value = JSON.stringify(s);
}

export async function save(): Promise<boolean> {
  saveError.value = "";
  saveErrorField = undefined;
  const name = (Object.keys(errors.value) as Field[])[0];
  if (name) {
    tab.value = fields[name].pane;
    await nextTick();
    document.getElementById(name)?.focus();
    saveError.value = fields[name].label + ": " + errors.value[name];
    saveErrorField = name;
    return false;
  }
  try {
    const s = payload.value;
    await Service.Save(s);
    saved.value = JSON.stringify(s);
    return true;
  } catch (err) {
    saveError.value = errorMessage(err);
    return false;
  }
}

watch(errors, (errors) => {
  if (saveErrorField && !errors[saveErrorField]) saveError.value = "";
});
