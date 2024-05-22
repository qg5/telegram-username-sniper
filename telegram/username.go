package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	FragmentBaseURL  = "https://fragment.com/"
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0"
)

// Checks if a username is available.
func IsUsernameAvailable(username string) bool {
	status, err := getUser(username)
	if err != nil {
		return false
	}

	return status == "Available"
}

func getUser(username string) (string, error) {
	req, err := http.NewRequest("GET", FragmentBaseURL+"username/"+username, nil)
	if err != nil {
		return "", err
	}
	setHeaders(req, username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return processResponse(resp)
}

func setHeaders(req *http.Request, username string) {
	referUrl := fmt.Sprintf("%s?query=%s", FragmentBaseURL, username)

	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("X-Aj-Referer", referUrl)
	req.Header.Set("Referer", referUrl)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Priority", "u=1")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
}

func processResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if rData, ok := response["r"].(string); ok && rData == "/" {
		return "Ratelimit", nil
	}

	hData, ok := response["h"].(string)
	if !ok {
		return "Available", nil
	}

	status := strings.Split(hData, "tm-section-header-status")[1]
	status = strings.Split(status, `">`)[0]
	status = strings.TrimSpace(status)

	statusMapping := map[string]string{
		"tm-status-taken":   "Taken",
		"tm-status-avail":   "Auctioned or for sale",
		"tm-status-unavail": "Sold",
	}

	if mappedStatus, exists := statusMapping[status]; exists {
		return mappedStatus, nil
	}

	return "Unknown", nil
}
