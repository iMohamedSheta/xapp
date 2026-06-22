<script setup lang="ts">
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubItem,
  useSidebar,
} from '@/components/ui/sidebar';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger
} from '@/components/ui/collapsible';
import { Button } from '@/components/ui/button';
import { ChevronRight } from 'lucide-vue-next';
import type { NavItem } from '@/types';
import { Link, usePage } from '@inertiajs/vue3';
import { openedSidebarMenus } from '@/composables/useAppearance';
import { lang } from '@/lib/lang';

withDefaults(defineProps<{
  items: NavItem[];
  showEditButton?: boolean;
}>(), {
  showEditButton: false,
});

const emit = defineEmits<{
  (event: 'edit-click'): void;
}>();

const page = usePage();
const { state } = useSidebar();
</script>

<template>
  <SidebarGroup
    class="px-2 py-0  scrollbar scrollbar-thumb-[#353535] scrollbar-track-[#1e1e1e] dark:scrollbar-thumb-[#353535] dark:scrollbar-track-[#1e1e1e]">
    <SidebarGroupLabel class="flex items-center justify-between">
      <span>
        {{ lang('sidebar.platform') }}
      </span>
      <Button v-if="showEditButton && state === 'expanded'" size="icon" variant="ghost" class="h-6 w-6"
        @click.stop="emit('edit-click')">
        <slot name="label-action">⚙</slot>
      </Button>
    </SidebarGroupLabel>
    <SidebarMenu>
      <template v-for="item in items" :key="item.title">
        <!-- Simple item -->
        <SidebarMenuItem v-if="!item.children">
          <SidebarMenuButton as-child :is-active="page.url.startsWith(item.activeKey ?? item.href)"
            :tooltip="item.title">
            <Link :href="item.href">
              <component :is="item.icon" />
              <span>{{ item.title }}</span>
            </Link>
          </SidebarMenuButton>
        </SidebarMenuItem>

        <!-- Collapsible group -->
        <Collapsible v-model:open="openedSidebarMenus[item.href]" class="group/collapsible" v-else>
          <SidebarMenuItem>
            <CollapsibleTrigger asChild :tooltip="item.title">
              <SidebarMenuButton>
                <component :is="item.icon" />
                <span>{{ item.title }}</span>
                <ChevronRight
                  class="rtl:mr-auto ltr:ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90" />
              </SidebarMenuButton>
            </CollapsibleTrigger>

            <CollapsibleContent>
              <SidebarMenuSub>
                <SidebarMenuSubItem v-for="child in item.children" :key="child.title">
                  <SidebarMenuButton as-child :is-active="child.href === page.url" :tooltip="child.title">
                    <Link :href="child.href">
                      <component :is="child.icon" />
                      <span dir="ltr">{{ child.title }}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>
      </template>
    </SidebarMenu>
  </SidebarGroup>
</template>
