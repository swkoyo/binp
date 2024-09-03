package cli

import (
	"binp/server"
	"binp/storage"
	"binp/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

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

		if !storage.IsValidExpiration(expiry) {
			fmt.Fprintln(os.Stderr, "Error: Invalid expiry. Valid values are 1h, 1d, 1w, 1m")
			os.Exit(1)
		}

		if !storage.IsValidLanguage(language) {
			fmt.Fprintln(os.Stderr, "Error: Invalid language. Valid values are plaintext, bash, css, docker, go, html, javascript, json, markdown, python, ruby, typescript")
			os.Exit(1)
		}

		snippetBody := &server.PostSnippetReq{
			Text:          text,
			BurnAfterRead: burnAfterRead,
			Expiry:        expiry,
			Language:      language,
		}

		postBody, err := json.Marshal(snippetBody)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			os.Exit(1)
		}
		req := bytes.NewBuffer(postBody)
		resp, err := util.HTTPPost("http://localhost:8080/snippet", req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			os.Exit(1)
		}

		createdSnippet := &storage.Snippet{}
		err = json.Unmarshal(resBody, createdSnippet)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			os.Exit(1)
		}

		fmt.Println(createdSnippet.ID)
		os.Exit(0)
	},
}

func init() {
	createCmd.Flags().StringP("language", "l", "plaintext", "The language of the snippet (plaintext, bash, css, docker, go, html, javascript, json, markdown, python, ruby, typescript)")
	createCmd.Flags().StringP("expiry", "e", "1h", "The expiry time of the snippet (1h, 1d, 1w, 1m)")
	createCmd.Flags().BoolP("burn-after-read", "b", false, "Burn the snippet after reading it once")
	rootCmd.AddCommand(createCmd)
}
