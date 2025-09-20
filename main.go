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
var widgetType = "cfb" // nba, nfl, cfb

func main() {
	// Load widget from file
	templatePath := fmt.Sprintf("widgets/%s.j2", widgetType)
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template: %v\n", err)
		return
	}

	tmpl, err := template.New(widgetType).Funcs(glance.GetTemplateFunctions()).Parse(string(templateBytes))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	var jsonResult gjson.Result
	if sourceType == "api" {
		jsonResult, err = getAPIData(widgetType)
	} else {
		jsonResult, err = getTestData(widgetType)
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

	resultPath := fmt.Sprintf("data/results/%s-widget.html", widgetType)
	err = os.WriteFile(resultPath, []byte(result.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing result to file: %v\n", err)
		return
	}

	fmt.Println(result.String())
}

func getTestData(widgetType string) (gjson.Result, error) {
	filePath := fmt.Sprintf("data/responses/%s-response.json", widgetType)
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("reading test data: %w", err)
	}
	return gjson.Parse(string(jsonBytes)), nil
}

func getAPIData(widgetType string) (gjson.Result, error) {
	var apiURL string
	switch widgetType {
	case "nba":
		apiURL = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"
	case "nfl":
		apiURL = "https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard"
	case "cfb":
		apiURL = "https://site.api.espn.com/apis/site/v2/sports/football/college-football/scoreboard"
	default:
		return gjson.Result{}, fmt.Errorf("unsupported widget type: %s", widgetType)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("fetching %s data: %w", widgetType, err)
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
	responsePath := fmt.Sprintf("data/responses/%s-response.json", widgetType)
	err = os.WriteFile(responsePath, jsonBytes, 0644)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("saving json response: %w", err)
	}

	return gjson.Parse(string(jsonBytes)), nil
}
