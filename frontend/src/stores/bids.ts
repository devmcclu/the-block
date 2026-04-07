import { ref } from "vue";
import { defineStore } from "pinia";
import { api } from "@/lib/api/client";
import type { components } from "@/lib/api/v1";

export type Bid = components["schemas"]["Bid"];

export const useBidsStore = defineStore("bids", () => {
  const bids = ref<Bid[]>([]);
  const loading = ref(false);

  async function fetchBids() {
    loading.value = true;
    try {
      const { data } = await api.GET("/bids/");
      if (data) {
        bids.value = data;
      }
    } finally {
      loading.value = false;
    }
  }

  return { bids, loading, fetchBids };
});
