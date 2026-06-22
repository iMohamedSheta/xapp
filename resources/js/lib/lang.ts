import { ref } from "vue"
import { translations, type Lang } from "@/lang"
import type { Paths } from "type-fest"

type TranslationKey = Paths<typeof translations["en"]>

export const currentLang = ref<Lang>("ar")

export function lang(key: TranslationKey, params?: Record<string, any>): string {
  const parts = key.split(".")
  let result: any = translations[currentLang.value]

  for (const part of parts) {
    result = result?.[part]
    if (result === undefined) break // Early exit if path doesn't exist
  }

  if (typeof result !== 'string') return key

  if (params) {
    return result.replace(/{{(\w+)}}/g, (_match: string, key: string) =>
      params[key]?.toString() ?? _match
    )
  }

  return result
}

export function setLang(lang: Lang) {
  currentLang.value = lang
}
