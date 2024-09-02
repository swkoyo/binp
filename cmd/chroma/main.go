package main

import (
	"binp/util"
)

func main() {
	util.InitLogger()
	logger := util.GetLogger()

	err := util.GenerateChromaCSS()

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate chroma.css")
	}
}
