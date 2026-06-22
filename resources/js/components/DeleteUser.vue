<script setup lang="ts">
import { useForm } from '@inertiajs/vue3';
import { ref } from 'vue';
import HeadingSmall from '@/components/HeadingSmall.vue';
import InputError from '@/components/InputError.vue';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { lang } from '@/lib/lang';

const passwordInput = ref<HTMLInputElement | null>(null);

const form = useForm({
  password: '',
});

const deleteUser = (e: Event) => {
  e.preventDefault();

  form.delete('profile/destroy', {
    preserveScroll: true,
    onSuccess: () => closeModal(),
    onError: () => passwordInput.value?.focus(),
    onFinish: () => form.reset(),
  });
};

const closeModal = () => {
  form.clearErrors();
  form.reset();
};
</script>

<template>
  <div class="space-y-6">
    <HeadingSmall :title="lang('settings.delete_user.title')" :description="lang('settings.delete_user.description')" />
    <div class="space-y-4 rounded-lg border border-red-100 bg-red-50 p-4 dark:border-red-200/10 dark:bg-red-700/10">
      <div class="relative space-y-0.5 text-red-600 dark:text-red-100">
        <p class="font-medium">
          {{ lang('main.warning') }}
        </p>
        <p class="text-sm">
          {{ lang('settings.delete_user.description') }}
        </p>
      </div>
      <Dialog>
        <DialogTrigger as-child>
          <Button variant="destructive">
            {{ lang('settings.delete_user.delete') }}
          </Button>
        </DialogTrigger>
        <DialogContent>
          <form class="space-y-6" @submit="deleteUser">
            <DialogHeader class="space-y-3">
              <DialogTitle>
                {{ lang('settings.delete_user.delete_confirm.title') }}
              </DialogTitle>
              <DialogDescription>
                {{ lang('settings.delete_user.delete_confirm.description') }}
              </DialogDescription>
            </DialogHeader>

            <div class="grid gap-2">
              <Label for="password" class="sr-only">
                {{ lang('settings.delete_user.delete_confirm.password') }}
              </Label>
              <Input id="password" type="password" name="password" ref="passwordInput" v-model="form.password"
                :placeholder="lang('settings.delete_user.delete_confirm.password')" />
              <InputError :message="form.errors.password" />
            </div>

            <DialogFooter class="gap-2">
              <DialogClose as-child>
                <Button variant="secondary" @click="closeModal">{{ lang('main.cancel') }}</Button>
              </DialogClose>

              <Button type="submit" variant="destructive" :disabled="form.processing">
                {{ lang('settings.delete_user.delete_confirm.button') }}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  </div>
</template>
