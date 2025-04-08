package external_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"stock-information/external"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAllRecommendations(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/production/swechallenge/list" {
			assert.Equal(t, "Bearer test_token", r.Header.Get("Authorization"))

			// Simular una respuesta exitosa con una página de datos y NextPage vacío
			response := external.ApiResponse{
				Items: []external.ApiItem{
					{
						Ticker:     "AAPL",
						Company:    "Apple Inc.",
						Brokerage:  "Goldman Sachs",
						Action:     "Buy",
						RatingFrom: "Neutral",
						RatingTo:   "Buy",
						TargetFrom: "$150.00",
						TargetTo:   "$180.50",
					},
					{
						Ticker:     "GOOG",
						Company:    "Alphabet Inc.",
						Brokerage:  "Morgan Stanley",
						Action:     "Hold",
						RatingFrom: "Buy",
						RatingTo:   "Hold",
						TargetFrom: "$2500",
						TargetTo:   "$2700.75",
					},
				},
				NextPage: "", // NextPage está vacío
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ruta no encontrada en el mock"})
	}))
	defer testServer.Close()

	os.Setenv("STOCK_API_TOKEN", "test_token")
	defer os.Unsetenv("STOCK_API_TOKEN")

	originalBaseURL := external.BaseURL
	external.BaseURL = testServer.URL + "/production/swechallenge/list"
	defer func() { external.BaseURL = originalBaseURL }()

	recommendations, err := external.FetchAllRecommendations()

	assert.NoError(t, err)
	assert.Len(t, recommendations, 2)

	assert.Equal(t, "AAPL", recommendations[0].Ticker)
	assert.Equal(t, "Apple Inc.", recommendations[0].Company)
	assert.Equal(t, 150.00, recommendations[0].TargetFrom)
	assert.Equal(t, 180.50, recommendations[0].TargetTo)

	assert.Equal(t, "GOOG", recommendations[1].Ticker)
	assert.Equal(t, "Alphabet Inc.", recommendations[1].Company)
	assert.Equal(t, 2500.00, recommendations[1].TargetFrom)
	assert.Equal(t, 2700.75, recommendations[1].TargetTo)

}

func TestFetchAllRecommendationsNoToken(t *testing.T) {
	os.Unsetenv("STOCK_API_TOKEN")
	recommendations, err := external.FetchAllRecommendations()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "STOCK_API_TOKEN no está definido en el entorno")
	assert.Nil(t, recommendations)
}

func TestFetchAllRecommendationsAPIError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer testServer.Close()

	os.Setenv("STOCK_API_TOKEN", "test_token")
	defer os.Unsetenv("STOCK_API_TOKEN")

	originalBaseURL := external.BaseURL
	external.BaseURL = testServer.URL
	defer func() { external.BaseURL = originalBaseURL }()

	recommendations, err := external.FetchAllRecommendations()

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, recommendations)
}

func TestFetchAllRecommendationsInvalidJSON(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items": [{"ticker": "AAPL",}`))
	}))
	defer testServer.Close()

	os.Setenv("STOCK_API_TOKEN", "test_token")
	defer os.Unsetenv("STOCK_API_TOKEN")

	originalBaseURL := external.BaseURL
	external.BaseURL = testServer.URL
	defer func() { external.BaseURL = originalBaseURL }()

	recommendations, err := external.FetchAllRecommendations()

	assert.Error(t, err)
	assert.Nil(t, recommendations)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestParseDollar(t *testing.T) {
	assert.Equal(t, 150.00, external.ParseDollar("$150.00"))
	assert.Equal(t, 2700.75, external.ParseDollar("$2700.75"))
	assert.Equal(t, 100.0, external.ParseDollar("100"))
	assert.Equal(t, 0.0, external.ParseDollar("$"))
	assert.Equal(t, 123.45, external.ParseDollar(" $123.45 "))
}
