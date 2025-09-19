package main

import (
	"fmt"
	"glance-widgets/glance"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

// var sourceType = "test-data"
var sourceType = "api"

func main() {

	// Load widget from file
	templateBytes, err := os.ReadFile("widgets/nba.tmpl")
	if err != nil {
		fmt.Printf("Error reading template: %v\n", err)
		return
	}

	tmpl, err := template.New("nba").Funcs(glance.GetTemplateFunctions()).Parse(string(templateBytes))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	var jsonResult gjson.Result
	if sourceType == "api" {
		jsonResult, err = getAPIData()
	} else {
		jsonResult, err = getTestData()
	}
	if err != nil {
		fmt.Printf("Error getting JSON data: %v\n", err)
		return
	}

	testData := map[string]interface{}{
		"JSON": glance.NewDecoratedGJSONResult(jsonResult),
	}

	// Execute the template
	var result strings.Builder
	err = tmpl.Execute(&result, testData)
	if err != nil {
		fmt.Printf("Template execution failed: %v\n", err)
		return
	}
	err = os.WriteFile("data/results/nba-widget.html", []byte(result.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing result to file: %v\n", err)
		return
	}

	fmt.Println(result.String())
}

func getTestData() (gjson.Result, error) {
	jsonBytes, err := os.ReadFile("data/mock/nba-response.json")
	if err != nil {
		return gjson.Result{}, fmt.Errorf("reading test data: %w", err)
	}
	return gjson.Parse(string(jsonBytes)), nil
}

func getAPIData() (gjson.Result, error) {
	// Fetch NBA API data
	apiURL := "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"

	resp, err := http.Get(apiURL)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("fetching NBA data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return gjson.Result{}, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("reading response: %w", err)
	}

	// Save the raw JSON response to file for inspection
	err = os.WriteFile("data/responses/nba-response.json", jsonBytes, 0644)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("saving json response: %w", err)
	}

	return gjson.Parse(string(jsonBytes)), nil

}
