package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server   *httptest.Server
	usersURL string
	reader   io.Reader
	mux      *http.ServeMux
)

func init() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	usersURL = fmt.Sprintf("%s/signup?username=apple&password=123&email=apple@mail.com", server.URL)
}

func TestTokenHandler(t *testing.T) {
	userJson := `{"access_token": "abc", "token_type": "bearer", "expires_in": 2592000}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("GET", usersURL, reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("%s Success expected: %s", usersURL, res) // res.StatusCode
	}
}
