package slug

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtemS18/ShortURL-API/internal/delivery/mocks"
	"github.com/ArtemS18/ShortURL-API/internal/entity"
	"github.com/ArtemS18/ShortURL-API/internal/usecase/dto"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateSlugHandler(t *testing.T) {
	t.Parallel()
	valiURL := "http://example.com"

	tests := []struct {
		name           string
		reqBody        map[string]any
		setupMock      func(mockUC *mocks.MockSlugUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			reqBody: map[string]any{
				"url": valiURL,
			},
			setupMock: func(mockUC *mocks.MockSlugUseCase) {
				mockUC.EXPECT().CreateSlug(gomock.Any(), &dto.CreateSlugRequest{URL: valiURL}).Return(&dto.CreateSlugResponse{SlugURL: "https://test.com/slug"}, nil)

			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing url",
			reqBody: map[string]any{
				"url": "",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error usecase",
			reqBody: map[string]any{
				"url": valiURL,
			},
			setupMock: func(mockUC *mocks.MockSlugUseCase) {
				mockUC.EXPECT().CreateSlug(gomock.Any(), &dto.CreateSlugRequest{URL: valiURL}).Return(nil, entity.ServiceError)

			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockSlugUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}
			handler := NewSlugHandler(mockUC)
			body, _ := json.Marshal(tt.reqBody)

			req := httptest.NewRequest(http.MethodPost, "/slugs", bytes.NewBufferString(string(body)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateSlugHandler(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	t.Parallel()
	valiURL := "http://example.com"
	validSlug := "slugslug12"

	tests := []struct {
		name           string
		slug           string
		redirect       bool
		setupMock      func(mockUC *mocks.MockSlugUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success with redirect",
			slug:     validSlug,
			redirect: true,
			setupMock: func(mockUC *mocks.MockSlugUseCase) {
				mockUC.EXPECT().GetURL(gomock.Any(), &dto.GetURLRequest{SlugURL: validSlug}).Return(&dto.GetURLResponse{URL: valiURL}, nil)

			},
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:     "success without redirect",
			slug:     validSlug,
			redirect: false,
			setupMock: func(mockUC *mocks.MockSlugUseCase) {
				mockUC.EXPECT().GetURL(gomock.Any(), &dto.GetURLRequest{SlugURL: validSlug}).Return(&dto.GetURLResponse{URL: valiURL}, nil)

			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error usecase",
			slug: validSlug,
			setupMock: func(mockUC *mocks.MockSlugUseCase) {
				mockUC.EXPECT().GetURL(gomock.Any(), &dto.GetURLRequest{SlugURL: validSlug}).Return(nil, entity.ServiceError)

			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockSlugUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}
			handler := NewSlugHandler(mockUC)
			url := "/" + tt.slug
			if tt.redirect {
				url += "?redirect=true"
			} else {
				url += "?redirect=false"
			}
			req := httptest.NewRequest(http.MethodPost, url, nil)
			req = mux.SetURLVars(req, map[string]string{
				"slug": tt.slug,
			})
			rec := httptest.NewRecorder()

			handler.GetURLHandler(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
