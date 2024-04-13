package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"

	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

// NewRouter creates a new handler for server routes
func NewRouter() (http.Handler, error) {
	// Load openapi spec
	spec, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load openapi spec: %w", err)
	}

	// Set server endpoint to the route. This ensures request validation
	// doesn't fail in the validation middleware.
	// See: https://github.com/deepmap/oapi-codegen/issues/1123
	spec.Servers = openapi3.Servers{&openapi3.Server{URL: "/"}}

	// Create router
	chiRouter := chi.NewRouter()

	// Use request validation middleware
	chiRouter.Use(
		oapimiddleware.OapiRequestValidatorWithOptions(
			spec,
			&oapimiddleware.Options{
				SilenceServersWarning: true,
				Options: openapi3filter.Options{
					AuthenticationFunc: func(ctx context.Context, authInput *openapi3filter.AuthenticationInput) error {
						// Skip auth
						return nil
					},
				},
			},
		),
	)

	// Create route handler for OpenAPI routes
	server := NewStrictHandler(NewServer(), nil)
	routeHandler := HandlerFromMux(server, chiRouter)

	return routeHandler, nil
}