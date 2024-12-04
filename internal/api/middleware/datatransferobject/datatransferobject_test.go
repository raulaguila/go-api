package datatransferobject

//
//import (
//	"github.com/gofiber/fiber/v2"
//	"reflect"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestNew(t *testing.T) {
//	tests := []struct {
//		name          string
//		config        Config
//		contextData   func(c *fiber.Ctx)
//		expectedError string
//		localKey      string
//	}{
//		{
//			name: "BodyParserSuccess",
//			config: Config{
//				OnLookup: Body,
//				Model:    &struct{ Name string }{},
//				ErrorHandler: func(c *fiber.Ctx, err error) error {
//					return err
//				},
//				ContextKey: "data",
//			},
//			contextData: func(c *fiber.Ctx) {
//				c.Request().SetBody([]byte(`{"Name":"Alice"}`))
//			},
//			expectedError: "",
//			localKey:      "data",
//		},
//		{
//			name: "BodyParserError",
//			config: Config{
//				OnLookup: Body,
//				Model:    &struct{ Age int }{},
//				ErrorHandler: func(c *fiber.Ctx, err error) error {
//					return err
//				},
//				ContextKey: "data",
//			},
//			contextData: func(c *fiber.Ctx) {
//				c.Request().SetBody([]byte(`{"Age":"not_an_int"}`))
//			},
//			expectedError: "Error mid: ",
//			localKey:      "data",
//		},
//		{
//			name: "QueryParserSuccess",
//			config: Config{
//				OnLookup: Query,
//				Model:    &struct{ Name string }{},
//				ErrorHandler: func(c *fiber.Ctx, err error) error {
//					return err
//				},
//				ContextKey: "data",
//			},
//			contextData: func(c *fiber.Ctx) {
//				c.Request().URI().SetQueryString("Name=Bob")
//			},
//			expectedError: "",
//			localKey:      "data",
//		},
//		{
//			name: "ParamsParserSuccess",
//			config: Config{
//				OnLookup: Params,
//				Model:    &struct{ ID string }{},
//				ErrorHandler: func(c *fiber.Ctx, err error) error {
//					return err
//				},
//				ContextKey: "data",
//			},
//			contextData: func(c *fiber.Ctx) {
//				c.Locals("ID", "123")
//			},
//			expectedError: "",
//			localKey:      "data",
//		},
//		{
//			name: "NextHandlerSkip",
//			config: Config{
//				Next: func(c *fiber.Ctx) bool {
//					return true
//				},
//				OnLookup: Params,
//				Model:    &struct{}{},
//				ErrorHandler: func(c *fiber.Ctx, err error) error {
//					return err
//				},
//			},
//			contextData:   func(c *fiber.Ctx) {},
//			expectedError: "",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			app.Use(New(tt.config))
//
//			app.Get("/", func(c *fiber.Ctx) error {
//				obj := c.Locals(tt.config.ContextKey)
//				if tt.expectedError != "" {
//					assert.Nil(t, obj)
//				} else {
//					assert.NotNil(t, obj)
//					assert.Equal(t, reflect.TypeOf(tt.config.Model), reflect.TypeOf(obj))
//				}
//				return nil
//			})
//
//			req := &fiber.Ctx{Request: app.AcquireCtx(nil).Request()}
//			tt.contextData(req)
//
//			err := app.Handler()(req)
//			if tt.expectedError != "" {
//				assert.Contains(t, err.Error(), tt.expectedError)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
