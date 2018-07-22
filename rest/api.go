package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/EatsLemons/fa_steam_shop/store"
	"github.com/gorilla/mux"
)

type ShopService interface {
	GetPriceForItem(itemName, game, currency string) (*store.Price, error)
}

type CacheService interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Exists(key string) (bool, error)
}

type Rest struct {
	ShopService ShopService
	Cache       CacheService

	httpServer *http.Server
}

func (rs *Rest) Run(port int) {
	log.Printf("[INFO] server started at :%d", port)

	r := mux.NewRouter()
	r.Use(rs.recoverWrap)
	r.HandleFunc("/api/v1/find", rs.findHandler).Methods("GET")

	rs.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	err := rs.httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)
}

func (rs *Rest) findHandler(w http.ResponseWriter, r *http.Request) {
	response := rs.newResponseItem()

	itemName := r.URL.Query().Get("name")
	if itemName == "" {
		response.Errors = append(response.Errors, ErrorRs{Message: "name is empty"})
		rs.makeJSONResponse(w, response)
		return
	}

	game := r.URL.Query().Get("game")
	if itemName == "" {
		response.Errors = append(response.Errors, ErrorRs{Message: "game is empty"})
		rs.makeJSONResponse(w, response)
		return
	}

	game = strings.ToUpper(game)

	itemPrice := &store.Price{}
	if ok, err := rs.Cache.Exists(itemName); ok && err == nil {
		itemInfoFromCache, err := rs.Cache.Get(itemName)
		if err == nil {
			err = json.Unmarshal(itemInfoFromCache, itemPrice)
			if err != nil {
				response.Errors = append(response.Errors, ErrorRs{Message: err.Error()})
				rs.makeJSONResponse(w, response)
				return
			}

			result := ItemInfo{
				Name: itemName,
				Price: &Money{
					Currency: "USD",
					Amount:   itemPrice.Cost,
				},
			}

			response.Result = &result
			rs.makeJSONResponse(w, response)
			return
		}
	}

	itemPrice, err := rs.ShopService.GetPriceForItem(itemName, game, "USD")
	if err != nil {
		response.Errors = append(response.Errors, ErrorRs{Message: err.Error()})
		rs.makeJSONResponse(w, response)
		return
	}

	itemBytes, err := json.Marshal(itemPrice)
	if err != nil {
		log.Println("[WARN] ", err.Error())
	}

	rs.Cache.Set(itemName, itemBytes)

	result := ItemInfo{
		Name: itemName,
		Price: &Money{
			Currency: "USD",
			Amount:   itemPrice.Cost,
		},
	}

	response.Result = &result
	rs.makeJSONResponse(w, response)
}

func (rs *Rest) recoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var err error
			if r := recover(); r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}

				response := rs.newResponseItem()
				response.Errors = append(response.Errors, ErrorRs{Message: err.Error()})
				log.Println("[WARN] handled error: ", err.Error())
				rs.makeJSONResponse(w, response)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (rs *Rest) makeJSONResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("[WARN] response marshaling fail %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (rs *Rest) newResponseItem() *ItemResponse {
	result := ItemResponse{
		Errors: make([]ErrorRs, 0),
	}

	return &result
}
