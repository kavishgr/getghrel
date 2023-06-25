package utils

import (
	"bufio"
	// "log"
	"os"
)

func ScanStdIn(apiUrl chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		apiUrl <- scanner.Text()
	}
	close(apiUrl)
}
