package jsonrespond_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ripple-cloud/cloud/jsonrespond"
)

func setupServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

type SampleObj struct {
	Key string "json: key"
	Val string "json: value"
}

type SampleObjs struct {
	Root []SampleObj "json: objects"
}

func TestRespond(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		payload := SampleObjs{
			[]SampleObj{
				{"foo", "bar"},
			},
		}
		jsonrespond.Respond(w, http.StatusOK, payload)
	}
	ts := setupServer(h)

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	// check the status
	if res.StatusCode != http.StatusOK {
		t.Error("expected status code %v, received %v", http.StatusOK, res.StatusCode)
	}

	// check the response
	output := &SampleObjs{}
	err = json.NewDecoder(res.Body).Decode(output)
	if err != nil {
		t.Fatal(err)
	}
	if output.Root[0].Key != "foo" || output.Root[0].Val != "bar" {
		t.Error("unexepected Response: %v", output)
	}
}

func TestRespondWithNilPointer(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		jsonrespond.OK(w, nil)
	}
	ts := setupServer(h)

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	// check the status
	if res.StatusCode != http.StatusOK {
		t.Error("expected status code %v, received %v", http.StatusOK, res.StatusCode)
	}

	// check the response
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if string(output) != "null" {
		t.Errorf("expected null response, but got %s", output)
	}
}
