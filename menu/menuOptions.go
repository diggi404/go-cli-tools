package menu

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

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

func MenuSelection(selectedOption int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP Address", "Open Ports"})
	if selectedOption == 0 {
		filePath, err := GenIP()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			AfterGenIP("", "", table)
		}
		var input string
		fmt.Print("Do you want to scan these IPs now? Y/n :> ")
		fmt.Scanln(&input)
		AfterGenIP(input, filePath, table)
	} else if selectedOption == 1 {
		ScanIPs(table)
	} else if selectedOption == 2 {
		SmtpCrack()
	}
	os.Exit(0)
}

func AfterGenIP(choice, filePath string, table *tablewriter.Table) {
	choice = strings.ToLower(choice)
	if choice == "yes" || choice == "y" {
		ScanIPs(table, filePath)
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
		ScanIPs(table)
	case 3:
		SmtpCrack()
	case 4:
		os.Exit(0)
	}
}
