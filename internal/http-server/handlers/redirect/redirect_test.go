package redirect_test

import (
	"github.com/vanya-egorov/url_shortener.git/internal/lib/api"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vanya-egorov/url_shortener.git/internal/http-server/handlers/redirect"
	"github.com/vanya-egorov/url_shortener.git/internal/http-server/handlers/redirect/mocks"
	"github.com/vanya-egorov/url_shortener.git/internal/lib/logger/handlers/slogdiscard"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},

		{
			name:  "Redirect to Long URL",
			alias: "long_url_alias",
			url:   "https://www.example.com/very/long/url/that/should/be/shortened",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			if tc.respError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.respError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.url, redirectedToURL)
			}
		})
	}
}
