<script setup lang="ts">
import { useToggle } from "@vueuse/core";
import { Browser } from "@wailsio/runtime";
import IconEye from "~icons/ph/eye";
import IconEyeSlash from "~icons/ph/eye-slash";
import SettingsGroup from "@/components/SettingsGroup.vue";
import SettingsRow from "@/components/SettingsRow.vue";
import {
  approvalURL,
  isConnecting,
  nocturneCode,
  nocturneConnected,
  nocturneStatus,
  nocturneSupported,
  testConnection,
  testStatus,
  testing,
  toggleNocturne,
} from "@/connection";
import { errors, settings } from "@/settings";

const [revealToken, toggleRevealToken] = useToggle(false);

function openApproval() {
  if (approvalURL.value) Browser.OpenURL(approvalURL.value);
}
</script>

<template>
  <settings-group>
    <settings-row label="URL" for="url">
      <input
        id="url"
        v-model="settings.url"
        type="url"
        placeholder="https://nightscout.example.com"
        spellcheck="false"
        autocomplete="off"
        :disabled="isConnecting"
        :aria-invalid="errors.url ? true : undefined"
        :title="errors.url"
      />
      <button type="button" :disabled="testing" @click="testConnection">Test</button>
    </settings-row>
    <p class="status" :class="testStatus.kind" aria-live="polite">{{ testStatus.text }}</p>
  </settings-group>

  <settings-group title="Authentication">
    <template v-if="nocturneSupported">
      <settings-row label="Nocturne" description="Link this app using Nocturne's Connected Apps.">
        <div class="nocturne-status">{{ nocturneStatus }}</div>
        <template #trailing>
          <button type="button" @click="toggleNocturne">
            {{ isConnecting ? "Cancel" : nocturneConnected ? "Disconnect" : "Connect…" }}
          </button>
        </template>
      </settings-row>
      <settings-row v-if="isConnecting" class="pending">
        <div class="grow">
          <div>Approve the app in your browser and confirm this code:</div>
          <div class="code">{{ nocturneCode }}</div>
          <a href="#" @click.prevent="openApproval">Open approval page</a>
        </div>
      </settings-row>
      <div class="or">or</div>
    </template>
    <settings-row
      label="Token"
      for="token"
      description="An access token is recommended instead of the API secret."
      :class="{ dimmed: nocturneConnected }"
    >
      <input
        id="token"
        v-model="settings.token"
        :type="revealToken ? 'text' : 'password'"
        autocomplete="off"
        spellcheck="false"
      />
      <template #trailing>
        <button
          type="button"
          class="icon"
          :title="revealToken ? 'Hide token' : 'Show token'"
          :aria-label="revealToken ? 'Hide token' : 'Show token'"
          :aria-pressed="revealToken"
          @click="toggleRevealToken()"
        >
          <icon-eye-slash v-if="revealToken" />
          <icon-eye v-else />
        </button>
      </template>
    </settings-row>
    <p v-if="nocturneConnected" class="status muted">Not used while connected to Nocturne.</p>
  </settings-group>
</template>

<style scoped>
.nocturne-status {
  line-height: 22px;
}
.dimmed {
  opacity: 0.5;
}

.pending {
  align-items: flex-start;
  line-height: 1.5;
}
.code {
  font:
    600 22px ui-monospace,
    "SF Mono",
    Menlo,
    Consolas,
    monospace;
  letter-spacing: 0.08em;
  margin: 4px 0;
  -webkit-user-select: text;
  user-select: text;
}

.or {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 10px;
  font-size: 11px;
  color: var(--secondary);
}
.or::before,
.or::after {
  content: "";
  flex: 1;
  border-top: 1px solid var(--separator);
}

button.icon {
  appearance: none;
  border: 0;
  background: none;
  padding: 2px;
  border-radius: 4px;
  color: var(--secondary);
  display: flex;
}
button.icon:hover {
  color: var(--label);
}
button.icon svg {
  width: 18px;
  height: 18px;
}
</style>
