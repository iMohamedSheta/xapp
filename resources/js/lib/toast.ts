import { toast } from "vue-sonner"

export type ToastOptions = {
  title: string
  type: 'success' | 'error'
  richColors: boolean
  description?: string
  classes: Record<string, string>
  closeButton: boolean
  cancelButtonStyle: Record<string, string>
  style: Record<string, string>
  position: 'bottom-center'
  duration: number
  invert: boolean
  important: boolean
}

/* ===================== SUCCESS ===================== */
export function ToastSuccessOptions(
  msg: string,
  description?: string
): ToastOptions {
  return {
    title: msg,
    type: 'success',
    richColors: true,
    description,
    classes: {
      description: 'text-xs !text-gray-400 px-2',
      title: 'px-2',
      closeButton: '!bg-[#1f1f1f] !outline-none !border-none !text-white',
    },
    closeButton: true,
    cancelButtonStyle: {
      backgroundColor: '#1f1f1f',
    },
    style: {
      backgroundColor: '#0a0a0a',
      borderColor: '#1f1f1f',
      direction: 'rtl',
      textAlign: 'right',
      fontFamily: 'Cairo, sans-serif',
      paddingRight: '10px',
    },
    position: 'bottom-center',
    duration: 5000,
    invert: true,
    important: false,
  }
}

/* ===================== ERROR ===================== */
export function ToastErrorOptions(
  msg: string,
  description?: string
): ToastOptions {
  return {
    title: msg,
    type: 'error',
    richColors: true,
    description,
    classes: {
      description: 'text-xs !text-gray-400 px-2',
      title: 'px-2 !text-red-500',
      icon: '!text-red-500',
      closeButton: '!bg-[#1f1f1f] !outline-none !border-none !text-white',
    },
    closeButton: true,
    cancelButtonStyle: {
      backgroundColor: '#1f1f1f',
    },
    style: {
      backgroundColor: '#0a0a0a',
      borderColor: '#1f1f1f',
      direction: 'rtl',
      textAlign: 'right',
      fontFamily: 'Cairo, sans-serif',
      paddingRight: '10px',
    },
    position: 'bottom-center',
    duration: 5000,
    invert: true,
    important: false,
  }
}


export function ToastError(msg: string, description?: string) {
    toast(msg,ToastErrorOptions(msg, description))
}

export function ToastSuccess(msg: string, description?: string) {
    toast(msg,ToastSuccessOptions(msg, description))
}