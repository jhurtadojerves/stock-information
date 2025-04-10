<template>
    <div>
      <div class="flex justify-between mb-4">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Buscar por compañía"
          class="border px-3 py-2 rounded w-full max-w-sm"
        />
        <button @click="fetchStockRecommendation" class="ml-4 bg-blue-600 text-white px-4 py-2 rounded">
          Recomendar
        </button>
      </div>
  
      <div v-if="loading" class="text-center py-4">Cargando...</div>
  
      <div v-else>
        <table class="min-w-full border divide-y divide-gray-200">
          <thead>
            <tr>
                <th @click="toggleSort" class="cursor-pointer px-4 py-2 select-none flex items-center space-x-1">
                    <span>Ticker</span>
                    <svg
                        v-if="sort === 'asc'"
                        class="w-4 h-4 text-gray-500"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        viewBox="0 0 24 24"
                    >
                        <path d="M5 15l7-7 7 7" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                    <svg
                        v-else
                        class="w-4 h-4 text-gray-500"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        viewBox="0 0 24 24"
                    >
                        <path d="M19 9l-7 7-7-7" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                    </th>
              <th class="px-4 py-2">Compañía</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="stock in stocks"
              :key="stock.ticker"
              class="hover:bg-gray-100 cursor-pointer"
              @click="goToDetail(stock.ticker)"
            >
              <td class="px-4 py-2 font-mono">{{ stock.ticker }}</td>
              <td class="px-4 py-2">{{ stock.company }}</td>
            </tr>
          </tbody>
        </table>
  
        <div class="mt-4 flex justify-center space-x-2">
          <button @click="prevPage" :disabled="page === 1">Anterior</button>
          <span>Página {{ page }}</span>
          <button @click="nextPage" :disabled="!hasNextPage">Siguiente</button>
        </div>
      </div>
  
      <RecommendModal v-if="showModal" :stock="recommendedStock" @close="closeModal" />
    </div>
  </template>
  
  <script lang="ts" setup>
  import { ref, onMounted, watch } from 'vue'
  import RecommendModal from './RecommendModal.vue'
  import { fetchStocks, fetchRecommendation } from '../utils/api'
    import type { Stock } from '../utils/api'
  
  const stocks = ref<Stock[]>([])
  const loading = ref(true)
  const page = ref(1)
  const sort = ref<'asc' | 'desc'>('asc')
  const searchQuery = ref('')
  const hasNextPage = ref(false)
  
  // Función principal para cargar datos
  async function loadStocks() {
  console.log('🔁 Ejecutando loadStocks...')
  loading.value = true
  try {
    const res = await fetchStocks(page.value, sort.value, searchQuery.value)
    console.log('📦 Respuesta de fetchStocks:', res)
    stocks.value = res.data
    hasNextPage.value = res.nextPage !== null
  } catch (err) {
    console.error('❌ Error al obtener acciones:', err)
  } finally {
    loading.value = false
  }
}
  
  // Al montar el componente
  onMounted(() => {
    loadStocks()
  })
  
  // Buscar al cambiar el texto
  watch(searchQuery, () => {
    page.value = 1
    loadStocks()
  })
  
  // Cambio de orden
  function toggleSort() {
    sort.value = sort.value === 'asc' ? 'desc' : 'asc'
    loadStocks()
  }
  
  // Navegación
  function nextPage() {
    if (hasNextPage.value) {
      page.value++
      loadStocks()
    }
  }
  
  function prevPage() {
    if (page.value > 1) {
      page.value--
      loadStocks()
    }
  }
  
  // Detalle
  function goToDetail(ticker: string) {
    window.location.href = `/stock/${ticker}`
  }
  
  // Modal de recomendación
  const showModal = ref(false)
  const recommendedStock = ref<Stock | null>(null)
  
  function closeModal() {
    showModal.value = false
  }
  
  async function fetchStockRecommendation() {
    try {
      const data = await fetchRecommendation()
      recommendedStock.value = data
      showModal.value = true
    } catch (err) {
      console.error('Error al obtener recomendación:', err)
    }
  }
  </script>
