import { computed } from 'vue'
import { currentLang } from './../lib/lang';

export function useDir() {
    const dir = computed(() => {
        return currentLang.value === 'ar' ? 'rtl' : 'ltr'
    })
    return { dir }
}