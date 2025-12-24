<script setup lang="ts">
import { reactive, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { userApi } from '@/api/services';
import { useToast } from 'vue-toastification';
import { isAdmin } from '@/composables/useAuth';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';
import type { User } from '@/models/User';

const router = useRouter();
const route = useRoute();
const toast = useToast();
const userId = route.params.id as string;

const state = reactive({
  user: null as User | null,
  isLoading: true,
  isEditing: false,
  isSaving: false,
  errors: {} as Record<string, string>,
});

const form = reactive({
  email: '',
  username: '',
  firstName: '',
  lastName: '',
  dob: '',
  phoneNumber: '',
  address: '',
  isAdmin: false,
});

const loadUser = async () => {
  try {
    state.isLoading = true;
    state.user = await userApi.getById(userId);
    
    // Populate form with user data
    form.email = state.user.email;
    form.username = state.user.username;
    form.firstName = state.user.first_name;
    form.lastName = state.user.last_name;
    form.dob = state.user.dob ? state.user.dob.split('T')[0] : '';
    form.phoneNumber = state.user.phone_number;
    form.address = state.user.address;
    form.isAdmin = state.user.is_admin;
  } catch (error: any) {
    console.error('Error fetching user:', error);
    toast.error('Failed to load user');
    router.push('/users');
  } finally {
    state.isLoading = false;
  }
};

const toggleEdit = () => {
  if (state.isEditing) {
    // Cancel editing - reset form
    if (state.user) {
      form.email = state.user.email;
      form.username = state.user.username;
      form.firstName = state.user.first_name;
      form.lastName = state.user.last_name;
      form.dob = state.user.dob ? state.user.dob.split('T')[0] : '';
      form.phoneNumber = state.user.phone_number;
      form.address = state.user.address;
      form.isAdmin = state.user.is_admin;
    }
    state.errors = {};
  }
  state.isEditing = !state.isEditing;
};

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

const handleSave = async () => {
  if (!validateForm()) {
    toast.error('Please fix the errors in the form');
    return;
  }
  
  state.isSaving = true;
  
  try {
    await userApi.update(userId, {
      email: form.email,
      username: form.username,
      firstName: form.firstName,
      lastName: form.lastName,
      dob: form.dob,
      phoneNumber: form.phoneNumber,
      address: form.address,
      isAdmin: form.isAdmin,
    });
    
    toast.success('User updated successfully!');
    state.isEditing = false;
    await loadUser(); // Reload user data
  } catch (error: any) {
    console.error('Error updating user:', error);
    const errorMessage = error.response?.data?.error || 'Failed to update user';
    toast.error(errorMessage);
  } finally {
    state.isSaving = false;
  }
};

const formatDate = (dateString: string) => {
  if (!dateString) return 'Not set';
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

const formatDateTime = (dateString: string) => {
  if (!dateString) return 'Not set';
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

onMounted(async () => {
  await loadUser();
});
</script>

<template>
  <section class="bg-blue-50 px-4 py-10 min-h-screen">
    <div class="container mx-auto max-w-4xl">
      <!-- Loading State -->
      <div v-if="state.isLoading" class="text-center py-20">
        <PulseLoader color="#2563eb" />
        <p class="mt-4 text-gray-600">Loading user profile...</p>
      </div>

      <!-- User Profile -->
      <div v-else-if="state.user" class="bg-white rounded-lg shadow-md overflow-hidden">
        <!-- Header -->
        <div class="bg-blue-700 text-white px-8 py-6">
          <div class="flex justify-between items-start">
            <div>
              <h1 class="text-3xl font-bold">{{ state.user.username }}</h1>
              <p class="text-blue-200 mt-1">{{ state.user.email }}</p>
              <div class="mt-3 flex items-center gap-3">
                <span v-if="state.user.is_admin" class="px-3 py-1 text-xs font-semibold rounded-full bg-purple-500">
                  <i class="fa fa-shield"></i> Administrator
                </span>
                <span :class="state.user.is_active ? 'bg-green-500' : 'bg-red-500'" class="px-3 py-1 text-xs font-semibold rounded-full">
                  <i class="fa" :class="state.user.is_active ? 'fa-check-circle' : 'fa-times-circle'"></i>
                  {{ state.user.is_active ? 'Active' : 'Inactive' }}
                </span>
              </div>
            </div>
            <div class="flex gap-2">
              <button
                @click="router.push('/users')"
                class="bg-blue-800 hover:bg-blue-900 text-white px-4 py-2 rounded-lg transition"
              >
                <i class="fa fa-arrow-left mr-2"></i>Back to Users
              </button>
              <button
                v-if="isAdmin && !state.isEditing"
                @click="toggleEdit"
                class="bg-white text-blue-700 hover:bg-blue-50 px-4 py-2 rounded-lg transition font-semibold"
              >
                <i class="fa fa-edit mr-2"></i>Edit Profile
              </button>
            </div>
          </div>
        </div>

        <!-- Profile Content -->
        <div class="p-8">
          <!-- View Mode -->
          <div v-if="!state.isEditing" class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Personal Information -->
            <div class="col-span-2">
              <h3 class="text-xl font-bold text-gray-800 mb-4 border-b pb-2">
                <i class="fa fa-user mr-2 text-blue-600"></i>Personal Information
              </h3>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">First Name</label>
              <p class="text-gray-900 text-lg">{{ state.user.first_name || 'Not set' }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Last Name</label>
              <p class="text-gray-900 text-lg">{{ state.user.last_name || 'Not set' }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Username</label>
              <p class="text-gray-900 text-lg">{{ state.user.username }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Email Address</label>
              <p class="text-gray-900 text-lg">{{ state.user.email }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Date of Birth</label>
              <p class="text-gray-900 text-lg">{{ formatDate(state.user.dob) }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Phone Number</label>
              <p class="text-gray-900 text-lg">{{ state.user.phone_number || 'Not set' }}</p>
            </div>

            <div class="col-span-2 space-y-1">
              <label class="text-sm font-semibold text-gray-600">Address</label>
              <p class="text-gray-900 text-lg">{{ state.user.address || 'Not set' }}</p>
            </div>

            <!-- System Information -->
            <div class="col-span-2 mt-6">
              <h3 class="text-xl font-bold text-gray-800 mb-4 border-b pb-2">
                <i class="fa fa-info-circle mr-2 text-blue-600"></i>System Information
              </h3>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">User ID</label>
              <p class="text-gray-900 text-sm font-mono">{{ state.user.id }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Keycloak ID</label>
              <p class="text-gray-900 text-sm font-mono">{{ state.user.keycloack_id }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Created At</label>
              <p class="text-gray-900">{{ formatDateTime(state.user.created_at) }}</p>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-semibold text-gray-600">Last Updated</label>
              <p class="text-gray-900">{{ formatDateTime(state.user.updated_at) }}</p>
            </div>
          </div>

          <!-- Edit Mode -->
          <form v-else @submit.prevent="handleSave" class="space-y-6">
            <!-- Personal Information -->
            <div>
              <h3 class="text-xl font-bold text-gray-800 mb-4 border-b pb-2">
                <i class="fa fa-user mr-2 text-blue-600"></i>Personal Information
              </h3>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- First Name -->
                <div>
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
                <div>
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

                <!-- Username -->
                <div>
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

                <!-- Email -->
                <div>
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

                <!-- Date of Birth -->
                <div>
                  <label for="dob" class="block text-gray-700 font-bold mb-2">
                    Date of Birth
                  </label>
                  <input
                    v-model="form.dob"
                    type="date"
                    id="dob"
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <!-- Phone Number -->
                <div>
                  <label for="phoneNumber" class="block text-gray-700 font-bold mb-2">
                    Phone Number
                  </label>
                  <input
                    v-model="form.phoneNumber"
                    type="tel"
                    id="phoneNumber"
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="+1 (555) 123-4567"
                  />
                </div>

                <!-- Address -->
                <div class="md:col-span-2">
                  <label for="address" class="block text-gray-700 font-bold mb-2">
                    Address
                  </label>
                  <textarea
                    v-model="form.address"
                    id="address"
                    rows="3"
                    class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="123 Main St, City, State, ZIP"
                  ></textarea>
                </div>

                <!-- Admin Checkbox -->
                <div class="md:col-span-2">
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
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex gap-4 pt-4 border-t">
              <button
                type="submit"
                :disabled="state.isSaving"
                class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg transition disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span v-if="!state.isSaving"><i class="fa fa-save mr-2"></i>Save Changes</span>
                <span v-else><i class="fa fa-spinner fa-spin mr-2"></i>Saving...</span>
              </button>
              <button
                type="button"
                @click="toggleEdit"
                :disabled="state.isSaving"
                class="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-3 px-6 rounded-lg transition disabled:opacity-50"
              >
                <i class="fa fa-times mr-2"></i>Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </section>
</template>
