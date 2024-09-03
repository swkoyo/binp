package cli

import (
	"binp/server"
	"binp/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		language, _ := cmd.Flags().GetString("language")
		expiry, _ := cmd.Flags().GetString("expiry")
		burnAfterRead, _ := cmd.Flags().GetBool("burn")
		text := args[0]

		snippetBody := &server.PostSnippetReq{
			Text:          text,
			BurnAfterRead: burnAfterRead,
			Expiry:        expiry,
			Language:      language,
		}

		client := &http.Client{}

		postBody, err := json.Marshal(snippetBody)
		req := bytes.NewBuffer(postBody)

		request, err := http.NewRequest("POST", "http://localhost:8080/snippet", req)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		resp, err := client.Do(request)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer resp.Body.Close()
		resBody, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error: ", err)
		}

		createdSnippet := &storage.Snippet{}

		err = json.Unmarshal(resBody, createdSnippet)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println(createdSnippet.ID)
	},
}

func init() {
	createCmd.Flags().StringP("language", "l", "plaintext", "The language of the snippet")
	createCmd.Flags().StringP("expiry", "e", "one_hour", "The expiry time of the snippet (one_hour, one_day, one_week, one_month)")
	createCmd.Flags().BoolP("burn-after-read", "b", false, "Burn the snippet after reading it once")
	rootCmd.AddCommand(createCmd)
}
