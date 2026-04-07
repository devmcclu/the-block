<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useVehiclesStore } from "@/stores/vehicles";
import { Skeleton } from "@/components/ui/skeleton";
import VehicleCard from "./VehicleCard.vue";

const store = useVehiclesStore();
const { vehicles, loading, error } = storeToRefs(store);
</script>

<template>
  <!-- Loading State -->
  <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
    <div v-for="i in 6" :key="i" class="space-y-3">
      <Skeleton class="aspect-4/3 w-full rounded-lg" />
      <Skeleton class="h-4 w-3/4" />
      <Skeleton class="h-3 w-1/2" />
    </div>
  </div>

  <!-- Error State -->
  <div v-else-if="error" class="flex flex-col items-center justify-center py-16 text-center">
    <p class="text-lg font-medium text-destructive">{{ error }}</p>
    <p class="text-sm text-muted-foreground mt-1">Try refreshing the page</p>
  </div>

  <!-- Empty State -->
  <div
    v-else-if="vehicles.length === 0"
    class="flex flex-col items-center justify-center py-16 text-center"
  >
    <p class="text-lg font-medium text-muted-foreground">No vehicles match your filters</p>
    <p class="text-sm text-muted-foreground mt-1">Try adjusting or clearing your filters</p>
  </div>

  <!-- Vehicle Grid -->
  <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
    <VehicleCard v-for="vehicle in vehicles" :key="vehicle.external_id" :vehicle="vehicle" />
  </div>
</template>
