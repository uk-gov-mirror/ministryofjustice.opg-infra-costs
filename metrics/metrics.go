package metrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"opg-infra-costs/costs"
)

const ApiEndpoint = "https://lz0rqicsng.execute-api.eu-west-1.amazonaws.com/prod/streams/TimestreamMetricsStream/records"

// SendToApi formats a put request to send to the api
func SendToApi(costs []costs.CostRow) {
	// convert from cost struct to json
	jsonData, _ := FromCosts(costs)

	client := &http.Client{}
	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, ApiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	fmt.Println(req)
	fmt.Println()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	fmt.Println()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
