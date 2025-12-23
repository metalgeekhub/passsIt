<script setup lang="ts">
import { authApi } from '@/api/services';
authApi.login();
</script>

<template>
    <div>Redirecting to Login...</div>
</template>