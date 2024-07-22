package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// Number of requests allowed in 5 seconds
	limiterRequests     = 5
	limiterCooldown     = 5
	limitClients        = 3
	clearClientsSeconds = 5
)

type Clients struct {
	m  map[string]struct{}
	mu sync.Mutex
}

func NewClients() *Clients {
	clients := &Clients{
		m:  make(map[string]struct{}),
		mu: sync.Mutex{},
	}

	// очищаем мапу, чтобы другие пользователи могли иметь шанс подключиться
	go func() {
		for {
			time.Sleep(clearClientsSeconds * time.Second)
			clients.mu.Lock()
			clients.m = make(map[string]struct{})
			clients.mu.Unlock()
		}
	}()

	return clients
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func applyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func middlewareLimitClients(clients *Clients) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			clients.mu.Lock()
			if _, ok := clients.m[ip]; !ok {
				if len(clients.m)+1 > limitClients {
					clients.mu.Unlock()
					// тут не совсем понимаю, что лучше отдавать, поэтому оставил StatusTooManyRequests
					// буду благодарен, если подскажите :)
					http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
					return
				}

				clients.m[ip] = struct{}{}
			}
			clients.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func middlewareLimitRequests(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	limiter := rate.NewLimiter(rate.Every(limiterCooldown*time.Second), limiterRequests)
	clients := NewClients()

	mux := http.NewServeMux()
	mux.Handle("/", applyMiddlewares(
		http.HandlerFunc(handler),
		middlewareLimitRequests(limiter),
		middlewareLimitClients(clients),
	))

	fmt.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error server running on port 8080: %s", err.Error())
	}
}
