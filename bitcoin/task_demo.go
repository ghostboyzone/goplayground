package main

import (
	"github.com/fatih/color"
	"log"
)

func main() {

	color.Set(color.FgRed, color.Bold)
	color.Set(color.FgYellow, color.Bold)
	color.Set(color.FgGreen, color.Bold)
	// color.Unset()
	log.Println("123")
}
