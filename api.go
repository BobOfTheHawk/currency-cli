// In api.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// apiKey should be replaced with your actual API key from exchangerate-api.com
const apiKey = "04f1d4f1e6e887bc8411157c"

// ApiResponse matches the structure of the JSON response from the API.
type ApiResponse struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

// CodesResponse matches the structure of the /codes endpoint response.
type CodesResponse struct {
	Result         string     `json:"result"`
	SupportedCodes [][]string `json:"supported_codes"`
}

// getCurrencyNames fetches the list of supported currency codes and their names.
func getCurrencyNames() (map[string]string, error) {
	if apiKey == "YOUR_API_KEY" || apiKey == "" {
		return nil, fmt.Errorf("error: API key is not set. Please replace 'YOUR_API_KEY' in api.go")
	}

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/codes", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request for codes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("codes API request failed with status: %s", resp.Status)
	}

	var codesResponse CodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&codesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode codes API response: %w", err)
	}

	if codesResponse.Result != "success" {
		return nil, fmt.Errorf("codes API returned an error: %s", codesResponse.Result)
	}

	// Convert the list of lists to a map for easy lookup.
	currencyNames := make(map[string]string)
	for _, codePair := range codesResponse.SupportedCodes {
		if len(codePair) == 2 {
			currencyNames[codePair[0]] = codePair[1]
		}
	}

	return currencyNames, nil
}

// getRates fetches the latest conversion rates for a given base currency.
func getRates(baseCurrency string) (map[string]float64, error) {
	// Check if the user has replaced the placeholder API key.
	if apiKey == "YOUR_API_KEY" || apiKey == "" {
		return nil, fmt.Errorf("error: API key is not set. Please replace 'YOUR_API_KEY' in api.go")
	}

	// Construct the API request URL.
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, baseCurrency)

	// Make the HTTP GET request.
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-successful status codes.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Decode the JSON response into our struct.
	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Check if the API call itself was successful.
	if apiResponse.Result != "success" {
		return nil, fmt.Errorf("API returned an error: %s", apiResponse.Result)
	}

	return apiResponse.ConversionRates, nil
}
