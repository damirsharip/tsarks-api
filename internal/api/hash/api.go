package hash

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

var semaphore = make(chan struct{}, 5)

func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
	return &redisCache{
		host:      host,
		db:        db,
		exp:       exp,
		semaphore: make(chan struct{}, 5),
	}
}

type redisCache struct {
	host      string
	db        int
	exp       time.Duration
	semaphore chan struct{}
}

func (i *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     i.host,
		Password: "",
		DB:       i.db,
	})
}

func (i *redisCache) Setup(router *httprouter.Router) {
	router.Handler(http.MethodPost, "/rest/hash/calc", i.CalcHash())
	router.Handler(http.MethodGet, "/rest/hash/result/:id", i.GetResult())
}

type CalcHashRequest struct {
	Word string `json:"word"`
}

func (r *CalcHashRequest) Validate() error {
	if r.Word == "" {
		return fmt.Errorf("word required")
	}
	return nil
}

type CalcHashResponse struct {
	Answer string `json:"answer"`
}

func (i *redisCache) CalcHash() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var req CalcHashRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		err := req.Validate()
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("bad request: %s", err))
			return
		}

		// Convert UUID to string
		uuidStr := uuid.New().String()

		c := i.getClient()

		err = c.Set(uuidStr, 0, 0).Err()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		go hasher(uuidStr, req.Word, c, i.semaphore)

		respondJSON(w, http.StatusOK, CalcHashResponse{Answer: uuidStr})
	})
}

func hasher(id string, s string, c *redis.Client, semaphore chan struct{}) {
	var res int
	// Вычисляем хэш
	hash := crc64.Checksum([]byte(s), crc64.MakeTable(crc64.ECMA)) // Вычисление CRC64 хэша

	semaphore <- struct{}{}
	var mutex sync.Mutex
	// Выполняем шаги 2-3 в течение минуты
	endTime := time.Now().Add(time.Minute)
	for time.Now().Before(endTime) {
		timestamp := getCurrentTimestamp(&mutex)

		// Выполняем операцию "И" между timestamp и хэшем
		result := timestamp & int64(hash)

		// Считаем количество единиц в двоичной записи числа
		count := countOnes(result)

		res += count

		time.Sleep(5 * time.Second)
	}

	err := c.Set(id, res, 0).Err()
	if err != nil {
		return
	}

	<-semaphore
}

func (i *redisCache) GetResult() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")

		if _, err := uuid.Parse(id); err != nil {
			respondError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		c := i.getClient()

		val, err := c.Get(id).Result()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if val == "0" {
			respondJSON(w, http.StatusOK, CalcHashResponse{Answer: "PENDING"})
		} else {
			respondJSON(w, http.StatusOK, CalcHashResponse{Answer: val})
		}
	})
}

func getCurrentTimestamp(mutex *sync.Mutex) int64 {
	mutex.Lock()
	defer mutex.Unlock()
	return time.Now().UnixNano()
}

func countOnes(num int64) int {
	count := 0
	for num != 0 {
		count += int(num & 1)
		num >>= 1
	}
	return count
}

func decodeJSONBody(_ http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
