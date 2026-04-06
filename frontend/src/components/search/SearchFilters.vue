<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useVehiclesStore } from '@/stores/vehicles'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Slider } from '@/components/ui/slider'
import { Button } from '@/components/ui/button'

const store = useVehiclesStore()
const { filterOptions } = storeToRefs(store)
const filters = store.filters

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
  )
})

const yearMin = computed(() => filterOptions.value.year_min ?? 2000)
const yearMax = computed(() => filterOptions.value.year_max ?? 2025)
const odometerMinOpt = computed(() => filterOptions.value.odometer_min ?? 0)
const odometerMaxOpt = computed(() => filterOptions.value.odometer_max ?? 300000)
const conditionMinOpt = computed(() => filterOptions.value.condition_min ?? 1)
const conditionMaxOpt = computed(() => filterOptions.value.condition_max ?? 5)

const yearRange = computed({
  get: () => [filters.yearMin ?? yearMin.value, filters.yearMax ?? yearMax.value],
  set: (val: number[]) => {
    filters.yearMin = val[0] === yearMin.value ? undefined : val[0]
    filters.yearMax = val[1] === yearMax.value ? undefined : val[1]
  },
})

const odometerRange = computed({
  get: () => [
    filters.odometerMin ?? odometerMinOpt.value,
    filters.odometerMax ?? odometerMaxOpt.value,
  ],
  set: (val: number[]) => {
    filters.odometerMin = val[0] === odometerMinOpt.value ? undefined : val[0]
    filters.odometerMax = val[1] === odometerMaxOpt.value ? undefined : val[1]
  },
})

const conditionRange = computed({
  get: () => [
    filters.conditionMin ?? conditionMinOpt.value,
    filters.conditionMax ?? conditionMaxOpt.value,
  ],
  set: (val: number[]) => {
    filters.conditionMin = val[0] === conditionMinOpt.value ? undefined : val[0]
    filters.conditionMax = val[1] === conditionMaxOpt.value ? undefined : val[1]
  },
})

const availableModels = computed(() => {
  if (!filterOptions.value.models) return []
  return filterOptions.value.models
})

function formatOdometer(km: number | undefined) {
  if (km == null) return ''
  return `${(km / 1000).toFixed(0)}k km`
}

type ArrayFilterKey =
  | 'makes'
  | 'models'
  | 'bodyStyles'
  | 'exteriorColors'
  | 'interiorColors'
  | 'transmissions'
  | 'drivetrains'
  | 'fuelTypes'
  | 'titleStatuses'

function isChecked(key: ArrayFilterKey, value: string) {
  return filters[key].includes(value)
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
      <div class="space-y-3">
        <Label class="text-sm font-medium">Year</Label>
        <Slider v-model="yearRange" :min="yearMin" :max="yearMax" :step="1" />
        <div class="flex justify-between text-xs text-muted-foreground">
          <span>{{ yearRange[0] }}</span>
          <span>{{ yearRange[1] }}</span>
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
              :model-value="isChecked('makes', make!)"
              @update:model-value="store.toggleFilter('makes', make!)"
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
              :model-value="isChecked('models', model!)"
              @update:model-value="store.toggleFilter('models', model!)"
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('bodyStyles', style!)"
              @update:model-value="store.toggleFilter('bodyStyles', style!)"
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('exteriorColors', color!)"
              @update:model-value="store.toggleFilter('exteriorColors', color!)"
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('interiorColors', color!)"
              @update:model-value="store.toggleFilter('interiorColors', color!)"
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('transmissions', trans!)"
              @update:model-value="store.toggleFilter('transmissions', trans!)"
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('drivetrains', dt!)"
              @update:model-value="store.toggleFilter('drivetrains', dt!)"
            />
            <span class="text-sm">{{ dt }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Odometer Range -->
      <div class="space-y-3">
        <Label class="text-sm font-medium">Odometer</Label>
        <Slider v-model="odometerRange" :min="odometerMinOpt" :max="odometerMaxOpt" :step="5000" />
        <div class="flex justify-between text-xs text-muted-foreground">
          <span v-if="odometerRange[0]">{{ formatOdometer(odometerRange[0]) }}</span>
          <span v-if="odometerRange[1]">{{ formatOdometer(odometerRange[1]) }}</span>
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('fuelTypes', fuel!)"
              @update:model-value="store.toggleFilter('fuelTypes', fuel!)"
            />
            <span class="text-sm capitalize">{{ fuel }}</span>
          </label>
        </div>
      </div>

      <Separator />

      <!-- Condition Grade Range -->
      <div class="space-y-3">
        <Label class="text-sm font-medium">Condition Grade</Label>
        <Slider
          v-model="conditionRange"
          :min="conditionMinOpt"
          :max="conditionMaxOpt"
          :step="0.1"
        />
        <div class="flex justify-between text-xs text-muted-foreground">
          <span>{{ conditionRange[0]?.toFixed(1) }}</span>
          <span>{{ conditionRange[1]?.toFixed(1) }}</span>
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
            class="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              :model-value="isChecked('titleStatuses', status!)"
              @update:model-value="store.toggleFilter('titleStatuses', status!)"
            />
            <span class="text-sm capitalize">{{ status }}</span>
          </label>
        </div>
      </div>
    </div>
  </ScrollArea>
</template>
