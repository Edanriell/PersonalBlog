package newsletter

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

func (h *Handler) GetAll(c echo.Context) error {
	var req SubscriberListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	subscribers, pagination, err := h.service.GetAll(req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to fetch subscribers", err)
	}

	return response.Paginated(c, subscribers, pagination, "Subscribers retrieved successfully")
}

func (h *Handler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid subscriber ID", err)
	}

	subscriber, err := h.service.GetByID(uint(id))
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Subscriber not found", err)
	}

	return response.Success(c, subscriber, "Subscriber retrieved successfully")
}

func (h *Handler) Subscribe(c echo.Context) error {
	var req SubscribeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	subscriber, err := h.service.Subscribe(req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to subscribe", err)
	}

	return response.Success(c, subscriber, "Subscription successful")
}

func (h *Handler) Unsubscribe(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid subscriber ID", err)
	}

	if err := h.service.Unsubscribe(uint(id)); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to unsubscribe", err)
	}

	return response.Success(c, nil, "Unsubscribed successfully")
}

func (h *Handler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid subscriber ID", err)
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to delete subscriber", err)
	}

	return response.Success(c, nil, "Subscriber deleted successfully")
}
