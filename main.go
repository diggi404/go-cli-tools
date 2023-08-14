package main

import (
	"fmt"
	"go_cli/menu"
	"log"
	"os"

	"github.com/fatih/color"

	"github.com/eiannone/keyboard"
)

func main() {
	options := []string{"Bulk Range IP Generator", "Mass IP Scanner", "SMTP Cracker", "Exit"}
	selectedIndex := 0
	var selectedOption int

	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

	for {
		printMenu(options, selectedIndex)

		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEnter {
			if selectedIndex == len(options)-1 {
				fmt.Println("Exiting Program...")
				os.Exit(0)
			} else {
				selectedOption = selectedIndex
				keyboard.Close()
				break
			}
		} else if key == keyboard.KeyArrowUp {
			selectedIndex = (selectedIndex - 1 + len(options)) % len(options)
		} else if key == keyboard.KeyArrowDown {
			selectedIndex = (selectedIndex + 1) % len(options)
		}
	}
	if selectedOption == 0 {
		menu.GenIP()
		os.Exit(0)
	}
	fmt.Printf("selectedOption: %v\n", selectedOption)

}

func printMenu(items []string, selectedIndex int) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Select an option using the arrow keys (Up/Down) and press Enter:")
	red := color.New(color.BgRed).PrintfFunc()
	for i, item := range items {
		if i == selectedIndex {
			red("* %s\n", item)
		} else {
			fmt.Printf("  %s\n", item)
		}
	}
}
