package rainbow

import (
	"fmt"
)

var colors = []string{
	"\033[38;5;196m",
	"\033[38;5;202m",
	"\033[38;5;208m",
	"\033[38;5;214m",
	"\033[38;5;220m",
	"\033[38;5;226m",
	"\033[38;5;190m",
	"\033[38;5;154m",
	"\033[38;5;118m",
	"\033[38;5;82m",
	"\033[38;5;46m",
	"\033[38;5;51m",
	"\033[38;5;45m",
	"\033[38;5;39m",
	"\033[38;5;33m",
	"\033[38;5;27m",
	"\033[38;5;93m",
	"\033[38;5;99m",
	"\033[38;5;105m",
	"\033[38;5;141m",
	"\033[38;5;147m",
	"\033[38;5;183m",
	"\033[38;5;219m",
}

var LastColorI int = 0

func Color(input string) string {
	var result string
	var colorI int
	if LastColorI != 0 {
		colorI = LastColorI
	} else {
		colorI = 0
	}
	for _, char := range input {
		if char != ' ' {
			color := colors[colorI%len(colors)]
			result += fmt.Sprintf("%s%c", color, char)
			colorI++
		} else {
			result += " "
		}
	}
	LastColorI = colorI
	result += "\033[0m"
	return result
}
