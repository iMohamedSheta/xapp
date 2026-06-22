import { computed, watch } from 'vue'
import { currentLang, setLang } from './../lib/lang';
import { type Lang } from "@/lang"
import { useSidebarAppearance } from './useAppearance'

const { updateSidebarProps } = useSidebarAppearance()

export interface Language {
  code: Lang
  name: string
  nativeName: string
  dir?: 'ltr' | 'rtl'
}

const languages: Language[] = [
  { code: 'en', name: 'English', nativeName: 'English', dir: 'ltr' },
  { code: 'ar', name: 'Arabic', nativeName: 'العربية', dir: 'rtl' },
  // { code: 'es', name: 'Spanish', nativeName: 'Español', dir: 'ltr' },
  // { code: 'fr', name: 'French', nativeName: 'Français', dir: 'ltr' },
  // { code: 'de', name: 'German', nativeName: 'Deutsch', dir: 'ltr' },
]

function getStoredLanguage(): Lang | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('app-language') as Lang | null
}

function setStoredLanguage(code: Lang): void {
  if (typeof window === 'undefined') return
  localStorage.setItem('app-language', code)
}

// Initialize from localStorage on load
if (typeof window !== 'undefined') {
  const storedLang = getStoredLanguage()
  if (storedLang && languages.some(l => l.code === storedLang)) {
    setLang(storedLang)

    // Initialize sidebar direction
    const lang = languages.find(l => l.code === storedLang)
    if (lang) {
      const side = lang.dir === 'rtl' ? 'right' : 'left'
      updateSidebarProps({ side })
    }
  }
}

export function useLang() {
  const currentLanguage = computed(() => {
    return languages.find(lang => lang.code === currentLang.value) || languages[0]
  })

  const availableLanguages = computed(() => languages)

  const setLanguage = (code: Lang) => {
    const lang = languages.find(l => l.code === code)
    if (!lang) {
      console.warn(`Language code "${code}" not found`)
      return
    }

    // Update your translation system
    setLang(code)
    setStoredLanguage(code)

    // Update sidebar direction based on language
    const side = lang.dir === 'rtl' ? 'right' : 'left'
    updateSidebarProps({ side })
  }

  const isCurrentLanguage = (code: Lang) => {
    return currentLang.value === code
  }

  // Watch for external language changes and update sidebar accordingly
  watch(currentLang, (newLang) => {
    const lang = languages.find(l => l.code === newLang)
    if (lang && typeof document !== 'undefined') {
      document.documentElement.lang = lang.code
      document.documentElement.dir = lang.dir || 'ltr'

      const side = lang.dir === 'rtl' ? 'right' : 'left'
      updateSidebarProps({ side })
    }
  })

  return {
    currentLanguage,
    currentLang,
    availableLanguages,
    setLanguage,
    isCurrentLanguage,
  }
}