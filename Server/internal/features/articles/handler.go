package articles

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
	var req ArticleListRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	articles, pagination, err := h.service.GetAll(req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to fetch articles", err)
	}

	return response.Paginated(c, articles, pagination, "Articles retrieved successfully")
}

func (h *Handler) GetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid article ID", err)
	}

	article, err := h.service.GetByID(uint(id))
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Article not found", err)
	}

	return response.Success(c, article, "Article retrieved successfully")
}

func (h *Handler) GetBySlug(c echo.Context) error {
	slug := c.Param("slug")

	article, err := h.service.GetBySlug(slug)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Article not found", err)
	}

	return response.Success(c, article, "Article retrieved successfully")
}

func (h *Handler) Create(c echo.Context) error {
	authorID := c.Get("user_id").(uint)

	var req CreateArticleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	article, err := h.service.Create(authorID, req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to create article", err)
	}

	return response.Success(c, article, "Article created successfully")
}

func (h *Handler) Update(c echo.Context) error {
	authorID := c.Get("user_id").(uint)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid article ID", err)
	}

	var req UpdateArticleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body", err)
	}

	article, err := h.service.Update(uint(id), authorID, req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to update article", err)
	}

	return response.Success(c, article, "Article updated successfully")
}

func (h *Handler) Delete(c echo.Context) error {
	authorID := c.Get("user_id").(uint)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid article ID", err)
	}

	if err := h.service.Delete(uint(id), authorID); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to delete article", err)
	}

	return response.Success(c, nil, "Article deleted successfully")
}
