<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <div data-slot="sidebar-content" data-sidebar="content" :class="cn(
    'flex min-h-0 flex-1 flex-col gap-2 overflow-auto',
    // Hide scrollbar by default, show on hover
    '[&::-webkit-scrollbar]:w-1.5',
    '[&::-webkit-scrollbar-track]:bg-transparent',
    '[&::-webkit-scrollbar-thumb]:bg-transparent',
    '[&::-webkit-scrollbar-thumb]:rounded-full',
    'hover:[&::-webkit-scrollbar-thumb]:bg-white/20',
    'transition-all duration-300',
    'group-data-[collapsible=icon]:overflow-hidden',
    props.class
  )">
    <slot />
  </div>
</template>
