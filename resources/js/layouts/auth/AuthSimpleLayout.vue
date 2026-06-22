<script setup lang="ts">
import AppLogoIcon from '@/components/AppLogoIcon.vue';
import { Link } from '@inertiajs/vue3';
import { WifiIcon, ZapIcon } from 'lucide-vue-next';

defineProps<{
  title?: string;
  description?: string;
}>();

</script>

<template>
  <div class="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
    <div class="w-full max-w-sm">
      <div class="flex flex-col gap-8">
        <div class="flex flex-col items-center gap-4">
          <Link href="/home" class="flex flex-col items-center gap-2 font-medium">
            <div class="mb-1 flex h-20 w-20 items-center justify-center rounded-md">
              <div
                class="group flex aspect-square size-20 items-center justify-center rounded-md text-sidebar-primary-foreground cursor-pointer">
                <!-- <WifiIcon class=" size-20 dark:text-yellow-300 text-yellow-500 fill-yellow-500 dark:fill-yellow-300 transition-all duration-300
             drop-shadow-[0_0_12px_rgba(168,85,247,0.8)] 
             animate-electric-float" /> -->
                <AppLogoIcon />
              </div>
            </div>
            <span class=" sr-only">{{ title }}</span>
          </Link>
          <div class="space-y-2 text-center">
            <h1 class="text-xl font-medium">{{ title }}</h1>
            <p class="text-center text-sm text-muted-foreground">{{ description }}</p>
          </div>
        </div>
        <slot />
      </div>
    </div>
  </div>
</template>
<style scoped>
/* Idle float effect */
@keyframes electric-float {

  0%,
  100% {
    transform: translateY(0px);
    filter: drop-shadow(0 0 12px rgba(227, 156, 39, 0.8));
  }

  25% {
    transform: translateY(-1px);
    filter: drop-shadow(0 0 16px rgba(227, 156, 39, 0.9));
  }

  50% {
    transform: translateY(-2px);
    filter: drop-shadow(0 0 20px rgba(227, 156, 39, 1)) drop-shadow(0 0 30px rgba(227, 156, 39, 0.6));
  }

  75% {
    transform: translateY(-1px);
    filter: drop-shadow(0 0 16px rgba(227, 156, 39, 0.9));
  }
}

/* Hover combo: explosion + surge */
@keyframes hover-combo {
  0% {
    transform: scale(1) rotate(0deg);
    filter: drop-shadow(0 0 12px rgba(227, 156, 39, 0.8));
  }

  25% {
    transform: scale(1.3) rotate(-5deg);
    filter: drop-shadow(0 0 25px rgba(227, 156, 39, 1)) drop-shadow(0 0 50px rgba(16, 185, 129, 0.9));
  }

  50% {
    transform: scale(1.6) rotate(10deg);
    filter: drop-shadow(0 0 40px rgba(227, 156, 39, 1)) drop-shadow(0 0 80px rgba(16, 185, 129, 0.8)) brightness(1.5);
  }

  75% {
    transform: scale(1.2) rotate(-3deg);
    filter: drop-shadow(0 0 30px rgba(227, 156, 39, 1)) drop-shadow(0 0 60px rgba(16, 185, 129, 0.6));
  }

  100% {
    transform: scale(1.05) rotate(0deg);
    filter: drop-shadow(0 0 18px rgba(227, 156, 39, 1)) drop-shadow(0 0 25px rgba(16, 185, 129, 0.3));
  }
}

/* Apply animations */
.animate-electric-float {
  animation: electric-float 3s ease-in-out infinite;
}

/* Override idle animation on hover */
.group:hover .animate-electric-float {
  animation: hover-combo 0.9s ease-out, electric-float 3s ease-in-out infinite;
}
</style>
