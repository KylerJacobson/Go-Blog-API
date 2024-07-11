package main

import (
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

func getRecentPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all recent posts")
}

func main() {

	// Retrieve environment variables
	dbConn := config.GetDBConn()

	postsApi := posts.New(postsRepo.New(dbConn))
	http.HandleFunc("GET /", handler)
	// POSTS

	http.HandleFunc("GET /api/posts", postsApi.Get)
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
