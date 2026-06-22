import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'
import 'dayjs/locale/ar'
import { QuotaUnit } from '@/types'

dayjs.extend(relativeTime)
dayjs.extend(utc)
dayjs.extend(timezone)
dayjs.locale('ar')
dayjs.tz.setDefault('Africa/Cairo')

export function formatRelativeTime(timestamp?: string): string {
  const tz = 'Africa/Cairo'

  try {
    if (!timestamp) return dayjs().tz(tz).fromNow()

    // Normalize Go formats (space instead of T, missing timezone)
    let normalized = timestamp.trim()

    if (/^\d{4}-\d{2}-\d{2} \d/.test(normalized)) {
      normalized = normalized.replace(' ', 'T') + '+03:00'
    }

    const date = dayjs.tz(normalized, tz)

    if (!date.isValid()) {
      return dayjs().tz(tz).fromNow()
    }

    return date.fromNow()
  } catch (err) {
    return dayjs().tz(tz).fromNow()
  }
}

export const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('ar-EG', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    // timeZone: 'Africa/Cairo'
  });
};

export const formatTime = (timeString: string, hour12: boolean = true) => {
  if (!timeString) return '';
  const date = new Date(`1970-01-01T${timeString}Z`);
  return date.toLocaleTimeString('ar-EG', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: hour12,
    //  timeZone: 'Africa/Cairo'
  });
};

export const formatCurrency = (amount: number, currency: string = 'EGP', showSign = false) => {
  if (!currency) return '';
  const formatted = new Intl.NumberFormat('ar-EG', {
    style: 'currency',
    currency: currency
  }).format(Math.abs(amount));

  if (showSign && amount !== 0) {
    return amount > 0 ? `+${formatted}` : `-${formatted}`;
  }
  return formatted;
};

export const formatNumber = (input: number | string): string => {
  if (typeof input === "number") {
    return input.toLocaleString("ar-EG", {
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    });
  }

  return input.replace(/\d+(\.\d+)?/g, (match) => {
    const num = parseFloat(match);
    if (isNaN(num)) return match;
    return num.toLocaleString("ar-EG", {
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    });
  });
};


/**
 * Format price in Egyptian Pounds
 */
export const formatPrice = (price: number) => {
  return new Intl.NumberFormat('ar-EG', {
    style: 'currency',
    currency: 'EGP',
  }).format(price);
};


export function formatSeconds(seconds: number) {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;

  const parts = [];
  if (h > 0) parts.push(`${h} ساعة`);
  if (m > 0) parts.push(`${m} دقيقة`);
  if (s > 0 || parts.length === 0) parts.push(`${s} ثانية`);

  return parts.join(' و ');
}

export const getTimeLeft = (deletedAt: string) => {
  if (!deletedAt) return '';
  const date = dayjs(deletedAt).add(15, 'day');
  return date.fromNow();
};
