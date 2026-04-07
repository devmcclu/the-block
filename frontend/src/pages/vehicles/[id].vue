<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { api } from "@/lib/api/client";
import type { Vehicle } from "@/stores/vehicles";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { Icon } from "@iconify/vue";
import { useAuctionTime, loadAuctionConfig } from "@/composables/useAuctionTime";
import { formatCurrency, formatOdometer } from "@/lib/format";

const route = useRoute("/vehicles/[id]");
const vehicle = ref<Vehicle | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

const auctionStart = computed(() => vehicle.value?.auction_start);
const { ended, timeRemaining } = useAuctionTime(auctionStart, true);

const reserveMet = computed(() => {
  if (!vehicle.value) return false;
  return (
    vehicle.value.reserve_price == null ||
    (vehicle.value.current_bid ?? 0) >= vehicle.value.reserve_price
  );
});

const bidLabel = computed(() => {
  if (!ended.value) return "Current Bid";
  return reserveMet.value ? "Final Price" : "Final Bid";
});

onMounted(async () => {
  await loadAuctionConfig();
  const id = route.params.id;
  const { data, error: fetchError } = await api.GET("/vehicles/{id}", {
    params: { path: { id } },
  });
  if (fetchError) {
    error.value = "Unable to connect to server";
  } else if (data) {
    vehicle.value = data;
  }
  loading.value = false;
});
</script>

<template>
  <main class="container mx-auto px-4 py-6 max-w-5xl">
    <RouterLink to="/">
      <Button variant="ghost" size="sm" class="mb-4">
        <Icon icon="hugeicons:arrow-left-01" class="h-4 w-4 mr-2" />
        Back to Search
      </Button>
    </RouterLink>

    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <Skeleton class="aspect-video w-full rounded-lg" />
      <Skeleton class="h-8 w-2/3" />
      <Skeleton class="h-4 w-1/3" />
    </div>

    <template v-else-if="vehicle">
      <!-- Image Gallery -->
      <Carousel class="w-full mb-6">
        <CarouselContent>
          <CarouselItem v-for="(image, index) in vehicle.images" :key="index">
            <div class="aspect-video overflow-hidden rounded-lg bg-muted">
              <img
                :src="image.url"
                :alt="`Photo ${index + 1}`"
                :loading="index === 0 ? 'eager' : 'lazy'"
                class="h-full w-full object-cover"
              />
            </div>
          </CarouselItem>
        </CarouselContent>
        <CarouselPrevious class="absolute left-4 top-1/2 -translate-y-1/2" />
        <CarouselNext class="absolute right-4 top-1/2 -translate-y-1/2" />
      </Carousel>

      <!-- Title & Badges -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold">
          {{ vehicle.year }} {{ vehicle.make }} {{ vehicle.model }} {{ vehicle.trim }}
        </h1>
        <p class="text-sm text-muted-foreground mt-1">
          {{ vehicle.city }}, {{ vehicle.province }} &middot; Lot {{ vehicle.lot }}
        </p>
        <div class="flex flex-wrap gap-1.5 mt-3">
          <Badge variant="secondary">{{ vehicle.body_style }}</Badge>
          <Badge variant="secondary">{{ vehicle.drivetrain }}</Badge>
          <Badge variant="secondary" class="capitalize">{{ vehicle.fuel_type }}</Badge>
          <Badge variant="secondary" class="capitalize">{{ vehicle.transmission }}</Badge>
          <Badge variant="secondary" class="capitalize">{{ vehicle.title_status }}</Badge>
        </div>
      </div>

      <div class="grid md:grid-cols-2 gap-8">
        <!-- Specs -->
        <div>
          <h2 class="text-lg font-semibold mb-3">Vehicle Details</h2>
          <dl class="space-y-2 text-sm">
            <div class="flex justify-between">
              <dt class="text-muted-foreground">VIN</dt>
              <dd class="font-mono text-xs">{{ vehicle.vin }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Engine</dt>
              <dd>{{ vehicle.engine }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Transmission</dt>
              <dd class="capitalize">{{ vehicle.transmission }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Drivetrain</dt>
              <dd>{{ vehicle.drivetrain }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Odometer</dt>
              <dd>{{ formatOdometer(vehicle.odometer_km) }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Exterior Color</dt>
              <dd class="capitalize">{{ vehicle.exterior_color }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Interior Color</dt>
              <dd class="capitalize">{{ vehicle.interior_color }}</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Condition Grade</dt>
              <dd>{{ vehicle.condition_grade?.toFixed(1) }} / 5.0</dd>
            </div>
            <Separator />
            <div class="flex justify-between">
              <dt class="text-muted-foreground">Selling Dealership</dt>
              <dd>{{ vehicle.selling_dealership }}</dd>
            </div>
          </dl>
        </div>

        <!-- Auction & Condition -->
        <div class="space-y-6">
          <div>
            <h2 class="text-lg font-semibold mb-3">Auction</h2>
            <dl class="space-y-2 text-sm">
              <div class="flex justify-between">
                <dt class="text-muted-foreground">Time Remaining</dt>
                <dd class="font-semibold" :class="ended ? 'text-destructive' : ''">
                  {{ timeRemaining }}
                </dd>
              </div>
              <Separator />
              <div class="flex justify-between">
                <dt class="text-muted-foreground">{{ bidLabel }}</dt>
                <dd class="font-semibold text-base">
                  {{ formatCurrency(vehicle.current_bid) }}
                </dd>
              </div>
              <template v-if="ended && !reserveMet">
                <p class="text-xs text-destructive">
                  Reserve price not met — auction ended without a sale.
                </p>
              </template>
              <Separator />
              <div class="flex justify-between">
                <dt class="text-muted-foreground">Bids</dt>
                <dd>{{ vehicle.bid_count }}</dd>
              </div>
              <Separator />
              <div class="flex justify-between">
                <dt class="text-muted-foreground">Starting Bid</dt>
                <dd>{{ formatCurrency(vehicle.starting_bid) }}</dd>
              </div>
              <Separator />
              <div v-if="vehicle.reserve_price" class="flex justify-between">
                <dt class="text-muted-foreground">Reserve Price</dt>
                <dd>{{ formatCurrency(vehicle.reserve_price) }}</dd>
              </div>
              <Separator v-if="vehicle.reserve_price" />
              <div v-if="!ended && vehicle.buy_now_price" class="flex justify-between">
                <dt class="text-muted-foreground">Buy Now Price</dt>
                <dd class="font-semibold">
                  {{ formatCurrency(vehicle.buy_now_price) }}
                </dd>
              </div>
            </dl>
          </div>

          <!-- Condition Report -->
          <div v-if="vehicle.condition_report">
            <h2 class="text-lg font-semibold mb-3">Condition Report</h2>
            <p class="text-sm text-muted-foreground">
              {{ vehicle.condition_report }}
            </p>
          </div>

          <!-- Damage Notes -->
          <div v-if="vehicle.damage_notes && vehicle.damage_notes.length > 0">
            <h2 class="text-lg font-semibold mb-3">Damage Notes</h2>
            <ul class="list-disc list-inside text-sm text-muted-foreground space-y-1">
              <li v-for="note in vehicle.damage_notes" :key="note.note">
                {{ note.note }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </template>

    <!-- Error State -->
    <div v-else-if="error" class="py-16 text-center">
      <p class="text-lg font-medium text-destructive">{{ error }}</p>
      <p class="text-sm text-muted-foreground mt-1">Try refreshing the page</p>
    </div>

    <!-- Not Found -->
    <div v-else class="py-16 text-center">
      <p class="text-lg font-medium text-muted-foreground">Vehicle not found</p>
    </div>
  </main>
</template>
