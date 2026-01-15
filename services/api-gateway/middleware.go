package main

import (
	"net/http"
	"strings"

	"github.com/baobei23/e-ticket/services/api-gateway/grpc_clients"
	"github.com/baobei23/e-ticket/shared/proto/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *grpc_clients.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		token := parts[1]

		// Validasi via gRPC
		resp, err := authClient.Client.ValidateToken(c.Request.Context(), &auth.ValidateTokenRequest{
			Token: token,
		})

		if err != nil || !resp.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set UserID ke Context
		c.Set("user_id", resp.UserId)
		c.Next()
	}
}
