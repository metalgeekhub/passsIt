<script setup lang="ts">
import { authApi } from '@/api/services';
authApi.logout();
</script>

<template>
    <div>Logged Out!</div>
</template>