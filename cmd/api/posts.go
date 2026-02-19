package main

import (
	"net/http"

	"github.com/iamabhishekch/Social/internal/store"
)

type CreatePostPaylaod struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	// we will receive data of type CreatePostPaylaod from user (r.Body)
	// which we will decode in payload variable
	// to push data in db we need store.Post (to fill all details)
	var payload CreatePostPaylaod
	if err := readJson(w, r, &payload); err != nil {
		writeJSONEroor(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: change after auth
		UserID: 1,
	}

	ctx := r.Context()

	// data push
	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONEroor(w, http.StatusInternalServerError, err.Error())
		return
	}
	// sending back to user

	if err := writeJson(w, http.StatusOK, post); err != nil {
		writeJSONEroor(w, http.StatusInternalServerError, err.Error())
		return
	}

}
