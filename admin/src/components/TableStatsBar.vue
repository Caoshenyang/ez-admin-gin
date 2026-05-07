<script setup lang="ts">
</script>

<template>
  <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-2.5">
    <div>
      <slot />
    </div>
    <div>
      <slot name="actions" />
    </div>
  </div>
</template>
