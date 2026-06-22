<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue';
import { goBack } from '@/lib/utils';
import { usePage } from '@inertiajs/vue3';




interface Props {
  message: string;
  error_code: string;
  error_status_code: string;
}
defineProps<Props>();


const page = usePage();
const props = page.props;

const contactSupport = () => {
  window.location.href = "https://chat.whatsapp.com/D20g7KXF9Y557uRzPDa9OW";
}

</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 flex items-center justify-center p-4">
    <div class="max-w-md w-full space-y-8">
      <!-- Error Icon -->
      <div class="text-center">
        <div class="mx-auto h-24 w-24 text-red-500 mb-6">
          <svg class="w-full h-full" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>

        <!-- Error Code -->
        <h1 class="text-6xl font-bold text-slate-900 mb-2">
          {{ props.error_status_code }}
        </h1>

        <!-- Error Title -->
        <h2 class="text-2xl font-semibold text-slate-800 mb-4">
          {{ props.error_code }}
        </h2>

        <!-- Error Description -->
        <p class="text-slate-600 mb-8 leading-relaxed">
          {{ props.message }}
        </p>
      </div>


      <!-- Action Buttons -->
      <div class="space-y-4">
        <!-- Primary Action Button -->
        <Button
          class="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-slate-900 hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-slate-500"
          @click="goBack" color="slate">
          الرجوع
        </Button>
      </div>

      <!-- Additional Help -->
      <div class="text-center pt-6 border-t border-slate-200">
        <p class="text-sm text-gray-800 mb-3">
          لازال لديك مشكلة؟
        </p>
        <Button @click="contactSupport" variant="link" class="text-gray-800">
          التواصل مع الدعم
        </Button>
      </div>
    </div>
  </div>
</template>
