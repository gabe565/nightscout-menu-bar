<script setup lang="ts">
defineProps<{ title?: string }>();
</script>

<template>
  <div class="settings-group">
    <h2 v-if="title">{{ title }}</h2>
    <div class="box">
      <slot />
    </div>
    <p v-if="$slots.note" class="note"><slot name="note" /></p>
  </div>
</template>

<style scoped>
.settings-group + .settings-group {
  margin-top: 16px;
}
h2 {
  font-size: 13px;
  font-weight: 600;
  margin: 0 0 6px 2px;
}
.box {
  background: var(--group);
  border: 1px solid var(--group-border);
  border-radius: 10px;
}
.note {
  font-size: 11px;
  color: var(--secondary);
  margin: 6px 2px 0;
}
</style>
