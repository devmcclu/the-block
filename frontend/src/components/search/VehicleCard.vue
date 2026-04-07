<script setup lang="ts">
import { useRouter } from "vue-router";
import type { Vehicle } from "@/stores/vehicles";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { computed } from "vue";
import { useAuctionTime } from "@/composables/useAuctionTime";
import { formatCurrency, formatOdometer } from "@/lib/format";

const router = useRouter();

const props = defineProps<{
  vehicle: Vehicle;
}>();

const displayName = computed(
  () => `${props.vehicle.year} ${props.vehicle.make} ${props.vehicle.model} ${props.vehicle.trim}`,
);
const { ended, timeRemaining, urgency } = useAuctionTime(props.vehicle.auction_start);

const timeClass = computed(() => {
  switch (urgency.value) {
    case "ended":
      return "text-destructive";
    case "urgent":
      return "text-destructive font-semibold";
    case "warning":
      return "text-warning font-semibold";
    default:
      return "text-muted-foreground";
  }
});


function navigateToVehicle() {
  router.push(`/vehicles/${props.vehicle.external_id}`);
}
</script>

<template>
  <Card
    class="overflow-hidden transition-shadow hover:shadow-md h-full cursor-pointer"
    @click="navigateToVehicle"
  >
    <!-- Image Carousel -->
    <div class="relative">
      <Carousel class="w-full">
        <CarouselContent>
          <CarouselItem v-for="(image, index) in vehicle.images" :key="index">
            <div class="aspect-4/3 overflow-hidden bg-muted">
              <img
                :src="image.url"
                :alt="`${displayName} - Photo ${index + 1}`"
                :loading="index === 0 ? 'eager' : 'lazy'"
                class="h-full w-full object-cover"
              />
            </div>
          </CarouselItem>
        </CarouselContent>
        <CarouselPrevious
          v-if="vehicle.images && vehicle.images.length > 1"
          class="absolute left-2 top-1/2 -translate-y-1/2"
          @click.stop
        />
        <CarouselNext
          v-if="vehicle.images && vehicle.images.length > 1"
          class="absolute right-2 top-1/2 -translate-y-1/2"
          @click.stop
        />
      </Carousel>
    </div>

    <CardContent class="p-4 space-y-2">
      <h3 class="font-semibold text-sm leading-tight line-clamp-1">
        {{ displayName }}
      </h3>

      <div class="flex flex-wrap gap-1">
        <Badge variant="secondary" class="text-xs">
          {{ vehicle.body_style }}
        </Badge>
        <Badge variant="secondary" class="text-xs">
          {{ vehicle.drivetrain }}
        </Badge>
        <Badge variant="secondary" class="text-xs capitalize">
          {{ vehicle.fuel_type }}
        </Badge>
      </div>

      <div class="text-xs text-muted-foreground space-y-0.5">
        <p>{{ formatOdometer(vehicle.odometer_km) }}</p>
        <p>{{ vehicle.city }}, {{ vehicle.province }}</p>
      </div>
    </CardContent>

    <CardFooter class="px-4 pb-4 pt-0 flex items-center justify-between">
      <div>
        <p class="text-sm font-semibold">
          {{ formatCurrency(vehicle.current_bid) }}
        </p>
        <p class="text-xs text-muted-foreground">
          {{ vehicle.bid_count }} {{ vehicle.bid_count === 1 ? "bid" : "bids" }}
        </p>
      </div>
      <div class="text-right">
        <p v-if="!ended && vehicle.buy_now_price" class="text-xs text-muted-foreground">
          Buy Now {{ formatCurrency(vehicle.buy_now_price) }}
        </p>
        <p class="text-xs" :class="timeClass">
          {{ timeRemaining }}
        </p>
      </div>
    </CardFooter>
  </Card>
</template>
