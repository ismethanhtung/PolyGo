package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/models"
	"github.com/polygo/internal/polymarket"
	"github.com/polygo/pkg/response"
)

// EventsHandler handles event-related endpoints
type EventsHandler struct {
	gamma *polymarket.GammaClient
}

// NewEventsHandler creates a new events handler
func NewEventsHandler(gamma *polymarket.GammaClient) *EventsHandler {
	return &EventsHandler{gamma: gamma}
}

// GetEvents godoc
// @Summary List all events
// @Description Get a list of events with optional filtering
// @Tags Events
// @Accept json
// @Produce json
// @Param limit query int false "Limit results" default(100)
// @Param cursor query string false "Pagination cursor"
// @Param active query bool false "Filter by active status"
// @Param closed query bool false "Filter by closed status"
// @Param archived query bool false "Filter by archived status"
// @Param slug query string false "Filter by slug"
// @Param tag query string false "Filter by tag"
// @Success 200 {object} response.Response{data=[]models.Event}
// @Failure 500 {object} response.Response
// @Router /api/v1/events [get]
func (h *EventsHandler) GetEvents(c *fiber.Ctx) error {
	params := &models.EventQueryParams{
		Limit:  c.QueryInt("limit", 100),
		Cursor: c.Query("cursor"),
		Slug:   c.Query("slug"),
		Tag:    c.Query("tag"),
	}
	
	// Handle bool pointers
	if c.Query("active") != "" {
		active := c.QueryBool("active")
		params.Active = &active
	}
	if c.Query("closed") != "" {
		closed := c.QueryBool("closed")
		params.Closed = &closed
	}
	if c.Query("archived") != "" {
		archived := c.QueryBool("archived")
		params.Archived = &archived
	}
	
	data, cacheHit, err := h.gamma.GetEvents(params)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetEvent godoc
// @Summary Get event by ID
// @Description Get detailed information about a specific event including its markets
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} response.Response{data=models.Event}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/events/{id} [get]
func (h *EventsHandler) GetEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Event ID is required")
	}
	
	data, cacheHit, err := h.gamma.GetEvent(id)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	if len(data) == 0 || string(data) == "null" {
		return response.NotFound(c, "Event not found")
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetEventBySlug godoc
// @Summary Get event by slug
// @Description Get event by its URL slug
// @Tags Events
// @Accept json
// @Produce json
// @Param slug path string true "Event slug"
// @Success 200 {object} response.Response{data=models.Event}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/events/slug/{slug} [get]
func (h *EventsHandler) GetEventBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.BadRequest(c, "Slug is required")
	}
	
	data, cacheHit, err := h.gamma.GetEventBySlug(slug)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// SearchEvents godoc
// @Summary Search events
// @Description Search events by query string
// @Tags Events
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Limit results" default(20)
// @Success 200 {object} response.Response{data=[]models.Event}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/events/search [get]
func (h *EventsHandler) SearchEvents(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.BadRequest(c, "Search query is required")
	}
	
	limit := c.QueryInt("limit", 20)
	
	data, cacheHit, err := h.gamma.SearchEvents(query, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}
