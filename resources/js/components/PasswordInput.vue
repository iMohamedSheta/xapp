<script setup lang="ts">
import { ref } from "vue"
import { Eye, EyeOff } from "lucide-vue-next"

defineProps<{
  modelValue: string
  label?: string
  placeholder?: string
}>()

defineEmits<{
  (e: "update:modelValue", value: string): void
}>()

const show = ref(false)
</script>

<template>
  <div class="space-y-1">
    <label v-if="label" class="block text-sm font-medium text-gray-700">
      {{ label }}
    </label>
    <div class="relative">
      <input :type="show ? 'text' : 'password'"
        class="w-full rounded-md border border-gray-300 pr-10 pl-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        :placeholder="placeholder" :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)" />
      <button type="button" class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
        @click="show = !show">
        <Eye v-if="!show" class="w-5 h-5" />
        <EyeOff v-else class="w-5 h-5" />
      </button>
    </div>
  </div>
</template>
