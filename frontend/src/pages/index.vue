<script setup lang="ts">
import { watch, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useVehiclesStore, type VehicleFilters } from "@/stores/vehicles";
import { useDebounceFn } from "@vueuse/core";
import SearchFilters from "@/components/search/SearchFilters.vue";
import VehicleGrid from "@/components/search/VehicleGrid.vue";
import MobileFilterSheet from "@/components/search/MobileFilterSheet.vue";
import { loadAuctionConfig } from "@/composables/useAuctionTime";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const route = useRoute();
const router = useRouter();
const store = useVehiclesStore();
const { vehicles, sort, loading } = storeToRefs(store);

const arrayKeys: { filter: keyof VehicleFilters; param: string }[] = [
  { filter: "makes", param: "make" },
  { filter: "models", param: "model" },
  { filter: "bodyStyles", param: "body_style" },
  { filter: "exteriorColors", param: "exterior_color" },
  { filter: "interiorColors", param: "interior_color" },
  { filter: "transmissions", param: "transmission" },
  { filter: "drivetrains", param: "drivetrain" },
  { filter: "fuelTypes", param: "fuel_type" },
  { filter: "titleStatuses", param: "title_status" },
];

const numericKeys: { filter: keyof VehicleFilters; param: string }[] = [
  { filter: "yearMin", param: "year_min" },
  { filter: "yearMax", param: "year_max" },
  { filter: "odometerMin", param: "odometer_min" },
  { filter: "odometerMax", param: "odometer_max" },
  { filter: "conditionMin", param: "condition_min" },
  { filter: "conditionMax", param: "condition_max" },
];

function loadFromQuery() {
  const q = route.query;
  for (const { filter, param } of arrayKeys) {
    const val = q[param];
    if (typeof val === "string" && val) {
      (store.filters[filter] as string[]).splice(0, Infinity, ...val.split(","));
    }
  }
  for (const { filter, param } of numericKeys) {
    const val = q[param];
    if (typeof val === "string" && val) {
      (store.filters as unknown as Record<string, number | undefined>)[filter] = Number(val);
    }
  }
  if (typeof q.sort === "string" && q.sort) {
    sort.value = q.sort;
  }
}

function syncToQuery() {
  const query: Record<string, string> = {};
  for (const { filter, param } of arrayKeys) {
    const arr = store.filters[filter] as string[];
    if (arr.length) query[param] = arr.join(",");
  }
  for (const { filter, param } of numericKeys) {
    const val = (store.filters as unknown as Record<string, number | undefined>)[filter];
    if (val != null) query[param] = String(val);
  }
  if (sort.value) query.sort = sort.value;
  router.replace({ query });
}

const debouncedFetch = useDebounceFn(() => {
  store.fetchVehicles();
  syncToQuery();
}, 300);

loadFromQuery();

onMounted(async () => {
  await Promise.allSettled([
    loadAuctionConfig(),
    store.fetchFilterOptions(),
    store.fetchVehicles(),
  ]);
});

watch(
  () => store.filters,
  () => debouncedFetch(),
  { deep: true },
);
watch(sort, () => {
  store.fetchVehicles();
  syncToQuery();
});
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
          <div class="flex items-center gap-2">
            <Select v-model="sort">
              <SelectTrigger class="w-44">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="price_asc">Price: Low to High</SelectItem>
                <SelectItem value="price_desc">Price: High to Low</SelectItem>
                <SelectItem value="year_desc">Newest</SelectItem>
                <SelectItem value="year_asc">Oldest</SelectItem>
                <SelectItem value="bids_desc">Most Bids</SelectItem>
                <SelectItem value="bids_asc">Fewest Bids</SelectItem>
                <SelectItem value="ending_soon">Ending Soon</SelectItem>
                <SelectItem value="ending_last">Ending Last</SelectItem>
              </SelectContent>
            </Select>
            <MobileFilterSheet />
          </div>
        </div>

        <VehicleGrid />
      </div>
    </div>
  </main>
</template>
