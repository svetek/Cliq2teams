package msteams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TeamCreationResponse struct {
	statusCode      int
	requestID       string `json:"Request-Id"`
	clientRequestID string `json:"Client-Request-Id"`
	contentLocation string `json:"Content-Location"`
	location        string `json:"Location"`
}

type ChannelCreationResponse struct {
	statusCode      int
	requestID       string `json:"Request-Id"`
	clientRequestID string `json:"Client-Request-Id"`
	contentLocation string `json:"Content-Location"`
	location        string `json:"Location"`

	OdataContext        string      `json:"@odata.context"`
	Id                  string      `json:"id"`
	CreatedDateTime     interface{} `json:"createdDateTime"`
	DisplayName         string      `json:"displayName"`
	Description         string      `json:"description"`
	IsFavoriteByDefault interface{} `json:"isFavoriteByDefault"`
	Email               interface{} `json:"email"`
	WebUrl              interface{} `json:"webUrl"`
	MembershipType      interface{} `json:"membershipType"`
	ModerationSettings  interface{} `json:"moderationSettings"`
}

type TeamsChannel struct {
	ODataContext string    `json:"@odata.context"`
	ODataCount   int       `json:"@odata.count"`
	Channels     []Channel `json:"value"`
}

type Channel struct {
	ODataID             string `json:"@odata.id"`
	ID                  string `json:"id"`
	CreatedDateTime     string `json:"createdDateTime"`
	DisplayName         string `json:"displayName"`
	Description         string `json:"description"`
	IsFavoriteByDefault string `json:"isFavoriteByDefault"`
	Email               string `json:"email"`
	TenantID            string `json:"tenantId"`
	WebURL              string `json:"webUrl"`
	MembershipType      string `json:"membershipType"`
}

func CreateTeamMigrate(accessToken string, teamName string, teamDescription string, createdDateTime string) (string, error) {

	url := "https://graph.microsoft.com/v1.0/teams"

	// Prepare the request payload
	payload := map[string]interface{}{
		"@microsoft.graph.teamCreationMode": "migration",
		"template@odata.bind":               "https://graph.microsoft.com/v1.0/teamsTemplates('standard')",
		"displayName":                       teamName,
		"description":                       teamDescription,
		"createdDateTime":                   createdDateTime,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err)
		return "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}

	defer resp.Body.Close()

	teamCreationResponse := &TeamCreationResponse{
		statusCode:      resp.StatusCode,
		requestID:       resp.Header.Get("Request-Id"),
		clientRequestID: resp.Header.Get("Client-Request-Id"),
		contentLocation: resp.Header.Get("Content-Location"),
		location:        resp.Header.Get("Location"),
	}

	fmt.Println("Team ID:", teamCreationResponse)

	return extractTeamID(teamCreationResponse.contentLocation), nil
}

func extractTeamID(teamURL string) string {
	// Find the starting position of the team ID
	start := strings.Index(teamURL, "('")
	if start == -1 {
		return "" // Team ID not found
	}

	// Find the ending position of the team ID
	end := strings.Index(teamURL, "')")
	if end == -1 {
		return "" // Team ID not found
	}

	// Extract the team ID
	teamID := teamURL[start+2 : end]
	return teamID
}

func CreateChannelMigrate(accessToken string, teamID string, displayName string, description string, createdDateTime string) (int, string, error) {

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%v/channels", teamID)

	// Prepare the request payload
	payload := map[string]interface{}{
		"@microsoft.graph.channelCreationMode": "migration",
		"displayName":                          displayName,
		"description":                          description,
		"membershipType":                       "standard",
		"createdDateTime":                      createdDateTime,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err)
		return 0, "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0, "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 0, "", err
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0, "", err
	}

	defer resp.Body.Close()

	ChannelCreationResponse := &ChannelCreationResponse{
		statusCode:      resp.StatusCode,
		requestID:       resp.Header.Get("Request-Id"),
		clientRequestID: resp.Header.Get("Client-Request-Id"),
		contentLocation: resp.Header.Get("Content-Location"),
		location:        resp.Header.Get("Location"),
	}

	err = json.Unmarshal(body, &ChannelCreationResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return 0, "", err
	}
	//fmt.Println("Response code:", ChannelCreationResponse.statusCode)
	if ChannelCreationResponse.statusCode != 201 {

		fmt.Println("Body:", string(body))
	}
	return ChannelCreationResponse.statusCode, ChannelCreationResponse.Id, nil

}

func PushMessageMigrate(accessToken string, teamID string, channelID string, payload map[string]interface{}) (int, string, error) {

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%v/channels/%v/messages", teamID, channelID)

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON payload:", err)
		return 0, "", err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0, "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return resp.StatusCode, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return resp.StatusCode, fmt.Sprintf("%v - %v\n", resp.Status, string(body)), nil
}

func CompleateChannelMigrate(accessToken string, teamID string, channelID string) {

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%v/channels/%v/completeMigration", teamID, channelID)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// Read the response body
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Error reading response body:", err)
	//	return
	//}

	defer resp.Body.Close()

	//fmt.Println(resp)
	//fmt.Println(string(body))

}

func CompleateTeamMigrate(accessToken string, teamID string) {

	// Close migration on all channels
	listCh, _ := ListChannels(accessToken, teamID)

	for _, channel := range listCh.Channels {
		fmt.Println("ID:", channel.ID)
		fmt.Println("Created Date and Time:", channel.CreatedDateTime)
		fmt.Println("Display Name:", channel.DisplayName)
		fmt.Println("Description:", channel.Description)
		fmt.Println("Tenant ID:", channel.TenantID)
		fmt.Println("Web URL:", channel.WebURL)
		fmt.Println("Membership Type:", channel.MembershipType)
		fmt.Println("--------------------")

		CompleateChannelMigrate(accessToken, teamID, channel.ID)
	}

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%v/completeMigration", teamID)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println(resp)
	fmt.Println(string(body))

}

func ListChannels(accessToken string, teamID string) (TeamsChannel, error) {

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%v/allChannels", teamID)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return TeamsChannel{}, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return TeamsChannel{}, err
	}

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return TeamsChannel{}, err
	}

	defer resp.Body.Close()

	var teamsChannels TeamsChannel
	if err := json.Unmarshal(body, &teamsChannels); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return TeamsChannel{}, err
	}

	return teamsChannels, nil

}
