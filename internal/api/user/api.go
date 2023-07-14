package user

import (
	"context"
	"encoding/json"
	"net/http"

	userentity "tech-tsarka/internal/storage/user/entity"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type User interface {
	Create(ctx context.Context, arg userentity.UserCreateInput) (userentity.User, error)
	Get(ctx context.Context, id string) (userentity.User, error)
	Update(ctx context.Context, id string, arg userentity.UserUpdateInput) error
	Delete(ctx context.Context, id string) error
}

func New(user User) *implementation {
	return &implementation{
		user: user,
	}
}

type implementation struct {
	user User
}

func (h *implementation) Setup(router *httprouter.Router) {
	router.Handler(http.MethodPost, "/rest/user", h.Create())
	router.Handler(http.MethodGet, "/rest/user/:id", h.Get())
	router.Handler(http.MethodPut, "/rest/user/:id", h.Update())
	router.Handler(http.MethodDelete, "/rest/user/:id", h.Delete())
}

type userCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type userCreateResponse struct {
	ID string `json:"id"`
}

func (i *implementation) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req userCreateRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		usr, err := i.user.Create(r.Context(), userentity.UserCreateInput{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		})

		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		respondJSON(w, http.StatusCreated, userCreateResponse{ID: usr.ID})
	})
}

type userGetResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (i *implementation) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")

		if _, err := uuid.Parse(id); err != nil {
			respondError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		usr, err := i.user.Get(r.Context(), id)

		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		respondJSON(w, http.StatusOK, userGetResponse{FirstName: usr.FirstName, LastName: usr.LastName})
	})
}

type userUpdateRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func (i *implementation) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")

		if _, err := uuid.Parse(id); err != nil {
			respondError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		var req userUpdateRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		err := i.user.Update(r.Context(), id, userentity.UserUpdateInput{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		})

		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"success": "ok",
		})
	})
}

func (i *implementation) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")

		if _, err := uuid.Parse(id); err != nil {
			respondError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		err := i.user.Delete(r.Context(), id)

		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"success": "ok",
		})
	})
}

func decodeJSONBody(_ http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
