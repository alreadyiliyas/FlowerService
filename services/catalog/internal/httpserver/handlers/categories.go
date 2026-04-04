package handlers

import (
	"errors"
	"net/http"
	"strings"

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
		switch {
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrNotFoundCategoryName), errors.Is(err, apperrors.ErrNotFoundCategorySlug):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
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
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrNotFoundCategoryName), errors.Is(err, apperrors.ErrNotFoundCategorySlug):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusOK, resp, "ok")
}

func (h *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, utils.CategoryMaxUploadSize)

	if err := r.ParseMultipartForm(utils.CategoryMaxUploadSize); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "category image is required")
		return
	}

	req := dto.CreateCategoryRequest{
		Name:        r.FormValue("name"),
		Slug:        r.FormValue("slug"),
		Description: r.FormValue("description"),
		Image:       file,
		ImageHeader: header,
	}

	resp, err := h.usecase.CreateCategory(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrDuplicateCategoryName), errors.Is(err, apperrors.ErrDuplicateCategorySlug):
			utils.Send(w, http.StatusConflict, nil, err.Error())
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

	r.Body = http.MaxBytesReader(w, r.Body, utils.CategoryMaxUploadSize)
	if err := r.ParseMultipartForm(utils.CategoryMaxUploadSize); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid multipart form")
		return
	}

	req := dto.UpdateCategoryRequest{
		Name:        r.FormValue("name"),
		Slug:        r.FormValue("slug"),
		Description: r.FormValue("description"),
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		req.Image = file
		req.ImageHeader = header
	} else if !strings.Contains(err.Error(), "no such file") {
		utils.Send(w, http.StatusBadRequest, nil, "invalid image")
		return
	}

	resp, err := h.usecase.UpdateCategory(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrNotFoundCategoryName), errors.Is(err, apperrors.ErrNotFoundCategorySlug):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		case errors.Is(err, apperrors.ErrDuplicateCategoryName), errors.Is(err, apperrors.ErrDuplicateCategorySlug):
			utils.Send(w, http.StatusConflict, nil, err.Error())
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
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrNotFoundCategoryName), errors.Is(err, apperrors.ErrNotFoundCategorySlug):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusOK, nil, "deleted")
}
