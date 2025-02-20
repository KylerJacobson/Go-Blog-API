package main

import (
	"github.com/KylerJacobson/Go-Blog-API/internal/services/azure"
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
	zapLogger, err := logger.NewLogger(env)
	if err != nil {
		log.Fatal(err)
	}
	defer zapLogger.Sync()
	dbPool := config.GetDBConn(zapLogger)
	defer dbPool.Close()

	azureClient := azure.NewAzureClient(zapLogger)

	session.Init()

	mux := http.NewServeMux()
	usersApi := users.New(usersRepo.New(dbPool, zapLogger), zapLogger)
	postsApi := posts.New(postsRepo.New(dbPool, zapLogger), zapLogger)

	sessionApi := session.New(usersRepo.New(dbPool, zapLogger), zapLogger)
	mediaApi := media.New(mediaRepo.New(dbPool, zapLogger), zapLogger, azureClient)

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
	mux.HandleFunc("PUT /api/user/{id}", usersApi.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", usersApi.DeleteUserById)

	// ---------------------------- Session ----------------------------

	mux.HandleFunc("POST /api/session", sessionApi.CreateSession)
	mux.HandleFunc("POST /api/verifyToken", authorization.VerifyToken)
	mux.HandleFunc("DELETE /api/session", sessionApi.DeleteSession)

	// ---------------------------- Media ----------------------------
	mux.HandleFunc("POST /api/media", mediaApi.UploadMedia)
	mux.HandleFunc("GET /api/media/{id}", mediaApi.GetMediaByPostId)

	zapLogger.Sugar().Infof("Logging level set to %s", env)
	zapLogger.Sugar().Infof("listening on port: %d", 8080)
	http.ListenAndServe(":8080", session.Manager.LoadAndSave(mux))
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
