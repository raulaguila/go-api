package handler

//
//import (
//	"errors"
//	"github.com/gofiber/fiber/v2"
//	"github.com/raulaguila/go-api/internal/pkg/filters"
//	"testing"
//
//	"github.com/raulaguila/go-api/internal/pkg/domain"
//	"github.com/raulaguila/go-api/internal/pkg/dto"
//	"github.com/raulaguila/go-api/pkg/filter"
//	"github.com/raulaguila/go-api/pkg/helper"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//// MockService is a mock implementation of the ProductService
//type MockService struct {
//	mock.Mock
//}
//
//func (m *MockService) GenerateProductOutputDTO(p *domain.Product) *dto.ProductOutputDTO {
//	return m.Called(p).Get(0).(*dto.ProductOutputDTO)
//}
//
//func (m *MockService) GetProductByID(ctx context.Context, id uint) (*dto.ProductOutputDTO, error) {
//	args := m.Called(ctx, id)
//	return args.Get(0).(*dto.ProductOutputDTO), args.Error(1)
//}
//
//func (m *MockService) GetProducts(ctx context.Context, f *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error) {
//	args := m.Called(ctx, f)
//	return args.Get(0).(*dto.ItemsOutputDTO[dto.ProductOutputDTO]), args.Error(1)
//}
//
//func (m *MockService) CreateProduct(ctx context.Context, dto *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
//	args := m.Called(ctx, dto)
//	return args.Get(0).(*dto.ProductOutputDTO), args.Error(1)
//}
//
//func (m *MockService) UpdateProduct(ctx context.Context, id uint, dto *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
//	args := m.Called(ctx, id, dto)
//	return args.Get(0).(*dto.ProductOutputDTO), args.Error(1)
//}
//
//func (m *MockService) DeleteProducts(ctx context.Context, ids []uint) error {
//	return m.Called(ctx, ids).Error(0)
//}
//
//func mockErrorHandler(c *fiber.Ctx, err error) error {
//	return err
//}
//
//func TestGetProducts(t *testing.T) {
//	tests := []struct {
//		name    string
//		setup   func(*MockService)
//		wantErr bool
//	}{
//		{
//			name: "Success",
//			setup: func(m *MockService) {
//				m.On("GetProducts", mock.Anything, mock.Anything).Return(new(dto.ItemsOutputDTO[dto.ProductOutputDTO]), nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "ServiceError",
//			setup: func(m *MockService) {
//				m.On("GetProducts", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			service := new(MockService)
//			tt.setup(service)
//			handler := &ProductHandler{
//				productService: service,
//				handlerError:   mockErrorHandler,
//			}
//
//			c := app.AcquireCtx(&fiber.Ctx{})
//			defer app.ReleaseCtx(c)
//			c.Locals(helper.LocalFilter, &filter.Filter{})
//
//			err := handler.getProducts(c)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestGetProductByID(t *testing.T) {
//	tests := []struct {
//		name    string
//		id      uint
//		setup   func(*MockService)
//		wantErr bool
//	}{
//		{
//			name: "Success",
//			id:   1,
//			setup: func(m *MockService) {
//				m.On("GetProductByID", mock.Anything, uint(1)).Return(new(dto.ProductOutputDTO), nil)
//			},
//			wantErr: false,
//		},
//		{
//			name:    "InvalidID",
//			id:      0,
//			setup:   func(m *MockService) {},
//			wantErr: true,
//		},
//		{
//			name: "NotFoundError",
//			id:   1,
//			setup: func(m *MockService) {
//				m.On("GetProductByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			service := new(MockService)
//			tt.setup(service)
//			handler := &ProductHandler{
//				productService: service,
//				handlerError:   mockErrorHandler,
//			}
//
//			c := app.AcquireCtx(&fiber.Ctx{})
//			defer app.ReleaseCtx(c)
//			c.Locals(helper.LocalID, &filters.IDFilter{ID: tt.id})
//
//			err := handler.getProductByID(c)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestCreateProduct(t *testing.T) {
//	tests := []struct {
//		name    string
//		setup   func(*MockService)
//		wantErr bool
//	}{
//		{
//			name: "Success",
//			setup: func(m *MockService) {
//				m.On("CreateProduct", mock.Anything, mock.Anything).Return(new(dto.ProductOutputDTO), nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "ServiceError",
//			setup: func(m *MockService) {
//				m.On("CreateProduct", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			service := new(MockService)
//			tt.setup(service)
//			handler := &ProductHandler{
//				productService: service,
//				handlerError:   mockErrorHandler,
//			}
//
//			c := app.AcquireCtx(&fiber.Ctx{})
//			defer app.ReleaseCtx(c)
//			c.Locals(helper.LocalDTO, &dto.ProductInputDTO{})
//
//			err := handler.createProduct(c)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUpdateProduct(t *testing.T) {
//	tests := []struct {
//		name    string
//		id      uint
//		setup   func(*MockService)
//		wantErr bool
//	}{
//		{
//			name: "Success",
//			id:   1,
//			setup: func(m *MockService) {
//				m.On("UpdateProduct", mock.Anything, uint(1), mock.Anything).Return(new(dto.ProductOutputDTO), nil)
//			},
//			wantErr: false,
//		},
//		{
//			name:    "InvalidID",
//			id:      0,
//			setup:   func(m *MockService) {},
//			wantErr: true,
//		},
//		{
//			name: "ServiceError",
//			id:   1,
//			setup: func(m *MockService) {
//				m.On("UpdateProduct", mock.Anything, uint(1), mock.Anything).Return(nil, errors.New("error"))
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			service := new(MockService)
//			tt.setup(service)
//			handler := &ProductHandler{
//				productService: service,
//				handlerError:   mockErrorHandler,
//			}
//
//			c := app.AcquireCtx(&fiber.Ctx{})
//			defer app.ReleaseCtx(c)
//			c.Locals(helper.LocalID, &filters.IDFilter{ID: tt.id})
//			c.Locals(helper.LocalDTO, &dto.ProductInputDTO{})
//
//			err := handler.updateProduct(c)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestDeleteProducts(t *testing.T) {
//	tests := []struct {
//		name    string
//		ids     []uint
//		setup   func(*MockService)
//		wantErr bool
//	}{
//		{
//			name: "Success",
//			ids:  []uint{1, 2, 3},
//			setup: func(m *MockService) {
//				m.On("DeleteProducts", mock.Anything, []uint{1, 2, 3}).Return(nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "ParseError",
//			ids:  nil,
//			setup: func(m *MockService) {
//				m.On("DeleteProducts", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			service := new(MockService)
//			tt.setup(service)
//			handler := &ProductHandler{
//				productService: service,
//				handlerError:   mockErrorHandler,
//			}
//
//			c := app.AcquireCtx(&fiber.Ctx{})
//			defer app.ReleaseCtx(c)
//			c.BodyParser(&dto.IDsInputDTO{IDs: tt.ids})
//
//			err := handler.deleteProducts(c)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
