import { ref, reactive } from "vue";
import { defineStore } from "pinia";
import { api } from "@/lib/api/client";
import type { components } from "@/lib/api/v1";

export type Vehicle = components["schemas"]["Vehicle"];
export type VehicleFilterOptions = components["schemas"]["VehicleFilterOptions"];

export interface VehicleFilters {
  yearMin?: number;
  yearMax?: number;
  makes: string[];
  models: string[];
  bodyStyles: string[];
  exteriorColors: string[];
  interiorColors: string[];
  transmissions: string[];
  drivetrains: string[];
  fuelTypes: string[];
  titleStatuses: string[];
  odometerMin?: number;
  odometerMax?: number;
  conditionMin?: number;
  conditionMax?: number;
}

function emptyFilters(): VehicleFilters {
  return {
    yearMin: undefined,
    yearMax: undefined,
    makes: [],
    models: [],
    bodyStyles: [],
    exteriorColors: [],
    interiorColors: [],
    transmissions: [],
    drivetrains: [],
    fuelTypes: [],
    titleStatuses: [],
    odometerMin: undefined,
    odometerMax: undefined,
    conditionMin: undefined,
    conditionMax: undefined,
  };
}

export const useVehiclesStore = defineStore("vehicles", () => {
  const vehicles = ref<Vehicle[]>([]);
  const filterOptions = ref<VehicleFilterOptions>({});
  const filters = reactive<VehicleFilters>(emptyFilters());
  const sort = ref<string>("");
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function fetchFilterOptions() {
    try {
      const { data, error: err } = await api.GET("/vehicles/filters");
      if (err) return;
      if (data) {
        filterOptions.value = data;
      }
    } catch {
      // filter fetch failure is non-critical; don't overwrite shared error
    }
  }

  let latestFetchId = 0;

  async function fetchVehicles() {
    const fetchId = ++latestFetchId;
    loading.value = true;
    error.value = null;

    const params: Record<string, unknown> = {};

    // Normalize inverted ranges by swapping min/max
    let yearMin = filters.yearMin;
    let yearMax = filters.yearMax;
    if (yearMin != null && yearMax != null && yearMin > yearMax) {
      [yearMin, yearMax] = [yearMax, yearMin];
    }
    if (yearMin != null) params.year_min = yearMin;
    if (yearMax != null) params.year_max = yearMax;

    if (filters.makes.length) params.make = filters.makes;
    if (filters.models.length) params.model = filters.models;
    if (filters.bodyStyles.length) params.body_style = filters.bodyStyles;
    if (filters.exteriorColors.length) params.exterior_color = filters.exteriorColors;
    if (filters.interiorColors.length) params.interior_color = filters.interiorColors;
    if (filters.transmissions.length) params.transmission = filters.transmissions;
    if (filters.drivetrains.length) params.drivetrain = filters.drivetrains;
    if (filters.fuelTypes.length) params.fuel_type = filters.fuelTypes;
    if (filters.titleStatuses.length) params.title_status = filters.titleStatuses;

    let odometerMin = filters.odometerMin;
    let odometerMax = filters.odometerMax;
    if (odometerMin != null && odometerMax != null && odometerMin > odometerMax) {
      [odometerMin, odometerMax] = [odometerMax, odometerMin];
    }
    if (odometerMin != null) params.odometer_min = odometerMin;
    if (odometerMax != null) params.odometer_max = odometerMax;

    let conditionMin = filters.conditionMin;
    let conditionMax = filters.conditionMax;
    if (conditionMin != null && conditionMax != null && conditionMin > conditionMax) {
      [conditionMin, conditionMax] = [conditionMax, conditionMin];
    }
    if (conditionMin != null) params.condition_min = String(conditionMin);
    if (conditionMax != null) params.condition_max = String(conditionMax);

    if (sort.value) params.sort = sort.value;

    try {
      const { data, error: err } = await api.GET("/vehicles/", {
        params: { query: params as Record<string, never> },
      });
      if (fetchId !== latestFetchId) return;
      if (err) {
        error.value = "Unable to connect to server";
      } else {
        vehicles.value = data ?? [];
      }
    } catch {
      if (fetchId !== latestFetchId) return;
      error.value = "Unable to connect to server";
    } finally {
      if (fetchId === latestFetchId) {
        loading.value = false;
      }
    }
  }

  function toggleFilter(
    key: keyof Pick<
      VehicleFilters,
      | "makes"
      | "models"
      | "bodyStyles"
      | "exteriorColors"
      | "interiorColors"
      | "transmissions"
      | "drivetrains"
      | "fuelTypes"
      | "titleStatuses"
    >,
    value: string,
  ) {
    const arr = filters[key];
    const idx = arr.indexOf(value);
    if (idx === -1) {
      arr.push(value);
    } else {
      arr.splice(idx, 1);
    }
  }

  function resetFilters() {
    Object.assign(filters, emptyFilters());
  }

  return {
    vehicles,
    filterOptions,
    filters,
    sort,
    loading,
    error,
    fetchFilterOptions,
    fetchVehicles,
    toggleFilter,
    resetFilters,
  };
});
