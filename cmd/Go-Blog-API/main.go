package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KylerJacobson/Go-Blog-API/internal/authorization"
	"github.com/KylerJacobson/Go-Blog-API/internal/db/config"

	mediaRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/media"
	postsRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
	usersRepo "github.com/KylerJacobson/Go-Blog-API/internal/db/users"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/media"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/posts"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/session"
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
	dbPool := config.GetDBConn()
	defer dbPool.Close()

	session.Init()

	mux := http.NewServeMux()

	postsApi := posts.New(postsRepo.New(dbPool))
	usersApi := users.New(usersRepo.New(dbPool))
	sessionApi := session.New(usersRepo.New(dbPool))
	mediaApi := media.New(mediaRepo.New(dbPool))

	// ---------------------------- Posts ----------------------------
	mux.HandleFunc("GET /api/posts", postsApi.GetPosts)
	mux.HandleFunc("GET /api/posts/recent", postsApi.GetRecentPosts)
	mux.HandleFunc("GET /api/posts/{id}", postsApi.GetPostById)
	mux.HandleFunc("DELETE /api/posts/{id}", postsApi.DeletePostById)
	mux.HandleFunc("POST /api/posts", postsApi.CreatePost)
	mux.HandleFunc("PUT /api/posts/{id}", postsApi.UpdatePost)

	// ---------------------------- Users ----------------------------
	mux.HandleFunc("POST /api/users", usersApi.CreateUser)
	mux.HandleFunc("GET /api/user", usersApi.GetUserFromSession)
	mux.HandleFunc("GET /api/users/{id}", usersApi.GetUserById)
	mux.HandleFunc("GET /api/users/list", usersApi.ListUsers)
	// http.HandleFunc("PUT /api/users/{id}", usersApi.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", usersApi.DeleteUserById)

	mux.HandleFunc("POST /api/session", sessionApi.CreateSession)
	mux.HandleFunc("POST /api/verifyToken", authorization.VerifyToken)
	mux.HandleFunc("DELETE /api/session", sessionApi.DeleteSession)

	// ---------------------------- Media ----------------------------
	mux.HandleFunc("POST /api/media", mediaApi.UploadMedia)
	mux.HandleFunc("GET /api/media/{id}", mediaApi.GetMediaByPostId)

	logger.Sugar.Infof("Logging level set to %s", env)
	logger.Sugar.Infof("listening on port: %d", 8080)
	http.ListenAndServe(":8080", session.Manager.LoadAndSave(mux))
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
