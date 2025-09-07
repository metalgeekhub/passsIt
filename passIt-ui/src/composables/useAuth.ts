import { ref } from 'vue'
import authAxios from '../api/axios';

export const isAuthenticated = ref(false)

export async function checkAuth() {
  try {
    await authAxios.get('/users')
    isAuthenticated.value = true
  } catch {
    isAuthenticated.value = false
  }
}