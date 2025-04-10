export interface Stock {
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_from: string;
  rating_to: string;
  target_from: number;
  target_to: number;
}

export interface PaginatedStocksResponse {
  data: Stock[];
  limit: number;
  nextPage: number | null;
  page: number;
  sort: "asc" | "desc";
  total: number;
  totalInPage: number;
}

const API_BASE = "http://localhost:8080/api/v1/stocks";

/**
 * Obtener lista paginada de stocks desde el backend.
 */
export async function fetchStocks(
  page = 1,
  sort: "asc" | "desc" = "asc",
  search = ""
): Promise<PaginatedStocksResponse> {
  const url = new URL(`${API_BASE}`);
  url.searchParams.set("page", String(page));
  url.searchParams.set("sort", sort);
  if (search) {
    url.searchParams.set("company", search);
  }

  return fetch(url.toString()).then((res) => {
    if (!res.ok) throw new Error("Error al obtener listado de acciones");
    return res.json();
  });
}

/**
 * Obtener una única acción recomendada.
 */
export async function fetchRecommendation(): Promise<Stock> {
  const res = await fetch(`${API_BASE}/recommendation`);
  if (!res.ok) throw new Error("Error al obtener recomendación");
  return res.json();
}

/**
 * Obtener el detalle de una acción por su ticker.
 */
export async function fetchStockByTicker(ticker: string): Promise<Stock> {
  const res = await fetch(`${API_BASE}/${ticker}`);
  if (!res.ok)
    throw new Error(`Error al obtener detalle para el ticker ${ticker}`);
  return res.json();
}
