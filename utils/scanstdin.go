package utils

import (
	"bufio"
	"os"
)

// scan StdIn and send each line to the apiUrl channel

func ScanStdIn(apiUrl chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		apiUrl <- scanner.Text()
	}
	close(apiUrl)
}
