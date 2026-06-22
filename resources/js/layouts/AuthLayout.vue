<script setup lang="ts">
import AppDebug from '@/components/debug/AppDebug.vue';
import Toast from '@/components/Toast.vue';
import AuthLayout from '@/layouts/auth/AuthSimpleLayout.vue';
import { usePage } from '@inertiajs/vue3';

defineProps<{
  title?: string;
  description?: string;
}>();

const page = usePage()
</script>

<template>
  <AuthLayout :title="title" :description="description">
    <slot />
    <AppDebug v-if="page.props.debug" />
    <Toast />
  </AuthLayout>
</template>
