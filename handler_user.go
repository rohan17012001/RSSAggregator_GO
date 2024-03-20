package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	// "github.com/rohan17012001/rssaggregator/internal/auth"
	"github.com/rohan17012001/rssaggregator/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v ", err))
		return
	}
	respondWithJSON(w, 201, databaseUsertoUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUsertoUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts for user: %v ", err))
		return
	}
	// fmt.Println(posts)
	respondWithJSON(w, 200, databasePoststoPosts(posts))
}
