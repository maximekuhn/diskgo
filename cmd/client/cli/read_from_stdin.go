package cli

import (
	"bufio"
	"os"
)

func ReadFromStdin(outCh chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		outCh <- input
	}
}
