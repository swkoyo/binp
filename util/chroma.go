package util

import (
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var logger = GetLogger()

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
