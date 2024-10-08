package cli

import (
	"bytes"
	"net/http"

	"github.com/spf13/cobra"
)

var client = &http.Client{}

func HTTPGet(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func HTTPPost(url string, body *bytes.Buffer) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

var rootCmd = &cobra.Command{
	Use:   "binp",
	Short: "A cli tool for the binp pastebin service",
	Long:  "A cli tool for the binp pastebin service",
}

func Execute() error {
	return rootCmd.Execute()
}
