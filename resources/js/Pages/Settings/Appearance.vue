<script setup lang="ts">
import { Head } from '@inertiajs/vue3';

import AppearanceTabs from '@/components/AppearanceTabs.vue';
import HeadingSmall from '@/components/HeadingSmall.vue';
import { type BreadcrumbItem } from '@/types';

import AppLayout from '@/layouts/AppLayout.vue';
import SettingsLayout from '@/layouts/settings/Layout.vue';
import { lang } from '@/lib/lang';

const breadcrumbItems: BreadcrumbItem[] = [
  {
    title: lang('settings.appearance.title'),
    href: '/settings/appearance',
  },
];
</script>

<template>
  <AppLayout :breadcrumbs="breadcrumbItems">

    <Head :title="lang('settings.appearance.title')" />
    <SettingsLayout>
      <div class="space-y-6">
        <HeadingSmall :title="lang('settings.appearance.title')"
          :description="lang('settings.appearance.description')" />
        <AppearanceTabs />
      </div>
    </SettingsLayout>
  </AppLayout>
</template>
