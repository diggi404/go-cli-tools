package main

import (
	"fmt"
	"go_cli/menu"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	fmt.Print("\033[H\033[2J")
	art := `
	███    ██████████████████████     ███████████████████    ████████    ████████████████ 
	████   ███       ██  ██   ███    ██   ███       ██       ██    ██    ███  ██  ██      
	██ ██  ██████    ██  ████████    █████████████  ██       ████████    ███  ██  █████   
	██  ██ ███       ██  ██   ███    ██   ██    ██  ██            ███    ███  ██  ██      
	██   ██████████  ██  ██████████████   ████████  ██       ███████████████  ██  ███████ 	
	
				[x] Created by @realdiggi [x]
	 `
	menu.SlowPrintArt(art, time.Millisecond*50)

	menuOpts := `
	1. Bulk Range IP Generator			2. Mass IP Scanner

	3. Mass Mailer					4. Email Bomber

	5. SMTP Checker					6. CPanel Checker

	7. Exit

	`
	menu.SlowPrintMenu(menuOpts, time.Millisecond*50)
	var choiceStr string
	// fmt.Print()
	color.New(color.FgHiBlue).Print("\n\nEnter your option :> ")
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println()
		color.New(color.FgRed).Println("invalid choice. Exiting Program...")
		return
	}
	menu.MenuSelection(choice)
}
