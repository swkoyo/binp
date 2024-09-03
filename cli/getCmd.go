package cli

import (
	"binp/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ID := args[0]

		client := &http.Client{}

		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/%s", ID), nil)
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

		snippet := &storage.Snippet{}

		jsonStr := string(resBody)

		err = json.Unmarshal(resBody, snippet)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		prettyPrint, _ := cmd.Flags().GetBool("pretty-print")
		jsonPrint, _ := cmd.Flags().GetBool("json")

		if jsonPrint {
			fmt.Printf("%+v\n", jsonStr)
		} else if !prettyPrint {
			fmt.Println(snippet.Text)
		} else {
			_, err = exec.LookPath("bat")
			if err != nil {
				fmt.Println("Pretty print requires bat. Please install bat.")
				return
			}
			batCmd := exec.Command("bat", "-l", snippet.Language)
			batCmd.Stdin = bytes.NewBufferString(snippet.Text)
			batCmd.Stdout = os.Stdout
			batCmd.Stderr = os.Stderr

			err = batCmd.Run()
			if err != nil {
				fmt.Println("BAT Error: ", err)
				return
			}
		}
	},
}

func init() {
	getCmd.Flags().BoolP("pretty-print", "p", false, "Pretty print snippet (requires bat)")
	getCmd.Flags().BoolP("json", "j", false, "Print snippet as JSON")
	rootCmd.AddCommand(getCmd)
}
