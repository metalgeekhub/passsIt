<script setup lang="ts">
import { RouterLink } from 'vue-router';
import Card from '@/components/Card.vue';
import { isAuthenticated, isAdmin } from '@/composables/useAuth';
</script>

<template>
    <section class="py-4">
      <div class="container-xl lg:container m-auto">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 rounded-lg">
          <Card v-if="isAdmin">
            <h2 class="text-2xl font-bold">For Event Organizers</h2>
            <p class="mt-2 mb-4">
              Create and manage events, sell tickets, and track attendees
            </p>
            <RouterLink
              to="/users"
              class="inline-block bg-black text-white rounded-lg px-4 py-2 hover:bg-gray-700"
            >
              Manage Users
            </RouterLink>
          </Card>
          <Card v-else>
            <h2 class="text-2xl font-bold">Discover Events</h2>
            <p class="mt-2 mb-4">
              Browse amazing events and purchase tickets instantly
            </p>
            <RouterLink
              v-if="isAuthenticated"
              to="/profile"
              class="inline-block bg-black text-white rounded-lg px-4 py-2 hover:bg-gray-700"
            >
              My Profile
            </RouterLink>
            <RouterLink
              v-else
              to="/login"
              class="inline-block bg-black text-white rounded-lg px-4 py-2 hover:bg-gray-700"
            >
              Sign In
            </RouterLink>
          </Card>
          <Card bg="bg-blue-100">
            <h2 class="text-2xl font-bold">{{ isAuthenticated ? 'Get Started' : 'Join PassIt' }}</h2>
            <p class="mt-2 mb-4">
              {{ isAuthenticated ? 'Explore events and manage your tickets' : 'Create an account and start attending events' }}
            </p>
            <RouterLink
              v-if="!isAuthenticated"
              to="/signup"
              class="inline-block bg-blue-500 text-white rounded-lg px-4 py-2 hover:bg-blue-600"
            >
              Create Account
            </RouterLink>
            <RouterLink
              v-else
              to="/profile"
              class="inline-block bg-blue-500 text-white rounded-lg px-4 py-2 hover:bg-blue-600"
            >
              View Profile
            </RouterLink>
          </Card>
        </div>
      </div>
    </section>
</template>