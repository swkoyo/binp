package util

import (
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var logger = GetLogger()

func GenerateChromaCSS() error {
	logger.Info().Msg("Generating chroma.css...")

	style := styles.Get("tokyonight-night")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))

	var buffer strings.Builder
	err := formatter.WriteCSS(&buffer, style)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to generate chroma.css")
		return err
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
		logger.Error().Err(err).Msg("Failed to write chroma.css")
		return err
	}

	logger.Info().Msg("Generated chroma.css")
	return nil
}

func HighlightCode(code, language string) (string, error) {
	logger.Debug().Str("language", language).Msg("Highlighting code")
	lexer := lexers.Get(language)
	if lexer == nil {
		logger.Warn().Str("language", language).Msg("No lexer found for language. Using fallback lexer.")
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	style := styles.Get("tokyonight-night")
	if style == nil {
		logger.Warn().Msg("No style found. Using fallback style.")
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true))
	var buf strings.Builder
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
