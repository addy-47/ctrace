package utils

import (
	"fmt"
)

// ANSI Color Codes
const (
	ColorCyan  = "\033[36m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
)

// PrintBanner prints the CTRACE ASCII art banner.
func PrintBanner() {
	banner := `
   ______   ______     ____     ___     ______   ______
  / ____/  /_  __/    / __ \   /   |   / ____/  / ____/
 / /        / /      / /_/ /  / /| |  / /      / __/   
/ /___     / /      / _, _/  / ___ | / /___   / /___   
\____/    /_/      /_/ |_|  /_/  |_| \____/  /_____/   
`
	// Applying Cyan color to the banner
	fmt.Println(ColorCyan + banner + ColorReset)
	fmt.Println(ColorBlue + "      Static Code Lifecycle Tracer" + ColorReset)
	fmt.Println()
}
