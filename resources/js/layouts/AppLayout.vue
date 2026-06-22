<script setup lang="ts">
import AppHeaderLayout from '@/layouts/app/AppHeaderLayout.vue'
import AppSidebarLayout from '@/layouts/app/AppSidebarLayout.vue'
import type { BreadcrumbItemType } from '@/types'
import Toast from '@/components/Toast.vue'
import { computed } from 'vue'
import { useSidebarAppearance } from '@/composables/useAppearance'
import AppDebug from '@/components/debug/AppDebug.vue'
import { usePage } from '@inertiajs/vue3'

const { sidebarProps } = useSidebarAppearance()

interface Props {
  breadcrumbs?: BreadcrumbItemType[]
  side?: 'left' | 'right'
}

const props = withDefaults(defineProps<Props>(), {
  breadcrumbs: () => [],
})

const side = computed(() => sidebarProps.value.side)
const dir = computed(() => (side.value === 'left' ? 'ltr' : 'rtl'))

const page = usePage()
const role = page.props.auth.user.role

// Pick layout component dynamically
const Layout = computed(() => {
  if (role === 'super_manager' || role === 'manager') return AppSidebarLayout
  if (role === 'super_admin' || role === 'admin' || role === 'distributor') return AppSidebarLayout
  if (role === 'client') return AppHeaderLayout
  return AppHeaderLayout // fallback
})
</script>

<template>
  <component :is="Layout" :breadcrumbs="breadcrumbs" :dir="dir">
    <div class="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
      <div
        class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/15 blur-[150px] rounded-full animate-pulse">
      </div>
      <div
        class="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] bg-blue-500/5 blur-[180px] rounded-full animate-pulse"
        style="animation-delay: 2s"></div>
    </div>
    <slot />
    <Toast />
    <AppDebug v-if="page.props.debug" />
  </component>
</template>
