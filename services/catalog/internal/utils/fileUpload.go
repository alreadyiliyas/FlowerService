package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
)

type UploadImageParams struct {
	File         multipart.File
	Header       *multipart.FileHeader
	Dir          string
	PublicPrefix string
	AllowedExt   []string
	FileNameSize int
}

type UploadedFile struct {
	FileName  string
	FullPath  string
	PublicURL string
}

type UploadImagesParams struct {
	Files        []multipart.File
	Headers      []*multipart.FileHeader
	Dir          string
	PublicPrefix string
	AllowedExt   []string
	FileNameSize int
}

func ValidateImageExtension(filename string, allowed []string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, item := range allowed {
		if ext == strings.ToLower(item) {
			return ext, nil
		}
	}

	return "", fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "такой формат картинки не поддерживается, должен быть: '.jpg', '.jpeg', '.png', '.webp'")
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func RandomFileName(size int, ext string) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf) + ext, nil
}

func SaveUploadedFile(file multipart.File, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		_ = dst.Close()
		_ = os.Remove(dstPath)
		return err
	}

	return nil
}

func UploadImage(params UploadImageParams) (*UploadedFile, error) {
	ext, err := ValidateImageExtension(params.Header.Filename, params.AllowedExt)
	if err != nil {
		return nil, err
	}

	if err := EnsureDir(params.Dir); err != nil {
		return nil, err
	}

	fileName, err := RandomFileName(params.FileNameSize, ext)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(params.Dir, fileName)
	if err := SaveUploadedFile(params.File, fullPath); err != nil {
		return nil, err
	}

	return &UploadedFile{
		FileName:  fileName,
		FullPath:  fullPath,
		PublicURL: strings.TrimRight(params.PublicPrefix, "/") + "/" + fileName,
	}, nil
}

func UploadImages(params UploadImagesParams) ([]UploadedFile, error) {
	if len(params.Files) == 0 {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "не переданы картинки")
	}
	if len(params.Files) != len(params.Headers) {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "количество файлов и заголовков картинок не совпадает")
	}

	uploaded := make([]UploadedFile, 0, len(params.Files))
	for i := range params.Files {
		file := params.Files[i]
		header := params.Headers[i]

		item, err := UploadImage(UploadImageParams{
			File:         file,
			Header:       header,
			Dir:          params.Dir,
			PublicPrefix: params.PublicPrefix,
			AllowedExt:   params.AllowedExt,
			FileNameSize: params.FileNameSize,
		})

		closeErr := file.Close()
		if err != nil {
			DeleteUploadedFiles(uploaded)
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
		if closeErr != nil {
			DeleteUploadedFiles(uploaded)
			DeleteFileIfExists(item.FullPath)
			return nil, closeErr
		}

		uploaded = append(uploaded, *item)
	}

	return uploaded, nil
}

func DeleteFileIfExists(path string) {
	if path == "" {
		return
	}
	_ = os.Remove(path)
}

func DeleteUploadedFiles(files []UploadedFile) {
	for _, file := range files {
		DeleteFileIfExists(file.FullPath)
	}
}

type ReplaceImageParams struct {
	File              multipart.File
	Header            *multipart.FileHeader
	ExistingPublicURL string
	AllowedExt        []string
}

func ReplaceImage(params ReplaceImageParams) (*UploadedFile, error) {
	if params.ExistingPublicURL == "" {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, "отсутствует текущий url картинки")
	}

	_, err := ValidateImageExtension(params.Header.Filename, params.AllowedExt)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.FromSlash(strings.TrimPrefix(params.ExistingPublicURL, "/"))
	if err := EnsureDir(filepath.Dir(fullPath)); err != nil {
		return nil, err
	}
	if err := SaveUploadedFile(params.File, fullPath); err != nil {
		return nil, err
	}

	return &UploadedFile{
		FileName:  filepath.Base(fullPath),
		FullPath:  fullPath,
		PublicURL: params.ExistingPublicURL,
	}, nil
}
