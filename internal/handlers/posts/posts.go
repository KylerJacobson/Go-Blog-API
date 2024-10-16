package posts

import (
	"encoding/json"
	"errors"
	"fmt"
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
	b, err := json.Marshal(posts)
	if err != nil {
		logger.Sugar.Errorf("error unmarshalling recent posts : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
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
	b, err := json.Marshal(posts)
	if err != nil {
		logger.Sugar.Errorf("error unmarshalling recent public posts : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
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
	b, err := json.Marshal(post)
	if err != nil {
		logger.Sugar.Errorf("error unmarshalling post (%d) : %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
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
	err := validatePost(post)
	if err != nil {
		logger.Sugar.Errorf("the post was not formatter correctly: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&post)
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

func (postsApi *postsApi) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post post_models.PostRequestBody
	id := r.PathValue("id")
	iId, err := strconv.Atoi(id)
	if err != nil {
		logger.Sugar.Errorf("UpdatePost parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logger.Sugar.Errorf("Error decoding the post request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	err = validatePost(post)
	if err != nil {
		logger.Sugar.Errorf("the post was not formatter correctly: %v", err)

		b, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}
	updatedPost, err := postsApi.postsRepository.UpdatePost(post, iId)
	if err != nil {
		logger.Sugar.Errorf("error updating post (%s) : %v", post.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(updatedPost)
	if err != nil {
		logger.Sugar.Errorf("error unmarshalling updated post (%s) : %v", post.Title, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func validatePost(post post_models.PostRequestBody) error {
	if len(post.Title) < 1 {
		return fmt.Errorf("post title must not be empty")
	}
	if len(post.Content) < 1 {
		return fmt.Errorf("post content must not be empty")
	}
	return nil
}
