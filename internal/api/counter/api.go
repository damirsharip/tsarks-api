package counter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
	return &redisCache{
		host: host,
		db:   db,
		exp:  exp,
	}
}

type redisCache struct {
	host string
	db   int
	exp  time.Duration
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Setup(router *httprouter.Router) {
	c := cache.getClient()

	err := c.Set("counter", 0, 0).Err()
	if err != nil {
		return
	}

	router.Handler(http.MethodPut, "/rest/counter/add/:number", cache.Add())
	router.Handler(http.MethodPut, "/rest/counter/sub/:number", cache.Sub())
	router.Handler(http.MethodGet, "/rest/counter/val", cache.GetCounter())
}

func (cache *redisCache) Add() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cache.getClient()
		number := httprouter.ParamsFromContext(r.Context()).ByName("number")

		numberInt, err := strconv.Atoi(number)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid number")
			return
		}

		val, err := c.Get("counter").Result()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		valInt, err := strconv.Atoi(val)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid number")
			return
		}

		err = c.Set("counter", valInt+numberInt, 0).Err()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"success": "ok",
		})
	})
}

func (cache *redisCache) Sub() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cache.getClient()
		number := httprouter.ParamsFromContext(r.Context()).ByName("number")

		numberInt, err := strconv.Atoi(number)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid number")
			return
		}

		val, err := c.Get("counter").Result()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		valInt, err := strconv.Atoi(val)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid number")
			return
		}

		err = c.Set("counter", valInt-numberInt, 0).Err()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"success": "ok",
		})
	})
}

type GetCounterResponce struct {
	Answer int `json:"answer"`
}

func (cache *redisCache) GetCounter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cache.getClient()
		val, err := c.Get("counter").Result()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		valInt, err := strconv.Atoi(val)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid number")
			return
		}

		respondJSON(w, http.StatusOK, GetCounterResponce{Answer: valInt})
	})
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
