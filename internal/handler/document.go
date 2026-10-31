package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"request-service/internal/model"
	"request-service/internal/service"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) Register(r gin.IRouter) {
	r.POST("/api/v1/requests/:id/documents", h.Create)
	r.GET("/api/v1/requests/:id/documents", h.List)
	r.GET("/api/v1/documents/:id", h.Get)
}

func (h *DocumentHandler) Create(c *gin.Context) {
	requestID, ok := pathID(c)
	if !ok {
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "document name is required")
		return
	}

	document, err := h.service.Create(c.Request.Context(), requestID, body.Name)
	if respondDocumentError(c, err, "request not found") {
		return
	}
	c.JSON(http.StatusCreated, document)
}

func (h *DocumentHandler) List(c *gin.Context) {
	requestID, ok := pathID(c)
	if !ok {
		return
	}

	documents, err := h.service.List(c.Request.Context(), requestID)
	if respondDocumentError(c, err, "request not found") {
		return
	}
	if documents == nil {
		documents = []model.Document{}
	}
	c.JSON(http.StatusOK, documents)
}

func (h *DocumentHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	document, err := h.service.Get(c.Request.Context(), id)
	if respondDocumentError(c, err, "document not found") {
		return
	}
	c.JSON(http.StatusOK, document)
}

func pathID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func respondDocumentError(c *gin.Context, err error, notFound string) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, service.ErrDocumentName) {
		respondError(c, http.StatusBadRequest, "document name is required")
		return true
	}
	if errors.Is(err, service.ErrNotFound) {
		respondError(c, http.StatusNotFound, notFound)
		return true
	}
	respondError(c, http.StatusInternalServerError, "internal error")
	return true
}
