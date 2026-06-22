<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <ul data-slot="sidebar-menu" data-sidebar="menu"
    :class="cn('flex w-full min-w-0 flex-col gap-1  scrollbar scrollbar-thumb-[#353535] scrollbar-track-[#1e1e1e] dark:scrollbar-thumb-[#353535] dark:scrollbar-track-[#1e1e1e]', props.class)">
    <slot />
  </ul>
</template>
