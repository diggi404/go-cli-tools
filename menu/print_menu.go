package menu

import (
	"strings"
	"time"

	"github.com/fatih/color"
)

func SlowPrintArt(text string, delay time.Duration) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		color.New(color.FgHiRed).Println(line)
		time.Sleep(delay)
	}
}
func SlowPrintMenu(text string, delay time.Duration) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		color.New(color.FgHiGreen).Println(line)
		time.Sleep(delay)
	}
}
