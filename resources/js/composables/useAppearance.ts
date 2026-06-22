import { SidebarProps } from '@/components/ui/sidebar';
import { onMounted, ref } from 'vue';
import { applyTheme, ThemeName } from './useTheme';

type Appearance = 'light' | 'dark' | 'system';

const defaultTheme = 'golden-white-dark';

export function updateTheme(value: Appearance) {
  if (typeof window === 'undefined') {
    return;
  }

  if (value === 'system') {
    const mediaQueryList = window.matchMedia('(prefers-color-scheme: dark)');
    const systemTheme = mediaQueryList.matches ? 'dark' : 'light';

    document.documentElement.classList.toggle('dark', systemTheme === 'dark');
  } else {
    document.documentElement.classList.toggle('dark', value === 'dark');
  }
}

const setCookie = (name: string, value: string, days = 365) => {
  if (typeof document === 'undefined') {
    return;
  }

  const maxAge = days * 24 * 60 * 60;

  document.cookie = `${name}=${value};path=/;max-age=${maxAge};SameSite=Lax`;
};

const mediaQuery = () => {
  if (typeof window === 'undefined') {
    return null;
  }

  return window.matchMedia('(prefers-color-scheme: dark)');
};

const getStoredAppearance = () => {
  if (typeof window === 'undefined') {
    return null;
  }

  return localStorage.getItem('appearance') as Appearance | null;
};

const handleSystemThemeChange = () => {
  const currentAppearance = getStoredAppearance();

  updateTheme(currentAppearance || 'system');
};

export function initializeTheme() {
  // Initialize color theme FIRST (WhatsApp, default, etc.)
  const storedColorTheme = (localStorage.getItem('app-theme') as ThemeName) || defaultTheme;

  // Initialize appearance (light/dark mode)
  const savedAppearance = getStoredAppearance();
  updateTheme(savedAppearance || 'dark');

  // Apply the color theme with current dark mode state
  const isDark = document.documentElement.classList.contains('dark');
  applyTheme(storedColorTheme, isDark);

  // Set up system theme change listener
  mediaQuery()?.addEventListener('change', handleSystemThemeChange);
}

const appearance = ref<Appearance>('dark');

export function useAppearance() {
  onMounted(() => {
    const savedAppearance = localStorage.getItem('appearance') as Appearance | null;

    if (savedAppearance) {
      appearance.value = savedAppearance;
    }
  });

  function updateAppearance(value: Appearance) {
    appearance.value = value;

    // Store in localStorage for client-side persistence...
    localStorage.setItem('appearance', value);

    // Store in cookie for SSR...
    setCookie('appearance', value);

    updateTheme(value);
  }

  return {
    appearance,
    updateAppearance,
  };
}


export const sidebarProps = ref<SidebarProps>((() => {
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('sidebarProps')
    if (saved) {
      try {
        const parsed = JSON.parse(saved) as Partial<SidebarProps>
        return {
          side: 'right',
          variant: 'inset',
          class: '',
          collapsible: 'icon',
          ...parsed
        }
      } catch { }
    }
  }

  return {
    side: 'right',
    variant: 'inset',
    class: '',
    collapsible: 'icon',
  }
})())

export const openedSidebarMenus = ref<Record<string, boolean>>({})

export function useSidebarAppearance() {
  function updateSidebarProps(props: Partial<SidebarProps>) {
    sidebarProps.value = { ...sidebarProps.value, ...props }
    localStorage.setItem('sidebarProps', JSON.stringify(sidebarProps.value))
  }

  return {
    sidebarProps,
    updateSidebarProps,
  }
}