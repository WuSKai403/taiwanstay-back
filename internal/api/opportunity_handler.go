package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
	"github.com/taiwanstay/taiwanstay-back/internal/repository"
	"github.com/taiwanstay/taiwanstay-back/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OpportunityHandler struct {
	oppService  service.OpportunityService
	hostService service.HostService
}

func NewOpportunityHandler(oppService service.OpportunityService, hostService service.HostService) *OpportunityHandler {
	return &OpportunityHandler{
		oppService:  oppService,
		hostService: hostService,
	}
}

func (h *OpportunityHandler) Create(c *gin.Context) {
	var opp domain.Opportunity
	if err := c.ShouldBindJSON(&opp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get User ID -> Host ID
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	host, err := h.hostService.GetHostByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a host"})
		return
	}
	opp.HostID = host.ID

	createdOpp, err := h.oppService.CreateOpportunity(c.Request.Context(), &opp)
	if err != nil {
		fmt.Printf("Error creating opportunity: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create opportunity: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOpp)
}

func (h *OpportunityHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	opp, err := h.oppService.GetOpportunityByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}
	c.JSON(http.StatusOK, opp)
}

func (h *OpportunityHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	filter := bson.M{}

	// Filter by Status (default ACTIVE if not specified?)
	// For now, let's allow listing all public ones
	// filter["status"] = domain.OpportunityStatusActive

	if hostID := c.Query("hostId"); hostID != "" {
		objID, _ := primitive.ObjectIDFromHex(hostID)
		filter["hostId"] = objID
	}

	// Exclude deleted opportunities
	if _, ok := filter["status"]; !ok {
		filter["status"] = bson.M{"$ne": domain.OpportunityStatusDeleted}
	}

	opps, err := h.oppService.ListOpportunities(c.Request.Context(), filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list opportunities"})
		return
	}

	c.JSON(http.StatusOK, opps)
}

func (h *OpportunityHandler) Search(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	distStr := c.Query("distance")
	var lat, lng, dist float64
	if latStr != "" && lngStr != "" {
		lat, _ = strconv.ParseFloat(latStr, 64)
		lng, _ = strconv.ParseFloat(lngStr, 64)
	}
	if distStr != "" {
		dist, _ = strconv.ParseFloat(distStr, 64)
	}

	filter := repository.OpportunityFilter{
		Query:     c.Query("q"),
		Type:      c.Query("type"),
		City:      c.Query("city"),
		Country:   c.Query("country"),
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
		Lat:       lat,
		Lng:       lng,
		Distance:  dist,
		Limit:     limit,
		Offset:    offset,
	}

	opps, total, err := h.oppService.SearchOpportunities(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search opportunities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  opps,
		"total": total,
	})
}

func (h *OpportunityHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req domain.Opportunity
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Get User ID -> Host ID
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	host, err := h.hostService.GetHostByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a host"})
		return
	}

	// 2. Get Existing Opportunity
	existingOpp, err := h.oppService.GetOpportunityByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}

	// 3. Check Ownership
	if existingOpp.HostID != host.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this opportunity"})
		return
	}

	// 4. Update
	req.ID = existingOpp.ID
	req.HostID = existingOpp.HostID

	err = h.oppService.UpdateOpportunity(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update opportunity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "opportunity updated"})
}

func (h *OpportunityHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// 1. Get User ID -> Host ID
	claims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	mapClaims := claims.(jwt.MapClaims)
	userID := mapClaims["sub"].(string)

	host, err := h.hostService.GetHostByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a host"})
		return
	}

	// 2. Get Existing Opportunity
	existingOpp, err := h.oppService.GetOpportunityByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}

	// 3. Check Ownership
	if existingOpp.HostID != host.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this opportunity"})
		return
	}

	// 4. Delete
	err = h.oppService.DeleteOpportunity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete opportunity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "opportunity deleted"})
}
