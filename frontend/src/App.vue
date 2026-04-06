<script setup lang="ts">
import { RouterView } from 'vue-router'
import ModeToggle from '@/components/ModeToggle.vue'
</script>

<template>
  <div class="min-h-screen bg-background font-sans antialiased">
    <header class="border-b border-border">
      <div class="container mx-auto flex h-14 items-center px-4">
        <RouterLink to="/" class="text-lg font-semibold tracking-tight"> The Block </RouterLink>
        <div class="ml-auto">
          <ModeToggle />
        </div>
      </div>
    </header>

    <RouterView />
  </div>
</template>
