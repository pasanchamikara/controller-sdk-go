package domains

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

const domainsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "created": "2014-01-01T00:00:00UTC",
            "domain": "example.example.com",
            "owner": "test",
            "updated": "2014-01-01T00:00:00UTC"
        }
    ]
}`

const domainFixture string = `
{
    "app": "example-go",
    "created": "2014-01-01T00:00:00UTC",
    "domain": "example.example.com",
    "owner": "test",
    "updated": "2014-01-01T00:00:00UTC"
}`

const domainCreateExpected string = `{"domain":"example.example.com"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/domains/" && req.Method == "GET" {
		res.Write([]byte(domainsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/domains/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != domainCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", domainCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(domainFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/domains/test.com" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte(domainsFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestDomainsList(t *testing.T) {
	t.Parallel()

	expected := api.Domains{
		{
			App:     "example-go",
			Created: "2014-01-01T00:00:00UTC",
			Domain:  "example.example.com",
			Owner:   "test",
			Updated: "2014-01-01T00:00:00UTC",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(deis, "example-go", 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestDomainsAdd(t *testing.T) {
	t.Parallel()

	expected := api.Domain{
		App:     "example-go",
		Created: "2014-01-01T00:00:00UTC",
		Domain:  "example.example.com",
		Owner:   "test",
		Updated: "2014-01-01T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := New(deis, "example-go", "example.example.com")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestDomainsRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(deis, "example-go", "test.com"); err != nil {
		t.Fatal(err)
	}
}
