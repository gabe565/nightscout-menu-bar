<script setup lang="ts">
// A font dropdown that shows each font's name in the font itself.
// The input shows the selected font's name, while the model is its file.
import { unrefElement } from "@vueuse/core";
import {
  AutocompleteInput,
  AutocompleteRoot,
  ComboboxAnchor,
  ComboboxContent,
  ComboboxItem,
  ComboboxPortal,
  ComboboxViewport,
  ComboboxVirtualizer,
} from "reka-ui";
import { type ComponentPublicInstance, computed, ref, watch } from "vue";
import type { Font } from "@/service";

defineOptions({ inheritAttrs: false });

const model = defineModel<string>({ required: true });
const props = defineProps<{ fonts: Font[]; disabled?: boolean }>();

const families = new Map<string, string>();

function previewFamily(file: string): string {
  if (!families.has(file)) {
    const family = `font-preview-${families.size}`;
    const face = new FontFace(family, `url("/font?file=${encodeURIComponent(file)}")`);
    document.fonts.add(face);
    face.load().catch(() => {});
    families.set(file, family);
  }
  return `"${families.get(file)}", system-ui`;
}

const embedded: Font = { file: "", name: "Embedded font" };
const allFonts = computed(() => [embedded, ...props.fonts]);
const nameOf = (file: string) => allFonts.value.find((f) => f.file === file)?.name ?? file;

const text = ref("");
const open = ref(false);
const query = ref("");
let highlighted = false;
const onHighlight = (item: unknown) => (highlighted = !!item);

const shown = computed(() => {
  const q = query.value.trim().toLowerCase();
  return allFonts.value.filter((f) => !q || f.name.toLowerCase().includes(q));
});

function sync() {
  text.value = model.value === "" ? "" : nameOf(model.value);
}
watch([model, allFonts], sync, { immediate: true });
watch(open, () => (query.value = ""));

function select(file: string) {
  open.value = false;
  model.value = file;
  sync();
}

// Accepts a typed font name, file name, or path
function commit() {
  const typed = text.value.trim();
  if (typed === nameOf(model.value) || (typed === "" && model.value === "")) return;
  const match = allFonts.value.find((f) => f.name.toLowerCase() === typed.toLowerCase());
  select(match ? match.file : typed);
}

// Reka selects the highlighted font. Otherwise, use the typed font, and let the form submit.
function onEnter() {
  if (!(open.value && highlighted)) commit();
}

// Close the list without closing the window
function onEscape(e: KeyboardEvent) {
  if (!open.value) return;
  e.stopPropagation();
  open.value = false;
  sync();
}

function measure(virtualizer: { measureElement: (el: Element | null) => void }, item: unknown) {
  virtualizer.measureElement(unrefElement(item as ComponentPublicInstance) ?? null);
}
</script>

<template>
  <autocomplete-root
    v-model="text"
    v-model:open="open"
    class="grow"
    :disabled="disabled"
    open-on-click
    ignore-filter
    @highlight="onHighlight"
  >
    <combobox-anchor class="combo">
      <autocomplete-input
        v-bind="$attrs"
        placeholder="Embedded font"
        spellcheck="false"
        autocomplete="off"
        :style="{ fontFamily: previewFamily(model) }"
        @input="query = text"
        @blur="commit"
        @keydown.enter="onEnter"
        @keydown.esc="onEscape"
      />
    </combobox-anchor>
    <combobox-portal>
      <combobox-content class="combo-list" position="popper" :side-offset="4">
        <combobox-viewport class="combo-viewport">
          <combobox-virtualizer
            v-slot="{ option, virtualizer }"
            :options="shown"
            :estimate-size="30"
            :text-content="(font) => font.name"
          >
            <combobox-item
              :ref="(item) => measure(virtualizer, item)"
              :value="option.name"
              class="combo-option"
              :style="{ fontFamily: previewFamily(option.file) }"
              @select.prevent="select(option.file)"
            >
              {{ option.name }}
            </combobox-item>
          </combobox-virtualizer>
        </combobox-viewport>
      </combobox-content>
    </combobox-portal>
  </autocomplete-root>
</template>

<style scoped>
/* The list is portaled into Reka's positioning wrapper, which gets the scope attribute */
.combo {
  display: flex;
}
.combo > input {
  flex: 1;
  min-width: 0;
}
:deep(.combo-list) {
  z-index: 10;
  width: var(--reka-combobox-trigger-width);
  background: var(--popup);
  border: 1px solid var(--group-border);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
  font-size: 13px;
}
:deep(.combo-viewport) {
  max-height: min(260px, var(--reka-combobox-content-available-height, 260px) - 2px);
  padding: 4px;
}
:deep(.combo-option) {
  width: 100%;
  padding: 5px 8px;
  border-radius: 5px;
  font-size: 15px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
:deep(.combo-option[data-highlighted]) {
  background: var(--accent);
  color: #fff;
}
</style>
