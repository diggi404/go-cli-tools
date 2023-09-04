package bomber

import (
	"os"

	"github.com/schollz/progressbar/v3"
)

func MakePgBar(numBombs int, description string) *progressbar.ProgressBar {
	pgBar := progressbar.NewOptions(numBombs,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	return pgBar
}
