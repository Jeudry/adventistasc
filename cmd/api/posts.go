package main

import (
	"context"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

type CreatePostCommentPayload struct {
	Comment string `json:"comment"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := &models.PostsModel{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  1,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.RetrieveCommentsByPostId(r.Context(), post.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createPostCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	idAsInt, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	var payload CreatePostCommentPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	comment := &models.CommentsModel{
		Content: payload.Comment,
		UserID:  1,
		PostID:  idAsInt,
	}

	err = app.store.Comments.CreatePostComment(ctx, comment)

	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload UpdatePostPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	postToBeUpdated := getPostFromCtx(r)

	if payload.Title != nil {
		postToBeUpdated.Title = *payload.Title
	}

	if payload.Content != nil {
		postToBeUpdated.Content = *payload.Content
	}

	if err := app.store.Posts.Update(ctx, postToBeUpdated); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, postToBeUpdated); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	idAsInt, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, idAsInt); err != nil {
		app.handleError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postId")
		idAsInt, err := strconv.ParseInt(idParam, 10, 64)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.RetrieveById(ctx, idAsInt)

		if err != nil {
			app.handleError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *models.PostsModel {
	post, _ := r.Context().Value(postCtx).(*models.PostsModel)
	return post
}
