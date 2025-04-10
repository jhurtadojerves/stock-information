<template>
    <div class="border p-6 rounded-lg shadow-md bg-white space-y-2">
      <template v-if="stock">
        <div><strong>Ticker:</strong> <span class="font-mono">{{ stock.ticker }}</span></div>
        <div><strong>Compañía:</strong> {{ stock.company }}</div>
        <div><strong>Broker:</strong> {{ stock.brokerage }}</div>
        <div><strong>Acción:</strong> {{ stock.action }}</div>
        <div><strong>Rating:</strong> {{ stock.rating_from }} → {{ stock.rating_to }}</div>
        <div><strong>Target:</strong> ${{ stock.target_from }} → ${{ stock.target_to }}</div>
      </template>
  
      <template v-else>
        <div class="space-y-2 animate-pulse">
          <div class="h-4 bg-gray-200 rounded w-1/3"></div>
          <div class="h-4 bg-gray-200 rounded w-2/3"></div>
          <div class="h-4 bg-gray-200 rounded w-1/2"></div>
          <div class="h-4 bg-gray-200 rounded w-3/4"></div>
          <div class="h-4 bg-gray-200 rounded w-2/4"></div>
          <div class="h-4 bg-gray-200 rounded w-1/4"></div>
        </div>
      </template>
    </div>
  </template>
  
  <script lang="ts" setup>
  import { ref, watchEffect } from 'vue'
  import type { Stock } from '../utils/api'
  
  const props = defineProps<{
    stockJson: string
  }>()
  
  const stock = ref<Stock | null>(null)
  
  watchEffect(() => {
    try {
      stock.value = JSON.parse(props.stockJson)
    } catch (e) {
      console.error('❌ Error parseando stockJson', e)
    }
  })
  </script>
