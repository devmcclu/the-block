import { ref, watch } from "vue";
import { defineStore } from "pinia";

export interface BidRecord {
  vehicleId: string;
  vehicleName: string;
  bidAmount: number;
  bidTime: string;
  isBuyNow: boolean;
}

const STORAGE_KEY = "the-block-bids";

function loadFromStorage(): BidRecord[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export const useBidsStore = defineStore("bids", () => {
  const bids = ref<BidRecord[]>(loadFromStorage());

  watch(bids, (val) => localStorage.setItem(STORAGE_KEY, JSON.stringify(val)), { deep: true });

  function addBid(record: BidRecord) {
    bids.value.unshift(record);
  }

  function getBidsForVehicle(vehicleId: string) {
    return bids.value.filter((b) => b.vehicleId === vehicleId);
  }

  return { bids, addBid, getBidsForVehicle };
});
