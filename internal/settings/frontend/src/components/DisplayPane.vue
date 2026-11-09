<script setup lang="ts">
import { computed } from "vue";
import FontPicker from "@/components/FontPicker.vue";
import SettingsGroup from "@/components/SettingsGroup.vue";
import SettingsRow from "@/components/SettingsRow.vue";
import { preview } from "@/preview";
import { Service } from "@/service";
import { errors, fonts, options, settings } from "@/settings";

const platform = options.platform;

function isLight(hex: string): boolean {
  const n = parseInt(hex.slice(1), 16);
  const [r, g, b] = [n >> 16, (n >> 8) & 0xff, n & 0xff];
  return 0.299 * r + 0.587 * g + 0.114 * b > 140;
}

const barClass = computed(() => {
  if (!preview.value?.icon) return "";
  return isLight(preview.value.iconColor ?? "#000000") ? "on-dark" : "on-light";
});

async function browseFont() {
  const path = await Service.BrowseFont();
  if (path) settings.value.fontFile = path;
}
</script>

<template>
  <!-- Windows only shows an icon in the tray, so there's no text to preview -->
  <settings-group v-if="platform !== 'windows' || settings.dynamicIcon" title="Preview">
    <div class="preview">
      <div class="menubar" :class="barClass">
        <img v-if="preview?.icon" :src="preview.icon" alt="" />
        <span v-else :class="{ stale: preview?.stale }">{{ preview?.text }}</span>
      </div>
    </div>
    <p v-if="preview?.error" class="status err">{{ preview.error }}</p>
    <p v-else-if="preview && !preview.live" class="status muted">
      Sample reading. Test the connection to preview your latest reading.
    </p>
  </settings-group>

  <settings-group>
    <settings-row label="Title" for="title">
      <input id="title" v-model="settings.title" type="text" />
    </settings-row>
    <settings-row label="Units" for="units">
      <select id="units" v-model="settings.units">
        <option v-for="unit in options.units" :key="unit" :value="unit">{{ unit }}</option>
      </select>
    </settings-row>
    <!-- Strikethrough for stale readings is only rendered on macOS -->
    <settings-row v-if="platform === 'darwin'" label="Stale after" for="staleThreshold">
      <input
        id="staleThreshold"
        v-model="settings.staleThreshold"
        type="text"
        class="short"
        placeholder="15m"
        :aria-invalid="errors.staleThreshold ? true : undefined"
        :title="errors.staleThreshold"
      />
      <span class="hint">Older readings are crossed out.</span>
    </settings-row>
  </settings-group>

  <!-- Windows doesn't show text next to tray icons -->
  <settings-group v-if="platform !== 'windows'" title="Menu Bar">
    <settings-row label="Show">
      <div class="grow checks">
        <label><input v-model="settings.showArrow" type="checkbox" /> Arrow</label>
        <label><input v-model="settings.showDelta" type="checkbox" /> Delta</label>
        <label><input v-model="settings.showTimeAgo" type="checkbox" /> Time since reading</label>
      </div>
    </settings-row>
  </settings-group>

  <settings-group title="Dynamic Icon">
    <settings-row label="Show the current reading as the tray icon" for="dynamicIcon" wide>
      <input id="dynamicIcon" v-model="settings.dynamicIcon" type="checkbox" switch />
    </settings-row>
    <fieldset :disabled="!settings.dynamicIcon">
      <settings-row label="Color" for="fontColor">
        <label class="inline">
          <input v-model="settings.autoColor" type="checkbox" /> Automatic
        </label>
        <input
          id="fontColor"
          v-model="settings.fontColor"
          type="color"
          :disabled="settings.autoColor"
        />
        <span v-if="!settings.autoColor" class="hint">{{ settings.fontColor }}</span>
      </settings-row>
      <settings-row label="Font" for="font-input">
        <font-picker
          id="font-input"
          v-model="settings.fontFile"
          :fonts="fonts"
          :disabled="!settings.dynamicIcon"
        />
        <button type="button" @click="browseFont">Browse…</button>
      </settings-row>
      <settings-row label="Max size" for="maxFontSize">
        <input
          id="maxFontSize"
          v-model.number="settings.maxFontSize"
          type="number"
          min="1"
          step="1"
          class="short"
          :aria-invalid="errors.maxFontSize ? true : undefined"
          :title="errors.maxFontSize"
        />
        <span class="hint">pt</span>
      </settings-row>
    </fieldset>
  </settings-group>
</template>

<style scoped>
.checks {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 14px;
}
.checks label,
label.inline {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.preview {
  padding: 12px 10px 8px;
}
.menubar {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 30px;
  border-radius: 6px;
  border: 1px solid var(--group-border);
  background: var(--menubar);
  color: var(--menubar-text);
  font-size: 13px;
}
.menubar.on-dark {
  background: #2a2a2a;
  color: #fff;
}
.menubar.on-light {
  background: #f2f2f2;
  color: #000;
}
.menubar img {
  height: 22px;
  width: 22px;
  image-rendering: auto;
}
.menubar .stale {
  text-decoration: line-through;
}
</style>
