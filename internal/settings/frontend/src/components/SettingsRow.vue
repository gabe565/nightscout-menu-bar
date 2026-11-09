<script setup lang="ts">
defineProps<{
  label?: string;
  for?: string;
  // The label fills the row, e.g. before a switch
  wide?: boolean;
  description?: string;
}>();
</script>

<template>
  <div class="row" :class="{ described: description }">
    <label v-if="label" :for="$props.for" :class="{ wide }">{{ label }}</label>
    <div v-if="description" class="grow">
      <slot />
      <div class="desc">{{ description }}</div>
    </div>
    <slot v-else />
    <slot name="trailing" />
  </div>
</template>

<style scoped>
.row {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 40px;
  padding: 6px 10px;
}
.row + .row,
fieldset > .row:first-child,
.row + fieldset > .row:first-child {
  border-top: 1px solid var(--separator);
}
label {
  flex: 0 0 76px;
}
label.wide {
  flex: 1;
}
.row > :slotted(input:not([type="checkbox"], [type="color"])),
.row > :slotted(select) {
  flex: 1;
  min-width: 0;
}
.row > :slotted(input.short) {
  flex: 0 0 72px;
}
.grow > :slotted(input) {
  width: 100%;
}

.described {
  align-items: flex-start;
}
.described > label {
  line-height: 22px;
}
.desc {
  font-size: 11px;
  color: var(--secondary);
  margin-top: 3px;
}
</style>
