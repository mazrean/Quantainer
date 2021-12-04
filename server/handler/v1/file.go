package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mazrean/Quantainer/domain/values"
	Openapi "github.com/mazrean/Quantainer/handler/v1/openapi"
	"github.com/mazrean/Quantainer/service"
)

type File struct {
	checker     *Checker
	fileService service.File
}

func NewFile(checker *Checker, fileService service.File) *File {
	return &File{
		checker:     checker,
		fileService: fileService,
	}
}

func (f *File) PostFile(c echo.Context) error {
	err := f.checker.check(c)
	if err != nil {
		return err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	reqFile, err := fileHeader.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to open file:%w", err))
	}
	defer reqFile.Close()

	file, err := f.fileService.Upload(c.Request().Context(), reqFile)
	if err != nil {
		return err
	}

	var fileType Openapi.FileType
	switch file.GetType() {
	case values.FileTypeJpeg:
		fileType = Openapi.FileTypeJpeg
	case values.FileTypePng:
		fileType = Openapi.FileTypePng
	case values.FileTypeWebP:
		fileType = Openapi.FileTypeWebp
	case values.FileTypeSvg:
		fileType = Openapi.FileTypeSvg
	case values.FileTypeGif:
		fileType = Openapi.FileTypeGif
	case values.FileTypeOther:
		fileType = Openapi.FileTypeOther
	default:
		log.Printf("error: unknown file type: %d", file.GetType())
		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected file type")
	}

	return c.JSON(http.StatusOK, Openapi.File{
		Id:        uuid.UUID(file.GetID()).String(),
		Type:      fileType,
		CreatedAt: file.GetCreatedAt(),
	})
}
