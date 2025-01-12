package main

import (
	"context"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type FollowUserToggle struct {
	UserID int64 `json:"user_id"`
}

type userKey string

const userCtx userKey = "user"

// @SUMMARY		Fetches a user profile
// @Description	Fetches a user profile by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	models.UsersModel
// @Failure		400	{object}	error
// @Failure		404	{object}	error
// @Failure		500	{object}	error
// @Security		ApiKeyAuth
// @Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// @SUMMARY		Follows a user
// @Description	Fetches a user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path		int		true	"User ID"
// @Success		204		{object}	string	"User followed"
// @Failure		400		{object}	string	"User payload failed"
// @Failure		404		{object}	error	"User not found"
// @Security		ApiKeyAuth
// @Router			/users/{id}/follow [put]
func (app *application) toggleFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)
	var payload FollowUserToggle

	if err := readJson(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.ToggleUserFollow(r.Context(), followerUser.ID, payload.UserID); err != nil {
		app.handleError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, followerUser); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)

		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.RetrieveById(ctx, userId)

		if err != nil {
			app.handleError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *models.UsersModel {
	user, _ := r.Context().Value(userCtx).(*models.UsersModel)
	return user
}
