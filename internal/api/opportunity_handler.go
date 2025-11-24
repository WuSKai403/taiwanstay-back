package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taiwanstay/taiwanstay-back/internal/domain"
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

	opps, err := h.oppService.ListOpportunities(c.Request.Context(), filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list opportunities"})
		return
	}

	c.JSON(http.StatusOK, opps)
}
