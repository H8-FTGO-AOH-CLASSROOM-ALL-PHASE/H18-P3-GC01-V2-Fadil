package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"H18-P3-GC01-V2-Fadil/repository"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

type productRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

func bindProduct(c echo.Context) (*productRequest, error) {
	var req productRequest
	if err := c.Bind(&req); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if strings.TrimSpace(req.Name) == "" || req.Price < 0 || req.Qty < 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "name is required; price and qty must not be negative")
	}
	return &req, nil
}

func parseID(c echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	return id, nil
}

func mapRepoError(err error, what string) error {
	if errors.Is(err, repository.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, what+" not found")
	}
	return err
}

// Create - POST /products
func (h *ProductHandler) Create(c echo.Context) error {
	req, err := bindProduct(c)
	if err != nil {
		return err
	}
	p := &repository.Product{Name: req.Name, Price: req.Price, Qty: req.Qty}
	if err := h.repo.Create(c.Request().Context(), p); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, p)
}

// GetAll - GET /products
func (h *ProductHandler) GetAll(c echo.Context) error {
	products, err := h.repo.FindAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, products)
}

// GetByID - GET /products/:id
func (h *ProductHandler) GetByID(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	p, err := h.repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return mapRepoError(err, "product")
	}
	return c.JSON(http.StatusOK, p)
}

// Update - PUT /products/:id
func (h *ProductHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	req, err := bindProduct(c)
	if err != nil {
		return err
	}
	p := &repository.Product{ID: id, Name: req.Name, Price: req.Price, Qty: req.Qty}
	if err := h.repo.Update(c.Request().Context(), p); err != nil {
		return mapRepoError(err, "product")
	}
	return c.JSON(http.StatusOK, p)
}

// Delete - DELETE /products/:id
func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return mapRepoError(err, "product")
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "product deleted successfully"})
}
