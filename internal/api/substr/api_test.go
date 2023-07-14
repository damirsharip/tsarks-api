package substr

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func TestFindSubstr(t *testing.T) {
	tests := []struct {
		name       string
		body       io.Reader
		wantStatus int
		wantRes    string
	}{
		{
			name: "OK",
			body: bytes.NewBuffer(
				[]byte(
					`{
						"word": "damird"
					}`,
				),
			),
			wantStatus: http.StatusOK,
			wantRes:    "{\"answer\":\"amird\"}\n",
		},
		{
			name: "invalid request payload",
			body: bytes.NewBuffer(
				[]byte(
					`{
						asd
					}`,
				),
			),
			wantStatus: http.StatusBadRequest,
			wantRes:    "{\"error\":\"invalid request payload\"}\n",
		},
		{
			name: "OK",
			body: bytes.NewBuffer(
				[]byte(
					`{
						"word": ""
					}`,
				),
			),
			wantStatus: http.StatusBadRequest,
			wantRes:    "{\"error\":\"bad request: word required\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := httprouter.New()
			router.Handler(http.MethodPost, "/rest/substr/find", FindSubstr())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/rest/substr/find", tt.body)

			router.ServeHTTP(rec, req)

			require.Equal(t, rec.Code, tt.wantStatus)
			require.Equal(t, rec.Body.String(), tt.wantRes)
		})
	}
}
