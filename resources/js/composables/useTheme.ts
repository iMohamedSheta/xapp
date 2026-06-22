// composables/useTheme.ts
import { ref, watch } from 'vue'

export interface TypographySettings {
  fontFamily: string
  fontSize: string
  fontWeight: string
}

export type ThemeName =
  | 'default'
  | 'connect'
  | 'whatsapp'
  | 'electric-blue'
  | 'golden-white-dark'
  | 'modern-luxe'
  | 'aurora-noir'
  | 'warm-stone'
  | 'violet-depths'
  | 'mint-crisp'
  | 'ember-glow'
  | 'steel-sky'

interface Theme {
  name: ThemeName
  label: string
  light: {
    background: string
    foreground: string
    card: string
    cardForeground: string
    popover: string
    popoverForeground: string
    primary: string
    primaryForeground: string
    secondary: string
    secondaryForeground: string
    muted: string
    mutedForeground: string
    accent: string
    accentForeground: string
    destructive: string
    destructiveForeground: string
    border: string
    input: string
    ring: string
    chart1: string
    chart2: string
    chart3: string
    chart4: string
    chart5: string
    sidebarBackground: string
    sidebarForeground: string
    sidebarPrimary: string
    sidebarPrimaryForeground: string
    sidebarAccent: string
    sidebarAccentForeground: string
    sidebarBorder: string
    sidebarRing: string
  }
  dark: {
    background: string
    foreground: string
    card: string
    cardForeground: string
    popover: string
    popoverForeground: string
    primary: string
    primaryForeground: string
    secondary: string
    secondaryForeground: string
    muted: string
    mutedForeground: string
    accent: string
    accentForeground: string
    destructive: string
    destructiveForeground: string
    border: string
    input: string
    ring: string
    chart1: string
    chart2: string
    chart3: string
    chart4: string
    chart5: string
    sidebarBackground: string
    sidebarForeground: string
    sidebarPrimary: string
    sidebarPrimaryForeground: string
    sidebarAccent: string
    sidebarAccentForeground: string
    sidebarBorder: string
    sidebarRing: string
  }
}

const themes: Theme[] = [
  {
    name: 'connect',
    label: 'Connect',
    light: {
      background: 'hsl(210 30% 98%)',           // softer, cleaner white with blue tint
      foreground: 'hsl(210 40% 12%)',           // deeper, more readable dark text
      card: 'hsl(0 0% 100%)',                   // pure white cards for contrast
      cardForeground: 'hsl(210 40% 12%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(210 40% 12%)',
      primary: 'hsl(210 55% 42%)',              // richer, more saturated blue
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(210 25% 95%)',            // subtle secondary
      secondaryForeground: 'hsl(210 50% 28%)',
      muted: 'hsl(210 25% 96%)',                // soft muted backgrounds
      mutedForeground: 'hsl(210 15% 48%)',      // better contrast for muted text
      accent: 'hsl(210 40% 94%)',               // more distinct accent
      accentForeground: 'hsl(210 55% 38%)',
      destructive: 'hsl(0 85% 58%)',            // vibrant red
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(210 20% 88%)',               // softer borders
      input: 'hsl(210 25% 96%)',                // subtle input background
      ring: 'hsl(210 55% 42%)',                 // matches primary
      chart1: 'hsl(210 60% 45%)',               // vibrant chart colors
      chart2: 'hsl(195 65% 50%)',
      chart3: 'hsl(185 70% 52%)',
      chart4: 'hsl(220 55% 52%)',
      chart5: 'hsl(175 60% 48%)',
      sidebarBackground: 'hsl(0 0% 99%)',       // ultra-light sidebar
      sidebarForeground: 'hsl(210 35% 18%)',
      sidebarPrimary: 'hsl(210 55% 42%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(210 35% 96%)',
      sidebarAccentForeground: 'hsl(210 55% 38%)',
      sidebarBorder: 'hsl(210 20% 92%)',
      sidebarRing: 'hsl(210 55% 42%)',
    },
    dark: {
      background: 'hsl(210 35% 8%)',            // deeper, richer dark blue
      foreground: 'hsl(210 30% 92%)',           // softer light text for less eye strain
      card: 'hsl(210 32% 11%)',                 // elevated card surface
      cardForeground: 'hsl(210 30% 92%)',
      popover: 'hsl(210 32% 11%)',
      popoverForeground: 'hsl(210 30% 92%)',
      primary: 'hsl(210 60% 52%)',              // brighter, more vivid primary
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(210 28% 18%)',            // better contrast secondary
      secondaryForeground: 'hsl(210 35% 78%)',
      muted: 'hsl(210 30% 14%)',                // subtle muted surfaces
      mutedForeground: 'hsl(210 15% 62%)',      // readable muted text
      accent: 'hsl(210 32% 20%)',               // distinct accent areas
      accentForeground: 'hsl(210 40% 82%)',
      destructive: 'hsl(0 75% 55%)',            // balanced destructive
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(210 28% 20%)',               // visible but subtle borders
      input: 'hsl(210 30% 18%)',                // clear input fields
      ring: 'hsl(210 60% 52%)',                 // matches primary
      chart1: 'hsl(210 65% 55%)',               // bright, distinguishable charts
      chart2: 'hsl(195 70% 58%)',
      chart3: 'hsl(185 75% 60%)',
      chart4: 'hsl(220 60% 65%)',
      chart5: 'hsl(175 65% 55%)',
      sidebarBackground: 'hsl(210 38% 10%)',    // distinct sidebar
      sidebarForeground: 'hsl(210 30% 90%)',
      sidebarPrimary: 'hsl(210 60% 50%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(210 32% 18%)',
      sidebarAccentForeground: 'hsl(210 40% 85%)',
      sidebarBorder: 'hsl(210 30% 18%)',
      sidebarRing: 'hsl(210 60% 52%)',
    },
  },
  {
    name: 'default',
    label: 'Default',
    light: {
      background: 'hsl(240 5% 96%)',
      foreground: 'hsl(0 0% 3.9%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(0 0% 3.9%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(0 0% 3.9%)',
      primary: 'hsl(0 0% 9%)',
      primaryForeground: 'hsl(0 0% 98%)',
      secondary: 'hsl(0 0% 92.1%)',
      secondaryForeground: 'hsl(0 0% 9%)',
      muted: 'hsl(0 0% 96.1%)',
      mutedForeground: 'hsl(0 0% 45.1%)',
      accent: 'hsl(0 0% 96.1%)',
      accentForeground: 'hsl(0 0% 9%)',
      destructive: 'hsl(0 84.2% 60.2%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(0 0% 90%)',
      input: 'hsl(0 0% 88%)',
      ring: 'hsl(0 0% 3.9%)',
      chart1: 'hsl(12 76% 61%)',
      chart2: 'hsl(173 58% 39%)',
      chart3: 'hsl(197 37% 24%)',
      chart4: 'hsl(43 74% 66%)',
      chart5: 'hsl(27 87% 67%)',
      sidebarBackground: 'hsl(240 10% 12%)',
      sidebarForeground: 'hsl(0 0% 98%)',
      sidebarPrimary: 'hsl(0 0% 98%)',
      sidebarPrimaryForeground: 'hsl(240 10% 12%)',
      sidebarAccent: 'hsl(240 10% 20%)',
      sidebarAccentForeground: 'hsl(0 0% 98%)',
      sidebarBorder: 'hsl(240 10% 18%)',
      sidebarRing: 'hsl(217.2 91.2% 59.8%)',
    },
    dark: {
      background: 'hsl(0 0% 2%)',
      foreground: 'hsl(0 0% 98%)',
      card: 'hsl(0 0% 4.5%)',
      cardForeground: 'hsl(0 0% 98%)',
      popover: 'hsl(0 0% 4.5%)',
      popoverForeground: 'hsl(0 0% 98%)',
      primary: 'hsl(0 0% 98%)',
      primaryForeground: 'hsl(0 0% 9%)',
      secondary: 'hsl(0 0% 14.9%)',
      secondaryForeground: 'hsl(0 0% 98%)',
      muted: 'hsl(0 0% 16.08%)',
      mutedForeground: 'hsl(0 0% 63.9%)',
      accent: 'hsl(0 0% 14.9%)',
      accentForeground: 'hsl(0 0% 98%)',
      destructive: 'hsl(0 84% 60%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(0 0% 14.9%)',
      input: 'hsl(0 0% 14.9%)',
      ring: 'hsl(0 0% 83.1%)',
      chart1: 'hsl(220 70% 50%)',
      chart2: 'hsl(160 60% 45%)',
      chart3: 'hsl(30 80% 55%)',
      chart4: 'hsl(280 65% 60%)',
      chart5: 'hsl(340 75% 55%)',
      sidebarBackground: 'hsl(0 0% 3%)',
      sidebarForeground: 'hsl(0 0% 95.9%)',
      sidebarPrimary: 'hsl(360, 100%, 100%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(0 0% 15.9%)',
      sidebarAccentForeground: 'hsl(240 4.8% 95.9%)',
      sidebarBorder: 'hsl(0 0% 15.9%)',
      sidebarRing: 'hsl(217.2 91.2% 59.8%)',
    },
  },
  {
    name: 'electric-blue',
    label: 'Electric Blue',
    light: {
      background: 'hsl(0 0% 100%)',
      foreground: 'hsl(0 0% 3.9%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(0 0% 3.9%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(0 0% 3.9%)',
      primary: 'hsl(217 91% 60%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(0 0% 92.1%)',
      secondaryForeground: 'hsl(0 0% 9%)',
      muted: 'hsl(0 0% 96.1%)',
      mutedForeground: 'hsl(0 0% 45.1%)',
      accent: 'hsl(217 91% 95%)',
      accentForeground: 'hsl(217 91% 40%)',
      destructive: 'hsl(0 84.2% 60.2%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(0 0% 92.8%)',
      input: 'hsl(0 0% 89.8%)',
      ring: 'hsl(217 91% 60%)',
      chart1: 'hsl(217 91% 60%)',
      chart2: 'hsl(217 91% 45%)',
      chart3: 'hsl(217 91% 70%)',
      chart4: 'hsl(217 50% 50%)',
      chart5: 'hsl(217 30% 60%)',
      sidebarBackground: 'hsl(0 0% 98%)',
      sidebarForeground: 'hsl(240 5.3% 26.1%)',
      sidebarPrimary: 'hsl(217 91% 60%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(217 91% 96%)',
      sidebarAccentForeground: 'hsl(217 91% 35%)',
      sidebarBorder: 'hsl(0 0% 91%)',
      sidebarRing: 'hsl(217 91% 60%)',
    },
    dark: {
      background: 'hsl(222 20% 14%)',
      foreground: 'hsl(0 0% 98%)',
      card: 'hsl(222 20% 16%)',
      cardForeground: 'hsl(0 0% 98%)',
      popover: 'hsl(222 20% 16%)',
      popoverForeground: 'hsl(0 0% 98%)',
      primary: 'hsl(217 91% 60%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(222 20% 22%)',
      secondaryForeground: 'hsl(0 0% 98%)',
      muted: 'hsl(222 20% 20%)',
      mutedForeground: 'hsl(220 10% 65%)',
      accent: 'hsl(222 20% 22%)',
      accentForeground: 'hsl(217 91% 70%)',
      destructive: 'hsl(0 84% 60%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(222 20% 24%)',
      input: 'hsl(222 20% 22%)',
      ring: 'hsl(217 91% 60%)',
      chart1: 'hsl(217 91% 60%)',
      chart2: 'hsl(217 91% 50%)',
      chart3: 'hsl(217 91% 70%)',
      chart4: 'hsl(217 60% 55%)',
      chart5: 'hsl(217 40% 65%)',
      sidebarBackground: 'hsl(222 20% 12%)',
      sidebarForeground: 'hsl(0 0% 90%)',
      sidebarPrimary: 'hsl(217 91% 60%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(222 20% 20%)',
      sidebarAccentForeground: 'hsl(217 91% 70%)',
      sidebarBorder: 'hsl(222 20% 22%)',
      sidebarRing: 'hsl(217 91% 60%)',
    },
  },
  {
    name: 'golden-white-dark',
    label: 'Golden Dark',
    light: {
      background: 'hsl(210 20% 92%)',
      foreground: 'hsl(224 71.4% 4.1%)',
      card: 'hsl(210 20% 96%)',
      cardForeground: 'hsl(240 10% 3.9%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(240 10% 3.9%)',
      primary: 'hsl(240 5.9% 10%)',
      primaryForeground: 'hsl(0 0% 98%)',
      secondary: 'hsl(240 4.8% 95.9%)',
      secondaryForeground: 'hsl(240 5.9% 10%)',
      muted: 'hsl(240 4.8% 95.9%)',
      mutedForeground: 'hsl(240 3.8% 46.1%)',
      accent: 'hsl(240 4.8% 95.9%)',
      accentForeground: 'hsl(240 5.9% 10%)',
      destructive: 'hsl(0 84.2% 60.2%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(240 5.9% 88%)',
      input: 'hsl(240 5.9% 85%)',
      ring: 'hsl(240 5.9% 10%)',
      chart1: 'hsl(240 5.9% 10%)',
      chart2: 'hsl(240 5.9% 20%)',
      chart3: 'hsl(240 5.9% 30%)',
      chart4: 'hsl(240 5.9% 40%)',
      chart5: 'hsl(240 5.9% 50%)',
      sidebarBackground: 'hsl(222 47% 8%)',
      sidebarForeground: 'hsl(210 40% 98%)',
      sidebarPrimary: 'hsl(210 40% 98%)',
      sidebarPrimaryForeground: 'hsl(222 47% 8%)',
      sidebarAccent: 'hsl(222 47% 15%)',
      sidebarAccentForeground: 'hsl(210 40% 98%)',
      sidebarBorder: 'hsl(222 47% 18%)',
      sidebarRing: 'hsl(210 20% 70%)',
    },
    dark: {
      background: 'hsl(240 10% 8%)',
      foreground: 'hsl(0 0% 98%)',
      card: 'hsl(240 10% 10%)',
      cardForeground: 'hsl(0 0% 98%)',
      popover: 'hsl(240 10% 10%)',
      popoverForeground: 'hsl(0 0% 98%)',
      primary: '#e7b733ff',
      primaryForeground: 'hsl(240, 7%, 6%)',
      secondary: 'hsl(240 10% 22%)',
      secondaryForeground: 'hsl(0 0% 98%)',
      muted: 'hsl(240 10% 15%)',
      mutedForeground: 'hsl(240 5% 65%)',
      accent: 'hsl(240 10% 20%)',
      accentForeground: 'hsl(45 85% 52%)',
      destructive: 'hsl(0 84% 60%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(240 10% 20%)',
      input: 'hsl(240 10% 18%)',
      ring: 'hsl(45 85% 45%)',
      chart1: 'hsl(45 85% 45%)',
      chart2: 'hsl(45 85% 40%)',
      chart3: 'hsl(45 80% 50%)',
      chart4: 'hsl(45 75% 42%)',
      chart5: 'hsl(45 70% 48%)',
      sidebarBackground: 'hsl(240 10% 6%)',
      sidebarForeground: 'hsl(0 0% 90%)',
      sidebarPrimary: '#ecbf43ff',
      sidebarPrimaryForeground: 'hsl(240 10% 8%)',
      sidebarAccent: 'hsl(240 10% 14%)',
      sidebarAccentForeground: '#ecbf43ff',
      sidebarBorder: 'hsl(240 10% 18%)',
      sidebarRing: 'hsl(45 85% 45%)',
    },
  },
  {
    name: 'whatsapp',
    label: 'WhatsApp',
    light: {
      background: 'hsl(140 20% 96%)',
      foreground: 'hsl(140 60% 8%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(0 0% 13%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(0 0% 13%)',
      primary: 'hsl(142 70% 40%)', // Deeper WhatsApp green
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(142 20% 94%)',
      secondaryForeground: 'hsl(142 70% 25%)',
      muted: 'hsl(0 0% 96%)',
      mutedForeground: 'hsl(0 0% 45%)',
      accent: 'hsl(142 50% 92%)',
      accentForeground: 'hsl(142 70% 30%)',
      destructive: 'hsl(0 84% 60%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(0 0% 90%)',
      input: 'hsl(0 0% 94%)',
      ring: 'hsl(142 70% 40%)',
      chart1: 'hsl(142 70% 40%)',
      chart2: 'hsl(175 60% 45%)',
      chart3: 'hsl(195 70% 50%)',
      chart4: 'hsl(142 50% 35%)',
      chart5: 'hsl(160 65% 40%)',
      sidebarBackground: 'hsl(0 0% 98%)',
      sidebarForeground: 'hsl(0 0% 20%)',
      sidebarPrimary: 'hsl(142 70% 40%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(142 30% 95%)',
      sidebarAccentForeground: 'hsl(142 70% 30%)',
      sidebarBorder: 'hsl(0 0% 92%)',
      sidebarRing: 'hsl(142 70% 40%)',
    },
    dark: {
      background: 'hsl(0 0% 7%)',
      foreground: 'hsl(0 0% 88%)',
      card: 'hsl(0 0% 11%)',
      cardForeground: 'hsl(0 0% 88%)',
      popover: 'hsl(0 0% 11%)',
      popoverForeground: 'hsl(0 0% 88%)',
      primary: 'hsl(142 70% 40%)', // Same deeper green
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(0 0% 15%)',
      secondaryForeground: 'hsl(142 50% 70%)',
      muted: 'hsl(0 0% 13%)',
      mutedForeground: 'hsl(0 0% 60%)',
      accent: 'hsl(0 0% 16%)',
      accentForeground: 'hsl(142 50% 70%)',
      destructive: 'hsl(0 72% 51%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(0 0% 18%)',
      input: 'hsl(0 0% 16%)',
      ring: 'hsl(142 70% 40%)',
      chart1: 'hsl(142 70% 40%)',
      chart2: 'hsl(175 60% 45%)',
      chart3: 'hsl(195 60% 50%)',
      chart4: 'hsl(142 50% 60%)',
      chart5: 'hsl(160 55% 50%)',
      sidebarBackground: 'hsl(0 0% 9%)',
      sidebarForeground: 'hsl(0 0% 88%)',
      sidebarPrimary: 'hsl(142 70% 40%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(0 0% 16%)',
      sidebarAccentForeground: 'hsl(142 50% 70%)',
      sidebarBorder: 'hsl(0 0% 16%)',
      sidebarRing: 'hsl(142 70% 40%)',
    },
  },
  {
    name: 'modern-luxe',
    label: 'Modern Luxe',
    light: {
      background: 'hsl(210 20% 92%)',
      foreground: 'hsl(224 71.4% 4.1%)',
      card: 'hsl(210 20% 96%)',
      cardForeground: 'hsl(240 10% 3.9%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(240 10% 3.9%)',
      primary: 'hsl(240 5.9% 10%)',
      primaryForeground: 'hsl(0 0% 98%)',
      secondary: 'hsl(240 4.8% 95.9%)',
      secondaryForeground: 'hsl(240 5.9% 10%)',
      muted: 'hsl(240 4.8% 95.9%)',
      mutedForeground: 'hsl(240 3.8% 46.1%)',
      accent: 'hsl(240 4.8% 95.9%)',
      accentForeground: 'hsl(240 5.9% 10%)',
      destructive: 'hsl(0 84.2% 60.2%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(240 5.9% 88%)',
      input: 'hsl(240 5.9% 85%)',
      ring: 'hsl(240 5.9% 10%)',
      chart1: 'hsl(240 5.9% 10%)',
      chart2: 'hsl(240 5.9% 20%)',
      chart3: 'hsl(240 5.9% 30%)',
      chart4: 'hsl(240 5.9% 40%)',
      chart5: 'hsl(240 5.9% 50%)',
      sidebarBackground: 'hsl(222 47% 8%)',
      sidebarForeground: 'hsl(210 40% 98%)',
      sidebarPrimary: 'hsl(210 40% 98%)',
      sidebarPrimaryForeground: 'hsl(222 47% 8%)',
      sidebarAccent: 'hsl(222 47% 15%)',
      sidebarAccentForeground: 'hsl(210 40% 98%)',
      sidebarBorder: 'hsl(222 47% 18%)',
      sidebarRing: 'hsl(210 20% 70%)',
    },
    dark: {
      background: 'hsl(222 47% 2%)',
      foreground: 'hsl(210 40% 98%)',
      card: 'hsl(222 47% 5%)',
      cardForeground: 'hsl(210 40% 98%)',
      popover: 'hsl(222 47% 5%)',
      popoverForeground: 'hsl(210 40% 98%)',
      primary: 'hsl(217 91% 64%)',
      primaryForeground: 'hsl(222 47% 6%)',
      secondary: 'hsl(222 47% 12%)',
      secondaryForeground: 'hsl(210 40% 98%)',
      muted: 'hsl(222 47% 10%)',
      mutedForeground: 'hsl(210 20% 65%)',
      accent: 'hsl(217 33% 15%)',
      accentForeground: 'hsl(210 40% 98%)',
      destructive: 'hsl(0 80% 60%)',
      destructiveForeground: 'hsl(210 40% 98%)',
      border: 'hsl(222 47% 12%)',
      input: 'hsl(222 47% 12%)',
      ring: 'hsl(217 91% 64%)',
      chart1: 'hsl(217 91% 64%)',
      chart2: 'hsl(199 89% 55%)',
      chart3: 'hsl(168 83% 50%)',
      chart4: 'hsl(35 91% 55%)',
      chart5: 'hsl(0 72% 56%)',
      sidebarBackground: 'hsl(222 47% 1%)',
      sidebarForeground: 'hsl(210 40% 90%)',
      sidebarPrimary: 'hsl(217 91% 64%)',
      sidebarPrimaryForeground: 'hsl(222 47% 6%)',
      sidebarAccent: 'hsl(222 47% 10%)',
      sidebarAccentForeground: 'hsl(217 91% 80%)',
      sidebarBorder: 'hsl(222 47% 8%)',
      sidebarRing: 'hsl(217 91% 64%)',
    },
  },
  {
    name: 'aurora-noir',
    label: 'Aurora Noir',
    light: {
      background: 'hsl(220 20% 97%)',
      foreground: 'hsl(220 30% 10%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(220 30% 10%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(220 30% 10%)',
      primary: 'hsl(188 90% 42%)',
      primaryForeground: 'hsl(220 30% 6%)',
      secondary: 'hsl(220 15% 94%)',
      secondaryForeground: 'hsl(220 30% 20%)',
      muted: 'hsl(220 15% 95%)',
      mutedForeground: 'hsl(220 10% 48%)',
      accent: 'hsl(270 60% 94%)',
      accentForeground: 'hsl(270 60% 40%)',
      destructive: 'hsl(0 85% 58%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(220 15% 88%)',
      input: 'hsl(220 15% 92%)',
      ring: 'hsl(188 90% 42%)',
      chart1: 'hsl(188 90% 42%)',
      chart2: 'hsl(270 70% 60%)',
      chart3: 'hsl(330 80% 62%)',
      chart4: 'hsl(210 80% 55%)',
      chart5: 'hsl(155 70% 45%)',
      sidebarBackground: 'hsl(222 40% 12%)',
      sidebarForeground: 'hsl(210 30% 88%)',
      sidebarPrimary: 'hsl(188 90% 45%)',
      sidebarPrimaryForeground: 'hsl(222 40% 8%)',
      sidebarAccent: 'hsl(222 40% 20%)',
      sidebarAccentForeground: 'hsl(188 80% 72%)',
      sidebarBorder: 'hsl(222 40% 20%)',
      sidebarRing: 'hsl(188 90% 42%)',
    },
    dark: {
      background: 'hsl(222 35% 7%)',
      foreground: 'hsl(210 30% 90%)',
      card: 'hsl(222 35% 10%)',
      cardForeground: 'hsl(210 30% 90%)',
      popover: 'hsl(222 35% 10%)',
      popoverForeground: 'hsl(210 30% 90%)',
      primary: 'hsl(188 90% 48%)',
      primaryForeground: 'hsl(222 35% 5%)',
      secondary: 'hsl(222 30% 16%)',
      secondaryForeground: 'hsl(210 30% 80%)',
      muted: 'hsl(222 30% 13%)',
      mutedForeground: 'hsl(210 15% 55%)',
      accent: 'hsl(270 40% 18%)',
      accentForeground: 'hsl(270 80% 78%)',
      destructive: 'hsl(0 75% 55%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(222 30% 18%)',
      input: 'hsl(222 30% 15%)',
      ring: 'hsl(188 90% 48%)',
      chart1: 'hsl(188 90% 52%)',
      chart2: 'hsl(270 75% 65%)',
      chart3: 'hsl(330 80% 65%)',
      chart4: 'hsl(210 80% 60%)',
      chart5: 'hsl(155 70% 50%)',
      sidebarBackground: 'hsl(222 40% 5%)',
      sidebarForeground: 'hsl(210 30% 88%)',
      sidebarPrimary: 'hsl(188 90% 48%)',
      sidebarPrimaryForeground: 'hsl(222 40% 5%)',
      sidebarAccent: 'hsl(222 35% 14%)',
      sidebarAccentForeground: 'hsl(188 80% 72%)',
      sidebarBorder: 'hsl(222 35% 14%)',
      sidebarRing: 'hsl(188 90% 48%)',
    },
  },
  {
    name: 'warm-stone',
    label: 'Warm Stone',
    light: {
      background: 'hsl(40 15% 97%)',
      foreground: 'hsl(20 20% 10%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(20 20% 10%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(20 20% 10%)',
      primary: 'hsl(32 90% 42%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(35 15% 92%)',
      secondaryForeground: 'hsl(20 20% 20%)',
      muted: 'hsl(35 15% 94%)',
      mutedForeground: 'hsl(30 8% 46%)',
      accent: 'hsl(32 80% 94%)',
      accentForeground: 'hsl(32 80% 30%)',
      destructive: 'hsl(0 84% 60%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(35 12% 86%)',
      input: 'hsl(35 12% 91%)',
      ring: 'hsl(32 90% 42%)',
      chart1: 'hsl(32 90% 42%)',
      chart2: 'hsl(20 80% 50%)',
      chart3: 'hsl(45 85% 50%)',
      chart4: 'hsl(10 70% 55%)',
      chart5: 'hsl(55 80% 48%)',
      sidebarBackground: 'hsl(35 15% 94%)',
      sidebarForeground: 'hsl(20 20% 15%)',
      sidebarPrimary: 'hsl(32 90% 42%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(35 20% 88%)',
      sidebarAccentForeground: 'hsl(32 80% 30%)',
      sidebarBorder: 'hsl(35 12% 86%)',
      sidebarRing: 'hsl(32 90% 42%)',
    },
    dark: {
      background: 'hsl(25 15% 9%)',
      foreground: 'hsl(35 15% 90%)',
      card: 'hsl(25 15% 12%)',
      cardForeground: 'hsl(35 15% 90%)',
      popover: 'hsl(25 15% 12%)',
      popoverForeground: 'hsl(35 15% 90%)',
      primary: 'hsl(32 90% 50%)',
      primaryForeground: 'hsl(25 15% 6%)',
      secondary: 'hsl(25 15% 18%)',
      secondaryForeground: 'hsl(35 15% 82%)',
      muted: 'hsl(25 15% 14%)',
      mutedForeground: 'hsl(30 8% 55%)',
      accent: 'hsl(32 50% 18%)',
      accentForeground: 'hsl(32 90% 72%)',
      destructive: 'hsl(0 75% 55%)',
      destructiveForeground: 'hsl(0 0% 98%)',
      border: 'hsl(25 15% 20%)',
      input: 'hsl(25 15% 17%)',
      ring: 'hsl(32 90% 50%)',
      chart1: 'hsl(32 90% 55%)',
      chart2: 'hsl(20 80% 58%)',
      chart3: 'hsl(45 85% 55%)',
      chart4: 'hsl(10 70% 60%)',
      chart5: 'hsl(55 80% 52%)',
      sidebarBackground: 'hsl(25 15% 7%)',
      sidebarForeground: 'hsl(35 15% 88%)',
      sidebarPrimary: 'hsl(32 90% 50%)',
      sidebarPrimaryForeground: 'hsl(25 15% 6%)',
      sidebarAccent: 'hsl(25 15% 15%)',
      sidebarAccentForeground: 'hsl(32 90% 70%)',
      sidebarBorder: 'hsl(25 15% 18%)',
      sidebarRing: 'hsl(32 90% 50%)',
    },
  },
  {
    name: 'violet-depths',
    label: 'Violet Depths',
    light: {
      background: 'hsl(260 20% 97%)',
      foreground: 'hsl(260 30% 8%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(260 30% 8%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(260 30% 8%)',
      primary: 'hsl(262 80% 55%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(260 15% 93%)',
      secondaryForeground: 'hsl(260 30% 20%)',
      muted: 'hsl(260 15% 95%)',
      mutedForeground: 'hsl(260 10% 48%)',
      accent: 'hsl(330 70% 93%)',
      accentForeground: 'hsl(330 70% 38%)',
      destructive: 'hsl(0 84% 58%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(260 15% 88%)',
      input: 'hsl(260 15% 92%)',
      ring: 'hsl(262 80% 55%)',
      chart1: 'hsl(262 80% 55%)',
      chart2: 'hsl(330 80% 60%)',
      chart3: 'hsl(196 80% 50%)',
      chart4: 'hsl(290 70% 58%)',
      chart5: 'hsl(220 75% 55%)',
      sidebarBackground: 'hsl(260 30% 10%)',
      sidebarForeground: 'hsl(260 20% 88%)',
      sidebarPrimary: 'hsl(262 80% 68%)',
      sidebarPrimaryForeground: 'hsl(260 30% 6%)',
      sidebarAccent: 'hsl(260 30% 18%)',
      sidebarAccentForeground: 'hsl(262 80% 80%)',
      sidebarBorder: 'hsl(260 30% 18%)',
      sidebarRing: 'hsl(262 80% 55%)',
    },
    dark: {
      background: 'hsl(260 25% 6%)',
      foreground: 'hsl(260 20% 90%)',
      card: 'hsl(260 25% 9%)',
      cardForeground: 'hsl(260 20% 90%)',
      popover: 'hsl(260 25% 9%)',
      popoverForeground: 'hsl(260 20% 90%)',
      primary: 'hsl(262 80% 62%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(260 25% 16%)',
      secondaryForeground: 'hsl(260 20% 82%)',
      muted: 'hsl(260 25% 12%)',
      mutedForeground: 'hsl(260 10% 54%)',
      accent: 'hsl(330 40% 16%)',
      accentForeground: 'hsl(330 80% 75%)',
      destructive: 'hsl(0 72% 52%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(260 25% 18%)',
      input: 'hsl(260 25% 14%)',
      ring: 'hsl(262 80% 62%)',
      chart1: 'hsl(262 80% 65%)',
      chart2: 'hsl(330 80% 65%)',
      chart3: 'hsl(196 80% 55%)',
      chart4: 'hsl(290 70% 62%)',
      chart5: 'hsl(220 75% 60%)',
      sidebarBackground: 'hsl(260 30% 4%)',
      sidebarForeground: 'hsl(260 20% 88%)',
      sidebarPrimary: 'hsl(262 80% 62%)',
      sidebarPrimaryForeground: 'hsl(260 30% 4%)',
      sidebarAccent: 'hsl(260 25% 13%)',
      sidebarAccentForeground: 'hsl(262 80% 80%)',
      sidebarBorder: 'hsl(260 25% 14%)',
      sidebarRing: 'hsl(262 80% 62%)',
    },
  },
  {
    name: 'mint-crisp',
    label: 'Mint Crisp',
    light: {
      background: 'hsl(150 20% 97%)',
      foreground: 'hsl(150 30% 8%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(150 30% 8%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(150 30% 8%)',
      primary: 'hsl(160 75% 35%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(150 15% 93%)',
      secondaryForeground: 'hsl(150 30% 18%)',
      muted: 'hsl(150 15% 95%)',
      mutedForeground: 'hsl(150 8% 46%)',
      accent: 'hsl(160 60% 92%)',
      accentForeground: 'hsl(160 75% 28%)',
      destructive: 'hsl(0 84% 58%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(150 12% 87%)',
      input: 'hsl(150 12% 92%)',
      ring: 'hsl(160 75% 35%)',
      chart1: 'hsl(160 75% 35%)',
      chart2: 'hsl(185 70% 42%)',
      chart3: 'hsl(140 65% 38%)',
      chart4: 'hsl(200 75% 45%)',
      chart5: 'hsl(130 60% 40%)',
      sidebarBackground: 'hsl(150 15% 94%)',
      sidebarForeground: 'hsl(150 25% 12%)',
      sidebarPrimary: 'hsl(160 75% 35%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(150 20% 88%)',
      sidebarAccentForeground: 'hsl(160 75% 28%)',
      sidebarBorder: 'hsl(150 12% 86%)',
      sidebarRing: 'hsl(160 75% 35%)',
    },
    dark: {
      background: 'hsl(155 20% 7%)',
      foreground: 'hsl(150 20% 90%)',
      card: 'hsl(155 20% 10%)',
      cardForeground: 'hsl(150 20% 90%)',
      popover: 'hsl(155 20% 10%)',
      popoverForeground: 'hsl(150 20% 90%)',
      primary: 'hsl(160 70% 42%)',
      primaryForeground: 'hsl(155 20% 5%)',
      secondary: 'hsl(155 20% 16%)',
      secondaryForeground: 'hsl(150 20% 82%)',
      muted: 'hsl(155 20% 12%)',
      mutedForeground: 'hsl(150 8% 52%)',
      accent: 'hsl(160 35% 16%)',
      accentForeground: 'hsl(160 70% 70%)',
      destructive: 'hsl(0 72% 52%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(155 20% 18%)',
      input: 'hsl(155 20% 14%)',
      ring: 'hsl(160 70% 42%)',
      chart1: 'hsl(160 70% 46%)',
      chart2: 'hsl(185 70% 50%)',
      chart3: 'hsl(140 65% 46%)',
      chart4: 'hsl(200 75% 52%)',
      chart5: 'hsl(130 60% 48%)',
      sidebarBackground: 'hsl(155 22% 5%)',
      sidebarForeground: 'hsl(150 20% 88%)',
      sidebarPrimary: 'hsl(160 70% 42%)',
      sidebarPrimaryForeground: 'hsl(155 22% 5%)',
      sidebarAccent: 'hsl(155 20% 13%)',
      sidebarAccentForeground: 'hsl(160 70% 68%)',
      sidebarBorder: 'hsl(155 20% 15%)',
      sidebarRing: 'hsl(160 70% 42%)',
    },
  },
  {
    name: 'ember-glow',
    label: 'Ember Glow',
    light: {
      background: 'hsl(30 40% 97%)',
      foreground: 'hsl(15 30% 10%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(15 30% 10%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(15 30% 10%)',
      primary: 'hsl(20 90% 48%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(25 25% 93%)',
      secondaryForeground: 'hsl(15 30% 20%)',
      muted: 'hsl(25 20% 95%)',
      mutedForeground: 'hsl(20 10% 46%)',
      accent: 'hsl(15 80% 92%)',
      accentForeground: 'hsl(15 80% 32%)',
      destructive: 'hsl(0 84% 58%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(25 20% 86%)',
      input: 'hsl(25 20% 91%)',
      ring: 'hsl(20 90% 48%)',
      chart1: 'hsl(20 90% 48%)',
      chart2: 'hsl(0 80% 52%)',
      chart3: 'hsl(45 90% 50%)',
      chart4: 'hsl(10 75% 55%)',
      chart5: 'hsl(35 85% 52%)',
      sidebarBackground: 'hsl(25 25% 93%)',
      sidebarForeground: 'hsl(15 25% 14%)',
      sidebarPrimary: 'hsl(20 90% 48%)',
      sidebarPrimaryForeground: 'hsl(0 0% 100%)',
      sidebarAccent: 'hsl(25 20% 87%)',
      sidebarAccentForeground: 'hsl(20 80% 32%)',
      sidebarBorder: 'hsl(25 20% 84%)',
      sidebarRing: 'hsl(20 90% 48%)',
    },
    dark: {
      background: 'hsl(18 20% 7%)',
      foreground: 'hsl(30 20% 90%)',
      card: 'hsl(18 20% 10%)',
      cardForeground: 'hsl(30 20% 90%)',
      popover: 'hsl(18 20% 10%)',
      popoverForeground: 'hsl(30 20% 90%)',
      primary: 'hsl(20 90% 55%)',
      primaryForeground: 'hsl(18 20% 5%)',
      secondary: 'hsl(18 20% 16%)',
      secondaryForeground: 'hsl(30 20% 82%)',
      muted: 'hsl(18 20% 12%)',
      mutedForeground: 'hsl(20 8% 52%)',
      accent: 'hsl(15 40% 16%)',
      accentForeground: 'hsl(20 90% 72%)',
      destructive: 'hsl(0 72% 52%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(18 20% 20%)',
      input: 'hsl(18 20% 16%)',
      ring: 'hsl(20 90% 55%)',
      chart1: 'hsl(20 90% 58%)',
      chart2: 'hsl(0 80% 60%)',
      chart3: 'hsl(45 90% 56%)',
      chart4: 'hsl(10 75% 60%)',
      chart5: 'hsl(35 85% 56%)',
      sidebarBackground: 'hsl(18 22% 5%)',
      sidebarForeground: 'hsl(30 20% 88%)',
      sidebarPrimary: 'hsl(20 90% 55%)',
      sidebarPrimaryForeground: 'hsl(18 22% 5%)',
      sidebarAccent: 'hsl(18 20% 13%)',
      sidebarAccentForeground: 'hsl(20 90% 70%)',
      sidebarBorder: 'hsl(18 20% 16%)',
      sidebarRing: 'hsl(20 90% 55%)',
    },
  },
  {
    name: 'steel-sky',
    label: 'Steel Sky',
    light: {
      background: 'hsl(214 40% 96%)',
      foreground: 'hsl(215 40% 8%)',
      card: 'hsl(0 0% 100%)',
      cardForeground: 'hsl(215 40% 8%)',
      popover: 'hsl(0 0% 100%)',
      popoverForeground: 'hsl(215 40% 8%)',
      primary: 'hsl(217 88% 50%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(214 25% 92%)',
      secondaryForeground: 'hsl(215 40% 18%)',
      muted: 'hsl(214 20% 94%)',
      mutedForeground: 'hsl(215 12% 46%)',
      accent: 'hsl(217 80% 93%)',
      accentForeground: 'hsl(217 88% 36%)',
      destructive: 'hsl(0 84% 58%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(214 20% 86%)',
      input: 'hsl(214 20% 91%)',
      ring: 'hsl(217 88% 50%)',
      chart1: 'hsl(217 88% 50%)',
      chart2: 'hsl(196 85% 48%)',
      chart3: 'hsl(240 70% 60%)',
      chart4: 'hsl(170 75% 42%)',
      chart5: 'hsl(200 80% 52%)',
      sidebarBackground: 'hsl(216 52% 18%)',
      sidebarForeground: 'hsl(210 30% 88%)',
      sidebarPrimary: 'hsl(217 88% 62%)',
      sidebarPrimaryForeground: 'hsl(216 52% 10%)',
      sidebarAccent: 'hsl(216 52% 26%)',
      sidebarAccentForeground: 'hsl(210 50% 85%)',
      sidebarBorder: 'hsl(216 52% 26%)',
      sidebarRing: 'hsl(217 88% 62%)',
    },
    dark: {
      background: 'hsl(216 30% 8%)',
      foreground: 'hsl(210 25% 90%)',
      card: 'hsl(216 30% 11%)',
      cardForeground: 'hsl(210 25% 90%)',
      popover: 'hsl(216 30% 11%)',
      popoverForeground: 'hsl(210 25% 90%)',
      primary: 'hsl(217 88% 58%)',
      primaryForeground: 'hsl(0 0% 100%)',
      secondary: 'hsl(216 30% 17%)',
      secondaryForeground: 'hsl(210 25% 82%)',
      muted: 'hsl(216 30% 13%)',
      mutedForeground: 'hsl(210 12% 52%)',
      accent: 'hsl(217 40% 18%)',
      accentForeground: 'hsl(217 80% 78%)',
      destructive: 'hsl(0 72% 52%)',
      destructiveForeground: 'hsl(0 0% 100%)',
      border: 'hsl(216 30% 19%)',
      input: 'hsl(216 30% 15%)',
      ring: 'hsl(217 88% 58%)',
      chart1: 'hsl(217 88% 60%)',
      chart2: 'hsl(196 85% 55%)',
      chart3: 'hsl(240 70% 65%)',
      chart4: 'hsl(170 75% 48%)',
      chart5: 'hsl(200 80% 58%)',
      sidebarBackground: 'hsl(216 35% 6%)',
      sidebarForeground: 'hsl(210 25% 88%)',
      sidebarPrimary: 'hsl(217 88% 58%)',
      sidebarPrimaryForeground: 'hsl(216 35% 5%)',
      sidebarAccent: 'hsl(216 30% 14%)',
      sidebarAccentForeground: 'hsl(217 80% 78%)',
      sidebarBorder: 'hsl(216 30% 16%)',
      sidebarRing: 'hsl(217 88% 58%)',
    },
  },
]

const THEME_STORAGE_KEY = 'app-theme'
const TYPOGRAPHY_STORAGE_KEY = 'app-typography'

const currentTheme = ref<ThemeName>(
  (localStorage.getItem(THEME_STORAGE_KEY) as ThemeName) || 'golden-white-dark'
)

const typography = ref<TypographySettings>(
  JSON.parse(localStorage.getItem(TYPOGRAPHY_STORAGE_KEY) || 'null') || {
    fontFamily: 'Tajawal',
    fontSize: '1.125rem',
    fontWeight: '500',
  }
)

export function applyTheme(themeName: ThemeName, isDark: boolean) {
  const theme = themes.find(t => t.name === themeName)
  if (!theme) return

  const colors = isDark ? theme.dark : theme.light
  const root = document.documentElement

  // Apply all CSS variables
  root.style.setProperty('--background', colors.background)
  root.style.setProperty('--foreground', colors.foreground)
  root.style.setProperty('--card', colors.card)
  root.style.setProperty('--card-foreground', colors.cardForeground)
  root.style.setProperty('--popover', colors.popover)
  root.style.setProperty('--popover-foreground', colors.popoverForeground)
  root.style.setProperty('--primary', colors.primary)
  root.style.setProperty('--primary-foreground', colors.primaryForeground)
  root.style.setProperty('--secondary', colors.secondary)
  root.style.setProperty('--secondary-foreground', colors.secondaryForeground)
  root.style.setProperty('--muted', colors.muted)
  root.style.setProperty('--muted-foreground', colors.mutedForeground)
  root.style.setProperty('--accent', colors.accent)
  root.style.setProperty('--accent-foreground', colors.accentForeground)
  root.style.setProperty('--destructive', colors.destructive)
  root.style.setProperty('--destructive-foreground', colors.destructiveForeground)
  root.style.setProperty('--border', colors.border)
  root.style.setProperty('--input', colors.input)
  root.style.setProperty('--ring', colors.ring)
  root.style.setProperty('--chart-1', colors.chart1)
  root.style.setProperty('--chart-2', colors.chart2)
  root.style.setProperty('--chart-3', colors.chart3)
  root.style.setProperty('--chart-4', colors.chart4)
  root.style.setProperty('--chart-5', colors.chart5)
  root.style.setProperty('--sidebar-background', colors.sidebarBackground)
  root.style.setProperty('--sidebar-foreground', colors.sidebarForeground)
  root.style.setProperty('--sidebar-primary', colors.sidebarPrimary)
  root.style.setProperty('--sidebar-primary-foreground', colors.sidebarPrimaryForeground)
  root.style.setProperty('--sidebar-accent', colors.sidebarAccent)
  root.style.setProperty('--sidebar-accent-foreground', colors.sidebarAccentForeground)
  root.style.setProperty('--sidebar-border', colors.sidebarBorder)
  root.style.setProperty('--sidebar-ring', colors.sidebarRing)

  // Typography is pre-loaded in root.html
  const fontValue = typography.value.fontFamily.includes(' ')
    ? `"${typography.value.fontFamily}", sans-serif`
    : `${typography.value.fontFamily}, sans-serif`

  // Ultimate "Nuclear" Fix: Dynamic Style Injection
  const styleId = 'typography-fix'
  let styleTag = document.getElementById(styleId) as HTMLStyleElement
  if (!styleTag) {
    styleTag = document.createElement('style')
    styleTag.id = styleId
    document.head.appendChild(styleTag)
  }

  const css = `
    * {
      font-family: ${fontValue}, "Cairo", sans-serif !important;
      font-weight: ${typography.value.fontWeight} !important;
      font-variation-settings: "wght" ${typography.value.fontWeight} !important;
    }
    h1, h2, h3, h4, h5, h6 {
      font-weight: 700 !important;
      font-variation-settings: "wght" 700 !important;
    }
  `
  styleTag.textContent = css

  console.log('Typography enforced via dynamic style tag:', {
    family: fontValue,
    weight: typography.value.fontWeight
  })

  // Set on root element as well for Tailwind/Component usage
  root.style.setProperty('--font-family', fontValue)
  root.style.setProperty('--font-size', typography.value.fontSize)
  root.style.setProperty('--font-weight', typography.value.fontWeight)

  // Scale the global font size
  root.style.fontSize = typography.value.fontSize

  // Also force on body
  document.body.style.setProperty('--font-family', fontValue)
  document.body.style.setProperty('--font-weight', typography.value.fontWeight)
  document.body.style.setProperty('font-variation-settings', `"wght" ${typography.value.fontWeight}`)
  document.body.style.fontSize = typography.value.fontSize
}

export function useTheme() {
  const updateTheme = (themeName: ThemeName) => {
    currentTheme.value = themeName
    localStorage.setItem(THEME_STORAGE_KEY, themeName)

    // Get current dark mode state
    const isDark = document.documentElement.classList.contains('dark')
    applyTheme(themeName, isDark)
  }

  const updateTypography = (settings: Partial<TypographySettings>) => {
    typography.value = { ...typography.value, ...settings }
    localStorage.setItem(TYPOGRAPHY_STORAGE_KEY, JSON.stringify(typography.value))

    const isDark = document.documentElement.classList.contains('dark')
    applyTheme(currentTheme.value, isDark)
  }

  // Watch for dark mode changes
  const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.attributeName === 'class') {
        const isDark = document.documentElement.classList.contains('dark')
        applyTheme(currentTheme.value, isDark)
      }
    })
  })

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class'],
  })

  return {
    currentTheme,
    typography,
    themes,
    updateTheme,
    updateTypography,
    applyTheme,
  }
}
