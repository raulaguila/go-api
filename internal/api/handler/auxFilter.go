package handler

import (
	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/helper"
)

// middlewareFilterDTO is a data transfer object middleware for handling query parameters using the Filter structure.
var middlewareFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &filter.Filter{},
})

// middlewareUserFilterDTO serves as a middleware for parsing UserFilter data from HTTP request queries into a struct.
// It utilizes the datatransferobject package to bind query parameters to a UserFilter model and stores it in context.
// The data is extracted from the query string, based on the OnLookup setting of the middleware configuration.
var middlewareUserFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &filters.UserFilter{},
})

// middlewareIDDTO is a middleware configuration object that parses ID parameters from URL paths in HTTP requests.
// It stores the parsed data into the context for further handling in the request lifecycle.
var middlewareIDDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalID,
	OnLookup:   datatransferobject.Params,
	Model:      &filters.IDFilter{},
})
