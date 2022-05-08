// Package services provides methods for managing an app's services.
package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	deis "github.com/pasanchamikara/controller-sdk-go"
	"github.com/pasanchamikara/controller-sdk-go/api"
)

// List services registered with an app.
func List(c *deis.Client, appID string) (api.Services, error) {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)
	res, reqErr := c.Request("GET", u, nil)

	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return []api.Service{}, reqErr
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []api.Service{}, err
	}

	r := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return []api.Service{}, err
	}

	out, err := json.Marshal(r["services"].([]interface{}))
	if err != nil {
		return []api.Service{}, err
	}

	var services []api.Service
	if err := json.Unmarshal([]byte(out), &services); err != nil {
		return []api.Service{}, err
	}

	return services, reqErr
}

// New adds a new service to an app. App should already exists.
// Service is the way to route some traffic matching given URL pattern to worker different than `web`
// procfileType - name of the process in Procfile (i.e. <process type> from the `<process type>: <command>`), e.g. `webhooks`
// for more about Procfile see this https://devcenter.heroku.com/articles/procfile
// pathPattern - one or several regexp patterns separated by comma, all request matching given regexp
// are routed to the procfileType workers. E.g. `/webhooks/notify,~ ^/users/[0-9]+/.*/webhooks/notify,/webhooks/rest`
// procfileType and pathPattern are mandatory and should have valid values.
func New(c *deis.Client, appID string, procfileType string, pathPattern string) (api.Service, error) {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)

	req := api.ServiceCreateUpdateRequest{ProcfileType: procfileType, PathPattern: pathPattern}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Service{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !deis.IsErrAPIMismatch(reqErr) {
		return api.Service{}, reqErr
	}
	defer res.Body.Close()

	d := api.Service{ProcfileType: procfileType, PathPattern: pathPattern}
	return d, reqErr
}

// Delete service from app
// If given service for the app doesn't exists then error returned
func Delete(c *deis.Client, appID string, procfileType string) error {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)

	req := api.ServiceDeleteRequest{ProcfileType: procfileType}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.Request("DELETE", u, body)

	return err
}
