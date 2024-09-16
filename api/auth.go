package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GET "/api/key"
func GetAPIkey(username, password string) (string, error) {
	request := fmt.Sprintf("https://%s:%s@neocities.org/api/key", username, password)
	client := http.Client{Timeout: 10 * time.Second}
	reply, err := client.Get(request)
	if err != nil {
		fmt.Println("Get Error")
		return "", err
	}
	defer reply.Body.Close()
	var response Response
	err = json.NewDecoder(reply.Body).Decode(&response)
	if err != nil {
		fmt.Println("Decode Error")
		body := reply.Body
		buf := make([]byte, 1024)
		_, err := body.Read(buf)
		fmt.Printf("Reply: %s\n", buf)
		return "", err
	}
	if response.APIKey == "" {
		return "", fmt.Errorf("invalid credentials")
	}
	return response.APIKey, nil
}
