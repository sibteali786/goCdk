package main

import "fmt"

type MyEvent struct {
	Username string `json:"username"` // its called tagging
}

// Take in a payload and do something wit it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty ")
	}

	return fmt.Sprintf("Scucessfully called by - %s\n", event.Username), nil
}
func main() {

}
