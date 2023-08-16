package menu

import (
	"fmt"
	menu "go_cli/menu/scanIPs"
	"os"
	"strings"

	"github.com/fatih/color"
)

// PrintMenu Main Menu Function... CLears the terminal before printing the Menu.
func PrintMenu(items []string, selectedIndex int) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Select an option using the arrow keys (Up/Down) and press Enter:")
	red := color.New(color.BgRed).PrintfFunc()
	for i, item := range items {
		if i == selectedIndex {
			red("* %s\n", item)
		} else {
			fmt.Printf("%d  %s\n", i+1, item)
		}
	}
}

// MenuSelection Handles user's selection from the Main Menu i.e. When the app is launched.
func MenuSelection(selectedOption int) {
	if selectedOption == 0 {
		filePath, err := GenIP()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			AfterGenIP("", "")
		}
		var input string
		fmt.Print("Do you want to scan these IPs now? Y/n :> ")
		fmt.Scanln(&input)
		AfterGenIP(input, filePath)
	} else if selectedOption == 1 {
		menu.ScanIPs()
	} else if selectedOption == 2 {
		SmtpCrack()
	}
	os.Exit(0)
}

// AfterGenIP This function handles user selection right after generating bulk IPs.
func AfterGenIP(choice, filePath string) {
	choice = strings.ToLower(choice)
	if choice == "yes" || choice == "y" {
		menu.ScanIPs(filePath)
		os.Exit(0)
	}
	fmt.Print("\n")
	var selectedOption int
	fmt.Print("Enter a menu option :> ")
	fmt.Scanln(&selectedOption)
	switch selectedOption {
	case 1:
		GenIP()
	case 2:
		menu.ScanIPs()
	case 3:
		SmtpCrack()
	case 4:
		os.Exit(0)
	}
}
