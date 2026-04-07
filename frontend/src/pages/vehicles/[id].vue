<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { api } from "@/lib/api/client";
import type { Vehicle } from "@/stores/vehicles";
import { useBidsStore } from "@/stores/bids";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Icon } from "@iconify/vue";
import { useAuctionTime, loadAuctionConfig, getMinBidIncrement } from "@/composables/useAuctionTime";
import { formatCurrency, formatOdometer } from "@/lib/format";
import { toast } from "vue-sonner";

const route = useRoute("/vehicles/[id]");
const vehicle = ref<Vehicle | null>(null);
const loading = ref(true);
const notFound = ref(false);
const error = ref<string | null>(null);

const auctionStart = computed(() => vehicle.value?.auction_start);
const { ended, timeRemaining } = useAuctionTime(auctionStart, true);
const minBidIncrement = getMinBidIncrement();
const bidsStore = useBidsStore();

const bidOpen = ref(false);
const bidAmount = ref<number | undefined>(undefined);
const bidError = ref<string | null>(null);
const bidSuccess = ref(false);
const bidLoading = ref(false);
const buyNowOpen = ref(false);
const buyNowLoading = ref(false);

function openBidDialog() {
  bidAmount.value = undefined;
  bidError.value = null;
  bidSuccess.value = false;
  bidOpen.value = true;
}

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

const minimumBid = computed(() => {
  if (!vehicle.value) return 0;
  if (!vehicle.value.bid_count || vehicle.value.bid_count === 0) {
    return vehicle.value.starting_bid ?? 0;
  }
  return (vehicle.value.current_bid ?? 0) + (minBidIncrement.value ?? 100);
});

const vehicleName = computed(() => {
  if (!vehicle.value) return "";
  return `${vehicle.value.year} ${vehicle.value.make} ${vehicle.value.model} ${vehicle.value.trim}`.trim();
});

async function placeBid() {
  if (!vehicle.value || !bidAmount.value) return;
  if (bidAmount.value < minimumBid.value) {
    bidError.value = `Bid must be at least ${formatCurrency(minimumBid.value)}`;
    return;
  }

  bidError.value = null;
  bidLoading.value = true;
  try {
    const id = vehicle.value.external_id!;
    const { data, error: fetchError } = await api.PUT("/vehicles/{id}", {
      params: { path: { id } },
      body: { bid_amount: bidAmount.value },
    });
    if (fetchError) {
      bidError.value = (fetchError as { detail?: string }).detail ?? "Failed to place bid";
    } else if (data) {
      vehicle.value = data;
      bidsStore.addBid({
        vehicleId: id,
        vehicleName: vehicleName.value,
        bidAmount: bidAmount.value,
        bidTime: new Date().toISOString(),
        isBuyNow: false,
      });
      bidSuccess.value = true;
      bidError.value = null;
    }
  } catch {
    bidError.value = "Unable to connect to server";
  } finally {
    bidLoading.value = false;
  }
}

async function confirmBuyNow() {
  if (!vehicle.value || vehicle.value.buy_now_price == null) return;
  const price = vehicle.value.buy_now_price;

  buyNowLoading.value = true;
  try {
    const id = vehicle.value.external_id!;
    const { data, error: fetchError } = await api.PUT("/vehicles/{id}", {
      params: { path: { id } },
      body: { bid_amount: price },
    });
    if (fetchError) {
      toast.error("Purchase failed", {
        description: (fetchError as { detail?: string }).detail ?? "Unable to complete purchase",
      });
    } else if (data) {
      vehicle.value = data;
      bidsStore.addBid({
        vehicleId: id,
        vehicleName: vehicleName.value,
        bidAmount: price,
        bidTime: new Date().toISOString(),
        isBuyNow: true,
      });
      toast.success("Purchase complete", {
        description: `You purchased this vehicle for ${formatCurrency(price)}.`,
      });
      buyNowOpen.value = false;
    }
  } catch {
    toast.error("Purchase failed", { description: "Unable to connect to server" });
  } finally {
    buyNowLoading.value = false;
  }
}

onMounted(async () => {
  try {
    await loadAuctionConfig();
    const id = route.params.id;
    if (!id) {
      error.value = "Vehicle ID is missing";
      return;
    }
    const {
      data,
      error: fetchError,
      response,
    } = await api.GET("/vehicles/{id}", {
      params: { path: { id } },
    });
    if (fetchError) {
      if (response?.status === 404) {
        notFound.value = true;
      } else {
        error.value = "Unable to connect to server";
      }
    } else if (data) {
      vehicle.value = data;
    }
  } catch {
    error.value = "Unable to connect to server";
  } finally {
    loading.value = false;
  }
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
              <div v-if="vehicle.reserve_price != null" class="flex justify-between">
                <dt class="text-muted-foreground">Reserve Price</dt>
                <dd>{{ formatCurrency(vehicle.reserve_price) }}</dd>
              </div>
              <Separator v-if="vehicle.reserve_price != null" />
              <div v-if="!ended && vehicle.buy_now_price != null" class="flex justify-between">
                <dt class="text-muted-foreground">Buy Now Price</dt>
                <dd class="font-semibold">
                  {{ formatCurrency(vehicle.buy_now_price) }}
                </dd>
              </div>
            </dl>

            <!-- Bid Actions -->
            <div v-if="!ended" class="mt-4 space-y-3">
              <Separator />
              <Button class="w-full" @click="openBidDialog">
                <Icon icon="hugeicons:auction" class="h-4 w-4 mr-2" />
                Place Bid
              </Button>

              <Button
                v-if="vehicle.buy_now_price != null"
                class="w-full"
                variant="secondary"
                @click="buyNowOpen = true"
              >
                <Icon icon="hugeicons:shopping-bag-02" class="h-4 w-4 mr-2" />
                Buy Now &mdash; {{ formatCurrency(vehicle.buy_now_price) }}
              </Button>
            </div>
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

    <!-- Not Found -->
    <div v-else-if="notFound" class="py-16 text-center">
      <p class="text-lg font-medium text-muted-foreground">Vehicle not found</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="py-16 text-center">
      <p class="text-lg font-medium text-destructive">{{ error }}</p>
      <p class="text-sm text-muted-foreground mt-1">Try refreshing the page</p>
    </div>

    <!-- Place Bid Dialog -->
    <Dialog v-model:open="bidOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Place a Bid</DialogTitle>
          <DialogDescription>
            Enter your bid for this vehicle.
          </DialogDescription>
        </DialogHeader>

        <div v-if="vehicle" class="space-y-4 py-2">
          <div class="flex items-center gap-4">
            <div
              v-if="vehicle.images && vehicle.images.length > 0"
              class="h-16 w-24 shrink-0 overflow-hidden rounded bg-muted"
            >
              <img
                :src="vehicle.images?.[0]?.url"
                :alt="vehicleName"
                class="h-full w-full object-cover"
              />
            </div>
            <div class="min-w-0">
              <p class="font-semibold truncate">{{ vehicleName }}</p>
              <p class="text-sm text-muted-foreground">
                Current bid: {{ formatCurrency(vehicle.current_bid) }}
                &middot; {{ vehicle.bid_count }} {{ vehicle.bid_count === 1 ? "bid" : "bids" }}
              </p>
            </div>
          </div>

          <Separator />

          <!-- Success State -->
          <div v-if="bidSuccess" class="text-center space-y-3 py-2">
            <Icon icon="hugeicons:checkmark-circle-02" class="mx-auto h-10 w-10 text-green-500" />
            <div>
              <p class="font-semibold">Bid placed successfully</p>
              <p class="text-sm text-muted-foreground">
                Your bid of {{ formatCurrency(bidAmount) }} has been recorded.
              </p>
            </div>
          </div>

          <!-- Bid Input State -->
          <div v-else class="space-y-3">
            <div class="space-y-2">
              <label class="text-sm font-medium">Your Bid</label>
              <Input
                v-model.number="bidAmount"
                type="number"
                :min="minimumBid"
                :placeholder="`$${minimumBid.toLocaleString()}`"
                @keyup.enter="placeBid"
              />
              <p class="text-xs text-muted-foreground">
                Minimum bid: {{ formatCurrency(minimumBid) }}
              </p>
              <p v-if="bidError" class="text-xs text-destructive">{{ bidError }}</p>
            </div>
          </div>
        </div>

        <DialogFooter class="gap-2 sm:gap-0">
          <template v-if="bidSuccess">
            <Button @click="bidOpen = false">Done</Button>
          </template>
          <template v-else>
            <Button variant="outline" @click="bidOpen = false">Cancel</Button>
            <Button :disabled="bidLoading || !bidAmount" @click="placeBid">
              <Icon
                v-if="bidLoading"
                icon="hugeicons:loading-03"
                class="h-4 w-4 mr-2 animate-spin"
              />
              Place Bid
            </Button>
          </template>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Buy Now Dialog -->
    <Dialog v-model:open="buyNowOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Buy Now</DialogTitle>
          <DialogDescription>
            Complete your purchase of this vehicle.
          </DialogDescription>
        </DialogHeader>

        <div v-if="vehicle" class="space-y-4 py-2">
          <div class="flex items-center gap-4">
            <div
              v-if="vehicle.images && vehicle.images.length > 0"
              class="h-16 w-24 shrink-0 overflow-hidden rounded bg-muted"
            >
              <img
                :src="vehicle.images?.[0]?.url"
                :alt="vehicleName"
                class="h-full w-full object-cover"
              />
            </div>
            <div class="min-w-0">
              <p class="font-semibold truncate">{{ vehicleName }}</p>
              <p class="text-sm text-muted-foreground">Lot {{ vehicle.lot }}</p>
            </div>
          </div>

          <Separator />

          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-muted-foreground">Purchase Price</span>
              <span class="text-lg font-bold">
                {{ formatCurrency(vehicle.buy_now_price) }}
              </span>
            </div>
          </div>

          <Separator />

          <div class="space-y-1">
            <p class="text-sm font-medium">Payment Method</p>
            <div class="flex items-center gap-2 rounded-md border p-3">
              <Icon icon="hugeicons:credit-card" class="h-5 w-5 text-muted-foreground" />
              <span class="text-sm">Visa ending in 4242</span>
            </div>
          </div>
        </div>

        <DialogFooter class="gap-2 sm:gap-0">
          <Button variant="outline" @click="buyNowOpen = false">Cancel</Button>
          <Button :disabled="buyNowLoading" @click="confirmBuyNow">
            <Icon
              v-if="buyNowLoading"
              icon="hugeicons:loading-03"
              class="h-4 w-4 mr-2 animate-spin"
            />
            Confirm Purchase
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>
