<script setup lang="ts">
import AppLayout from '@/layouts/AppLayout.vue';
import type { BreadcrumbItem, Event } from '@/types'
import { Head } from '@inertiajs/vue3';

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: 'NAME',
    href: 'ROUTE',
  },
];

</script>

<template>

  <Head title="NAME" />

  <AppLayout :breadcrumbs="breadcrumbs">

    <div class="flex h-full flex-1 flex-col gap-6 rounded-xl p-4 overflow-x-auto">

    </div>
  </AppLayout>
</template>