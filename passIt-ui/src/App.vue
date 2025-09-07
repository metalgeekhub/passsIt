<script lang="ts" setup>
import Navbar from './components/Navbar.vue';
import { RouterView } from 'vue-router';
import { checkAuth } from '@/composables/useAuth'

checkAuth()
</script>

<template >
  <Navbar />
  <RouterView />
</template>