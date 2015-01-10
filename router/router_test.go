package router_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ripple-cloud/cloud/router"
)

func setupServer(tmpDir string) (*httptest.Server, error) {
	r := router.New()

	// directly access an endpoint
	r.GET("/endpoint", func(w http.ResponseWriter, r *http.Request, c router.Context) {
		w.Write([]byte("OK"))
	})

	// using handlers to modify the context
	r.GET("/with_modifiers",
		func(w http.ResponseWriter, r *http.Request, c router.Context) {
			c.Meta["test"] = "OK"
			// call next handler
			c.Next(w, r, c)
		},
		func(w http.ResponseWriter, r *http.Request, c router.Context) {
			o := c.Meta["test"].(string)
			w.Write([]byte(o))
		})

	// check named params and request form params
	r.GET("/with_params/:param", func(w http.ResponseWriter, r *http.Request, c router.Context) {
		fmt.Fprintf(w, "%s %s", c.Params.ByName("param"), r.FormValue("param"))
	})

	// respond early
	r.GET("/respond_early",
		func(w http.ResponseWriter, r *http.Request, c router.Context) {
			w.Write([]byte("foo"))
		},
		func(w http.ResponseWriter, r *http.Request, c router.Context) {
			w.Write([]byte("bar"))
		})

	// default handlers
	r.Default(func(w http.ResponseWriter, r *http.Request, c router.Context) {
		c.Meta["default"] = "hello world"
		// call next handler
		c.Next(w, r, c)
	})

	r.GET("/use_default_handler",
		func(w http.ResponseWriter, r *http.Request, c router.Context) {
			o := c.Meta["default"].(string)
			w.Write([]byte(o))
		})

	// serving static files
	err := ioutil.WriteFile(tmpDir+"/test.txt", []byte("foo bar"), 0644)
	if err != nil {
		return nil, err
	}
	r.ServeFiles("/public/*filepath", http.Dir(tmpDir))

	return httptest.NewServer(r), nil
}

func TestRouter(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "router-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	ts, err := setupServer(tmpDir)
	if err != nil {
		panic(err)
	}
	defer ts.Close()

	type testCase struct {
		path string
		body string
	}

	tCases := []testCase{
		{"/endpoint", "OK"},
		{"/with_modifiers", "OK"},
		{"/with_params/foo?param=bar", "foo bar"},
		{"/respond_early", "foo"},
		{"/use_default_handler", "hello world"},
		{"/public/test.txt", "foo bar"},
	}

	for _, tc := range tCases {
		res, err := http.Get(ts.URL + tc.path)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if body := string(b); body != tc.body {
			t.Errorf("%s - Expected %s, Got %s", tc.path, tc.body, body)
		}
	}
}
