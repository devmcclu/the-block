<script setup lang="ts">
import { RouterView } from "vue-router";
import ModeToggle from "@/components/ModeToggle.vue";
import "vue-sonner/style.css";
import { Toaster } from "@/components/ui/sonner";
</script>

<template>
  <div class="min-h-screen bg-background font-sans antialiased">
    <header class="border-b border-border">
      <div class="container mx-auto flex h-14 items-center px-4">
        <RouterLink to="/" class="text-lg font-semibold tracking-tight"> The Block </RouterLink>
        <nav class="ml-6 flex items-center gap-4">
          <RouterLink
            to="/bids"
            class="text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            My Bids
          </RouterLink>
        </nav>
        <div class="ml-auto">
          <ModeToggle />
        </div>
      </div>
    </header>

    <RouterView />
    <Toaster />
  </div>
</template>
