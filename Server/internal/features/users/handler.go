package users

import (
	"net/http"
	"strconv"

	"Server/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID", err)
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		return response.Error(c, http.StatusNotFound, "User not found", err)
	}

	return response.Success(c, user, "User retrieved successfully")
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	loginResp, err := h.service.Register(req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Registration failed", err)
	}

	return response.Success(c, loginResp, "User registered successfully")
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	loginResp, err := h.service.Login(req)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Login failed", err)
	}

	return response.Success(c, loginResp, "Login successful")
}

func (h *Handler) Update(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	user, err := h.service.Update(userID, req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Update failed", err)
	}

	return response.Success(c, user, "User updated successfully")
}
