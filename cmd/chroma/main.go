package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	style := styles.Get("tokyonight-night")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))

	var buffer strings.Builder
	err := formatter.WriteCSS(&buffer, style)
	if err != nil {
		panic(err)
	}

	css := buffer.String()

	re := regexp.MustCompile(`background-color:\s*[^;]+;`)
	css = re.ReplaceAllString(css, "")

	re = regexp.MustCompile(`[^\{\}]+\{\s*\}`)
	css = re.ReplaceAllString(css, "")

	css = strings.ReplaceAll(css, "  ", " ")

	err = os.WriteFile("static/css/chroma.css", []byte(css), 0644)
	if err != nil {
		panic(err)
	}
}
