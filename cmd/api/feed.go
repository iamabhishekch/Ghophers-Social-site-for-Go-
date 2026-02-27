package main

import (
	"net/http"

	"github.com/iamabhishekch/Social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// pagination, filter, sort

	// giving default values to pagination properperties
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	// parsing values send by user
	fq, err := fq.Parse(r)
	if err != nil{
		app.badRequestResponse(w,r,err)
		return
	}

	// validating recevied values from user
	if err := Validate.Struct(fq); err != nil{
		app.badRequestResponse(w,r,err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(23), fq)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
