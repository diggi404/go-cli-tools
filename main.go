package main

import (
	"fmt"
	"go_cli/menu"
	"log"
	"os"

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
		menu.PrintMenu(options, selectedIndex)

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
		var scanOption string
		fmt.Print("Do you want to scan these IPs now? Y/n :> ")
		fmt.Scanln(&scanOption)
		menu.Selections(scanOption)
	}
	fmt.Printf("selectedOption: %v\n", selectedOption)
}
