package shop

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/EatsLemons/fa_steam_shop/store"
)

type SteamAPI struct {
	address string

	httpClient *http.Client
}

const (
	CurrencyUSD = "1"
)

const (
	GameCSGO = "730"
)

func NewSteamAPI(address string) *SteamAPI {
	s := SteamAPI{
		address: address,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	return &s
}

func (s *SteamAPI) GetPriceForItem(itemName, game, currency string) (*store.Price, error) {
	if game == "CSGO" {
		game = GameCSGO
	}
	if currency == "USD" {
		currency = CurrencyUSD
	}

	requestURI := "/market/priceoverview/?appid=" + game + "&currency=" + currency + "&market_hash_name=" + itemName
	response := PriceResponse{}

	err := s.makeGetRequest(requestURI, &response)
	if err != nil {
		return nil, err
	}

	tmpPrice := strings.Replace(response.MedianPrice, "$", "", 1)
	cost, err := strconv.ParseFloat(tmpPrice, 64)
	if err != nil {
		return nil, err
	}

	result := store.Price{
		Item:      itemName,
		Currency:  currency,
		Cost:      cost,
		ActualFor: time.Now(),
	}

	return &result, nil
}

func (s *SteamAPI) makeGetRequest(url string, result interface{}) error {
	requestString := s.address + url
	requestString = strings.Replace(requestString, " ", "%20", -1)
	r, err := s.httpClient.Get(requestString)
	if err != nil {
		log.Printf("[WARN] request to steam has failed %s for query %s", err, url)
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(result)
}

type PriceResponse struct {
	MedianPrice string `json:"median_price"`
}
