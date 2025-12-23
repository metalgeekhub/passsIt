<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { userApi } from '@/api/services';
import { useToast } from 'vue-toastification';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

const router = useRouter();
const route = useRoute();
const toast = useToast();
const userId = route.params.id as string;

const form = reactive({
  email: '',
  username: '',
  firstName: '',
  lastName: '',
  isAdmin: false,
});

const state = reactive({
  isLoading: true,
  isSubmitting: false,
  errors: {} as Record<string, string>,
});

onMounted(async () => {
  try {
    const user = await userApi.getById(userId);
    form.email = user.email;
    form.username = user.username;
    form.firstName = user.first_name;
    form.lastName = user.last_name;
    form.isAdmin = user.is_admin;
  } catch (error: any) {
    console.error('Error fetching user:', error);
    toast.error('Failed to load user');
    router.push('/users');
  } finally {
    state.isLoading = false;
  }
});

const validateForm = () => {
  state.errors = {};
  
  if (!form.email) state.errors.email = 'Email is required';
  else if (!/\S+@\S+\.\S+/.test(form.email)) state.errors.email = 'Email is invalid';
  
  if (!form.username) state.errors.username = 'Username is required';
  else if (form.username.length < 3) state.errors.username = 'Username must be at least 3 characters';
  
  if (!form.firstName) state.errors.firstName = 'First name is required';
  if (!form.lastName) state.errors.lastName = 'Last name is required';
  
  return Object.keys(state.errors).length === 0;
};

const handleSubmit = async () => {
  if (!validateForm()) {
    toast.error('Please fix the errors in the form');
    return;
  }
  
  state.isSubmitting = true;
  
  try {
    await userApi.update(userId, {
      email: form.email,
      username: form.username,
      firstName: form.firstName,
      lastName: form.lastName,
      isAdmin: form.isAdmin,
    });
    
    toast.success('User updated successfully!');
    router.push('/users');
  } catch (error: any) {
    console.error('Error updating user:', error);
    const errorMessage = error.response?.data?.error || 'Failed to update user';
    toast.error(errorMessage);
  } finally {
    state.isSubmitting = false;
  }
};
</script>

<template>
  <section class="bg-blue-50 px-4 py-10 min-h-screen">
    <div class="container mx-auto max-w-2xl">
      <!-- Loading State -->
      <div v-if="state.isLoading" class="text-center py-20">
        <PulseLoader color="#2563eb" />
        <p class="mt-4 text-gray-600">Loading user...</p>
      </div>

      <!-- Edit Form -->
      <div v-else class="bg-white rounded-lg shadow-md p-8">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-blue-700">Edit User</h2>
          <p class="text-gray-600 mt-2">Update user information</p>
        </div>

        <form @submit.prevent="handleSubmit">
          <!-- Email -->
          <div class="mb-4">
            <label for="email" class="block text-gray-700 font-bold mb-2">
              Email <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.email"
              type="email"
              id="email"
              class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': state.errors.email }"
              placeholder="user@example.com"
            />
            <p v-if="state.errors.email" class="text-red-500 text-sm mt-1">{{ state.errors.email }}</p>
          </div>

          <!-- Username -->
          <div class="mb-4">
            <label for="username" class="block text-gray-700 font-bold mb-2">
              Username <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.username"
              type="text"
              id="username"
              class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': state.errors.username }"
              placeholder="johndoe"
            />
            <p v-if="state.errors.username" class="text-red-500 text-sm mt-1">{{ state.errors.username }}</p>
          </div>

          <!-- First Name -->
          <div class="mb-4">
            <label for="firstName" class="block text-gray-700 font-bold mb-2">
              First Name <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.firstName"
              type="text"
              id="firstName"
              class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': state.errors.firstName }"
              placeholder="John"
            />
            <p v-if="state.errors.firstName" class="text-red-500 text-sm mt-1">{{ state.errors.firstName }}</p>
          </div>

          <!-- Last Name -->
          <div class="mb-4">
            <label for="lastName" class="block text-gray-700 font-bold mb-2">
              Last Name <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.lastName"
              type="text"
              id="lastName"
              class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': state.errors.lastName }"
              placeholder="Doe"
            />
            <p v-if="state.errors.lastName" class="text-red-500 text-sm mt-1">{{ state.errors.lastName }}</p>
          </div>

          <!-- Admin Checkbox -->
          <div class="mb-6">
            <label class="flex items-center cursor-pointer">
              <input
                v-model="form.isAdmin"
                type="checkbox"
                id="isAdmin"
                class="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
              />
              <span class="ml-3 text-gray-700 font-bold">Administrator privileges</span>
            </label>
            <p class="text-gray-600 text-sm mt-1 ml-8">Administrators have full access to all features</p>
          </div>

          <!-- Buttons -->
          <div class="flex gap-4">
            <button
              type="submit"
              :disabled="state.isSubmitting"
              class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="!state.isSubmitting">Update User</span>
              <span v-else>Updating...</span>
            </button>
            <button
              type="button"
              @click="router.push('/users')"
              class="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-3 px-6 rounded-lg transition"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </section>
</template>
