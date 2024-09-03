package cli

import (
	"binp/storage"
	"binp/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prettyPrint, _ := cmd.Flags().GetBool("pretty-print")
		jsonPrint, _ := cmd.Flags().GetBool("json")
		ID := args[0]

		if prettyPrint && jsonPrint {
			fmt.Fprintln(os.Stderr, "Error: Cannot use both pretty-print and json flags")
			os.Exit(1)
		}

		if prettyPrint {
			if _, err := exec.LookPath("bat"); err != nil {
				fmt.Fprintln(os.Stderr, "Error: Pretty print requires bat. Please install bat.")
				os.Exit(1)
			}
		}

		resp, err := util.HTTPGet(fmt.Sprintf("http://localhost:8080/%s", ID))
		defer resp.Body.Close()

		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if resp.StatusCode != 200 {
			if resp.StatusCode == 404 {
				fmt.Fprintln(os.Stderr, "Error: Snippet not found")
			} else {
				fmt.Fprintln(os.Stderr, "Error: ", string(resBody))
			}
			os.Exit(1)
		}

		if jsonPrint {
			fmt.Println(string(resBody))
			os.Exit(0)
		}

		snippet := &storage.Snippet{}
		err = json.Unmarshal(resBody, snippet)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			os.Exit(1)
		}

		if !prettyPrint {
			fmt.Println(snippet.Text)
			os.Exit(0)
		}

		batCmd := exec.Command("bat", "-l", snippet.Language, "--file-name", snippet.ID)
		batCmd.Stdin = bytes.NewBufferString(snippet.Text)
		batCmd.Stdout = os.Stdout
		batCmd.Stderr = os.Stderr

		err = batCmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	getCmd.Flags().BoolP("pretty-print", "p", false, "Pretty print snippet (requires bat)")
	getCmd.Flags().BoolP("json", "j", false, "Print snippet as JSON")
	rootCmd.AddCommand(getCmd)
}
