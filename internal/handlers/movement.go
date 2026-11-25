package handlers

import (
	"net/http"
	"strconv"

	"inventrack/internal/models"
	"inventrack/internal/repository"
)

type MovementHandler struct {
	repo *repository.MovementRepo
}

func NewMovementHandler(repo *repository.MovementRepo) *MovementHandler {
	return &MovementHandler{repo: repo}
}

func (h *MovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.MovementInput
	if err := decode(w, r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := in.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	m := &models.Movement{ProductID: in.ProductID, Type: in.Type, Quantity: in.Quantity}
	if err := h.repo.Create(r.Context(), m); err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

// List soporta el filtro opcional ?product_id=1
func (h *MovementHandler) List(w http.ResponseWriter, r *http.Request) {
	var productID int64
	if raw := r.URL.Query().Get("product_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id <= 0 {
			writeError(w, http.StatusBadRequest, "invalid product_id")
			return
		}
		productID = id
	}

	movements, err := h.repo.List(r.Context(), productID)
	if err != nil {
		handleRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movements)
}
