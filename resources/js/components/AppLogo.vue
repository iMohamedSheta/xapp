<script setup lang="ts">
import { WifiIcon, ZapIcon } from 'lucide-vue-next';
import { useSidebar } from './ui/sidebar';

const { state } = useSidebar();
const props = defineProps({
  title: String,
});
</script>

<template>
  <div
    class="group flex aspect-square size-10 items-center justify-center rounded-md text-sidebar-primary-foreground cursor-pointer">
    <WifiIcon class="dark:text-yellow-400 text-yellow-600 fill-yellow-600 dark:fill-yellow-400 transition-all duration-300
             drop-shadow-[0_0_12px_rgba(34,197,94,0.8)] 
             animate-electric-float" :class="state === 'collapsed' ? 'size-7' : 'size-10'" />
  </div>
  <div v-if="props.title" class="ml-1 grid flex-1 text-sm">
    <span class="mb-0.5 truncate leading-tight font-semibold text-center">
      {{ props.title }}
    </span>
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