<script setup lang="ts">
import WebsiteHeader from './website/WebsiteHeader.vue';
import { CalendarIcon } from 'lucide-vue-next';
import Button from '@/components/ui/button/Button.vue';
import AppDebug from '@/components/debug/AppDebug.vue';
import { router } from '@inertiajs/vue3';
import Input from '@/components/ui/input/Input.vue';
</script>

<template>
  <div class="min-h-svh overflow-hidden bg-gradient-to-r from-[#5a4279c0] to-[#d6c1f8c0]">
    <WebsiteHeader />
    <div class="mt-4 lg:mt-16">
      <slot />
    </div>
    <div class="md:mx-8 mx-2 mb-4">
      <!-- Call to Action Section -->
      <div class="mt-12  bg-gradient-to-r  bg-background from-primary/10 to-purple-500/10 rounded-4xl p-12 text-center">
        <h3 class="text-2xl font-bold mb-4">لديك فعالية؟</h3>
        <p class="text-muted-foreground mb-6">شارك فعاليتك مع الآلاف من المهتمين وابني مجتمعك المهني</p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <Button @click="router.get('/login')" size="lg" class="bg-primary hover:bg-primary/90">
            <CalendarIcon class="w-5 h-5 ml-2" />
            سجل فعاليتك الان
          </Button>
          <Button size="lg" variant="outline">
            تعرف على المزيد
          </Button>
        </div>
      </div>
    </div>
  </div>
  <!-- Footer -->
  <footer class="pt-6 bg-gradient-to-r from-[#5a4279c0] to-[#d6c1f8c0]" dir="rtl">
    <div class="bg-[#0d0b11] text-gray-300 border-t border-[#2f2640]">
      <div class=" max-w-7xl mx-auto px-4 py-16 grid grid-cols-1 md:grid-cols-4 gap-8 text-sm">
        <div class="text-right space-y-2">
          <h5 class="font-semibold text-white">عن المنصة</h5>
          <p class="text-gray-400">منصة لاكتشاف وحجز أفضل الفعاليات.</p>
        </div>
        <div class="text-right space-y-2">
          <h5 class="font-semibold text-white">روابط</h5>
          <ul class="space-y-1">
            <li><a href="/events" class="hover:underline">الفعاليات</a></li>
            <li><a href="/login" class="hover:underline">تسجيل الدخول</a></li>
            <li><a href="/register" class="hover:underline">إنشاء حساب</a></li>
          </ul>
        </div>
        <div class="text-right space-y-2">
          <h5 class="font-semibold text-white">الموارد</h5>
          <ul class="space-y-1">
            <li><a href="#" class="hover:underline">الأسئلة الشائعة</a></li>
            <li><a href="#" class="hover:underline">سياسة الخصوصية</a></li>
            <li><a href="#" class="hover:underline">الشروط والأحكام</a></li>
          </ul>
        </div>
        <div class="text-right space-y-3">
          <h5 class="font-semibold text-white">النشرة البريدية</h5>
          <p class="text-gray-400">اشترك ليصلك جديد الفعاليات.</p>
          <div class="flex gap-2 justify-end">
            <Button class="bg-primary hover:bg-primary/90">اشترك</Button>
            <Input dir="rtl" class="max-w-[220px] text-right" placeholder="بريدك الإلكتروني" />
          </div>
        </div>
      </div>
      <div class="border-t border-[#2f2640] py-4 text-center text-gray-500 text-xs">
        © {{ new Date().getFullYear() }} جميع الحقوق محفوظة @MohamedSheta
      </div>
    </div>
  </footer>
  <AppDebug />
</template>