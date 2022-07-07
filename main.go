package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	pPath := flag.String("path", ".env", "Path to .env file")
	flag.Parse()

	keys, err := parseEnvFile(*pPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	errors := checkEnvVars(keys)
	if len(errors) != 0 {
		fmt.Printf("Errors found: %d\n", len(errors))
		for i, err := range errors {
			fmt.Printf("%d. %v\n", i+1, err)
		}
		os.Exit(1)
	}

	fmt.Printf("All %d env vars found: OK\n", len(keys))
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parseEnvFile(path string) ([]string, error) {
	vars := []string{}
	lines, err := readLines(path)
	if err != nil {
		return vars, fmt.Errorf("error reading env file: %v", err)
	}

	re, err := regexp.Compile(`(\w+)=?`)
	if err != nil {
		return vars, fmt.Errorf("error parsing env regex: %v", err)
	}

	for i, line := range lines {
		if line == "" {
			continue
		}

		lineNumber := i + 1

		if !re.MatchString(line) {
			return vars, fmt.Errorf("error parsing env, line %d is invalid: '%s'", lineNumber, line)
		}

		matches := re.FindStringSubmatch(line)
		minMatches := 2
		if len(matches) < minMatches {
			return vars, fmt.Errorf(
				"error parsing env, line %d has not enough matches: '%s'",
				lineNumber,
				line,
			)
		}

		varName := matches[1]
		vars = append(vars, varName)
	}

	return vars, nil
}

func checkEnvVars(keys []string) []error {
	errors := []error{}
	for _, key := range keys {
		_, exists := os.LookupEnv(key)
		if !exists {
			errors = append(errors, fmt.Errorf("env var '%s' not found", key))
		}
	}

	return errors
}
