<script setup lang="ts">
import InputError from "@/components/InputError.vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import AuthBase from "@/layouts/AuthLayout.vue";
import { Head, router, useForm } from "@inertiajs/vue3";
import { LoaderCircle, InfoIcon, Globe } from "lucide-vue-next";
import { lang } from "@/lib/lang";
import { onMounted } from "vue";
import { ToastError } from "@/lib/toast";

defineProps<{
  status?: string;
  canResetPassword: boolean;
}>();

const form = useForm({
  identifier: "",
  email: "",
  username: "",
  password: "",
  remember: false,
});

const submit = () => {
  // Detect if identifier is email or username
  const isEmail = form.identifier.includes('@');

  // Transform the form data before submission
  form.transform((data) => {
    const payload: any = {
      password: data.password,
      remember: data.remember,
    };

    if (isEmail) {
      payload.email = data.identifier;
    } else {
      payload.username = data.identifier;
    }

    return payload;
  }).post("/login", {
    onFinish: () => form.reset("password"),
  });
};
onMounted(async () => {
  const params = new URLSearchParams(window.location.search);
  const error = params.get("error");

  if (error === "user_blocked") {
    ToastError("مرفوض", "حسابك محظور. يرجى التواصل مع الدعم الفني.");
  } else if (error === "invalid_credentials") {
    ToastError("خطأ", "بيانات الدخول غير صحيحة.");
  }

  // Remove error from URL if it exists
  if (error) {
    setTimeout(() => {
      params.delete("error");
      const newParams = params.toString();
      const newUrl = newParams ? `${window.location.pathname}?${newParams}` : window.location.pathname;
      window.history.replaceState(null, "", newUrl);
    }, 40);
  }
});
</script>

<template>
  <AuthBase title="بوابة المشتركين" description="تسجيل الدخول لبوابة المشتركين الخاصة بك">

    <Head title="بوابة المشتركين - تسجيل الدخول" />

    <div v-if="status" class="mb-4 text-center text-sm font-medium text-green-600">
      {{ status }}
    </div>

    <form @submit.prevent="submit" class="flex flex-col gap-6" dir="rtl">
      <div class="grid gap-6">
        <div class="grid gap-2">
          <Label for="identifier">{{ lang('auth.login.username') }}</Label>
          <Input id="identifier" type="text" autofocus :tabindex="1" autocomplete="username" v-model="form.identifier"
            placeholder="username" />
          <InputError :message="form.errors.username || form.errors.email" />
        </div>

        <div class="grid gap-2">
          <div class="flex items-center justify-between">
            <Label for="password">{{ lang('auth.login.password') }}</Label>
            <TextLink v-if="canResetPassword" :href="'/password/reset'" class="text-sm" :tabindex="5">
              {{ lang('auth.login.forgot_password') }}
            </TextLink>
          </div>
          <Input id="password" type="password" required :tabindex="2" autocomplete="current-password"
            v-model="form.password" placeholder="كلمة المرور" />
          <InputError :message="form.errors.password" />
        </div>

        <div class="flex items-center justify-between">
          <Label for="remember" class="flex items-center space-x-3">
            <Checkbox id="remember" v-model="form.remember" :tabindex="3" />
            <span>{{ lang('auth.login.remember') }}</span>
          </Label>
        </div>

        <Button type="submit" class="w-full" :tabindex="4" :disabled="form.processing">
          <LoaderCircle v-if="form.processing" class="h-4 w-4 animate-spin" />
          {{ form.processing ? lang('auth.login.login_processing') : lang('auth.login.login_button') }}
        </Button>
      </div>
      <div class="mt-4 rounded-xl border border-primary/20 bg-primary/10 p-4 relative overflow-hidden">
        <div class="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent pointer-events-none"></div>
        <div class="relative z-10">
          <div class="mb-3 flex items-center justify-center gap-2 text-primary">
            <Globe class="h-4 w-4" />
            <p class="font-semibold text-sm">بوابة خدمات الإنترنت</p>
          </div>
          <div class="text-center text-sm text-muted-foreground leading-relaxed">
            يمكنك من خلال بوابة المشتركين معرفة استهلاكك الشهري، تجديد اشتراكك، والتواصل مع الدعم الفني بشكل مباشر.
          </div>
        </div>
      </div>
    </form>
  </AuthBase>
</template>
