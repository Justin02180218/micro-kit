package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

type BASEJSONRequest struct {
	V interface{} `json:"v"`
}

type BASEJSONResponse struct {
	V interface{} `json:"v"`
}

func EncodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func EncodeJSONRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func DecodeJSONResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response BASEJSONResponse
	var v interface{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	response.V = v
	return response, nil
}

func DecodeJSONRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req BASEJSONRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.V, nil
}
