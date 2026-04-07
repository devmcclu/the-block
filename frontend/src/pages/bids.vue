<script setup lang="ts">
import { onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useBidsStore } from "@/stores/bids";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Icon } from "@iconify/vue";
import { formatCurrency } from "@/lib/format";

const bidsStore = useBidsStore();
const { bids, loading } = storeToRefs(bidsStore);

function formatRelativeTime(iso: string | undefined) {
  if (!iso) return "";
  const diff = Date.now() - new Date(iso).getTime();
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return "Just now";
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

onMounted(() => bidsStore.fetchBids());
</script>

<template>
  <main class="container mx-auto px-4 py-6 max-w-3xl">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">My Bids</h1>
      <RouterLink to="/">
        <Button variant="ghost" size="sm">
          <Icon icon="hugeicons:arrow-left-01" class="h-4 w-4 mr-2" />
          Browse Vehicles
        </Button>
      </RouterLink>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <Skeleton v-for="i in 4" :key="i" class="h-16 w-full rounded-md" />
    </div>

    <div v-else-if="bids.length === 0" class="py-16 text-center">
      <Icon icon="hugeicons:invoice-03" class="mx-auto h-12 w-12 text-muted-foreground/50" />
      <p class="mt-4 text-lg font-medium text-muted-foreground">No bids yet</p>
      <p class="text-sm text-muted-foreground mt-1">
        Browse vehicles and place your first bid to see it here.
      </p>
      <RouterLink to="/">
        <Button class="mt-4">Browse Vehicles</Button>
      </RouterLink>
    </div>

    <div v-else class="space-y-1">
      <div
        v-for="(bid, index) in bids"
        :key="`${bid.vehicle_external_id}-${bid.bid_time}`"
      >
        <RouterLink
          :to="`/vehicles/${bid.vehicle_external_id}`"
          class="flex items-center justify-between py-3 px-2 -mx-2 rounded-md hover:bg-muted/50 transition-colors"
        >
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <p class="font-medium truncate">{{ bid.vehicle_name }}</p>
              <Badge v-if="bid.is_buy_now" variant="secondary" class="shrink-0">Buy Now</Badge>
            </div>
            <p class="text-sm text-muted-foreground mt-0.5">
              {{ formatRelativeTime(bid.bid_time) }}
            </p>
          </div>
          <div class="text-right shrink-0 ml-4">
            <p class="font-semibold">{{ formatCurrency(bid.bid_amount) }}</p>
          </div>
        </RouterLink>
        <Separator v-if="index < bids.length - 1" />
      </div>
    </div>
  </main>
</template>
