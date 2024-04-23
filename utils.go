package main

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens, nil
}
