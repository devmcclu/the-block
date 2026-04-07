import { describe, it, expect, vi, beforeEach } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useBidsStore } from "@/stores/bids";

vi.mock("@/lib/api/client", () => ({
  api: {
    GET: vi.fn(),
  },
}));

import { api } from "@/lib/api/client";

const mockGet = vi.mocked(api.GET);

describe("useBidsStore", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("starts with empty bids and not loading", () => {
    const store = useBidsStore();
    expect(store.bids).toEqual([]);
    expect(store.loading).toBe(false);
  });

  it("fetches bids from the API", async () => {
    const mockBids = [
      {
        vehicle_external_id: "abc-123",
        vehicle_name: "2021 Honda Civic EX",
        bid_amount: 1500,
        is_buy_now: false,
        bid_time: "2026-04-07T12:00:00Z",
      },
      {
        vehicle_external_id: "def-456",
        vehicle_name: "2020 BMW M4 Competition",
        bid_amount: 25000,
        is_buy_now: true,
        bid_time: "2026-04-07T11:00:00Z",
      },
    ];

    mockGet.mockResolvedValueOnce({
      data: mockBids,
      error: undefined,
      response: new Response(),
    } as never);

    const store = useBidsStore();
    await store.fetchBids();

    expect(mockGet).toHaveBeenCalledWith("/bids/");
    expect(store.bids).toEqual(mockBids);
    expect(store.loading).toBe(false);
  });

  it("sets loading to true while fetching", async () => {
    let resolvePromise: (value: unknown) => void;
    const pending = new Promise((resolve) => {
      resolvePromise = resolve;
    });

    mockGet.mockReturnValueOnce(pending as never);

    const store = useBidsStore();
    const fetchPromise = store.fetchBids();

    expect(store.loading).toBe(true);

    resolvePromise!({ data: [], error: undefined, response: new Response() });
    await fetchPromise;

    expect(store.loading).toBe(false);
  });

  it("sets loading to false even on error", async () => {
    mockGet.mockRejectedValueOnce(new Error("network error"));

    const store = useBidsStore();
    await expect(store.fetchBids()).rejects.toThrow("network error");

    expect(store.loading).toBe(false);
  });

  it("keeps existing bids if API returns no data", async () => {
    mockGet.mockResolvedValueOnce({
      data: [{ vehicle_external_id: "abc", bid_amount: 100 }],
      error: undefined,
      response: new Response(),
    } as never);

    const store = useBidsStore();
    await store.fetchBids();
    expect(store.bids).toHaveLength(1);

    mockGet.mockResolvedValueOnce({
      data: undefined,
      error: { detail: "server error" },
      response: new Response(),
    } as never);

    await store.fetchBids();
    expect(store.bids).toHaveLength(1);
  });
});
