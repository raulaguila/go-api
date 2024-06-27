package handler

import (
	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/helper"
)

var middlewareFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &filter.Filter{},
})

var middlewareUserFilterDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalFilter,
	OnLookup:   datatransferobject.Query,
	Model:      &filters.UserFilter{},
})

var middlewareIDDTO = datatransferobject.New(datatransferobject.Config{
	ContextKey: helper.LocalID,
	OnLookup:   datatransferobject.Params,
	Model:      &filters.IDFilter{},
})
