<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { userApi } from '@/api/services';
import { isAdmin } from '@/composables/useAuth';
import type { User } from '@/models/User';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

const state = reactive({
  users: [] as User[],
  isLoading: true,
  error: null as string | null,
});

onMounted(async () => {
  try {
    state.users = await userApi.getAll();
  } catch (error: any) {
    console.error('Error fetching users:', error);
    state.error = error.response?.data?.error || 'Failed to fetch users';
  } finally {
    state.isLoading = false;
  }
});
</script>

<template>
  <section class="bg-blue-50 px-4 py-10 min-h-screen">
    <div class="container mx-auto max-w-6xl">
      <div class="flex justify-between items-center mb-8">
        <h2 class="text-3xl font-bold text-blue-700">User Management</h2>
        <RouterLink
          v-if="isAdmin"
          to="/users/add"
          class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition"
        >
          <i class="fa fa-plus mr-2"></i>Add User
        </RouterLink>
      </div>

      <!-- Loading State -->
      <div v-if="state.isLoading" class="text-center py-20">
        <PulseLoader color="#2563eb" />
        <p class="mt-4 text-gray-600">Loading users...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="state.error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        <p>{{ state.error }}</p>
      </div>

      <!-- Users Table -->
      <div v-else class="bg-white rounded-lg shadow-md overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-blue-700 text-white">
            <tr>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap">Username</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden sm:table-cell">Email</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden md:table-cell">First Name</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden md:table-cell">Last Name</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap">Role</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden lg:table-cell">Created At</th>
              <th class="px-3 sm:px-6 py-3 text-right text-xs font-medium uppercase tracking-wider whitespace-nowrap">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="user in state.users" :key="user.id" class="hover:bg-blue-50 transition">
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ user.username }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-600 hidden sm:table-cell">{{ user.email }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-600 hidden md:table-cell">{{ user.first_name }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-600 hidden md:table-cell">{{ user.last_name }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm">
                <span v-if="user.is_admin" class="px-2 py-1 text-xs font-semibold rounded-full bg-purple-100 text-purple-800">
                  <i class="fa fa-shield"></i> Admin
                </span>
                <span v-else class="px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-600">
                  User
                </span>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-600 hidden lg:table-cell">
                {{ new Date(user.created_at).toLocaleDateString() }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-1 sm:space-x-2">
                <RouterLink
                  v-if="isAdmin"
                  :to="`/users/edit/${user.id}`"
                  class="text-blue-600 hover:text-blue-900 font-semibold text-xs sm:text-sm"
                >
                  <i class="fa fa-edit"></i> Edit
                </RouterLink>
                <span v-else class="text-gray-400 text-xs sm:text-sm">-</span>
              </td>
            </tr>
            <tr v-if="state.users.length === 0">
              <td colspan="7" class="px-6 py-8 text-center text-gray-500">
                No users found. Click "Add User" to create one.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>
