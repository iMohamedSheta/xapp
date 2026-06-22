<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { toast } from 'vue-sonner'
import { usePage } from '@inertiajs/vue3'
import { Toaster } from '@/components/ui/sonner'
import 'vue-sonner/style.css'
import { FlashMessage, LottieFlash } from '@/types'
import Swal from 'sweetalert2'
import '/resources/css/sweetalert_overwrite.css'
import { DotLottieVue } from '@lottiefiles/dotlottie-vue'

const page = usePage()
const showLottie = ref(false)
const lottieData = ref<LottieFlash | null>(null)

function showFlashToast(flash: FlashMessage) {
  const toastData = flash?.toast
  if (Array.isArray(toastData)) {
    toastData.forEach(toastr => {
      if (toastr?.title) toast(toastr.title, toastr)
    })
  }

  const alertData = flash?.sweetalert
  if (Array.isArray(alertData)) {
    alertData.forEach(alert => {
      if (alert?.title) Swal.fire(alert)
    })
  }

  const lottie = flash?.lottie as LottieFlash
  if (lottie) {
    lottieData.value = lottie
    showLottie.value = true

    // Auto-close by timer
    if (lottie.timer) {
      setTimeout(() => {
        showLottie.value = false
      }, lottie.timer)
    }
  }
}

onMounted(async () => {
  await nextTick()
  setTimeout(() => {
    showFlashToast(page.props.flash!)
  }, 50)

  // listen for global clicks
  window.addEventListener('click', handleClose)
})

onBeforeUnmount(() => {
  window.removeEventListener('click', handleClose)
})

function handleClose() {
  if (lottieData.value?.closeOnClick) {
    showLottie.value = false
  }
}

function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape' && lottieData.value?.closeOnEscape !== false) {
    showLottie.value = false
  }
}

onMounted(() => {
  setTimeout(() => {
    showFlashToast(page.props.flash!)
  }, 50)
  window.addEventListener('click', handleClose)
  window.addEventListener('keydown', handleEscape)
})

onBeforeUnmount(() => {
  window.removeEventListener('click', handleClose)
  window.removeEventListener('keydown', handleEscape)
})


watch(
  () => page.props.flash,
  (flash) => {
    showFlashToast(flash!)
  },
  { deep: true }
)
</script>

<template>
  <Toaster />
  <DotLottieVue v-if="showLottie && lottieData" :src="lottieData.src" :class="lottieData.class" v-bind="{
    autoplay: lottieData.lottieProps?.autoplay,
    loop: lottieData.lottieProps?.loop,
    speed: lottieData.lottieProps?.speed,
    playOnHover: lottieData.lottieProps?.playOnHover,
    backgroundColor: lottieData.lottieProps?.backgroundColor,
    animationId: lottieData.lottieProps?.animationId,
    segment: lottieData.lottieProps?.segment,
    width: lottieData.lottieProps?.width,
    height: lottieData.lottieProps?.height,
  }" />
</template>
