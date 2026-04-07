import { type MaybeRefOrGetter, toValue, ref, computed, onUnmounted } from 'vue'

const MAX_AUCTION_DURATION_HOURS = ref<number | null>(null)
let configLoaded = false

export async function loadAuctionConfig() {
  if (configLoaded) return
  const { api } = await import('@/lib/api/client')
  const { data } = await api.GET('/vehicles/config')
  if (data?.max_auction_duration_hours) {
    MAX_AUCTION_DURATION_HOURS.value = data.max_auction_duration_hours
  }
  configLoaded = true
}

export function useAuctionTime(
  auctionStart: MaybeRefOrGetter<string | undefined>,
  live = false,
) {
  const now = ref(Date.now())
  let timer: ReturnType<typeof setInterval> | null = null

  if (live) {
    timer = setInterval(() => {
      now.value = Date.now()
    }, 1000)
    onUnmounted(() => {
      if (timer) clearInterval(timer)
    })
  }

  const auctionEnd = computed(() => {
    const start = toValue(auctionStart)
    if (!start || MAX_AUCTION_DURATION_HOURS.value == null) return null
    return new Date(start).getTime() + MAX_AUCTION_DURATION_HOURS.value * 3600 * 1000
  })

  const ended = computed(() => {
    if (auctionEnd.value == null) return false
    return now.value >= auctionEnd.value
  })

  const timeRemaining = computed(() => {
    if (auctionEnd.value == null) return ''
    if (ended.value) return 'Ended'

    const diff = auctionEnd.value - now.value
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    const seconds = Math.floor((diff % (1000 * 60)) / 1000)

    if (days > 0) return `${days}d ${hours}h ${minutes}m`
    if (hours > 0) return `${hours}h ${minutes}m ${seconds}s`
    return `${minutes}m ${seconds}s`
  })

  return { ended, timeRemaining }
}
