package v1

import (
	"bytes"
	"errors"
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
	session     *Session
	checker     *Checker
	fileService service.File
}

func NewFile(session *Session, checker *Checker, fileService service.File) *File {
	return &File{
		session:     session,
		checker:     checker,
		fileService: fileService,
	}
}

func (f *File) PostFile(c echo.Context) error {
	err := f.checker.check(c)
	if err != nil {
		return err
	}

	session, err := getSession(c)
	if err != nil {
		log.Printf("error: failed to get session: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get session")
	}

	authSession, err := f.session.getAuthSession(session)
	if err != nil {
		log.Printf("error: failed to get user: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no file")
	}

	reqFile, err := fileHeader.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to open file:%w", err))
	}
	defer reqFile.Close()

	fileInfo, err := f.fileService.Upload(c.Request().Context(), authSession, reqFile)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to upload file:%w", err))
	}

	var fileType Openapi.FileType
	switch fileInfo.File.GetType() {
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
		log.Printf("error: unknown file type: %d", fileInfo.File.GetType())
		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected file type")
	}

	return c.JSON(http.StatusOK, Openapi.File{
		Id:        uuid.UUID(fileInfo.File.GetID()).String(),
		Type:      fileType,
		Creator:   string(fileInfo.Creator.GetName()),
		CreatedAt: fileInfo.File.GetCreatedAt(),
	})
}

func (f *File) GetFile(c echo.Context, strFileID Openapi.FileIDInPath) error {
	err := f.checker.check(c)
	if err != nil {
		return err
	}

	uuidFileID, err := uuid.Parse(string(strFileID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file id")
	}

	buf := bytes.NewBuffer(nil)
	file, err := f.fileService.Download(c.Request().Context(), values.NewFileIDFromUUID(uuidFileID), buf)
	if errors.Is(err, service.ErrNoFile) {
		return echo.NewHTTPError(http.StatusNotFound, "file not found")
	}
	if err != nil {
		log.Printf("error: failed to get file: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to download file")
	}

	var mime string
	switch file.GetType() {
	case values.FileTypeJpeg:
		mime = "image/jpeg"
	case values.FileTypePng:
		mime = "image/png"
	case values.FileTypeWebP:
		mime = "image/webp"
	case values.FileTypeSvg:
		mime = "image/svg+xml"
	case values.FileTypeGif:
		mime = "image/gif"
	case values.FileTypeOther:
		mime = "application/octet-stream"
	default:
		log.Printf("error: unknown file type: %d", file.GetType())
		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected file type")
	}

	return c.Stream(http.StatusOK, mime, buf)
}
