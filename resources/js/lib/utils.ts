import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { router } from '@inertiajs/vue3';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function goBack() {
  window.addEventListener('popstate', () => {
    router.reload()
  }, { once: true })

  if (window.history.length > 2) {
    window.history.back()
    return
  }

  router.visit('/dashboard')
}
