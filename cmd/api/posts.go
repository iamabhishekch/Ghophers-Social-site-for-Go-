package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iamabhishekch/Social/internal/store"
)

// not relying on the string relaying on the type for
// for : postContextMiddleware() ::: ctx = context.WithValue(ctx, "post", post)
type postKey string

const postCtx postKey = "post"

type CreatePostPaylaod struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

// createPostHandler will create new post

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	// we will receive data of type CreatePostPaylaod from user (r.Body)
	// which we will decode in payload variable
	// to push data in db we need store.Post (to fill all details)
	var payload CreatePostPaylaod
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// vallidate struct
	if err := Validate.Struct(payload); err != nil {
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
		app.internalServerError(w, r, err)
		return
	}
	// sending back to user

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// getPostHandler will return post by ID
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	// fetching comments of the post
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comment = comments

	// returning to user
	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// deletePostHandler will delete the post by id

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// parse the Postid to be deleted from url

	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// delete the post
	ctx := r.Context()
	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// payload for the UpdatePostHandler
type UpdatePostPayLoad struct {
	Title   *string `json:"title" validate:"omitempty,max=1000"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePostHandler will update the post
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayLoad
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// updating values of post with payload values
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err!= nil{
		app.internalServerError(w, r, err)
		return
	} 

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// reducing repeatative code block
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parsing id
		idParam := chi.URLParam(r, "postID")
		// converting string val to int64
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		//db layer
		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			// switching based on the error if not found 400, server error 500, or others
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
