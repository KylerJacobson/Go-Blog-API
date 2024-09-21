package posts

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	post_models "github.com/KylerJacobson/Go-Blog-API/internal/api/types/posts"
	posts_repo "github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	v5 "github.com/jackc/pgx/v5"
)

type PostsApi interface {
	GetRecentPosts(w http.ResponseWriter, r *http.Request)
	GetRecentPublicPosts(w http.ResponseWriter, r *http.Request)
	GetPostById(w http.ResponseWriter, r *http.Request)
	DeletePostById(w http.ResponseWriter, r *http.Request)
	CreatePost(w http.ResponseWriter, r *http.Request)
}

type postsApi struct {
	postsRepository posts_repo.PostsRepository
}

func New(postsRepo posts_repo.PostsRepository) *postsApi {
	return &postsApi{
		postsRepository: postsRepo,
	}
}

func (postsApi *postsApi) GetRecentPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postsApi.postsRepository.GetRecentPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func (postsApi *postsApi) GetRecentPublicPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postsApi.postsRepository.GetRecentPublicPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (postsApi *postsApi) GetPostById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	val, err := strconv.Atoi(id)
	if err != nil {
		logger.Sugar.Errorf("GetPostId parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	post, err := postsApi.postsRepository.GetPostById(val)
	if err != nil {
		if errors.Is(err, v5.ErrNoRows) {
			logger.Sugar.Infof("Post %v does not exist in the database", val)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (postsApi *postsApi) DeletePostById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	val, err := strconv.Atoi(id)
	if err != nil {
		logger.Sugar.Errorf("DeletePostById parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	err = postsApi.postsRepository.DeletePostById(val)
	if err != nil {
		if errors.Is(err, v5.ErrNoRows) {
			logger.Sugar.Infof("Post %v does not exist in the database", val)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (postsApi *postsApi) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post post_models.PostRequestBody
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logger.Sugar.Errorf("Error decoding the post request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	err = postsApi.postsRepository.CreatePost(post)
	if err != nil {
		logger.Sugar.Errorf("error creating post (%s) : %v", post.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
