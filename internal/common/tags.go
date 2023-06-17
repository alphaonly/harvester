package common

import "fmt"

func PrintBuildTags(version string, date string, commit string) {
	fmt.Printf("Build version: %v\n", version)
	fmt.Printf("Build date: %v\n", date)
	fmt.Printf("Build commit: %v\n", commit)
}
