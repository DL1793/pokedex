package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		fields := cleanInput(text)
		if len(fields) > 0 {
			fmt.Println("Your command was: " + fields[0])
		}
	}
}
