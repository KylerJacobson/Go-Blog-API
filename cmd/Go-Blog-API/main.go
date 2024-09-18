package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/KylerJacobson/Go-Blog-API/internal/db/config"
	postsRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/posts"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from %s", r.URL.Path[1:])
}

func main() {

	// Retrieve environment variables
	dbConn := config.GetDBConn()
	defer dbConn.Close(context.Background())

	postsApi := posts.New(postsRepo.New(dbConn))
	http.HandleFunc("GET /", handler)
	// POSTS
	http.HandleFunc("GET /api/posts/public", postsApi.GetRecentPublicPosts)
	http.HandleFunc("GET /api/posts/recent", postsApi.GetRecentPosts)
	http.HandleFunc("GET /api/media/:id", postsApi.GetRecentPosts)
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
