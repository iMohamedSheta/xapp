<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div data-slot="sidebar-group" data-sidebar="group"
    :class="cn('relative flex w-full min-w-0 scrollbar scrollbar-thumb-[#353535] scrollbar-track-[#1e1e1e] dark:scrollbar-thumb-[#353535] dark:scrollbar-track-[#1e1e1e] flex-col p-2', props.class)">
    <slot />
  </div>
</template>
