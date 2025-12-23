<script setup lang="ts">
import { reactive, onMounted, ref } from 'vue';
import { authApi, userApi } from '@/api/services';
import { useToast } from 'vue-toastification';
import { currentUser } from '@/composables/useAuth';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

const toast = useToast();
const loading = ref(true);

const form = reactive({
  email: '',
  username: '',
  firstName: '',
  lastName: '',
  password: '',
  confirmPassword: '',
});

const state = reactive({
  isSubmitting: false,
  errors: {} as Record<string, string>,
});

const loadProfile = async () => {
  try {
    loading.value = true;
    const user = await authApi.getCurrentUser();
    form.email = user.email;
    form.username = user.username;
    form.firstName = user.first_name || '';
    form.lastName = user.last_name || '';
  } catch (error) {
    toast.error('Failed to load profile');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const validateForm = () => {
  state.errors = {};
  
  if (!form.email) state.errors.email = 'Email is required';
  else if (!/\S+@\S+\.\S+/.test(form.email)) state.errors.email = 'Email is invalid';
  
  if (!form.username) state.errors.username = 'Username is required';
  else if (form.username.length < 3) state.errors.username = 'Username must be at least 3 characters';
  
  if (!form.firstName) state.errors.firstName = 'First name is required';
  if (!form.lastName) state.errors.lastName = 'Last name is required';
  
  // Password is optional, but if provided, validate it
  if (form.password) {
    if (form.password.length < 8) {
      state.errors.password = 'Password must be at least 8 characters';
    } else if (form.password !== form.confirmPassword) {
      state.errors.confirmPassword = 'Passwords do not match';
    }
  }
  
  return Object.keys(state.errors).length === 0;
};

const handleSubmit = async () => {
  if (!validateForm()) {
    toast.error('Please fix the errors in the form');
    return;
  }
  
  if (!currentUser.value?.id) {
    toast.error('User ID not found');
    return;
  }
  
  state.isSubmitting = true;
  
  try {
    const updateData: any = {
      email: form.email,
      username: form.username,
      firstName: form.firstName,
      lastName: form.lastName,
    };

    // Only include password if it was changed
    if (form.password) {
      updateData.password = form.password;
    }

    await userApi.update(currentUser.value.id, updateData);
    
    toast.success('Profile updated successfully');
    
    // Clear password fields
    form.password = '';
    form.confirmPassword = '';
    
    // Reload profile to get updated data
    await loadProfile();
  } catch (error: any) {
    toast.error(error.response?.data?.error || 'Failed to update profile');
    console.error(error);
  } finally {
    state.isSubmitting = false;
  }
};

onMounted(() => {
  loadProfile();
});
</script>

<template>
  <div class="container mx-auto px-4 py-8 max-w-2xl">
    <div class="bg-white shadow-md rounded-lg px-8 pt-6 pb-8 mb-4">
      <h2 class="text-2xl font-bold mb-6 text-gray-800">My Profile</h2>
      
      <div v-if="loading" class="flex justify-center py-8">
        <PulseLoader color="#3b82f6" />
      </div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-4">
        <!-- Username -->
        <div>
          <label for="username" class="block text-gray-700 text-sm font-bold mb-2">
            Username *
          </label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500"
            :class="{ 'border-red-500': state.errors.username }"
          />
          <p v-if="state.errors.username" class="text-red-500 text-xs italic mt-1">
            {{ state.errors.username }}
          </p>
        </div>

        <!-- Email -->
        <div>
          <label for="email" class="block text-gray-700 text-sm font-bold mb-2">
            Email *
          </label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500"
            :class="{ 'border-red-500': state.errors.email }"
          />
          <p v-if="state.errors.email" class="text-red-500 text-xs italic mt-1">
            {{ state.errors.email }}
          </p>
        </div>

        <!-- First Name -->
        <div>
          <label for="firstName" class="block text-gray-700 text-sm font-bold mb-2">
            First Name *
          </label>
          <input
            id="firstName"
            v-model="form.firstName"
            type="text"
            required
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500"
            :class="{ 'border-red-500': state.errors.firstName }"
          />
          <p v-if="state.errors.firstName" class="text-red-500 text-xs italic mt-1">
            {{ state.errors.firstName }}
          </p>
        </div>

        <!-- Last Name -->
        <div>
          <label for="lastName" class="block text-gray-700 text-sm font-bold mb-2">
            Last Name *
          </label>
          <input
            id="lastName"
            v-model="form.lastName"
            type="text"
            required
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500"
            :class="{ 'border-red-500': state.errors.lastName }"
          />
          <p v-if="state.errors.lastName" class="text-red-500 text-xs italic mt-1">
            {{ state.errors.lastName }}
          </p>
        </div>

        <!-- Password (Optional) -->
        <div class="pt-4 border-t">
          <h3 class="text-lg font-semibold mb-4 text-gray-700">Change Password (Optional)</h3>
          
          <div class="mb-4">
            <label for="password" class="block text-gray-700 text-sm font-bold mb-2">
              New Password
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              minlength="8"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': state.errors.password }"
            />
            <p v-if="state.errors.password" class="text-red-500 text-xs italic mt-1">
              {{ state.errors.password }}
            </p>
            <p class="text-gray-500 text-xs mt-1">Leave blank to keep current password</p>
          </div>

          <div>
            <label for="confirmPassword" class="block text-gray-700 text-sm font-bold mb-2">
              Confirm New Password
            </label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500"
              :class="{ 'border-red-500': state.errors.confirmPassword }"
            />
            <p v-if="state.errors.confirmPassword" class="text-red-500 text-xs italic mt-1">
              {{ state.errors.confirmPassword }}
            </p>
          </div>
        </div>

        <!-- Submit Button -->
        <div class="flex items-center justify-end pt-4">
          <button
            type="submit"
            :disabled="state.isSubmitting"
            class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded focus:outline-none focus:shadow-outline disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ state.isSubmitting ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
