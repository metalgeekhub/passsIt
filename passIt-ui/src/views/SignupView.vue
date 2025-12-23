<template>
  <div class="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md mx-auto">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-extrabold text-gray-900">Create your account</h2>
        <p class="mt-2 text-sm text-gray-600">
          Join PassIt to start attending amazing events
        </p>
      </div>

      <div class="bg-white py-8 px-6 shadow-lg rounded-lg">
        <form @submit.prevent="handleSignup" class="space-y-6">
          <!-- Username -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700">
              Username *
            </label>
            <input
              id="username"
              v-model="formData.username"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.username }"
            />
            <p v-if="errors.username" class="mt-1 text-sm text-red-600">{{ errors.username }}</p>
          </div>

          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">
              Email *
            </label>
            <input
              id="email"
              v-model="formData.email"
              type="email"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.email }"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
          </div>

          <!-- First Name -->
          <div>
            <label for="firstName" class="block text-sm font-medium text-gray-700">
              First Name
            </label>
            <input
              id="firstName"
              v-model="formData.firstName"
              type="text"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <!-- Last Name -->
          <div>
            <label for="lastName" class="block text-sm font-medium text-gray-700">
              Last Name
            </label>
            <input
              id="lastName"
              v-model="formData.lastName"
              type="text"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <!-- Password -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700">
              Password *
            </label>
            <input
              id="password"
              v-model="formData.password"
              type="password"
              required
              minlength="8"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.password }"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
            <p class="mt-1 text-xs text-gray-500">Must be at least 8 characters</p>
          </div>

          <!-- Confirm Password -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
              Confirm Password *
            </label>
            <input
              id="confirmPassword"
              v-model="formData.confirmPassword"
              type="password"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              :class="{ 'border-red-500': errors.confirmPassword }"
            />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{ errors.confirmPassword }}</p>
          </div>

          <!-- Error Message -->
          <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="ml-3">
                <p class="text-sm text-red-800">{{ errorMessage }}</p>
              </div>
            </div>
          </div>

          <!-- Buttons -->
          <div class="flex gap-4">
            <button
              type="submit"
              :disabled="loading"
              class="flex-1 flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading">Creating account...</span>
              <span v-else>Sign Up</span>
            </button>
            <button
              type="button"
              @click="$router.push('/')"
              class="flex-1 py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Cancel
            </button>
          </div>

          <!-- Login Link -->
          <div class="text-center">
            <p class="text-sm text-gray-600">
              Already have an account?
              <router-link to="/login" class="font-medium text-blue-600 hover:text-blue-500">
                Sign in
              </router-link>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/services'
import { toast } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'

const router = useRouter()
const loading = ref(false)
const errorMessage = ref('')

const formData = reactive({
  username: '',
  email: '',
  firstName: '',
  lastName: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validateForm = (): boolean => {
  // Reset errors
  errors.username = ''
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''
  errorMessage.value = ''

  let isValid = true

  if (!formData.username.trim()) {
    errors.username = 'Username is required'
    isValid = false
  }

  if (!formData.email.trim()) {
    errors.email = 'Email is required'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
    errors.email = 'Please enter a valid email address'
    isValid = false
  }

  if (!formData.password) {
    errors.password = 'Password is required'
    isValid = false
  } else if (formData.password.length < 8) {
    errors.password = 'Password must be at least 8 characters'
    isValid = false
  }

  if (!formData.confirmPassword) {
    errors.confirmPassword = 'Please confirm your password'
    isValid = false
  } else if (formData.password !== formData.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match'
    isValid = false
  }

  return isValid
}

const handleSignup = async () => {
  if (!validateForm()) {
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    // Create the account
    await authApi.signup({
      username: formData.username,
      email: formData.email,
      firstName: formData.firstName,
      lastName: formData.lastName,
      password: formData.password
    })

    toast.success('Account created successfully! Redirecting to login...')
    
    // Redirect to OAuth login
    setTimeout(() => {
      authApi.login()
    }, 1000)
  } catch (error: any) {
    console.error('Signup error:', error)
    errorMessage.value = error.response?.data?.error || 'Failed to create account. Please try again.'
    toast.error(errorMessage.value)
  } finally {
    loading.value = false
  }
}
</script>
