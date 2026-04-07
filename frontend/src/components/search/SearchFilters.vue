<script setup lang="ts">
import { computed } from "vue";
import { storeToRefs } from "pinia";
import { useVehiclesStore } from "@/stores/vehicles";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const store = useVehiclesStore();
const { filterOptions } = storeToRefs(store);
const filters = store.filters;

const hasActiveFilters = computed(() => {
  return (
    filters.makes.length > 0 ||
    filters.models.length > 0 ||
    filters.bodyStyles.length > 0 ||
    filters.exteriorColors.length > 0 ||
    filters.interiorColors.length > 0 ||
    filters.transmissions.length > 0 ||
    filters.drivetrains.length > 0 ||
    filters.fuelTypes.length > 0 ||
    filters.titleStatuses.length > 0 ||
    filters.yearMin != null ||
    filters.yearMax != null ||
    filters.odometerMin != null ||
    filters.odometerMax != null ||
    filters.conditionMin != null ||
    filters.conditionMax != null
  );
});

const availableModels = computed(() => {
  if (!filterOptions.value.models) return [];
  return filterOptions.value.models;
});

function parseNumber(value: string, min?: number, max?: number): number | undefined {
  const n = Number(value);
  if (value === "" || Number.isNaN(n)) return undefined;
  if (min != null && n < min) return min;
  if (max != null && n > max) return max;
  return n;
}

type ArrayFilterKey =
  | "makes"
  | "models"
  | "bodyStyles"
  | "exteriorColors"
  | "interiorColors"
  | "transmissions"
  | "drivetrains"
  | "fuelTypes"
  | "titleStatuses";

function isChecked(key: ArrayFilterKey, value: string) {
  return filters[key].includes(value);
}
</script>

<template>
  <ScrollArea class="h-full">
    <div class="space-y-4 p-4">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">Filters</h2>
        <Button v-if="hasActiveFilters" variant="ghost" size="sm" @click="store.resetFilters()">
          Clear All
        </Button>
      </div>

      <Separator />

      <!-- Year Range -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Year</Label>
        <div class="flex items-center gap-2">
          <Input
            lazy
            type="number"
            min="1900"
            max="2100"
            aria-label="Year minimum"
            :placeholder="String(filterOptions.year_min ?? 'Min')"
            :model-value="filters.yearMin"
            @update:model-value="filters.yearMin = parseNumber(String($event), 1900, 2100)"
          />
          <span class="text-xs text-muted-foreground">to</span>
          <Input
            lazy
            type="number"
            min="1900"
            max="2100"
            aria-label="Year maximum"
            :placeholder="String(filterOptions.year_max ?? 'Max')"
            :model-value="filters.yearMax"
            @update:model-value="filters.yearMax = parseNumber(String($event), 1900, 2100)"
          />
        </div>
      </div>

      <Separator />

      <!-- Make -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Make</Label>
        <div class="space-y-1.5 max-h-48 overflow-y-auto">
          <label
            v-for="make in filterOptions.makes"
            :key="make"
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('makes', make)"
              @update:model-value="store.toggleFilter('makes', make)"
            />
            <span class="text-sm">{{ make }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Model -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Model</Label>
        <div class="space-y-1.5 max-h-48 overflow-y-auto">
          <label
            v-for="model in availableModels"
            :key="model"
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('models', model)"
              @update:model-value="store.toggleFilter('models', model)"
            />
            <span class="text-sm">{{ model }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Body Style -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Body Style</Label>
        <div class="space-y-1.5">
          <label
            v-for="style in filterOptions.body_styles"
            :key="style"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('bodyStyles', style)"
              @update:model-value="store.toggleFilter('bodyStyles', style)"
            />
            <span class="text-sm">{{ style }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Exterior Color -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Exterior Color</Label>
        <div class="space-y-1.5 max-h-48 overflow-y-auto">
          <label
            v-for="color in filterOptions.exterior_colors"
            :key="color"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('exteriorColors', color)"
              @update:model-value="store.toggleFilter('exteriorColors', color)"
            />
            <span class="text-sm">{{ color }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Interior Color -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Interior Color</Label>
        <div class="space-y-1.5 max-h-48 overflow-y-auto">
          <label
            v-for="color in filterOptions.interior_colors"
            :key="color"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('interiorColors', color)"
              @update:model-value="store.toggleFilter('interiorColors', color)"
            />
            <span class="text-sm">{{ color }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Transmission -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Transmission</Label>
        <div class="space-y-1.5">
          <label
            v-for="trans in filterOptions.transmissions"
            :key="trans"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('transmissions', trans)"
              @update:model-value="store.toggleFilter('transmissions', trans)"
            />
            <span class="text-sm capitalize">{{ trans }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Drivetrain -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Drivetrain</Label>
        <div class="space-y-1.5">
          <label
            v-for="dt in filterOptions.drivetrains"
            :key="dt"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('drivetrains', dt)"
              @update:model-value="store.toggleFilter('drivetrains', dt)"
            />
            <span class="text-sm">{{ dt }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Odometer Range -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Odometer (km)</Label>
        <div class="flex items-center gap-2">
          <Input
            lazy
            type="number"
            min="0"
            aria-label="Odometer minimum"
            :placeholder="String(filterOptions.odometer_min ?? 'Min')"
            :model-value="filters.odometerMin"
            @update:model-value="filters.odometerMin = parseNumber(String($event), 0)"
          />
          <span class="text-xs text-muted-foreground">to</span>
          <Input
            lazy
            type="number"
            min="0"
            aria-label="Odometer maximum"
            :placeholder="String(filterOptions.odometer_max ?? 'Max')"
            :model-value="filters.odometerMax"
            @update:model-value="filters.odometerMax = parseNumber(String($event), 0)"
          />
        </div>
      </div>

      <Separator />

      <!-- Fuel Type -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Fuel Type</Label>
        <div class="space-y-1.5">
          <label
            v-for="fuel in filterOptions.fuel_types"
            :key="fuel"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('fuelTypes', fuel)"
              @update:model-value="store.toggleFilter('fuelTypes', fuel)"
            />
            <span class="text-sm capitalize">{{ fuel }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Condition Grade Range -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Condition Grade</Label>
        <div class="flex items-center gap-2">
          <Input
            lazy
            type="number"
            step="0.1"
            min="0"
            max="5"
            aria-label="Condition grade minimum"
            :placeholder="String(filterOptions.condition_min ?? 'Min')"
            :model-value="filters.conditionMin"
            @update:model-value="filters.conditionMin = parseNumber(String($event), 0, 5)"
          />
          <span class="text-xs text-muted-foreground">to</span>
          <Input
            lazy
            type="number"
            step="0.1"
            min="0"
            max="5"
            aria-label="Condition grade maximum"
            :placeholder="String(filterOptions.condition_max ?? 'Max')"
            :model-value="filters.conditionMax"
            @update:model-value="filters.conditionMax = parseNumber(String($event), 0, 5)"
          />
        </div>
      </div>

      <Separator />

      <!-- Title Status -->
      <div class="space-y-2">
        <Label class="text-sm font-medium">Title Status</Label>
        <div class="space-y-1.5">
          <label
            v-for="status in filterOptions.title_statuses"
            :key="status"
            class="flex items-center gap-2 cursor-pointer capitalize"
          >
            <Checkbox
              :model-value="isChecked('titleStatuses', status)"
              @update:model-value="store.toggleFilter('titleStatuses', status)"
            />
            <span class="text-sm capitalize">{{ status }}</span>
          </label>
        </div>
      </div>
    </div>
  </ScrollArea>
</template>
