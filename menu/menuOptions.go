package menu

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func PrintMenu(items []string, selectedIndex int) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Select an option using the arrow keys (Up/Down) and press Enter:")
	red := color.New(color.BgRed).PrintfFunc()
	for i, item := range items {
		if i == selectedIndex {
			red("> %s\n", item)
		} else {
			fmt.Printf("%d  %s\n", i+1, item)
		}
	}
}

func Selections(choice string) {
	choice = strings.ToLower(choice)
	if choice == "yes" || choice == "y" {
		ScanIPs()
		os.Exit(0)
	}
	var selectedOption int
	fmt.Print("Enter a menu option :> ")
	fmt.Scanln(&selectedOption)
	switch selectedOption {
	case 1:
		GenIP()
	case 2:
		ScanIPs()
	case 3:
		SmtpCrack()
	case 4:
		os.Exit(0)
	}
}
