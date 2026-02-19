package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		app.badRequestResponse(w, r, err)
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
		app.internalServerError(w,r, err)
		return
	}
	// sending back to user

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w,r, err)
		return
	}

}

// getPostHandler will return post by ID
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	// parsing id
	idParam := chi.URLParam(r, "postID")
	// converting string val to int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w,r, err)
		return
	}

	// db layer
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		// switching based on the error if not found 400, server error 500, or others
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w,r, err)
		}
		return
	}
	// returning to user
	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w,r, err)
		return
	}
	

}
