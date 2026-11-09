<script setup lang="ts">
import { onKeyStroke, useEventListener } from "@vueuse/core";
import { Browser, Dialogs, Events, Window } from "@wailsio/runtime";
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "reka-ui";
import { type Component, onMounted, onScopeDispose, useTemplateRef } from "vue";
import IconGeneral from "~icons/ph/gear";
import IconDisplay from "~icons/ph/palette";
import IconConnection from "~icons/ph/plug";
import IconAdvanced from "~icons/ph/sliders-horizontal";
import AdvancedPane from "@/components/AdvancedPane.vue";
import ConnectionPane from "@/components/ConnectionPane.vue";
import DisplayPane from "@/components/DisplayPane.vue";
import GeneralPane from "@/components/GeneralPane.vue";
import { testConnection } from "@/connection";
import { Service } from "@/service";
import { fonts, isDirty, options, save, saveError, settings, tab } from "@/settings";
import type { Pane } from "@/validation";

const tabs: { pane: Pane; label: string; icon: Component; component: Component }[] = [
  { pane: "connection", label: "Connection", icon: IconConnection, component: ConnectionPane },
  { pane: "display", label: "Display", icon: IconDisplay, component: DisplayPane },
  { pane: "general", label: "General", icon: IconGeneral, component: GeneralPane },
  { pane: "advanced", label: "Advanced", icon: IconAdvanced, component: AdvancedPane },
];

// Like macOS settings windows, tabs aren't tab stops.
// Reka only respects tabindex on as-child elements.
const tabStop = options.platform === "darwin" ? { tabindex: -1 } : {};

// Sizes the window to fit the tallest tab without scrolling
const main = useTemplateRef<HTMLElement>("main");

async function fitWindow(main: HTMLElement) {
  const panes = [...main.querySelectorAll<HTMLElement>('[role="tabpanel"]')];
  const visible = panes.find((p) => !p.hidden);
  let tallest = 0;
  for (const pane of panes) {
    pane.hidden = false;
    for (const other of panes) if (other !== pane) other.hidden = true;
    tallest = Math.max(tallest, main.scrollHeight);
  }
  for (const pane of panes) pane.hidden = pane !== visible;

  const chrome = document.body.scrollHeight - main.clientHeight;
  await Window.SetSize(window.innerWidth, Math.ceil(chrome + tallest));
}

// The service also routes the close button and Cmd+Q here
let closing = false;

async function requestClose() {
  if (closing) return;
  closing = true;
  try {
    if (isDirty.value) {
      const choice = await Dialogs.Question({
        Title: "Do you want to save the changes to your settings?",
        Message: "Your changes will be lost if you don't save them.",
        Buttons: [
          { Label: "Save", IsDefault: true },
          { Label: "Don't Save" },
          { Label: "Cancel", IsCancel: true },
        ],
      });
      if (choice !== "Save" && choice !== "Don't Save") return;
      if (choice === "Save" && !(await save())) return;
    }
    await Service.Close();
  } finally {
    closing = false;
  }
}

async function submit() {
  if (await save()) Service.Close();
}

const command = (e: KeyboardEvent) => e.metaKey || e.ctrlKey;

onKeyStroke(
  (e) => e.key === "Escape" || (command(e) && e.key === "w"),
  (e) => {
    e.preventDefault();
    requestClose();
  },
);

onKeyStroke(
  tabs.map((_, i) => String(i + 1)),
  (e) => {
    if (!command(e)) return;
    e.preventDefault();
    tab.value = tabs[Number(e.key) - 1]!.pane;
  },
);

useEventListener(document, "click", (e) => {
  const link = (e.target as Element).closest<HTMLAnchorElement>('a[href^="http"]');
  if (!link) return;
  e.preventDefault();
  Browser.OpenURL(link.href);
});

onScopeDispose(Events.On("settings:close-requested", requestClose));

onMounted(async () => {
  Service.GuardClose();
  await fitWindow(main.value!);

  // Slow or network-bound work after the window is ready
  Service.Fonts().then((list) => (fonts.value = list ?? []));
  if (settings.value.url) testConnection();
});
</script>

<template>
  <!-- Panes stay mounted so fitWindow can measure them -->
  <tabs-root v-model="tab" class="tabs" :unmount-on-hide="false">
    <header>
      <tabs-list as-child>
        <nav v-bind="tabStop">
          <tabs-trigger v-for="t in tabs" :key="t.pane" :value="t.pane" as-child>
            <button type="button" v-bind="tabStop">
              <component :is="t.icon" />
              {{ t.label }}
            </button>
          </tabs-trigger>
        </nav>
      </tabs-list>
    </header>

    <form novalidate @submit.prevent="submit">
      <main ref="main">
        <!-- Reka makes panes tab stops by default -->
        <tabs-content v-for="t in tabs" :key="t.pane" :value="t.pane" as="section" tabindex="-1">
          <component :is="t.component" />
        </tabs-content>
      </main>

      <footer>
        <p class="status err grow" aria-live="polite">{{ saveError }}</p>
        <button type="button" @click="requestClose">Cancel</button>
        <button type="button" :disabled="!isDirty" @click="save">Apply</button>
        <button type="submit" class="default">Save</button>
      </footer>
    </form>
  </tabs-root>
</template>

<style scoped>
.tabs {
  display: contents;
}

header {
  --wails-draggable: drag;
  flex: none;
  display: flex;
  justify-content: center;
  padding: 6px 0;
  border-bottom: 1px solid var(--separator);
}
/* Clear the traffic lights in macOS's hidden title bar */
[data-platform="darwin"] header {
  padding-top: 30px;
}
nav {
  display: flex;
  gap: 2px;
}
[role="tab"] {
  --wails-draggable: no-drag;
  appearance: none;
  border: 0;
  background: none;
  font-size: 11px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 72px;
  padding: 5px 8px 4px;
  border-radius: 6px;
  color: var(--secondary);
}
[role="tab"] svg {
  width: 22px;
  height: 22px;
}
[role="tab"]:hover {
  color: var(--label);
}
[role="tab"][aria-selected="true"] {
  color: var(--accent);
  background: var(--group);
}
[role="tab"]:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: -2px;
}

form {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
main {
  flex: 1;
  overflow: auto;
  overscroll-behavior: contain;
  padding: 18px 20px 4px;
}

footer {
  flex: none;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px 16px;
}
footer .status {
  margin: 0;
}
/* Keep the error message as a spacer so the buttons stay right-aligned */
footer .status:empty {
  display: block;
}
footer button {
  min-width: 72px;
}
button.default {
  appearance: default-button;
}
</style>
