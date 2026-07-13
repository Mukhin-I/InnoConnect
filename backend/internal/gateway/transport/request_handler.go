package transport

import (
	"net/http"
	"strconv"

	"innoconnect/internal/gateway/entity"
	"innoconnect/internal/gateway/usecase"
	"innoconnect/pkg/logger"
	requestpb "innoconnect/pkg/pb/request"
	userpb "innoconnect/pkg/pb/user"

	"innoconnect/pkg/pb/chat"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateRequest godoc
// @Summary Create a request
// @Description Creates a new request for the authenticated user
// @Tags requests
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request body entity.CreateRequestRequest true "Request data"
// @Success 201 {object} entity.CreateRequestResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /requests [post]
func (h *Handler) CreateRequest(c *gin.Context) {
	var req entity.CreateRequestRequest

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		logger.Error("missing authorization header")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid JSON": err.Error()})
		logger.Error("Invalid JSON in request creation")
		return
	}

	logger.Info("Sending gRPC request to create request")

	userID, name, err := usecase.GetUserFromToken(c)

	request, err := h.requestClient.CreateRequest(
		c.Request.Context(),
		&requestpb.CreateRequestRequest{
			CreatorId:        userID,
			CreatorName:      name,
			Title:            req.Title,
			Description:      req.Description,
			RequesterAddress: req.RequesterAddress,
			Type:             req.Type,
			Deadline:         req.Deadline,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error("Failed to get answer from request service via gRPC: " + err.Error())
		return
	}

	_, err = h.userClient.IncrementCreatedRequestsCount(
		c.Request.Context(),
		&userpb.IncrementUserRequestsCountRequest{
			UserId: userID,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error("Failed to increment created requests count: " + err.Error())
		return
	}

	_, err = h.chatClient.GetOrCreateRequestChat(
		c.Request.Context(),
		&chat.GetOrCreateRequestChatRequest{
			RequestId: request.Id,
			ChatName:  request.Title,
			UserId:    userID,
		},
	)

	c.JSON(http.StatusCreated, entity.CreateRequestResponse{
		ID: request.Id,
		Creator: entity.User{
			ID:   request.Creator.Id,
			Name: request.Creator.Name,
		},
		Title:            request.Title,
		Description:      request.Description,
		RequesterAddress: request.RequesterAddress,
		Type:             request.Type,
		Deadline:         request.Deadline,
	})
}

// GetRequests godoc
// @Summary Get all requests
// @Description Returns a list of requests
// @Tags requests
// @Produce json
// @Success 200 {object} entity.GetRequestsResponse
// @Failure 500 {object} map[string]interface{}
// @Router /requests [get]
func (h *Handler) GetRequests(c *gin.Context) {
	logger.Info("Getting requests from request service")

	resp, err := h.requestClient.GetRequests(
		c.Request.Context(),
		&requestpb.GetRequestsRequest{},
	)

	if err != nil {
		logger.Error("Failed to get requests from request service: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
		return
	}

	var requests []entity.RequestShort

	for _, r := range resp.Requests {
		requests = append(requests, entity.RequestShort{
			RequestID: r.Id,
			CreatorID: r.CreatorId,
			Title:     r.Title,
			Type:      r.Type,
			Deadline:  r.Deadline,
		})
	}

	c.JSON(http.StatusOK, entity.GetRequestsResponse{Requests: requests})
}

// GetRequest godoc
// @Summary Get request by ID
// @Description Returns detailed information about a request
// @Tags requests
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} entity.RequestFull
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /requests/{id} [get]
func (h *Handler) GetRequest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	logger.Info("Getting request by id: " + strconv.FormatInt(id, 10))

	resp, err := h.requestClient.GetRequest(
		c.Request.Context(),
		&requestpb.GetRequestRequest{Id: id},
	)

	if err != nil {
		logger.Error("Failed to get request from request service: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get request"})
		return
	}

	c.JSON(http.StatusOK, entity.RequestFull{
		RequestID: resp.Id,
		Creator: entity.User{
			ID:   resp.Creator.Id,
			Name: resp.Creator.Name,
		},
		Title:            resp.Title,
		Description:      resp.Description,
		RequesterAddress: resp.RequesterAddress,
		Type:             resp.Type,
		Deadline:         resp.Deadline,
	})
}

// DeleteRequest godoc
// @Summary Delete a request
// @Description Deletes a request if the authenticated user is the creator
// @Tags requests
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /requests/{id} [delete]
func (h *Handler) DeleteRequest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	creatorID, _, err := usecase.GetUserFromToken(c)
	if err != nil {
		logger.Error("Failed to get user from token: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	_, err = h.requestClient.DeleteRequest(
		c.Request.Context(),
		&requestpb.DeleteRequestRequest{
			Id:        id,
			CreatorId: creatorID,
		},
	)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		case codes.PermissionDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		default:
			logger.Error("Failed to delete request: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete request"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request deleted successfully"})
}

// ApplyToRequest godoc
// @Summary Apply to a request
// @Description Applies the authenticated user to a request
// @Tags requests
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /requests/{id}/apply [post]
func (h *Handler) ApplyToRequest(c *gin.Context) {
	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request id",
		})
		return
	}

	userID, userName, err := usecase.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: remove uneccesarry return from this method
	_, err = h.requestClient.ApplyToRequest(
		c.Request.Context(),
		&requestpb.ApplyToRequestRequest{
			RequestId: requestID,
			UserId:    userID,
			UserName:  userName,
		},
	)

	if err != nil {
		logger.Error(
			"failed to apply to request: " + err.Error(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to apply to request",
		})
		return
	}

	_, err = h.chatClient.AddToChat(
		c.Request.Context(),
		&chat.AddToChatRequest{
			// TODO rename on user
			CreatorId:   userID,
			CreatorName: userName,
			ChatType:    "REQUEST",
			RelaterId:   requestID,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// CancelRequestApplication godoc
// @Summary Cancel request application
// @Description Cancels the authenticated user's application to a request
// @Tags requests
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param request_id path int true "Request ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /requests/{request_id}/applications/{user_id} [delete]
func (h *Handler) CancelRequestApplication(c *gin.Context) {

	requestID, err := strconv.ParseInt(c.Param("request_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request id",
		})
		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	creatorID, _, err := usecase.GetUserFromToken(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = h.requestClient.CancelRequestApplication(
		c.Request.Context(),
		&requestpb.CancelRequestApplicationRequest{
			RequestId: requestID,
			UserId:    userID,
			CreatorId: creatorID,
		},
	)

	if err != nil {
		logger.Error(
			"failed to cancel application: " + err.Error(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to cancel application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
