package utils

import (
	"fmt"
)

func Print(v ...any) {
	for _, val := range v {
		fmt.Println(val)
	}
}

func PrintErr(v ...any) {
	for _, val := range v {
		fmt.Println(Red + fmt.Sprint(val) + Reset)
	}
}

func PrintSuccess(v ...any) {
	for _, val := range v {
		fmt.Println(Green + fmt.Sprint(val) + Reset)
	}
}

func PrintWarning(v ...any) {
	for _, val := range v {
		fmt.Println(Yellow + fmt.Sprint(val) + Reset)
	}
}

func PrintInfo(v ...any) {
	for _, val := range v {
		fmt.Println(Blue + fmt.Sprint(val) + Reset)
	}
}
