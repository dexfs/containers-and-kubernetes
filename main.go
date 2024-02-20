package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

var ctx = context.Background()

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
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	s := Server{rdb}
	s.Run()
}
