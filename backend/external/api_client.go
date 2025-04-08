package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"stock-information/internal/domain"
	"strconv"
	"strings"
)

var BaseURL = "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"

type ApiResponse struct {
	Items    []ApiItem `json:"items"`
	NextPage string    `json:"next_page"`
}

type ApiItem struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
}

func FetchAllRecommendations() ([]domain.Recommendation, error) {
	token := os.Getenv("STOCK_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("STOCK_API_TOKEN no está definido en el entorno")
	}

	var allItems []domain.Recommendation
	nextPage := ""

	for {
		url := BaseURL
		if nextPage != "" {
			url += "?next_page=" + nextPage
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var apiResp ApiResponse

		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, err
		}

		for _, item := range apiResp.Items {
			rec := domain.Recommendation{
				Ticker:     item.Ticker,
				Company:    item.Company,
				Brokerage:  item.Brokerage,
				Action:     item.Action,
				RatingFrom: item.RatingFrom,
				RatingTo:   item.RatingTo,
				TargetFrom: ParseDollar(item.TargetFrom),
				TargetTo:   ParseDollar(item.TargetTo),
			}
			allItems = append(allItems, rec)
		}

		if apiResp.NextPage == "" {
			break
		}
		nextPage = apiResp.NextPage
	}

	return allItems, nil
}

func ParseDollar(input string) float64 {
	trimmedInput := strings.TrimSpace(input) // Eliminar espacios en blanco al inicio y al final
	s := strings.ReplaceAll(trimmedInput, "$", "")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
