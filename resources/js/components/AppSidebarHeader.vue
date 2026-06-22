<script setup lang="ts">
import Breadcrumbs from '@/components/Breadcrumbs.vue';
import { SidebarTrigger } from '@/components/ui/sidebar';
import type { BreadcrumbItemType } from '@/types';
import { usePage } from '@inertiajs/vue3';


withDefaults(
  defineProps<{
    breadcrumbs?: BreadcrumbItemType[];
  }>(),
  {
    breadcrumbs: () => [],
  },
);
// Write your TypeScript code here.
const page = usePage()
const isImpersonating = page.props.auth.is_impersonating as boolean;
</script>

<template>
  <header
    class="flex h-12 bg-card! shrink-0 items-center gap-2 border-b border-sidebar-border/20 px-6 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 md:px-4">
    <div class="flex items-center gap-2 flex-1">
      <SidebarTrigger class="-ml-1" />
      <template v-if="breadcrumbs && breadcrumbs.length > 0">
        <Breadcrumbs class="hidden sm:flex" :breadcrumbs="breadcrumbs" />
      </template>
    </div>

    <Search />
    <Notify />
    <LangSwitcher />
    <FullscreenSwitch />
    <BackToImpersonator :isImpersonating="isImpersonating" v-if="isImpersonating" />
  </header>
</template>