<script setup lang="ts">
defineProps<{
    variant: 'neutral' | 'success' | 'error' | 'info';
}>();
</script>

<template>
    <div class="flex items-start gap-3 rounded-lg px-3.5 py-3 text-[13px]" :class="{
        'bg-muted/50 border border-border/50': variant === 'neutral',
        'border': true,
        'bg-green-50 border-green-200/70 dark:bg-green-950/30 dark:border-green-800/40': variant === 'success',
        'bg-red-50 border-red-200/70 dark:bg-red-950/30 dark:border-red-800/40': variant === 'error',
        'bg-blue-50 border-blue-200/70 dark:bg-blue-950/30 dark:border-blue-800/40': variant === 'info',
    }">
        <div class="mt-0.5 w-6 h-6 rounded-full flex items-center justify-center shrink-0" :class="{
            'bg-muted text-muted-foreground': variant === 'neutral',
            'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300': variant === 'success',
            'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300': variant === 'error',
            'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300': variant === 'info',
        }">
            <slot name="icon" />
        </div>

        <div class="min-w-0">
            <p class="font-medium leading-snug" :class="{
                'text-foreground': variant === 'neutral',
                'text-green-800 dark:text-green-300': variant === 'success',
                'text-red-800 dark:text-red-300': variant === 'error',
                'text-blue-800 dark:text-blue-300': variant === 'info',
            }">
                <slot name="title" />
            </p>
            <p class="mt-0.5 leading-relaxed" :class="{
                'text-muted-foreground': variant === 'neutral',
                'text-green-700/80 dark:text-green-400/80': variant === 'success',
                'text-red-700/80 dark:text-red-400/80': variant === 'error',
                'text-blue-700/80 dark:text-blue-400/80': variant === 'info',
            }">
                <slot name="desc" />
            </p>
        </div>
    </div>
</template>