import en from "./en"
import ar from "./ar"

export const translations = {
  en,
  ar,
} as const

export type Lang = keyof typeof translations
