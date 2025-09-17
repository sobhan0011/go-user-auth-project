package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	useruc "dekamond/internal/usecase/user"
)

type UserHandler struct {
	uuc *useruc.UserUsecase
}

func NewUserHandler(uuc *useruc.UserUsecase) *UserHandler {
	return &UserHandler{uuc: uuc}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	u, err := h.uuc.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case useruc.ErrInvalidID:
			WriteJSON(w, http.StatusBadRequest, ApiResponse{Error: "invalid_id"})
		case useruc.ErrNotFound:
			WriteJSON(w, http.StatusNotFound, ApiResponse{Error: "not_found"})
		default:
			WriteJSON(w, http.StatusInternalServerError, ApiResponse{Error: "server_error"})
		}
		return
	}
	WriteJSON(w, http.StatusOK, ApiResponse{Data: u})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Printf("server listening on :%s", r.URL.Query().Get("phone"))
	q := useruc.ListQuery{
		Phone: r.URL.Query().Get("phone"),
		Page:  atoiDefault(r.URL.Query().Get("page"), 0),
		Limit: atoiDefault(r.URL.Query().Get("limit"), 0),
	}

	page, err := h.uuc.List(r.Context(), q)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiResponse{Error: "list_failed"})
		return
	}
	WriteJSON(w, http.StatusOK, ApiResponse{Data: page})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}