<script setup lang="ts">
import { useAppearance, useSidebarAppearance } from '@/composables/useAppearance'
import { useModernLayout } from '@/composables/useModernLayout'
import { useTheme } from '@/composables/useTheme'
import { Monitor, Moon, Sun, Palette, LayoutDashboard, Type, TextSelect } from 'lucide-vue-next'

const { appearance, updateAppearance } = useAppearance()
const { sidebarProps, updateSidebarProps } = useSidebarAppearance()
const { currentTheme, typography, themes, updateTheme, updateTypography } = useTheme()
const { isModern, toggleModern } = useModernLayout()

function handleDirChange(value: 'ltr' | 'rtl') {
  updateSidebarProps({ side: value === 'ltr' ? 'left' : 'right' })
}

function handleVariantChange(value: 'sidebar' | 'floating' | 'inset') {
  updateSidebarProps({ variant: value })
}

function handleCollapsibleChange(value: 'offcanvas' | 'none' | 'icon') {
  updateSidebarProps({ collapsible: value })
}

const appearanceTabs = [
  { value: 'light', Icon: Sun, label: 'فاتح' },
  { value: 'dark', Icon: Moon, label: 'غامق' },
  { value: 'system', Icon: Monitor, label: 'تلقائي' },
] as const

const dirTabs = [
  { value: 'ltr', label: 'يسار إلى يمين' },
  { value: 'rtl', label: 'يمين إلى يسار' },
] as const

const variantTabs = [
  { value: 'sidebar', label: 'شريط جانبي' },
  { value: 'floating', label: 'عائم' },
  { value: 'inset', label: 'داخلي' },
] as const

const collapsibleTabs = [
  { value: 'offcanvas', label: 'تغطية' },
  { value: 'none', label: 'ثابت' },
  { value: 'icon', label: 'أيقونات' },
] as const

const layoutStyleTabs = [
  { value: false, label: 'كلاسيكي' },
  { value: true, label: 'عصري' },
] as const

const fontFamilies = [
  { value: 'Cairo', label: 'Cairo (عربي)' },
  { value: 'Changa', label: 'Changa (عصري)' },
  { value: 'Almarai', label: 'Almarai (نظيف)' },
  { value: 'Amiri', label: 'Amiri (كلاسيكي)' },
  { value: 'Reem Kufi', label: 'Reem Kufi (كوفي)' },
  { value: 'Tajawal', label: 'Tajawal (ناعم)' },
  { value: 'Vazirmatn', label: 'Vazirmatn (تقني)' },
  { value: 'IBM Plex Sans Arabic', label: 'IBM Plex (هندسي)' },
  { value: 'Cascadia Code', label: 'Cascadia (برمجة)' },
  { value: 'Inter', label: 'Inter' },
  { value: 'Instrument Sans', label: 'Instrument' },
]

const fontSizes = [
  { value: '0.875rem', label: 'صغير' },
  { value: '1rem', label: 'عادي' },
  { value: '1.125rem', label: 'كبير' },
]

const fontWeights = [
  { value: '300', label: 'خفيف' },
  { value: '400', label: 'عادي' },
  { value: '500', label: 'متوسط' },
  { value: '600', label: 'شبه سميك' },
  { value: '700', label: 'سميك' },
  { value: '800', label: 'عريض جداً' },
]
</script>

<template>
  <div class="space-y-4">
    <!-- Theme Selector -->
    <div class="space-y-2">
      <div class="flex items-center gap-2 text-sm font-medium">
        <Palette class="h-4 w-4" />
        <span>السمة الللونية</span>
      </div>
      <div class="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-4">
        <button v-for="theme in themes" :key="theme.name" @click="updateTheme(theme.name)" :class="[
          'flex items-center justify-center rounded-lg border-2 px-4 py-2.5 text-sm font-medium transition-all',
          currentTheme === theme.name
            ? 'border-primary bg-primary/5 text-primary'
            : 'border-border hover:border-primary/50 hover:bg-accent',
        ]">
          {{ theme.label }}
        </button>
      </div>
    </div>

    <!-- Typography Section -->
    <div class="space-y-4 pt-2 border-t border-border">
      <div class="flex items-center gap-2 text-sm font-medium">
        <Type class="h-4 w-4" />
        <span>الخطوط والتنسيق</span>
      </div>

      <!-- Font Family -->
      <div class="space-y-2">
        <div class="text-xs text-muted-foreground uppercase font-semibold">نوع الخط</div>
        <div class="inline-flex flex-wrap gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
          <button v-for="font in fontFamilies" :key="font.value" @click="updateTypography({ fontFamily: font.value })"
            :class="[
              'flex items-center rounded-md px-3 py-1.5 transition-colors text-sm',
              typography.fontFamily === font.value
                ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100 font-medium'
                : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
            ]">
            {{ font.label }}
          </button>
        </div>
      </div>

      <!-- Font Size -->
      <div class="space-y-2">
        <div class="text-xs text-muted-foreground uppercase font-semibold">حجم الخط</div>
        <div class="inline-flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
          <button v-for="size in fontSizes" :key="size.value" @click="updateTypography({ fontSize: size.value })"
            :class="[
              'flex items-center rounded-md px-3 py-1.5 transition-colors text-sm',
              typography.fontSize === size.value
                ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100 font-medium'
                : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
            ]">
            {{ size.label }}
          </button>
        </div>
      </div>

      <!-- Font Weight -->
      <div class="space-y-2">
        <div class="text-xs text-muted-foreground uppercase font-semibold">سمك الخط</div>
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
          <button v-for="weight in fontWeights" :key="weight.value"
            @click="updateTypography({ fontWeight: weight.value })" :class="[
              'flex items-center justify-center rounded-lg border-2 px-3 py-2 text-sm transition-all',
              typography.fontWeight === weight.value
                ? 'border-primary bg-primary/5 text-primary font-bold'
                : 'border-transparent bg-neutral-100 text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:bg-neutral-800 dark:text-neutral-400 dark:hover:bg-neutral-700/60',
            ]">
            {{ weight.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Appearance Switch -->
    <div class="space-y-2 pt-2 border-t border-border">
      <div class="text-sm font-medium">وضع العرض</div>
      <div class="inline-flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
        <button v-for="{ value, Icon, label } in appearanceTabs" :key="value" @click="updateAppearance(value)" :class="[
          'flex items-center rounded-md px-3.5 py-1.5 transition-colors',
          appearance === value
            ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100'
            : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
        ]">
          <component :is="Icon" class="-ml-1 h-4 w-4" />
          <span class="ml-1.5 text-sm">{{ label }}</span>
        </button>
      </div>
    </div>

    <!-- Sidebar Settings Section -->
    <div class="space-y-4 pt-2 border-t border-border">
      <div class="flex items-center gap-2 text-sm font-medium">
        <TextSelect class="h-4 w-4" />
        <span>إعدادات الشريط الجانبي</span>
      </div>

      <!-- Direction Switch -->
      <div class="space-y-2">
        <div class="text-xs text-muted-foreground uppercase font-semibold">الاتجاه</div>
        <div class="inline-flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
          <button v-for="{ value, label } in dirTabs" :key="value" @click="handleDirChange(value)" :class="[
            'flex items-center rounded-md px-3.5 py-1.5 transition-colors',
            sidebarProps.side === (value === 'ltr' ? 'left' : 'right')
              ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100'
              : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
          ]">
            <span class="text-sm">{{ label }}</span>
          </button>
        </div>
      </div>

      <!-- Variant Switch -->
      <div class="space-y-2">
        <div class="text-xs text-muted-foreground uppercase font-semibold">الشكل</div>
        <div class="inline-flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
          <button v-for="{ value, label } in variantTabs" :key="value" @click="handleVariantChange(value)" :class="[
            'flex items-center rounded-md px-3.5 py-1.5 transition-colors',
            sidebarProps.variant === value
              ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100'
              : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
          ]">
            <span class="text-sm">{{ label }}</span>
          </button>
        </div>
      </div>

      <!-- Collapsible Switch -->
      <div class="space-y-2">
        <div class="text-xs text-muted-foreground uppercase font-semibold">قابلية الطي</div>
        <div class="inline-flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
          <button v-for="{ value, label } in collapsibleTabs" :key="value" @click="handleCollapsibleChange(value)"
            :class="[
              'flex items-center rounded-md px-3.5 py-1.5 transition-colors',
              sidebarProps.collapsible === value
                ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100'
                : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
            ]">
            <span class="text-sm">{{ label }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Layout Style Switch -->
    <div class="space-y-2 pt-2 border-t border-border">
      <div class="flex items-center gap-2 text-sm font-medium">
        <LayoutDashboard class="h-4 w-4" />
        <span>نمط التخطيط</span>
      </div>
      <div class="inline-flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
        <button v-for="{ value, label } in layoutStyleTabs" :key="String(value)" @click="toggleModern(value)" :class="[
          'flex items-center rounded-md px-3.5 py-1.5 transition-colors',
          isModern === value
            ? 'bg-white shadow-xs dark:bg-neutral-700 dark:text-neutral-100'
            : 'text-neutral-500 hover:bg-neutral-200/60 hover:text-black dark:text-neutral-400 dark:hover:bg-neutral-700/60',
        ]">
          <span class="text-sm">{{ label }}</span>
        </button>
      </div>
    </div>
  </div>
</template>
