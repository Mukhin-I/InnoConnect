package transport

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    userpb "innoconnect/pkg/pb/user"
)

func (h *Handler) Register(c *gin.Context) {
    var req userpb.RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
        return
    }

    resp, err := h.userClient.Register(c.Request.Context(), &userpb.RegisterRequest{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": resp.GetMessage()})
}

func (h *Handler) Login(c *gin.Context) {
    var req userpb.LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request payload"})
        return
    }

    resp, err := h.userClient.Login(c.Request.Context(), &userpb.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token":     resp.GetToken(),
        "type":      resp.GetType(),
        "expiresIn": resp.GetExpiresIn(),
    })
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
        return
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
        return
    }

    resp, err := h.userClient.GetCurrentUser(c.Request.Context(), &userpb.GetCurrentUserRequest{})
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":    resp.GetId(),
        "email": resp.GetEmail(),
        "name":  resp.GetName(),
    })
}
