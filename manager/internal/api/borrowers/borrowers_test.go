package borrowers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/csokviktor/lib_manager/internal/api"
	"github.com/csokviktor/lib_manager/internal/api/books"
	"github.com/csokviktor/lib_manager/internal/repository"
	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	"github.com/csokviktor/lib_manager/internal/service"
	"github.com/csokviktor/lib_manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mocks.NewMockBorrowers(ctrl)

	type args struct {
		borrowersService service.Borrowers
	}
	tests := []struct {
		name string
		args args
		want Routes
	}{
		{
			name: "happy",
			args: args{
				borrowersService: serviceMock,
			},
			want: &routes{
				borrowersService: serviceMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoutes(tt.args.borrowersService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routes_CreateBorrower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		body []byte
	}
	type fields struct {
		borrowersService func() service.Borrowers
	}
	tests := []struct {
		name     string
		args     args
		fields   fields
		want     func() []byte
		wantCode int
	}{
		{
			name: "happy",
			fields: fields{
				borrowersService: func() service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						Create(gomock.Any(), &dbmodels.Borrower{Name: "John Doe"}).
						Return(&dbmodels.Borrower{ID: "mockid", Name: "John Doe"}, nil)
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"name": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(
					APIBorrower{
						ID:    "mockid",
						Name:  "John Doe",
						Books: []books.APIBook{},
					},
				)
				return jsonBody
			},
			wantCode: 201,
		},
		{
			name: "parse_error",
			fields: fields{
				borrowersService: func() service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"name_unknown": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "Key: 'BorrowerCreateValidator.Name' Error:Field validation for 'Name' failed on the 'required' tag",
				})
				return jsonBody
			},
			wantCode: 400,
		},
		{
			name: "internal_error",
			fields: fields{
				borrowersService: func() service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						Create(gomock.Any(), gomock.Any()).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				body: func() []byte {
					jsonBody, _ := json.Marshal(map[string]string{
						"name": "John Doe",
					})
					return jsonBody
				}(),
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "error",
				})
				return jsonBody
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRoutes(tt.fields.borrowersService())

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodPost, "/borrowers", bytes.NewBuffer(tt.args.body))
			handler.CreateBorrower(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want())
		})
	}
}

func Test_routes_GetBorrower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockID := "ab7c7712-d10c-412a-87fd-7340b8361e31"

	type fields struct {
		borrowersService func(id string) service.Borrowers
	}
	type args struct {
		borrowerID string
	}
	tests := []struct {
		name     string
		args     args
		fields   fields
		want     func(id string) []byte
		wantCode int
	}{
		{
			name: "happy",
			fields: fields{
				borrowersService: func(id string) service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						GetByID(gomock.Any(), id).
						Return(&dbmodels.Borrower{ID: id, Name: "John Doe"}, nil)
					return mock
				},
			},
			args: args{
				borrowerID: mockID,
			},
			want: func(id string) []byte {
				jsonBody, _ := json.Marshal(
					APIBorrower{
						ID:    id,
						Name:  "John Doe",
						Books: []books.APIBook{},
					},
				)
				return jsonBody
			},
			wantCode: 200,
		},
		{
			name: "invalid_uuid",
			fields: fields{
				borrowersService: func(id string) service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					return mock
				},
			},
			args: args{
				borrowerID: "invalid_uuid",
			},
			want: func(_ string) []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "invalid UUID length: 12",
				})
				return jsonBody
			},
			wantCode: 400,
		},
		{
			name: "not_found",
			fields: fields{
				borrowersService: func(id string) service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						GetByID(gomock.Any(), id).
						Return(nil, repository.NotFoundError{})
					return mock
				},
			},
			args: args{
				borrowerID: mockID,
			},
			want: func(_ string) []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "not found",
				})
				return jsonBody
			},
			wantCode: 404,
		},
		{
			name: "internal_error",
			fields: fields{
				borrowersService: func(id string) service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						GetByID(gomock.Any(), id).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			args: args{
				borrowerID: mockID,
			},
			want: func(_ string) []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "error",
				})
				return jsonBody
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRoutes(tt.fields.borrowersService(tt.args.borrowerID))

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = gin.Params{{Key: "id", Value: tt.args.borrowerID}}
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/borrower"+tt.args.borrowerID, nil)
			handler.GetBorrower(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want(tt.args.borrowerID))
		})
	}
}

func Test_routes_GetBorrowers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		borrowersService func() service.Borrowers
	}
	tests := []struct {
		name     string
		fields   fields
		want     func() []byte
		wantCode int
	}{
		{
			name: "happy",
			fields: fields{
				borrowersService: func() service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						GetAll(gomock.Any()).
						Return(
							[]dbmodels.Borrower{{ID: "mockid", Name: "John Doe"}},
							nil,
						)
					return mock
				},
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal([]APIBorrower{
					{
						ID:    "mockid",
						Name:  "John Doe",
						Books: []books.APIBook{},
					},
				},
				)
				return jsonBody
			},
			wantCode: 200,
		},
		{
			name: "internal_error",
			fields: fields{
				borrowersService: func() service.Borrowers {
					mock := mocks.NewMockBorrowers(ctrl)
					mock.EXPECT().
						GetAll(gomock.Any()).
						Return(nil, errors.New("error"))
					return mock
				},
			},
			want: func() []byte {
				jsonBody, _ := json.Marshal(api.ErrorResponse{
					Error: "error",
				})
				return jsonBody
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRoutes(tt.fields.borrowersService())

			// Setup Gin context
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest(http.MethodGet, "/borrowers", nil)
			handler.GetBorrowers(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, w.Body.Bytes(), tt.want())
		})
	}
}
