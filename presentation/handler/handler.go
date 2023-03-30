package handler

import (
	"encoding/json"
	"github.com/jpdel518/go-onionarch/usecase"
	"github.com/rs/cors"
	"net/http"
)

// NewHandler create handler
func NewHandler(usecase usecase.ArticleUsecase) error {
	articleHandler := NewArticleHandler(usecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/article/fetch", articleHandler.Fetch)
	mux.HandleFunc("/article/show-by-id/", articleHandler.ShowById)
	mux.HandleFunc("/article/get-by-title", articleHandler.ShowByTitle)
	mux.HandleFunc("/article/update", articleHandler.Update)
	mux.HandleFunc("/article/store", articleHandler.Store)
	mux.HandleFunc("/article/delete", articleHandler.Delete)

	// CORS
	c := cors.Default()
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:3000", "http://foo.com"},
	// 	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// })
	handler := c.Handler(mux)

	return http.ListenAndServe(":8080", handler)
}

// ApiRequestResponse response json
type ApiRequestResponse struct {
	Code int
	Data interface{}
}

// CreateResponseJson create response data as json
func CreateResponseJson(a *ApiRequestResponse) ([]byte, error) {
	js, err := json.Marshal(*a)
	if err != nil {
		return nil, err
	}
	return js, nil
}
