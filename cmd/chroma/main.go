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

	chromaRule := regexp.MustCompile(`\.chroma\s*\{[^\}]*\}`)
	if chromaRule.MatchString(css) {
		css = chromaRule.ReplaceAllStringFunc(css, func(match string) string {
			if strings.Contains(match, "white-space:") {
				return strings.Replace(match, "white-space: nowrap;", "white-space: pre-wrap;", 1)
			}
			return strings.TrimSuffix(match, "}") + " white-space: pre-wrap;}"
		})
	} else {
		css += "\n.chroma { white-space: pre-wrap; }"
	}

	// Add word-wrap: break-word to .chroma class
	if chromaRule.MatchString(css) {
		css = chromaRule.ReplaceAllStringFunc(css, func(match string) string {
			if !strings.Contains(match, "word-wrap:") {
				return strings.TrimSuffix(match, "}") + " word-wrap: break-word;}"
			}
			return match
		})
	} else {
		css += "\n.chroma { word-wrap: break-word; }"
	}

	err = os.WriteFile("static/css/chroma.css", []byte(css), 0644)
	if err != nil {
		panic(err)
	}
}
