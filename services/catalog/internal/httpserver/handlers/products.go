package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	usecase "github.com/ilyas/flower/services/catalog/internal/usecase/products"
	"github.com/ilyas/flower/services/catalog/internal/utils"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductsHandler(pu usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: pu}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	filter := utils.ParseProductFilter(r)
	resp, err := h.usecase.ListProducts(r.Context(), filter)
	if err != nil {
		utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		return
	}
	utils.Send(w, http.StatusOK, resp, "ok")
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint64Path(r, "id")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	resp, err := h.usecase.GetProduct(r.Context(), id)
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

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.CreateProduct(r.Context(), req)
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

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint64Path(r, "id")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var req dto.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resp, err := h.usecase.UpdateProduct(r.Context(), id, req)
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

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUint64Path(r, "id")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	if err := h.usecase.DeleteProduct(r.Context(), id); err != nil {
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
