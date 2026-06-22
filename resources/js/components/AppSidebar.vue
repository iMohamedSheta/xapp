<script setup lang="ts">
import NavFooter from '@/components/NavFooter.vue';
import NavMain from '@/components/NavMain.vue';
import NavUser from '@/components/NavUser.vue';
import { Button } from '@/components/ui/button';
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupLabel, SidebarHeader, SidebarMenu, SidebarMenuButton, SidebarMenuItem, useSidebar } from '@/components/ui/sidebar';
import { Link, usePage } from '@inertiajs/vue3';
import { applySidebarMenuSettings, getSidebarFooterMenu, getSidebarMenu, normalizeSidebarMenuSettings, type SidebarMenuSettings } from '@/lib/sidebar';
import AppLogo from './AppLogo.vue';
import { useModernLayout } from '@/composables/useModernLayout';
import type { NavItem } from '@/types';
import { Eye, EyeOff, GripVertical, Settings2, RotateCcw, Save, X } from 'lucide-vue-next';
import { computed, onMounted, ref } from 'vue';
import AppLogoIcon from './AppLogoIcon.vue';
import { lang } from '@/lib/lang';

const { state } = useSidebar();
const { isModern } = useModernLayout();

const page = usePage();
const authUser = page.props.auth.user;
const role = authUser.role;
const permissions = page.props.auth.permissions;
const baseMainNavItems = computed(() => getSidebarMenu(role, permissions));
const footerNavItems = computed(() => getSidebarFooterMenu(role, permissions));

interface SidebarEditableItem extends NavItem {
  isHidden: boolean;
}

const isEditorOpen = ref(false);
const menuSettings = ref<SidebarMenuSettings>({ order: [], hidden: [] });
const editableItems = ref<SidebarEditableItem[]>([]);

const sidebarSettingsStorageKey = computed(() => `sidebar.menu.settings.${authUser.id}.${role}`);
const mainNavItems = computed(() => applySidebarMenuSettings(baseMainNavItems.value, menuSettings.value));

function buildEditableItems(): SidebarEditableItem[] {
  const baseItems = baseMainNavItems.value;
  const normalized = normalizeSidebarMenuSettings(baseItems, menuSettings.value);
  const byHref = new Map(baseItems.map((item) => [item.href, item]));
  const hiddenSet = new Set(normalized.hidden);
  const orderedItems: NavItem[] = [];

  for (const href of normalized.order) {
    const item = byHref.get(href);
    if (!item) continue;
    orderedItems.push(item);
    byHref.delete(href);
  }

  for (const item of baseItems) {
    if (!byHref.has(item.href)) continue;
    orderedItems.push(item);
  }

  return orderedItems.map((item) => ({
    ...item,
    isHidden: hiddenSet.has(item.href),
  }));
}

function openSidebarEditor(): void {
  editableItems.value = buildEditableItems();
  isEditorOpen.value = true;
}

function closeSidebarEditor(): void {
  isEditorOpen.value = false;
  editableItems.value = [];
}

function toggleItemHidden(index: number): void {
  const item = editableItems.value[index];
  if (!item) return;
  item.isHidden = !item.isHidden;
}

const draggedIndex = ref<number | null>(null);
const dragOverIndex = ref<number | null>(null);

function onDragStart(index: number): void {
  draggedIndex.value = index;
}

function onDragEnter(index: number): void {
  const sourceIndex = draggedIndex.value;
  if (sourceIndex === null || sourceIndex === index) return;

  dragOverIndex.value = index;

  const clone = [...editableItems.value];
  const [movedItem] = clone.splice(sourceIndex, 1);
  clone.splice(index, 0, movedItem);
  editableItems.value = clone;
  draggedIndex.value = index;
}

function onDragLeave(index: number): void {
  if (dragOverIndex.value === index) {
    dragOverIndex.value = null;
  }
}

function onDragEnd(): void {
  draggedIndex.value = null;
  dragOverIndex.value = null;
}

function onDrop(index: number): void {
  draggedIndex.value = null;
  dragOverIndex.value = null;
}

function persistMenuSettings(): void {
  if (typeof window === 'undefined') return;
  window.localStorage.setItem(sidebarSettingsStorageKey.value, JSON.stringify(menuSettings.value));
}

function saveSidebarMenuSettings(): void {
  const order = editableItems.value.map((item) => item.href);
  const hidden = editableItems.value.filter((item) => item.isHidden).map((item) => item.href);
  menuSettings.value = normalizeSidebarMenuSettings(baseMainNavItems.value, { order, hidden });
  persistMenuSettings();
  closeSidebarEditor();
}

function resetSidebarMenuSettings(): void {
  menuSettings.value = { order: [], hidden: [] };
  if (typeof window !== 'undefined') {
    window.localStorage.removeItem(sidebarSettingsStorageKey.value);
  }
  if (isEditorOpen.value) {
    editableItems.value = buildEditableItems();
  }
}

onMounted(() => {
  if (typeof window === 'undefined') return;

  const raw = window.localStorage.getItem(sidebarSettingsStorageKey.value);
  if (!raw) return;

  try {
    const parsed = JSON.parse(raw) as Partial<SidebarMenuSettings>;
    menuSettings.value = normalizeSidebarMenuSettings(baseMainNavItems.value, parsed);
  } catch {
    window.localStorage.removeItem(sidebarSettingsStorageKey.value);
  }
});
</script>

<template>
  <Sidebar collapsible="icon" variant="inset"
    class="sidebar scrollbar scrollbar-thumb-[#353535] scrollbar-track-[#1e1e1e] dark:scrollbar-thumb-[#353535] dark:scrollbar-track-[#1e1e1e]">
    <SidebarHeader :class="isModern ? 'pt-4 px-2' : ''">
      <SidebarMenu>
        <SidebarMenuItem>
          <!-- <NavAgency v-if="role === 'super_admin' || role === 'admin' || role === 'distributor'"></NavAgency> -->
          <div class="flex items-center justify-center gap-3">
            <Link class="data-[state=open]:h-20" href="/">
              <div class="flex items-center justify-center  gap-3  w-full">
                <AppLogoIcon />
              </div>
            </Link>
          </div>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <SidebarGroup v-if="isEditorOpen && state === 'expanded'"
        class="px-2 py-0  scrollbar scrollbar-thumb-[#353535] scrollbar-track-[#1e1e1e] dark:scrollbar-thumb-[#353535] dark:scrollbar-track-[#1e1e1e]">
        <SidebarGroupLabel class="flex items-center justify-between">
          <span>
            {{ lang('sidebar.platform') }}
          </span>
          <div class="flex items-center gap-1">
            <Button size="icon" variant="ghost" class="h-6 w-6" @click.stop="resetSidebarMenuSettings">
              <RotateCcw class="size-4" />
            </Button>
            <Button size="icon" variant="ghost" class="h-6 w-6" @click.stop="closeSidebarEditor">
              <X class="size-4" />
            </Button>
            <Button size="icon" class="h-6 w-6" @click.stop="saveSidebarMenuSettings">
              <Save class="size-4" />
            </Button>
          </div>
        </SidebarGroupLabel>
        <SidebarMenu>
          <SidebarMenuItem v-for="(item, index) in editableItems" :key="item.href" draggable="true" :class="{
            'opacity-50': item.isHidden,
            'transition-all duration-150': true,
            'scale-[0.99]': draggedIndex === index,
          }" @dragstart="onDragStart(index)" @dragenter.prevent="onDragEnter(index)" @dragleave="onDragLeave(index)"
            @dragover.prevent @dragend="onDragEnd" @drop="onDrop(index)">
            <SidebarMenuButton class="cursor-grab active:cursor-grabbing border border-transparent" :class="{
              'bg-muted/60': draggedIndex === index,
              'border-primary bg-primary/10 shadow-sm': dragOverIndex === index,
            }">
              <GripVertical class="size-4 text-muted-foreground" />
              <component :is="item.icon" />
              <span class="ltr:text-left rtl:text-right">{{ item.title }}</span>
              <Button size="icon" variant="ghost" class="ltr:ml-auto rtl:mr-auto h-6 w-6"
                @click.stop="toggleItemHidden(index)">
                <EyeOff v-if="!item.isHidden" class="size-4" />
                <Eye v-else class="size-4" />
              </Button>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroup>
      <NavMain v-else :items="mainNavItems" :show-edit-button="true" @edit-click="openSidebarEditor">
        <template #label-action>
          <Settings2 class="size-4" />
        </template>
      </NavMain>
    </SidebarContent>
    <SidebarFooter>
      <NavFooter :items="footerNavItems" />
      <NavUser />
    </SidebarFooter>
  </Sidebar>
  <slot />
</template>
