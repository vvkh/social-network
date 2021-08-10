package profiles_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/mocks"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestHandle(t *testing.T) {
	type mockIn struct {
		firstName string
		lastName  string
		limit     int
	}
	tests := []struct {
		name         string
		url          string
		mockWantIn   mockIn
		mockResponse []entity.Profile
		wantBody     []string
	}{
		{
			name: "all_profiles_are_displayed",
			url:  "/profiles/",
			mockWantIn: mockIn{
				limit: 10,
			},
			mockResponse: []entity.Profile{
				{
					ID:        2,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        3,
					FirstName: "Topsy",
					LastName:  "Cret",
				},
			},
			wantBody: []string{
				`<a href="/profiles/2/">John Doe</a>`,
				`<a href="/profiles/3/">Topsy Cret</a>`,
			},
		},
		{
			name: "search_box_are_displayed_empty",
			url:  "/profiles/",
			mockWantIn: mockIn{
				limit: 10,
			},
			mockResponse: []entity.Profile{},
			wantBody: []string{
				`<input id="first_name" name="first_name" type="text" placeholder="John" value="" />`,
				`<input id="last_name" name="last_name" type="text" placeholder="Doe" value="" />`,
			},
		},
		{
			name: "custom_limit_is_used",
			url:  "/profiles/?limit=20",
			mockWantIn: mockIn{
				limit: 20,
			},
			mockResponse: []entity.Profile{},
			wantBody: []string{
				`<input id="first_name" name="first_name" type="text" placeholder="John" value="" />`,
				`<input id="last_name" name="last_name" type="text" placeholder="Doe" value="" />`,
			},
		},
		{
			name: "search_by_name_form",
			url:  "/profiles/?first_name=john&last_name=doe",
			mockWantIn: mockIn{
				firstName: "john",
				lastName:  "doe",
				limit:     10,
			},
			mockResponse: []entity.Profile{
				{
					ID:        2,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        3,
					FirstName: "Johny",
					LastName:  "Doewan",
				},
			},
			wantBody: []string{
				`<a href="/profiles/2/">John Doe</a>`,
				`<a href="/profiles/3/">Johny Doewan</a>`,
			},
		},
		{
			name: "search_by_first_name_form",
			url:  "/profiles/?first_name=john",
			mockWantIn: mockIn{
				firstName: "john",
				limit:     10,
			},
			mockResponse: []entity.Profile{
				{
					ID:        2,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        3,
					FirstName: "Johny",
					LastName:  "Doewan",
				},
			},
			wantBody: []string{
				`<a href="/profiles/2/">John Doe</a>`,
				`<a href="/profiles/3/">Johny Doewan</a>`,
			},
		},
		{
			name: "search_by_first_name_form",
			url:  "/profiles/?last_name=doe",
			mockWantIn: mockIn{
				lastName: "doe",
				limit:    10,
			},
			mockResponse: []entity.Profile{
				{
					ID:        2,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        3,
					FirstName: "Johny",
					LastName:  "Doewan",
				},
			},
			wantBody: []string{
				`<a href="/profiles/2/">John Doe</a>`,
				`<a href="/profiles/3/">Johny Doewan</a>`,
			},
		},
		{
			name: "form_is_filled_with_search_params",
			url:  "/profiles/?first_name=john&last_name=doe",
			mockWantIn: mockIn{
				firstName: "john",
				lastName:  "doe",
				limit:     10,
			},
			wantBody: []string{
				`<input id="first_name" name="first_name" type="text" placeholder="John" value="john" />`,
				`<input id="last_name" name="last_name" type="text" placeholder="Doe" value="doe" />`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := zap.NewNop().Sugar()
			ctrl := gomock.NewController(t)
			profilesUseCase := mocks.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().ListProfiles(gomock.Any(), test.mockWantIn.firstName, test.mockWantIn.lastName, test.mockWantIn.limit).Return(test.mockResponse, false, nil)
			s := server.New(log, ":80", "../../../templates", nil, profilesUseCase, nil)
			request := httptest.NewRequest("GET", test.url, nil)

			self := entity.Profile{
				ID:        1,
				UserID:    2,
				FirstName: "John",
				LastName:  "Doe",
			}
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), self))
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)
			response := responseWriter.Result()
			require.Equal(t, http.StatusOK, response.StatusCode)
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			for _, bodyPart := range test.wantBody {
				require.Contains(t, string(body), bodyPart)
			}
		})
	}
}
