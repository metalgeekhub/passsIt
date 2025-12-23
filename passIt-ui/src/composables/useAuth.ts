import { ref } from 'vue'
import { authApi } from '@/api/services';

export const isAuthenticated = ref(false)
export const isAdmin = ref(false)
export const currentUser = ref<any>(null)

export async function checkAuth() {
  try {
    const user = await authApi.getCurrentUser();
    isAuthenticated.value = true;
    isAdmin.value = user.is_admin || false;
    currentUser.value = user;
    return true;
  } catch {
    isAuthenticated.value = false;
    isAdmin.value = false;
    currentUser.value = null;
    return false;
  }
}