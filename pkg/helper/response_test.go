package helper

//
//import (
//	"github.com/gofiber/fiber/v2"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestNewHTTPResponse(t *testing.T) {
//	app := fiber.New()
//
//	tests := []struct {
//		name           string
//		status         int
//		message        string
//		expectedCode   int
//		expectedResult string
//	}{
//		{"Success", 200, "OK", 200, `{"Code":200,"Message":"OK"}`},
//		{"NotFound", 404, "Not Found", 404, `{"Code":404,"Message":"Not Found"}`},
//		{"ServerError", 500, "Internal Server Error", 500, `{"Code":500,"Message":"Internal Server Error"}`},
//		{"BadRequest", 400, "Bad Request", 400, `{"Code":400,"Message":"Bad Request"}`},
//		{"EmptyMessage", 200, "", 200, `{"Code":200,"Message":""}`},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c, _ := app.Test(httptest.NewRequest(fiber.MethodPost, "/", nil))
//
//			err := NewHTTPResponse(c, tt.status, tt.message)
//			assert.Nil(t, err)
//			assert.Equal(t, tt.expectedCode, c.Response().StatusCode())
//			assert.JSONEq(t, tt.expectedResult, c.Response().Body())
//		})
//	}
//}
