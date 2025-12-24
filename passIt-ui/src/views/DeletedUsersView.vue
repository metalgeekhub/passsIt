<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { RouterLink, useRouter } from 'vue-router';
import { userApi } from '@/api/services';
import type { User } from '@/models/User';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';

const router = useRouter();

const state = reactive({
  users: [] as User[],
  isLoading: true,
  error: null as string | null,
});

const goToUserProfile = (userId: string) => {
  router.push(`/users/${userId}`);
};

onMounted(async () => {
  try {
    state.users = await userApi.getInactive();
  } catch (error: any) {
    console.error('Error fetching deleted users:', error);
    state.error = error.response?.data?.error || 'Failed to fetch deleted users';
  } finally {
    state.isLoading = false;
  }
});
</script>

<template>
  <section class="bg-red-50 px-4 py-10 min-h-screen">
    <div class="container mx-auto max-w-6xl">
      <div class="flex justify-between items-center mb-8">
        <h2 class="text-3xl font-bold text-red-700">
          <i class="fa fa-trash mr-2"></i>Deleted Users
        </h2>
        <RouterLink
          to="/users"
          class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition"
        >
          <i class="fa fa-arrow-left mr-2"></i>Back to Active Users
        </RouterLink>
      </div>

      <!-- Loading State -->
      <div v-if="state.isLoading" class="text-center py-20">
        <PulseLoader color="#dc2626" />
        <p class="mt-4 text-gray-600">Loading deleted users...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="state.error" class="bg-red-100 border border-red-300 text-red-700 px-4 py-3 rounded">
        <p>{{ state.error }}</p>
      </div>

      <!-- Users Table -->
      <div v-else class="bg-white rounded-lg shadow-md overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-red-700 text-white">
            <tr>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap">Username</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden sm:table-cell">Email</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden md:table-cell">First Name</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden md:table-cell">Last Name</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap">Role</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden lg:table-cell">Created At</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium uppercase tracking-wider whitespace-nowrap hidden lg:table-cell">Deleted At</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr 
              v-for="user in state.users" 
              :key="user.id" 
              @click="goToUserProfile(user.id)"
              class="bg-red-50 opacity-75 cursor-pointer hover:opacity-100 transition"
            >
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                <i class="fa fa-ban text-red-500 mr-2"></i>{{ user.username }}
              </td>
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
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-600 hidden lg:table-cell">
                {{ new Date(user.updated_at).toLocaleDateString() }}
              </td>
            </tr>
            <tr v-if="state.users.length === 0">
              <td colspan="7" class="px-6 py-8 text-center text-gray-500">
                <i class="fa fa-check-circle text-green-500 text-3xl mb-2"></i>
                <p class="text-lg">No deleted users. All users are active.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>
