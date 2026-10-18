package handler

import (
	"net/http"
	"strings"

	"H18-P3-GC01-V2-Fadil/repository"

	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	repo *repository.TransactionRepository
}

func NewTransactionHandler(repo *repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

type transactionRequest struct {
	Name string `json:"name"`
}

func bindTransaction(c echo.Context) (*transactionRequest, error) {
	var req transactionRequest
	if err := c.Bind(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	return &req, nil
}

func success(c echo.Context, code int, data any) error {
	return c.JSON(code, echo.Map{"code": code, "status": "success", "data": data})
}

// Create - POST /transactions
func (h *TransactionHandler) Create(c echo.Context) error {
	req, err := bindTransaction(c)
	if err != nil {
		return err
	}
	t := &repository.Transaction{Name: req.Name}
	if err := h.repo.Create(c.Request().Context(), t); err != nil {
		return err
	}
	return success(c, http.StatusCreated, t)
}

// GetAll - GET /transactions
func (h *TransactionHandler) GetAll(c echo.Context) error {
	transactions, err := h.repo.FindAll(c.Request().Context())
	if err != nil {
		return err
	}
	return success(c, http.StatusOK, transactions)
}

// GetByID - GET /transactions/:id
func (h *TransactionHandler) GetByID(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	t, err := h.repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return mapRepoError(err, "transaction")
	}
	return success(c, http.StatusOK, t)
}

// Update - PUT /transactions/:id
func (h *TransactionHandler) Update(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	req, err := bindTransaction(c)
	if err != nil {
		return err
	}
	t := &repository.Transaction{ID: id, Name: req.Name}
	if err := h.repo.Update(c.Request().Context(), t); err != nil {
		return mapRepoError(err, "transaction")
	}
	return success(c, http.StatusOK, t)
}

// Delete - DELETE /transactions/:id
func (h *TransactionHandler) Delete(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return mapRepoError(err, "transaction")
	}
	return success(c, http.StatusOK, echo.Map{"message": "transaction deleted successfully"})
}
