<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { api } from "@/lib/api/client";
import type { Vehicle } from "@/stores/vehicles";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
import {
  useAuctionTime,
  loadAuctionConfig,
  getMinBidIncrement,
} from "@/composables/useAuctionTime";
import { formatCurrency, formatOdometer } from "@/lib/format";
import { toast } from "vue-sonner";

const route = useRoute("/vehicles/[id]");
const vehicle = ref<Vehicle | null>(null);
const loading = ref(true);
const notFound = ref(false);
const error = ref<string | null>(null);

const auctionStart = computed(() => vehicle.value?.auction_start);
const { ended, timeRemaining, urgency } = useAuctionTime(auctionStart, true);

const timeClass = computed(() => {
  switch (urgency.value) {
    case "ended":
      return "text-destructive";
    case "urgent":
      return "text-destructive";
    case "warning":
      return "text-warning";
    default:
      return "";
  }
});

const minBidIncrement = getMinBidIncrement();

const bidOpen = ref(false);
const bidAmount = ref<number | undefined>(undefined);
const bidError = ref<string | null>(null);
const bidSuccess = ref(false);
const bidLoading = ref(false);
const buyNowOpen = ref(false);
const buyNowLoading = ref(false);
const buyNowSuccess = ref(false);

function openBuyNowDialog() {
  buyNowSuccess.value = false;
  buyNowOpen.value = true;
}

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
  return (vehicle.value.current_bid ?? 0) + (minBidIncrement.value ?? 0);
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

  buyNowLoading.value = true;
  try {
    const id = vehicle.value.external_id!;
    const { data, error: fetchError } = await api.POST("/vehicles/{id}/buy", {
      params: { path: { id } },
    });
    if (fetchError) {
      toast.error("Purchase failed", {
        description: (fetchError as { detail?: string }).detail ?? "Unable to complete purchase",
      });
    } else if (data) {
      vehicle.value = data;
      buyNowSuccess.value = true;
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
  <main class="container mx-auto px-4 py-6 max-w-6xl">
    <RouterLink to="/">
      <Button variant="ghost" size="sm" class="mb-4">
        <Icon icon="hugeicons:arrow-left-01" class="h-4 w-4 mr-2" />
        Back to Search
      </Button>
    </RouterLink>

    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <Skeleton class="aspect-video w-full rounded-xl" />
      <div class="grid lg:grid-cols-5 gap-6">
        <div class="lg:col-span-3 space-y-4">
          <Skeleton class="h-8 w-2/3" />
          <Skeleton class="h-4 w-1/3" />
          <Skeleton class="h-48 w-full rounded-xl" />
        </div>
        <div class="lg:col-span-2">
          <Skeleton class="h-72 w-full rounded-xl" />
        </div>
      </div>
    </div>

    <template v-else-if="vehicle">
      <!-- Title -->
      <div class="mb-4">
        <h1 class="text-2xl font-bold tracking-tight">
          {{ vehicle.year }} {{ vehicle.make }} {{ vehicle.model }} {{ vehicle.trim }}
        </h1>
        <p class="text-sm text-muted-foreground mt-1">
          {{ vehicle.city }}, {{ vehicle.province }} &middot; Lot {{ vehicle.lot }}
        </p>
      </div>

      <!-- Image Gallery -->
      <Carousel class="w-full mb-8">
        <CarouselContent>
          <CarouselItem v-for="(image, index) in vehicle.images" :key="index">
            <div class="aspect-video overflow-hidden rounded-xl bg-muted">
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

      <div class="grid lg:grid-cols-5 gap-6 items-start">
        <!-- Left Column: Details + Condition -->
        <div class="lg:col-span-3 lg:col-start-1 lg:row-start-1 space-y-6 order-2 lg:order-none">
          <!-- Vehicle Details Card -->
          <Card>
            <CardHeader>
              <CardTitle class="text-sm uppercase tracking-wide text-muted-foreground">
                Vehicle Details
              </CardTitle>
            </CardHeader>
            <CardContent>
              <dl class="grid grid-cols-2 gap-x-8 gap-y-3 text-sm">
                <div>
                  <dt class="text-muted-foreground text-xs">VIN</dt>
                  <dd class="font-mono text-xs mt-0.5">{{ vehicle.vin }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Engine</dt>
                  <dd class="mt-0.5">{{ vehicle.engine }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Transmission</dt>
                  <dd class="capitalize mt-0.5">{{ vehicle.transmission }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Drivetrain</dt>
                  <dd class="mt-0.5">{{ vehicle.drivetrain }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Odometer</dt>
                  <dd class="mt-0.5">{{ formatOdometer(vehicle.odometer_km) }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Condition Grade</dt>
                  <dd class="mt-0.5">{{ vehicle.condition_grade?.toFixed(1) }} / 5.0</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Body Style</dt>
                  <dd class="capitalize mt-0.5">{{ vehicle.body_style }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Fuel Type</dt>
                  <dd class="capitalize mt-0.5">{{ vehicle.fuel_type }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Exterior Color</dt>
                  <dd class="capitalize mt-0.5">{{ vehicle.exterior_color }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Interior Color</dt>
                  <dd class="capitalize mt-0.5">{{ vehicle.interior_color }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Title Status</dt>
                  <dd class="capitalize mt-0.5">{{ vehicle.title_status }}</dd>
                </div>
                <div>
                  <dt class="text-muted-foreground text-xs">Selling Dealership</dt>
                  <dd class="mt-0.5">{{ vehicle.selling_dealership }}</dd>
                </div>
              </dl>
            </CardContent>
          </Card>

          <!-- Condition Report -->
          <Card v-if="vehicle.condition_report">
            <CardHeader>
              <CardTitle class="text-sm uppercase tracking-wide text-muted-foreground">
                Condition Report
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p class="text-sm leading-relaxed">
                {{ vehicle.condition_report }}
              </p>
            </CardContent>
          </Card>

          <!-- Damage Notes -->
          <Card
            v-if="vehicle.damage_notes && vehicle.damage_notes.length > 0"
            class="border-destructive/20 bg-destructive/5"
          >
            <CardHeader>
              <CardTitle
                class="text-sm uppercase tracking-wide text-destructive/80 flex items-center gap-2"
              >
                <Icon icon="hugeicons:alert-02" class="h-4 w-4" />
                Damage Notes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ul class="space-y-2 text-sm">
                <li
                  v-for="note in vehicle.damage_notes"
                  :key="note.note"
                  class="flex items-start gap-2"
                >
                  <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-destructive/40" />
                  <span>{{ note.note }}</span>
                </li>
              </ul>
            </CardContent>
          </Card>
        </div>

        <!-- Right Column: Auction Panel (sticky) -->
        <div class="order-1 lg:order-none lg:col-span-2 lg:col-start-4 lg:row-start-1 lg:sticky lg:top-6">
          <Card class="gap-0 py-0">
            <!-- Auction Timer -->
            <div class="flex items-center justify-between border-b px-5 py-3">
              <div class="flex items-center gap-2">
                <Icon icon="hugeicons:time-04" class="h-4 w-4 text-muted-foreground" />
                <span class="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                  {{ ended ? "Auction Ended" : "Time Remaining" }}
                </span>
              </div>
              <span class="font-semibold tabular-nums" :class="timeClass">
                {{ timeRemaining }}
              </span>
            </div>

            <!-- Current Bid Hero -->
            <div class="px-5 py-4 bg-muted/30">
              <p class="text-xs text-muted-foreground">{{ bidLabel }}</p>
              <p class="text-3xl font-bold tracking-tight mt-0.5">
                {{ formatCurrency(vehicle.current_bid) }}
              </p>
              <p class="text-xs text-muted-foreground mt-1">
                {{ vehicle.bid_count }} {{ vehicle.bid_count === 1 ? "bid" : "bids" }}
              </p>
              <p v-if="ended && !reserveMet" class="text-xs text-destructive mt-1.5">
                Reserve price not met — auction ended without a sale.
              </p>
            </div>

            <!-- Price Details -->
            <div class="px-5 py-4 space-y-4">
              <dl class="space-y-2.5 text-sm">
                <div class="flex justify-between">
                  <dt class="text-muted-foreground">Starting Bid</dt>
                  <dd>{{ formatCurrency(vehicle.starting_bid) }}</dd>
                </div>
                <div v-if="vehicle.reserve_price != null" class="flex justify-between">
                  <dt class="text-muted-foreground">Reserve Price</dt>
                  <dd>{{ formatCurrency(vehicle.reserve_price) }}</dd>
                </div>
                <div v-if="!ended && vehicle.buy_now_price != null" class="flex justify-between">
                  <dt class="text-muted-foreground">Buy Now Price</dt>
                  <dd class="font-semibold">{{ formatCurrency(vehicle.buy_now_price) }}</dd>
                </div>
              </dl>

              <!-- Bid Actions -->
              <div v-if="!ended" class="space-y-2.5">
                <Separator />
                <Button class="w-full" size="lg" @click="openBidDialog">
                  <Icon icon="hugeicons:auction" class="h-4 w-4 mr-2" />
                  Place Bid
                </Button>

                <Button
                  v-if="vehicle.buy_now_price != null"
                  class="w-full"
                  size="lg"
                  variant="outline"
                  @click="openBuyNowDialog"
                >
                  <Icon icon="hugeicons:shopping-bag-02" class="h-4 w-4 mr-2" />
                  Buy Now - {{ formatCurrency(vehicle.buy_now_price) }}
                </Button>
              </div>
            </div>
          </Card>
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
          <DialogDescription> Enter your bid for this vehicle. </DialogDescription>
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
                Current bid: {{ formatCurrency(vehicle.current_bid) }} &middot;
                {{ vehicle.bid_count }} {{ vehicle.bid_count === 1 ? "bid" : "bids" }}
              </p>
            </div>
          </div>

          <Separator />

          <!-- Success State -->
          <div v-if="bidSuccess" class="text-center space-y-3 py-2">
            <Icon icon="hugeicons:checkmark-circle-02" class="mx-auto h-10 w-10 text-success" />
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
          <DialogDescription> Complete your purchase of this vehicle. </DialogDescription>
        </DialogHeader>

        <div v-if="vehicle" class="space-y-4 py-2">
          <!-- Success State -->
          <div v-if="buyNowSuccess" class="text-center space-y-3 py-4">
            <Icon icon="hugeicons:checkmark-circle-02" class="mx-auto h-12 w-12 text-green-500" />
            <div>
              <p class="text-lg font-semibold">Purchase Complete</p>
              <p class="text-sm text-muted-foreground mt-1">
                You purchased this vehicle for
                {{ formatCurrency(vehicle.buy_now_price) }}.
              </p>
            </div>
            <p class="text-sm text-muted-foreground">{{ vehicleName }}</p>
          </div>

          <!-- Checkout Form -->
          <template v-else>
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
          </template>
        </div>

        <DialogFooter class="gap-2 sm:gap-0">
          <template v-if="buyNowSuccess">
            <Button @click="buyNowOpen = false">Done</Button>
          </template>
          <template v-else>
            <Button variant="outline" @click="buyNowOpen = false">Cancel</Button>
            <Button :disabled="buyNowLoading" @click="confirmBuyNow">
              <Icon
                v-if="buyNowLoading"
                icon="hugeicons:loading-03"
                class="h-4 w-4 mr-2 animate-spin"
              />
              Confirm Purchase
            </Button>
          </template>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>
