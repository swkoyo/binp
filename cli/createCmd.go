package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type PostSnippetReq struct {
	Text          string `json:"text"`
	BurnAfterRead bool   `json:"burn_after_read"`
	Language      string `json:"language"`
	Expiry        string `json:"expiry"`
}

type Snippet struct {
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	BurnAfterRead bool      `json:"burn_after_read"`
	IsRead        bool      `json:"is_read"`
	Language      string    `json:"language"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		language, _ := cmd.Flags().GetString("language")
		expiry, _ := cmd.Flags().GetString("expiry")
		burnAfterRead, _ := cmd.Flags().GetBool("burn")
		text := args[0]

		snippetBody := &PostSnippetReq{
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
		resp, err := HTTPPost("http://localhost:8080/snippet", req)
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

		createdSnippet := &Snippet{}
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
	createCmd.Flags().StringP("language", "l", "txt", "The language of the snippet")
	createCmd.Flags().StringP("expiry", "e", "1m", "The expiry time of the snippet. Valid values: %v")
	createCmd.Flags().BoolP("burn-after-read", "b", false, "Burn the snippet after reading it once")
	rootCmd.AddCommand(createCmd)
}
