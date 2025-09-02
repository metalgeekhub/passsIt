<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { defineProps, ref, computed } from 'vue';
import { type User } from '@/models/User'

const props = defineProps({
    user: Object as () => User
});

</script>

<template>
    <div class="bg-white rounded-xl shadow-md relative">
            <div class="p-4">
              <div class="mb-6">
                <div class="text-gray-600 my-2">{{ user?.username }}</div>
                <h3 class="text-xl font-bold">{{ user?.email }}</h3>
              </div>

              <div class="mb-5">
                <div>
                  {{ user?.address }}
                </div>
              </div>

              <h3 class="text-blue-500 mb-2">{{ user?.dob }}</h3>

              <div class="border border-gray-100 mb-5"></div>

              <div class="flex flex-col lg:flex-row justify-between mb-4">
                <div class="text-orange-700 mb-3">
                  <i class="pi pi-map-marker text-orange-700"></i>
                  {{ user?.phone_number }}
                </div>
                <RouterLink
                  :to="'/users/' + user?.id"
                  class="h-[36px] bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-center text-sm"
                >
                  Read More
                </RouterLink>
              </div>
            </div>
          </div>
</template>