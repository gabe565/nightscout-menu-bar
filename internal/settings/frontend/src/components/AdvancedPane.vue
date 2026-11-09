<script setup lang="ts">
import SettingsGroup from "@/components/SettingsGroup.vue";
import SettingsRow from "@/components/SettingsRow.vue";
import { errors, options, settings } from "@/settings";
</script>

<template>
  <settings-group>
    <settings-row label="Fetch delay" for="fetchDelay">
      <input
        id="fetchDelay"
        v-model="settings.fetchDelay"
        type="text"
        class="short"
        placeholder="30s"
        :aria-invalid="errors.fetchDelay ? true : undefined"
        :title="errors.fetchDelay"
      />
      <span class="hint">Wait after a reading is expected before fetching.</span>
    </settings-row>
    <settings-row label="Fallback" for="fallbackInterval">
      <input
        id="fallbackInterval"
        v-model="settings.fallbackInterval"
        type="text"
        class="short"
        placeholder="30s"
        :aria-invalid="errors.fallbackInterval ? true : undefined"
        :title="errors.fallbackInterval"
      />
      <span class="hint">Interval when the next reading can't be estimated.</span>
    </settings-row>
    <settings-row label="Round reading age up to the nearest minute" for="roundAge" wide>
      <input id="roundAge" v-model="settings.roundAge" type="checkbox" switch />
    </settings-row>
    <settings-row label="Log level" for="logLevel">
      <select id="logLevel" v-model="settings.logLevel">
        <option v-for="level in options.logLevels" :key="level" :value="level">{{ level }}</option>
      </select>
    </settings-row>
  </settings-group>
</template>
