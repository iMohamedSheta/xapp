<script setup lang="ts">
import InputError from "@/components/InputError.vue";
import TextLink from "@/components/TextLink.vue";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import AuthBase from "@/layouts/AuthLayout.vue";
import { Head, router, useForm } from "@inertiajs/vue3";
import { FacebookIcon, GithubIcon, LoaderCircle, InfoIcon } from "lucide-vue-next";
import { lang } from "@/lib/lang";
import { onMounted } from "vue";
import { ToastError } from "@/lib/toast";

const props = defineProps<{
  status?: string;
  canResetPassword: boolean;
  demo: boolean;
  username?: string;
  password?: string;
}>();

const form = useForm({
  identifier: props.demo ? props.username || "" : "",
  email: "",
  username: "",
  password: props.demo ? props.password || "" : "",
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
    ToastError("Blocked", "Your account is blocked. Please contact support.");
  } else if (error === "email_already_taken") {
    ToastError("Failed", "This email is already taken by another user.");
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

const loginWithProvider = (provider: string) => {
  window.location.href = `/auth/${provider}/redirect`;
};
</script>

<template>
  <AuthBase :title="lang('auth.login.title')" :description="lang('auth.login.description')">

    <Head :title="lang('auth.login.title')" />

    <div v-if="status" class="mb-4 text-center text-sm font-medium text-green-600">
      {{ status }}
    </div>

    <form @submit.prevent="submit" class="flex flex-col gap-6" dir="rtl">
      <div class="grid gap-6">
        <div class="grid gap-2">
          <Label for="identifier">{{ lang('auth.login.email') }} أو اسم المستخدم</Label>
          <Input id="identifier" type="text" autofocus :tabindex="1" autocomplete="username" v-model="form.identifier"
            placeholder="email@example.com أو username" />
          <InputError :message="form.errors.email || form.errors.username" />
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
      <!-- Divider -->
      <div class="relative">
        <div class="absolute inset-0 flex items-center">
          <span class="w-full border-t"></span>
        </div>
        <div class="relative flex justify-center text-xs uppercase">
          <span class="bg-background px-2 text-muted-foreground">أو التسجيل عن طريق</span>
        </div>
      </div>

      <!-- Social Login Buttons -->
      <div class="grid gap-3">
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
          <bdi>
            التسجيل عن طريق
            <bdi>
              GitHub
            </bdi>
          </bdi>
        </Button>
      </div>


      <div class="text-center text-sm text-muted-foreground">
        {{ lang('auth.login.no_account') }}
        <TextLink :href="'/register'" :tabindex="5">{{ lang('auth.login.sign_up') }}</TextLink>
      </div>
    </form>
  </AuthBase>
</template>
