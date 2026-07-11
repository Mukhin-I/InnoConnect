package transport

import (
	"net/http"
	"strconv"

	"innoconnect/internal/gateway/entity"
	"innoconnect/internal/gateway/usecase"
	"innoconnect/pkg/logger"
	requestpb "innoconnect/pkg/pb/request"

	"github.com/gin-gonic/gin"
)

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
			CreatorId:      userID,
			CreatorName:    name,
			Title:          req.Title,
			Description:    req.Description,
			RequesterAddress: req.RequesterAddress,
			Type:           req.Type,
			Deadline:       req.Deadline,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Error("Failed to get answer from request service via gRPC: " + err.Error())
		return
	}

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


	resp, err := h.requestClient.ApplyToRequest(
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


	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot apply to request",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

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


	resp, err := h.requestClient.CancelRequestApplication(
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


	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot cancel application",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}