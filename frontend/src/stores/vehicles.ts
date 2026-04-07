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
    const { data, error: err } = await api.GET("/vehicles/filters");
    if (err) {
      error.value = "Unable to connect to server";
      return;
    }
    if (data) {
      filterOptions.value = data;
    }
  }

  async function fetchVehicles() {
    loading.value = true;
    error.value = null;

    const params: Record<string, unknown> = {};

    // Only include range params when min <= max (skip the pair if inverted)
    const yearOk =
      filters.yearMin == null || filters.yearMax == null || filters.yearMin <= filters.yearMax;
    if (yearOk) {
      if (filters.yearMin != null) params.year_min = filters.yearMin;
      if (filters.yearMax != null) params.year_max = filters.yearMax;
    }

    if (filters.makes.length) params.make = filters.makes;
    if (filters.models.length) params.model = filters.models;
    if (filters.bodyStyles.length) params.body_style = filters.bodyStyles;
    if (filters.exteriorColors.length) params.exterior_color = filters.exteriorColors;
    if (filters.interiorColors.length) params.interior_color = filters.interiorColors;
    if (filters.transmissions.length) params.transmission = filters.transmissions;
    if (filters.drivetrains.length) params.drivetrain = filters.drivetrains;
    if (filters.fuelTypes.length) params.fuel_type = filters.fuelTypes;
    if (filters.titleStatuses.length) params.title_status = filters.titleStatuses;

    const odometerOk =
      filters.odometerMin == null ||
      filters.odometerMax == null ||
      filters.odometerMin <= filters.odometerMax;
    if (odometerOk) {
      if (filters.odometerMin != null) params.odometer_min = filters.odometerMin;
      if (filters.odometerMax != null) params.odometer_max = filters.odometerMax;
    }

    const conditionOk =
      filters.conditionMin == null ||
      filters.conditionMax == null ||
      filters.conditionMin <= filters.conditionMax;
    if (conditionOk) {
      if (filters.conditionMin != null) params.condition_min = String(filters.conditionMin);
      if (filters.conditionMax != null) params.condition_max = String(filters.conditionMax);
    }

    if (sort.value) params.sort = sort.value;

    try {
      const { data, error: err } = await api.GET("/vehicles/", {
        params: { query: params as Record<string, never> },
      });
      if (err) {
        error.value = "Unable to connect to server";
      } else {
        vehicles.value = data ?? [];
      }
    } catch {
      error.value = "Unable to connect to server";
    } finally {
      loading.value = false;
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
