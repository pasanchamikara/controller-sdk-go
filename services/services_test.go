package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	deis "github.com/pasanchamikara/controller-sdk-go"
	"github.com/pasanchamikara/controller-sdk-go/api"
)

const servicesFixture string = `
{
  "services": [
    {
        "procfile_type": "web",
        "path_pattern":  "/abc/xyz"
    },
    {
        "procfile_type": "worker",
        "path_pattern":  "/x1z/a2c"
    }
  ]
}`

const serviceFixture string = `
{
    "procfile_type": "web",
    "path_pattern":  "/abc/xyz"
}`

const serviceCreateExpected string = `{"procfile_type":"web","path_pattern":"/abc/xyz"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "GET" {
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != serviceCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", serviceCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(serviceFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte(servicesFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestServicesList(t *testing.T) {
	t.Parallel()

	expected := api.Services{
		{
			ProcfileType: "web",
			PathPattern:  "/abc/xyz",
		},
		{
			ProcfileType: "worker",
			PathPattern:  "/x1z/a2c",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := List(deis, "example-go")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestServicesAdd(t *testing.T) {
	t.Parallel()

	expected := api.Service{
		ProcfileType: "web",
		PathPattern:  "/abc/xyz",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := New(deis, "example-go", "web", "/abc/xyz")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestServicesRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(deis, "example-go", "web"); err != nil {
		t.Fatal(err)
	}
}
