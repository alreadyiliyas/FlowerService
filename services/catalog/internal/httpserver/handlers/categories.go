package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	usecase "github.com/ilyas/flower/services/catalog/internal/usecase/categories"
	"github.com/ilyas/flower/services/catalog/internal/utils"
)

type CategoriesHandler struct {
	usecase usecase.UsecaseCategories
}

func NewCategoriesHandler(uc usecase.UsecaseCategories) *CategoriesHandler {
	return &CategoriesHandler{usecase: uc}
}

func (h *CategoriesHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.ListCategories(r.Context())
	if err != nil {
		utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		return
	}
	utils.Send(w, http.StatusOK, resp, "ok")
}

func (h *CategoriesHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint64Path(r, "id")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	resp, err := h.usecase.GetCategory(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusOK, resp, "ok")
}

func (h *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.CreateCategory(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusCreated, resp, "created")
}

func (h *CategoriesHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint64Path(r, "id")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var req dto.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.UpdateCategory(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrNotFound):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusOK, resp, "updated")
}

func (h *CategoriesHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint64Path(r, "id")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	if err := h.usecase.DeleteCategory(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusOK, nil, "deleted")
}
