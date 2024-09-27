package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/KylerJacobson/Go-Blog-API/internal/db/config"
	postsRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
	usersRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/users"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/posts"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/users"
	"github.com/KylerJacobson/Go-Blog-API/logger"
)

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
	usersApi := users.New(usersRepo.New(dbConn))

	http.HandleFunc("GET /api/posts/public", postsApi.GetRecentPublicPosts)
	http.HandleFunc("GET /api/posts/recent", postsApi.GetRecentPosts)
	http.HandleFunc("GET /api/post/{id}", postsApi.GetPostById)
	http.HandleFunc("DELETE /api/post/{id}", postsApi.DeletePostById)
	http.HandleFunc("POST /api/post", postsApi.CreatePost)
	http.HandleFunc("PUT /api/post/{id}", postsApi.UpdatePost)
	http.HandleFunc("GET /api/users/{id}", usersApi.GetUserById)
	http.HandleFunc("DELETE /api/users/{id}", usersApi.DeleteUserById)
	http.HandleFunc("POST /api/users", usersApi.CreateUser)
	http.HandleFunc("POST /api/session", usersApi.LoginUser)
	logger.Sugar.Infof("Logging level set to %s", env)
	logger.Sugar.Infof("listening on port: %d", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
