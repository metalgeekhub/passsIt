<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router';
import logo from '@/assets/images/logo.png';
import { isAuthenticated, isAdmin } from '@/composables/useAuth';
import { onMounted } from 'vue';

const isActive = (path: string) => {
  const route = useRoute();
  return route.path === path;
}

onMounted(async () => {
  console.log(isAuthenticated)
});

</script>

<template>
    <nav class="bg-blue-700 border-b border-blue-500">
      <div class="mx-auto max-w-7xl px-2 sm:px-4 lg:px-8">
        <div class="flex h-16 sm:h-20 items-center justify-between">
          <div
            class="flex flex-1 items-center justify-center md:items-stretch md:justify-start"
          >
            <!-- Logo -->
            <RouterLink class="flex flex-shrink-0 items-center mr-2 sm:mr-4" to="/">
              <img class="h-8 sm:h-10 w-auto" :src="logo" alt="PassIt Logo" />
              <span class="hidden sm:block text-white text-xl sm:text-2xl font-bold ml-2"
                >PassIt</span
              >
            </RouterLink>
            <div class="md:ml-auto">
              <div class="flex space-x-1 sm:space-x-2">
                <RouterLink
                  to="/"
                  :class="[isActive('/') ? 'bg-blue-900' : 'hover:bg-gray-900 hover:text-white', 'text-white', 'rounded-md', 'px-2 sm:px-3', 'py-1.5 sm:py-2', 'text-sm sm:text-base']"
                  >Home</RouterLink>
                <RouterLink
                  v-if="isAuthenticated"
                  to="/profile"
                  :class="[isActive('/profile') ? 'bg-blue-900' : 'hover:bg-gray-900 hover:text-white', 'text-white', 'rounded-md', 'px-2 sm:px-3', 'py-1.5 sm:py-2', 'text-sm sm:text-base']"
                  >Profile</RouterLink>
                <RouterLink
                  v-if="isAuthenticated && isAdmin"
                  to="/users"
                  :class="[isActive('/users') ? 'bg-blue-900' : 'hover:bg-gray-900 hover:text-white', 'text-white', 'rounded-md', 'px-2 sm:px-3', 'py-1.5 sm:py-2', 'text-sm sm:text-base']"
                  >Users</RouterLink>
                <RouterLink
                  v-if="!isAuthenticated"
                  to="/login"
                  :class="[isActive('/login') ? 'bg-blue-50' : 'hover:bg-blue-100 hover:text-blue-900', 'bg-blue-50 text-blue-700 font-bold shadow-sm transition', 'rounded-md', 'px-2 sm:px-3', 'py-1.5 sm:py-2', 'text-sm sm:text-base']"
                  >Login</RouterLink>
                <RouterLink
                  v-if="!isAuthenticated"
                  to="/signup"
                  :class="[isActive('/signup') ? 'bg-blue-800' : 'hover:bg-blue-800', 'bg-blue-600 text-white font-bold shadow-sm transition', 'rounded-md', 'px-2 sm:px-3', 'py-1.5 sm:py-2', 'text-sm sm:text-base']"
                  >Sign Up</RouterLink>
                <RouterLink
                  v-if="isAuthenticated"
                  to="/logout"
                  :class="[isActive('/logout') ? 'bg-blue-50' : 'hover:bg-blue-100 hover:text-blue-900', 'bg-blue-50 text-blue-700 font-bold shadow-sm transition', 'rounded-md', 'px-2 sm:px-3', 'py-1.5 sm:py-2', 'text-sm sm:text-base']"
                  >Logout</RouterLink>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
</template>