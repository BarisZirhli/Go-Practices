package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

var store sync.Map

func generateShortCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Content-Type kontrolü
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	type Request struct {
		URL string `json:"url"`
	}
	type Response struct {
		Short string `json:"short"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}


	var code string
	for {
		code = generateShortCode(6)
		if _, exists := store.Load(code); !exists {
			break
		}
	}

	store.Store(code, req.URL)

	w.Header().Set("Content-Type", "application/json")
	resp := Response{Short: fmt.Sprintf("http://%s/%s", r.Host, code)}
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.Error(w, "URL code is missing", http.StatusBadRequest)
		return
	}

	url, ok := store.Load(code)
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url.(string), http.StatusFound)
}

func helper() {
	

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Server started at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
