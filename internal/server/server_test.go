package server_test

import (
	"github.com/a-clap/dictionary/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleTranslateEng(t *testing.T) {
	srv := server.New()

	t.Run("translate brain from english", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/translate/en/brain", nil)
		response := httptest.NewRecorder()

		srv.ServeHTTP(response, request)
		got := response.Body.String()
		want := "mózg"

		if got != want {
			t.Errorf("%s: Got %#v, want %q", t.Name(), got, want)
		}
	})
}
