package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var ctx = context.Background()

func init() {
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6773")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/app-config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

}

type Server struct {
	redis *redis.Client
}

func (s *Server) Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := s.redis.Set(ctx, "name", r.URL.Query().Get("name"), 0).Err()
			if err != nil {
				log.Println("ERROR: ", err)
			}

		case http.MethodGet:
			name, err := s.redis.Get(ctx, "name").Result()
			if err != nil {
				log.Println("ERROR: ", err)
			}

			w.Write([]byte(name))
		}
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
	})

	s := Server{rdb}
	s.Run()
}
