<script setup lang="ts">
import InputError from '@/components/InputError.vue';
import TextLink from '@/components/TextLink.vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import AuthBase from '@/layouts/AuthLayout.vue';
import { Head, useForm } from '@inertiajs/vue3';
import { FacebookIcon, GithubIcon, LoaderCircle } from 'lucide-vue-next';

const form = useForm({
  tenant_name: '',
  name: '',
  username: '',
  email: '',
  password: '',
  password_confirmation: '',
});

const submit = () => {
  form.post('/register', {
    onFinish: () => form.reset('password', 'password_confirmation'),
  });
};

const loginWithProvider = (provider: string) => {
  window.location.href = `/auth/${provider}/redirect`;
};
</script>

<template>
  <AuthBase title="إنشاء حساب جديد" description="أدخل بياناتك أدناه لإنشاء حسابك">

    <Head title="التسجيل" />

    <form @submit.prevent="submit" class="flex flex-col gap-6" dir="rtl">
      <div class="grid gap-6">

        <!-- Form Fields -->
        <div class="grid gap-2">
          <Label for="tenant_name">اسم الشبكة</Label>
          <Input id="tenant_name" type="text" required autofocus :tabindex="1" autocomplete="tenant_name"
            v-model="form.tenant_name" placeholder="اسم الشبكة" />
          <InputError :message="form.errors.tenant_name" />
        </div>

        <div class="grid gap-2">
          <Label for="name">الاسم الكامل</Label>
          <Input id="name" type="text" required :tabindex="2" autocomplete="name" v-model="form.name"
            placeholder="الاسم الكامل" />
          <InputError :message="form.errors.name" />
        </div>

        <div class="grid gap-2">
          <Label for="username">اسم المستخدم</Label>
          <Input id="username" type="username" required :tabindex="3" autocomplete="username" v-model="form.username"
            placeholder="username101" />
          <InputError :message="form.errors.username" />
        </div>

        <div class="grid gap-2">
          <Label for="email">البريد الإلكتروني</Label>
          <Input id="email" type="email" required :tabindex="3" autocomplete="email" v-model="form.email"
            placeholder="email@example.com" />
          <InputError :message="form.errors.email" />
        </div>

        <div class="grid gap-2">
          <Label for="password">كلمة المرور</Label>
          <Input id="password" type="password" required :tabindex="4" autocomplete="new-password"
            v-model="form.password" placeholder="كلمة المرور" />
          <InputError :message="form.errors.password" />
        </div>

        <div class="grid gap-2">
          <Label for="password_confirmation">تأكيد كلمة المرور</Label>
          <Input id="password_confirmation" type="password" required :tabindex="5" autocomplete="new-password"
            v-model="form.password_confirmation" placeholder="تأكيد كلمة المرور" />
          <InputError :message="form.errors.password_confirmation" />
        </div>

        <Button type="submit" class="mt-2 w-full" :tabindex="6" :disabled="form.processing">
          <LoaderCircle v-if="form.processing" class="h-4 w-4 animate-spin" />
          إنشاء الحساب
        </Button>
      </div>

      <!-- Divider -->
      <!-- <div class="relative">
        <div class="absolute inset-0 flex items-center">
          <span class="w-full border-t"></span>
        </div>
        <div class="relative flex justify-center text-xs uppercase">
          <span class="bg-background px-2 text-muted-foreground">أو التسجيل عن طريق</span>
        </div>
      </div> -->

      <!-- Social Login Buttons -->
      <!-- <div class="grid gap-3">
        <Button type="button" variant="outline" class="w-full" @click="loginWithProvider('google')">
          <svg class="h-4 w-4 mr-2" viewBox="0 0 24 24">
            <path fill="currentColor"
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
            <path fill="currentColor"
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
            <path fill="currentColor"
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
            <path fill="currentColor"
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
          </svg>
          التسجيل عن طريق جوجل
        </Button>

        <Button type="button" variant="outline" class="w-full" @click="loginWithProvider('facebook')">
          <FacebookIcon class="h-4 w-4 mr-2" />
          التسجيل عن طريق فيسبوك
        </Button>

        <Button type="button" variant="outline" class="w-full" @click="loginWithProvider('github')">
          <GithubIcon class="h-4 w-4 mr-2" />
          التسجيل عن طريق GitHub
        </Button>
      </div> -->

      <div class="text-center text-sm text-muted-foreground">
        لديك حساب بالفعل؟
        <TextLink :href="'/login'" class="underline underline-offset-4" :tabindex="7">تسجيل الدخول</TextLink>
      </div>
    </form>
  </AuthBase>
</template>