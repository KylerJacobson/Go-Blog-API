package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KylerJacobson/Go-Blog-API/internal/db/config"
	postsRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/posts"
	"github.com/KylerJacobson/Go-Blog-API/logger"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world from %s", r.URL.Path[1:])
}

func main() {
	env := os.Getenv("ENVIRONMENT")
	err := logger.Init(env)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	dbConn := config.GetDBConn()
	defer dbConn.Close(context.Background())

	postsApi := posts.New(postsRepo.New(dbConn))
	http.HandleFunc("GET /", handler)
	// POSTS
	http.HandleFunc("GET /api/posts/public", postsApi.GetRecentPublicPosts)
	http.HandleFunc("GET /api/posts/recent", postsApi.GetRecentPosts)
	http.HandleFunc("GET /api/media/:id", postsApi.GetRecentPosts)
	logger.Sugar.Infof("Logging level set to %s", env)
	logger.Sugar.Infof("listening on port: %d", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
