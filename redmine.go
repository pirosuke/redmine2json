package redmine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

/*
TicketFetchParams describes the format for ticket fetch parameters.
*/
type TicketFetchParams struct {
	Offset int
	Limit  int
}

/*
BasicAuth describes the format for basic authentication.
*/
type BasicAuth struct {
	UserName string
	Password string
}

/*
FetchTickets fetches ticket list from redmine.
*/
func FetchTickets(baseURL string, apiKey string, params TicketFetchParams, auth BasicAuth) ([]byte, error) {

	values := url.Values{}

	if params.Offset > 0 {
		values.Add("offset", strconv.Itoa(params.Offset))
	}

	if params.Limit > 0 {
		values.Add("limit", strconv.Itoa(params.Limit))
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL+"/issues.json?"+values.Encode(), nil)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Redmine-API-Key", apiKey)

	if auth.UserName != "" {
		req.SetBasicAuth(auth.UserName, auth.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
