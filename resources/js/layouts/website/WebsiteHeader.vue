<template>
  <!-- <nav dir="rtl" class="fixed top-4  right-1/2 translate-x-1/2 z-50 
           bg-gradient-to-r from-yellow-100  to-amber-300
           border border-gray-200 shadow-2xl w-[95%] max-w-8xl 
           px-6 py-5 flex items-center justify-between 
           rounded-2xl backdrop-blur-md"> -->
  <transition enter-active-class="transition transform duration-500" enter-from-class="-translate-y-full opacity-0"
    enter-to-class="translate-y-0 opacity-100" leave-active-class="transition transform duration-500"
    leave-from-class="translate-y-0 opacity-100" leave-to-class="-translate-y-full opacity-0">

    <nav dir="rtl" :class="[
      'w-[95%] max-w-8xl m-auto px-8 py-4 z-1 flex items-center justify-between rounded-b-2xl backdrop-blur-md transition-all duration-500 ease-in-out',
      stuck
        ? ' z-50 fixed w-full rounded-b-none top-0 left-0 right-0 bg-[#0d0a11c0] shadow-2xl '
        : 'relative bg-[#0d0a11d7] w-full  rounded-b-none'
    ]">
      <!-- الشعار -->
      <div class="font-extrabold text-xl text-gray-200 tracking-tight flex items-center m-0 p-0">
        <!-- <TicketIcon /> -->
        <Button @click="openMainPage" variant="link"
          class="font-extrabold text-xl text-gray-200 tracking-tight flex items-center">
          <Logo />
          <span>
            <!-- TicketIZY -->
            <!-- TicketEasy -->
            <span class="text-2xl">
              Connect4ar
            </span>
            <!-- - تيك ايزي -->
            <!-- تيك ايزي - Tickizy -->
          </span>
        </Button>
      </div>

      <!-- شريط البحث (ديسكتوب) -->
      <div class="hidden md:flex flex-1 mx-8">
        <div class="relative w-full max-w-md mx-auto">
          <Input v-model="searchQuery" @keyup.enter="handleSearch" type="text" placeholder="ابحث عن فعاليات..."
            class="w-full  py-6 text-gray-200 !ring-0 !border-0 rounded-full px-10" />
          <span class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 cursor-pointer"
            @click="handleSearch">
            <SearchIcon class="w-5 h-5" />
          </span>
        </div>
      </div>

      <!-- القائمة اليمنى (ديسكتوب) -->
      <div class="hidden md:flex items-center space-x-6 text-sm font-medium text-gray-700  ">
        <!-- <a href="#" class="hover:text-amber-600 transition">الفعاليات</a>
        <a href="#" class="hover:text-amber-600 transition">المنافذ</a>
        <a href="#" class="hover:text-amber-600 transition">الدعم والتواصل</a> -->
        <NavigationMenu>
          <NavigationMenuList dir="rtl">
            <NavigationMenuItem>
              <NavigationMenuLink class="!bg-transparent !text-gray-200 underline-offset-4 hover:underline "
                href="/events/list" :class="navigationMenuTriggerStyle()">
                الفعاليات
              </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NavigationMenuLink class="!bg-transparent !text-gray-200 underline-offset-4 hover:underline "
                href="/docs/introduction" :class="navigationMenuTriggerStyle()">
                الدعم والتواصل
              </NavigationMenuLink>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>

        <!-- اللغة -->
        <Button variant="link" size="icon" class="text-gray-300 hover:text-purple-300">
          <LanguagesIcon />
        </Button>

        <!-- الحساب -->
        <Button v-if="!isLoggedIn" @click="openLoginPage" asChild variant="ghost" size="icon"
          class="relative size-10 w-auto rounded-full p-1 focus-within:ring-2 focus-within:ring-primary">
          <Avatar class="size-10 overflow-hidden rounded-full">
            <AvatarFallback
              class="flex items-center justify-center rounded-full bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
              <User class="w-5 h-5" /> <!-- Fallback user icon -->
            </AvatarFallback>
          </Avatar>
        </Button>

        <DropdownMenu v-else>
          <DropdownMenuTrigger :as-child="true">
            <Button variant="ghost" size="icon"
              class="relative size-10 w-auto rounded-full p-1 focus-within:ring-2 focus-within:ring-primary">
              <Avatar class="size-8 overflow-hidden rounded-full">
                <AvatarImage v-if="page.props.auth?.user?.avatar" :src="page.props.auth?.user?.avatar"
                  :alt="page.props.auth?.user?.name" />
                <AvatarFallback
                  class="flex items-center justify-center rounded-lg bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
                  {{ getInitials(page.props.auth?.user?.name) }}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" class="w-56">
            <UserMenuContent :user="page.props.auth.user" />
          </DropdownMenuContent>
        </DropdownMenu>
        <Button variant="secondary" size="lg" @click="router.get('/register')" class="border-none">
          <CalendarIcon class="w-5 h-5" />
          نظم فعاليتك الان!
        </Button>
      </div>

      <!-- زر القائمة (موبايل) -->
      <Button variant="ghost" size="icon" class="md:hidden" @click="menuOpen = !menuOpen">
        <MenuIcon v-if="!menuOpen" class="w-6 h-6 text-gray-200" />
        <XIcon v-else class="w-6 h-6 text-gray-200" />
      </Button>

      <!-- القائمة المنسدلة (موبايل) -->
      <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0" leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 -translate-y-2">
        <div v-if="menuOpen" class="absolute top-full right-0 w-full 
                 bg-[#211829]
                  p-4  text-gray-200
                 flex flex-col space-y-3 md:hidden">
          <div class="relative w-full max-w-md mx-auto">
            <Input v-model="searchQuery" @keyup.enter="handleSearch" type="text" placeholder="ابحث عن فعاليات..."
              class="w-full  py-6 text-gray-200 !ring-0 !border-0 rounded-full px-10" />
            <span class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 cursor-pointer"
              @click="handleSearch">
              <SearchIcon class="w-5 h-5" />
            </span>
          </div>
          <a href="#" class="hover:text-purple-200 transition">الفعاليات</a>
          <a href="#" class="hover:text-purple-200 transition">الدعم والتواصل</a>

          <div class="flex space-x-4 pt-2 justify-end">
            <!-- اللغة -->
            <Button variant="link" size="icon" class="text-gray-300 hover:text-purple-300">
              <LanguagesIcon />
            </Button>

            <!-- الحساب -->
            <Button v-if="!isLoggedIn" @click="openLoginPage" asChild variant="ghost" size="icon"
              class="relative size-10 w-auto rounded-full p-1 focus-within:ring-2 focus-within:ring-primary">
              <Avatar class="size-10 overflow-hidden rounded-full">
                <AvatarFallback
                  class="flex items-center justify-center rounded-full bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
                  <User class="w-5 h-5" /> <!-- Fallback user icon -->
                </AvatarFallback>
              </Avatar>
            </Button>
            <DropdownMenu v-else>
              <DropdownMenuTrigger :as-child="true">
                <Button variant="ghost" size="icon"
                  class="relative size-10 w-auto rounded-full p-1 focus-within:ring-2 focus-within:ring-primary">
                  <Avatar class="size-8 overflow-hidden rounded-full">
                    <AvatarImage v-if="page.props.auth?.user?.avatar" :src="page.props.auth?.user?.avatar"
                      :alt="page.props.auth?.user?.name" />
                    <AvatarFallback
                      class="flex items-center justify-center rounded-lg bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
                      {{ getInitials(page.props.auth?.user?.name) }}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start" class="w-56">
                <UserMenuContent :user="page.props.auth?.user" />
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </transition>
    </nav>
  </Transition>

</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from "vue"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { SearchIcon, MenuIcon, XIcon, TicketIcon, LanguagesIcon, UserIcon, ChevronsUpDown, User, CalendarIcon } from "lucide-vue-next"
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"
import { router, usePage } from "@inertiajs/vue3";
import DropdownMenu from "@/components/ui/dropdown-menu/DropdownMenu.vue"
import DropdownMenuTrigger from "@/components/ui/dropdown-menu/DropdownMenuTrigger.vue"
import DropdownMenuContent from "@/components/ui/dropdown-menu/DropdownMenuContent.vue"
import UserMenuContent from "@/components/UserMenuContent.vue"
import Avatar from "@/components/ui/avatar/Avatar.vue"
import AvatarImage from "@/components/ui/avatar/AvatarImage.vue"
import AvatarFallback from "@/components/ui/avatar/AvatarFallback.vue"
import { getInitials } from "@/composables/useInitials"
import Logo from "@/components/Logo.vue"

const page = usePage();
const isLoggedIn = computed(() => page.props.auth?.user != null);
const searchQuery = ref('');

const openLoginPage = () => {
  router.get(`/login`, {}, {
    preserveScroll: true,
    preserveState: true,
  });
};

const openMainPage = () => {
  router.get(`/`, {}, {
    preserveScroll: true,
    preserveState: true,
  });
};

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.get('/events/list', { search: searchQuery.value.trim() }, {
      preserveScroll: true,
      preserveState: true,
    });
  }
};

const stuck = ref(false)

const handleScroll = () => {
  stuck.value = window.scrollY > 4
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})


const menuOpen = ref(false)
</script>
