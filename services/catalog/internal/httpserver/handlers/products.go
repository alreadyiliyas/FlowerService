package handlers

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/dto"
	"github.com/ilyas/flower/services/catalog/internal/httpserver/middleware"
	usecase "github.com/ilyas/flower/services/catalog/internal/usecase/products"
	"github.com/ilyas/flower/services/catalog/internal/utils"
)

const productMaxUploadSize int64 = 20 << 20

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
	r.Body = http.MaxBytesReader(w, r.Body, productMaxUploadSize)
	if err := r.ParseMultipartForm(productMaxUploadSize); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid multipart form")
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		utils.Send(w, http.StatusBadRequest, nil, "payload is required")
		return
	}

	typeUserID, idOK := middleware.UserIDFromContext(r.Context())
	typeRole, roleOK := middleware.RoleFromContext(r.Context())
	if !idOK || !roleOK {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	var product dto.ProductPayload
	if err := json.Unmarshal([]byte(payload), &product); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid payload")
		return
	}

	mainImage, mainHeader, err := r.FormFile("main_image")
	if err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "main_image is required")
		return
	}

	images, headers, err := utils.OpenMultipartFiles(r.MultipartForm, "images")
	if err != nil {
		_ = mainImage.Close()
		utils.Send(w, http.StatusBadRequest, nil, "invalid images")
		return
	}

	req := dto.CreateProductRequest{
		Product:         product,
		TypeUserID:      typeUserID,
		TypeRole:        typeRole,
		MainImage:       mainImage,
		MainImageHeader: mainHeader,
		Images:          images,
		ImageHeaders:    headers,
	}

	resp, err := h.usecase.CreateProduct(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrForbidden):
			utils.Send(w, http.StatusForbidden, nil, err.Error())
		case errors.Is(err, apperrors.ErrConflict):
			utils.Send(w, http.StatusConflict, nil, err.Error())
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrNotFoundCategoryName):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
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

	r.Body = http.MaxBytesReader(w, r.Body, productMaxUploadSize)
	if err := r.ParseMultipartForm(productMaxUploadSize); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid multipart form")
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		utils.Send(w, http.StatusBadRequest, nil, "payload is required")
		return
	}

	typeUserID, idOK := middleware.UserIDFromContext(r.Context())
	typeRole, roleOK := middleware.RoleFromContext(r.Context())
	if !idOK || !roleOK {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	var product dto.ProductPayload
	if err := json.Unmarshal([]byte(payload), &product); err != nil {
		utils.Send(w, http.StatusBadRequest, nil, "invalid payload")
		return
	}

	var mainImage multipart.File
	var mainHeader *multipart.FileHeader
	mainImage, mainHeader, err = r.FormFile("main_image")
	if err != nil {
		mainImage = nil
		mainHeader = nil
	}

	images, headers, err := utils.OpenMultipartFiles(r.MultipartForm, "images")
	if err != nil {
		if mainImage != nil {
			_ = mainImage.Close()
		}
		utils.Send(w, http.StatusBadRequest, nil, "invalid images")
		return
	}

	req := dto.UpdateProductRequest{
		Product:         product,
		TypeUserID:      typeUserID,
		TypeRole:        typeRole,
		MainImage:       mainImage,
		MainImageHeader: mainHeader,
		Images:          images,
		ImageHeaders:    headers,
	}

	resp, err := h.usecase.UpdateProduct(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidInput):
			utils.Send(w, http.StatusBadRequest, nil, err.Error())
		case errors.Is(err, apperrors.ErrForbidden):
			utils.Send(w, http.StatusForbidden, nil, err.Error())
		case errors.Is(err, apperrors.ErrConflict):
			utils.Send(w, http.StatusConflict, nil, err.Error())
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

	typeUserID, idOK := middleware.UserIDFromContext(r.Context())
	typeRole, roleOK := middleware.RoleFromContext(r.Context())
	if !idOK || !roleOK {
		utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
		return
	}

	req := dto.DeleteProductRequest{
		TypeUserID: typeUserID,
		TypeRole:   typeRole,
	}

	if err := h.usecase.DeleteProduct(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			utils.Send(w, http.StatusNotFound, nil, err.Error())
		case errors.Is(err, apperrors.ErrForbidden):
			utils.Send(w, http.StatusForbidden, nil, err.Error())
		default:
			utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
		}
		return
	}
	utils.Send(w, http.StatusOK, nil, "deleted")
}
