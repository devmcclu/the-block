<script setup lang="ts">
import { computed } from 'vue'
import { useVehiclesStore } from '@/stores/vehicles'
import { Icon } from '@iconify/vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import SearchFilters from './SearchFilters.vue'

const store = useVehiclesStore()
const filters = store.filters

const activeFilterCount = computed(() => {
  let count = 0
  count += filters.makes.length
  count += filters.models.length
  count += filters.bodyStyles.length
  count += filters.exteriorColors.length
  count += filters.interiorColors.length
  count += filters.transmissions.length
  count += filters.drivetrains.length
  count += filters.fuelTypes.length
  count += filters.titleStatuses.length
  if (filters.yearMin != null || filters.yearMax != null) count++
  if (filters.odometerMin != null || filters.odometerMax != null) count++
  if (filters.conditionMin != null || filters.conditionMax != null) count++
  return count
})
</script>

<template>
  <Sheet>
    <SheetTrigger as-child>
      <Button variant="outline" size="sm" class="lg:hidden relative">
        <Icon icon="hugeicons:preference-horizontal" class="h-4 w-4 mr-2" />
        Filters
        <Badge
          v-if="activeFilterCount > 0"
          class="ml-2 h-5 w-5 rounded-full p-0 flex items-center justify-center text-[10px]"
        >
          {{ activeFilterCount }}
        </Badge>
      </Button>
    </SheetTrigger>
    <SheetContent side="left" class="w-80 p-0">
      <SheetHeader class="sr-only">
        <SheetTitle>Filters</SheetTitle>
      </SheetHeader>
      <SearchFilters />
    </SheetContent>
  </Sheet>
</template>
