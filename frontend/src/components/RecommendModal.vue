<template>
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-lg p-6 w-full max-w-md relative">
        <button @click="$emit('close')" class="absolute top-2 right-2 text-gray-500 hover:text-black">
          ✕
        </button>
        <h2 class="text-xl font-bold mb-4">Mejor Opción Recomendad​a</h2>
  
        <div class="space-y-2">
          <div><strong>Ticker:</strong> <span class="font-mono">{{ stock.ticker }}</span></div>
          <div><strong>Compañía:</strong> {{ stock.company }}</div>
          <div><strong>Broker:</strong> {{ stock.brokerage }}</div>
          <div><strong>Acción:</strong> {{ stock.action }}</div>
          <div><strong>Rating:</strong> {{ stock.rating_from }} → {{ stock.rating_to }}</div>
          <div><strong>Target:</strong> ${{ stock.target_from }} → ${{ stock.target_to }}</div>
        </div>
  
        <div class="mt-6 text-right">
          <button @click="$emit('close')" class="bg-blue-600 text-white px-4 py-2 rounded">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script lang="ts" setup>
  import { defineProps } from 'vue'
  import type { Stock } from '../utils/api'
  
  defineProps<{
    stock: Stock
  }>()
  </script>
