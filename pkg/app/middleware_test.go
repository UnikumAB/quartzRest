package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_commonMiddleware(t *testing.T) {
	t.Parallel()
	type args struct {
		path string
	}
	type want struct {
		status      int
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "check normal path",
			args: args{
				path: "/test",
			},
			want: want{status: http.StatusOK, contentType: "application/json"},
		},
		{
			name: "check metrics",
			args: args{
				path: "/metrics",
			},
			want: want{status: http.StatusOK, contentType: "text/plain; version=0.0.4; charset=utf-8"},
		},
		{
			name: "check notFound",
			args: args{
				path: "/notFound",
			},
			want: want{status: http.StatusNotFound, contentType: "text/plain; charset=utf-8"},
		},
	}
	router := Router(App{})
	router.HandleFunc("/test", dummyHandler())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", tt.args.path, nil)
			w := httptest.NewRecorder()
			handler := dummyHandler()
			commonMiddleware(handler)
			router.ServeHTTP(w, r)
			assert.Equal(t, tt.want.status, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header().Get("Content-Type"))
		})
	}
}
func dummyHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
