<script setup lang="ts">
defineProps({
  title: {
    type: String,
    default: 'PassIt - Event Ticketing Platform',
  },
  subtitle: {
    type: String,
    default: 'Create events, sell tickets, and manage attendees all in one place',
  },
});
</script>
<template>
    <section class="bg-blue-700 py-12 sm:py-16 md:py-20 mb-4">
      <div
        class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col items-center"
      >
        <div class="text-center">
          <h1
            class="text-2xl sm:text-3xl md:text-4xl lg:text-5xl xl:text-6xl font-extrabold text-white leading-tight"
          >
            {{ title }}
          </h1>
          <p class="my-3 sm:my-4 text-base sm:text-lg md:text-xl text-white px-4">
            {{ subtitle }}
          </p>
        </div>
      </div>
    </section>
</template>