package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// func init() {
// 	viper.SetDefault("redis.host", "localhost")
// 	viper.SetDefault("redis.port", "6773")
// 	viper.SetDefault("redis.password", "")
// 	viper.SetDefault("redis.db", 0)

// 	viper.SetConfigName("config")
// 	viper.SetConfigType("toml")
// 	viper.AddConfigPath("/app-config")
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("fatal error config file: %w", err))
// 	}
// }

type Server struct {
	redis *redis.Client
}

func (s *Server) Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			log.Printf(http.MethodPost + " was requested with " + r.URL.Query().Get("name"))
			err := s.redis.Set(ctx, "name", r.URL.Query().Get("name"), 0).Err()
			if err != nil {
				log.Println("ERROR: ", err)
			}

		case http.MethodGet:
			name, err := s.redis.Get(ctx, "name").Result()
			if err != nil {
				log.Println("ERROR: ", err)
			}
			log.Printf(http.MethodGet + " was requested and return " + name)
			w.Write([]byte(name))
		}
	})
	log.Printf("Listening on" + os.Getenv("SERVER_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil))
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	s := Server{rdb}
	s.Run()
}
