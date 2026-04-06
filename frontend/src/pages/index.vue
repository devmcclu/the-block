<script setup lang="ts">
import { watch, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useVehiclesStore } from '@/stores/vehicles'
import { useDebounceFn } from '@vueuse/core'
import SearchFilters from '@/components/search/SearchFilters.vue'
import VehicleGrid from '@/components/search/VehicleGrid.vue'
import MobileFilterSheet from '@/components/search/MobileFilterSheet.vue'

const store = useVehiclesStore()
const { vehicles, loading } = storeToRefs(store)

const debouncedFetch = useDebounceFn(() => store.fetchVehicles(), 300)

onMounted(async () => {
  await Promise.all([store.fetchFilterOptions(), store.fetchVehicles()])
})

watch(() => store.filters, () => debouncedFetch(), { deep: true })
</script>

<template>
  <main class="container mx-auto px-4 py-6">
    <div class="flex gap-6">
      <!-- Desktop Sidebar -->
      <aside class="hidden lg:block w-72 shrink-0">
        <div class="sticky top-6 h-[calc(100vh-7rem)] border rounded-lg">
          <SearchFilters />
        </div>
      </aside>

      <!-- Main Content -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between mb-4">
          <p class="text-sm text-muted-foreground">
            <span v-if="!loading">{{ vehicles.length }} vehicles</span>
            <span v-else>Loading...</span>
          </p>
          <MobileFilterSheet />
        </div>

        <VehicleGrid />
      </div>
    </div>
  </main>
</template>
