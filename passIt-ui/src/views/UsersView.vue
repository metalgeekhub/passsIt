<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { RouterLink, useRouter } from 'vue-router';
import { userApi } from '@/api/services';
import { isAdmin } from '@/composables/useAuth';
import type { User } from '@/models/User';
import PulseLoader from 'vue-spinner/src/PulseLoader.vue';
import ConfirmDialog from '@/components/ConfirmDialog.vue';

const router = useRouter();

const state = reactive({
  users: [] as User[],
  isLoading: true,
  error: null as string | null,
  deletingUserId: null as string | null,
  showConfirmDialog: false,
  userToDelete: null as User | null,
});

const loadUsers = async () => {
  try {
    state.isLoading = true;
    state.users = await userApi.getAll();
  } catch (error: any) {
    console.error('Error fetching users:', error);
    state.error = error.response?.data?.error || 'Failed to fetch users';
  } finally {
    state.isLoading = false;
  }
};

const showDeleteConfirmation = (user: User, event: Event) => {
  event.stopPropagation(); // Prevent row click when clicking delete
  state.userToDelete = user;
  state.showConfirmDialog = true;
};

const goToUserProfile = (userId: string) => {
  router.push(`/users/${userId}`);
};

const cancelDelete = () => {
  state.showConfirmDialog = false;
  state.userToDelete = null;
};

const confirmDelete = async () => {
  if (!state.userToDelete) return;

  try {
    state.deletingUserId = state.userToDelete.id;
    state.showConfirmDialog = false;
    
    await userApi.delete(state.userToDelete.id);
    // Remove user from list
    state.users = state.users.filter(u => u.id !== state.userToDelete!.id);
    state.userToDelete = null;
  } catch (error: any) {
    console.error('Error deleting user:', error);
    state.error = error.response?.data?.error || 'Failed to delete user';
  } finally {
    state.deletingUserId = null;
  }
};

onMounted(async () => {
  await loadUsers();
});
</script>

<template>
  <section class="bg-blue-50 px-4 py-10 min-h-screen">
    <div class="container mx-auto max-w-6xl">
      <div class="flex justify-between items-center mb-8">
        <h2 class="text-3xl font-bold text-blue-700">User Management</h2>
        <div class="flex gap-3">
          <RouterLink
            to="/users/deleted"
            class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition"
          >
            <i class="fa fa-trash mr-2"></i>Deleted Users
          </RouterLink>
          <RouterLink
            v-if="isAdmin"
            to="/users/add"
            class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg shadow-md transition"
          >
            <i class="fa fa-plus mr-2"></i>Add User
          </RouterLink>
        </div>
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
            <tr 
              v-for="user in state.users" 
              :key="user.id" 
              @click="goToUserProfile(user.id)"
              class="hover:bg-blue-50 transition cursor-pointer"
            >
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
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  v-if="isAdmin"
                  @click="showDeleteConfirmation(user, $event)"
                  :disabled="state.deletingUserId === user.id"
                  class="text-red-600 hover:text-red-900 font-semibold text-xs sm:text-sm disabled:opacity-50"
                >
                  <i class="fa" :class="state.deletingUserId === user.id ? 'fa-spinner fa-spin' : 'fa-trash'"></i>
                  Delete
                </button>
                <span v-if="!isAdmin" class="text-gray-400 text-xs sm:text-sm">-</span>
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

    <!-- Confirmation Dialog -->
    <ConfirmDialog
      :show="state.showConfirmDialog"
      title="Deactivate User"
      :message="`Are you sure you want to deactivate ${state.userToDelete?.username}? This will disable their access.`"
      confirm-text="Deactivate"
      cancel-text="Cancel"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />
  </section>
</template>
