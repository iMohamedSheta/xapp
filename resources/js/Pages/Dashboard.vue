<script setup lang="ts">
import AppLayout from '@/layouts/AppLayout.vue';
import type { BreadcrumbItem } from '@/types'
import { Head } from '@inertiajs/vue3';
import { lang } from '@/lib/lang';

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: lang('dashboard.title'),
    href: '/dashboard',
  },
];

</script>

<template>

  <Head :title="lang('dashboard.title')" />

  <AppLayout :breadcrumbs="breadcrumbs">

    <div class="flex h-full flex-1 flex-col gap-6 rounded-xl p-4 overflow-x-auto">

    </div>
  </AppLayout>
</template>