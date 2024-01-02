package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func getCoordinates(address string, key string) (float64, float64, error) {
	// Create the URL
	url := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		address,
		key,
	)

	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("Error getting geocode: %d", resp.StatusCode)
	}

	// Read the HTTP response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	// Unmarshal the JSON response body
	var response struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, 0, err
	}

	// Return the coordinates
	return response.Results[0].Geometry.Location.Lat, response.Results[0].Geometry.Location.Lng, nil
}

func StartUpConnection(Location string) (float64, float64) {
	// Set the Google Maps API key
	key := "AIzaSyAcmeawBhpWjHpmEan-60vWy1p-L3YYjRI"

	// Get the coordinates
	param := url.PathEscape(Location)
	latitude, longitude, err := getCoordinates(param, key)
	if err != nil {
		log.Fatal(err)
	}

	// Print the coordinates
	fmt.Println(latitude, longitude)
	//fmt.Println(url.PathEscape(param), "\n\n")

	return latitude, longitude
	// Print the coordinates

	//fmt.Println(url.PathEscape(param1), "\n\n")

}
