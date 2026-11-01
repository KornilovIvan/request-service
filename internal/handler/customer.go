package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"request-service/internal/model"
	"request-service/internal/service"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Register(r gin.IRouter) {
	r.POST("/api/v1/customers", h.Create)
	r.GET("/api/v1/customers", h.List)
	r.GET("/api/v1/customers/:id", h.Get)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	customer, err := h.service.Get(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		respondError(c, http.StatusNotFound, "customer not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "name is required")
		return
	}

	customer, err := h.service.Create(c.Request.Context(), body.Name)
	if errors.Is(err, service.ErrNameRequired) {
		respondError(c, http.StatusBadRequest, "name is required")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) List(c *gin.Context) {
	customers, err := h.service.List(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	if customers == nil {
		customers = []model.Customer{}
	}
	c.JSON(http.StatusOK, customers)
}
