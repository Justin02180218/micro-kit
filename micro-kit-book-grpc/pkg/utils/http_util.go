package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

func EncodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
