package main

import (
	"context"
	"log"
	"net/http"

	"tech-tsarka/internal/api/counter"
	"tech-tsarka/internal/api/find"
	"tech-tsarka/internal/api/hash"
	"tech-tsarka/internal/api/substr"
	"tech-tsarka/internal/api/user"
	"tech-tsarka/internal/config"
	"tech-tsarka/internal/pkg/pgx"
	"tech-tsarka/internal/service/userservice"
	"tech-tsarka/internal/storage/user/pguserstorage"

	"github.com/julienschmidt/httprouter"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.ReadConfigYaml(); err != nil {
		log.Println("Failed to initialize config: ", err.Error())
		return
	}
	cfg := config.Get()

	pool, err := pgx.NewPool(ctx, cfg.PgxPool, cfg.Database.Dsn)
	if err != nil {
		log.Println("Failed to open a postgres pool: ", err.Error())
		return
	}
	defer pool.Close()

	router := httprouter.New()

	// task 1
	router.Handler(http.MethodPost, "/rest/substr/find", substr.FindSubstr())

	// task 2
	router.Handler(http.MethodPost, "/rest/email/check", find.FindEmail())
	router.Handler(http.MethodPost, "/rest/iin/check", find.FindIIN())

	// task 3 - counter service
	counterClient := counter.NewRedisCache("localhost:6379", 0, 1)
	counterClient.Setup(router)

	// task 4 - user service
	userStorage := pguserstorage.NewStorage(pool)
	userService := userservice.NewService(userStorage)
	userHandler := user.New(userService)
	userHandler.Setup(router)

	// task 5
	hashClient := hash.NewRedisCache("localhost:6379", 0, 1)
	hashClient.Setup(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
