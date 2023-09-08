package main

import (
	"fmt"
	"go_cli/menu"
	"strconv"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	art := `
	_____  ___   _______ ___________ _______  ___           __       ________ ___________       ________ ____  ____  __ ___________ _______  
	(\"   \|"  \ /"     "("     _   "|   _  "\|"  |         /""\     /"       ("     _   ")     /"       ("  _||_ " ||" ("     _   "/"     "| 
	|.\\   \    (: ______))__/  \\__/(. |_)  :||  |        /    \   (:   \___/ )__/  \\__/     (:   \___/|   (  ) : |||  )__/  \\__(: ______) 
	|: \.   \\  |\/    |     \\_ /   |:     \/|:  |       /' /\  \   \___  \      \\_ /         \___  \  (:  |  | . )|:  |  \\_ /   \/    |   
	|.  \    \. |// ___)_    |.  |   (|  _  \\ \  |___   //  __'  \   __/  \\     |.  |          __/  \\  \\ \__/ // |.  |  |.  |   // ___)_  
	|    \    \ (:      "|   \:  |   |: |_)  :( \_|:  \ /   /  \\  \ /" \   :)    \:  |         /" \   :) /\\ __ //\ /\  |\ \:  |  (:      "| 
	 \___|\____\)\_______)    \__|   (_______/ \_______(___/    \___(_______/      \__|        (_______/ (__________(__\_|_) \__|   \_______) 																																																																							
	 `
	color.New(color.FgHiRed).Println(art)

	menuOpts := `
	1.	Bulk Range IP Generator				2.	Mass IP Scanner

	3.	Mass Mailer					4.	Email Bomber

	5.	SMTP Cracker					6.	CPanel Cracker

	7.	Exit

	`
	color.New(color.FgGreen).Println(menuOpts)
	var choiceStr string
	fmt.Print("\n\nEnter your option :> ")
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("invalid choice!")
		return
	}
	menu.MenuSelection(choice)
}
