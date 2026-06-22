<script setup lang="ts">
import Heading from '@/components/Heading.vue';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { lang } from '@/lib/lang';
import { type NavItem } from '@/types';
import { Link, usePage } from '@inertiajs/vue3';
import { CreditCard, DatabaseIcon, IdCardIcon, Lock, Palette, PhoneIncoming, RefreshCw, ServerIcon, User, WifiHighIcon, WifiIcon } from 'lucide-vue-next';
import { computed } from 'vue';

const props = defineProps<{
  wide?: boolean
}>();

const page = usePage();
const role = page.props.auth.user.role;

const sidebarNavItems = computed<NavItem[]>(() => {
  const items: NavItem[] = [
    {
      title: lang('settings.nav.profile'),
      href: '/settings/profile',
      icon: User,
    },
    {
      title: lang('settings.nav.password'),
      href: '/settings/password',
      icon: Lock,
    },
    {
      title: lang('settings.nav.appearance'),
      href: '/settings/appearance',
      icon: Palette,
    },
  ];

  return items;
});

const currentPath = window.location.pathname;
</script>

<template>
  <div class="px-4 py-6">
    <Heading :title="lang('settings.settings')" :description="lang('settings.description')" />

    <div class="flex flex-col lg:flex-row lg:space-x-12">
      <aside class="w-full max-w-xl lg:w-48">
        <nav class="flex flex-col space-y-1 space-x-0">
          <Button v-for="item in sidebarNavItems" :key="item.href" variant="ghost"
            :class="['w-full justify-start', { 'bg-muted': currentPath === item.href }]" as-child>
            <Link :href="item.href">
              <component :is="item.icon"></component>
              {{ item.title }}
            </Link>
          </Button>
        </nav>
      </aside>

      <Separator class="my-6 lg:hidden" />

      <div :class="['flex-1', props.wide ? 'md:max-w-5xl' : 'md:max-w-2xl']">
        <section :class="['space-y-12', props.wide ? 'w-full' : 'max-w-xl']">
          <slot />
        </section>
      </div>
    </div>
  </div>
</template>
