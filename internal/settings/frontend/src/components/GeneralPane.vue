<script setup lang="ts">
import SettingsGroup from "@/components/SettingsGroup.vue";
import SettingsRow from "@/components/SettingsRow.vue";
import { settings } from "@/settings";
</script>

<template>
  <settings-group>
    <settings-row label="Start on login" for="startOnLogin" wide>
      <input id="startOnLogin" v-model="settings.startOnLogin" type="checkbox" switch />
    </settings-row>
    <settings-row label="Check for updates" for="checkUpdates" wide>
      <input id="checkUpdates" v-model="settings.checkUpdates" type="checkbox" switch />
    </settings-row>
  </settings-group>

  <settings-group title="Local Socket">
    <settings-row label="Expose readings over a local socket" for="socket" wide>
      <input id="socket" v-model="settings.socket" type="checkbox" switch />
    </settings-row>
    <template #note>
      Lets other apps read the latest reading, e.g.
      <a
        href="https://github.com/gabe565/nightscout-menu-bar/blob/main/contrib/powerlevel10k/README.md"
        target="_blank"
        >a shell prompt</a
      >.
    </template>
  </settings-group>
</template>
