package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	pPath := flag.String("env", "env.json", "Path to .env file")
	flag.Parse()

	requiredKeys, err := parseEnvFile(*pPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	missingKeys := checkEnvVars(requiredKeys)
	if len(missingKeys) != 0 {
		msg := fmt.Sprintf("ERROR: some variables are undefined: %d\n", len(missingKeys))
		for i, key := range missingKeys {
			msg = fmt.Sprintf("%s%d. %v\n", msg, i+1, key)
		}
		fmt.Println(msg)
		os.Exit(1)
	}

	fmt.Printf("All environment variables found: %d\n", len(requiredKeys))
}

func parseEnvFile(path string) ([]string, error) {
	vars := []string{}
	bs, err := os.ReadFile(path)
	if err != nil {
		return vars, fmt.Errorf("error reading env file: %v", err)
	}

	err = json.Unmarshal(bs, &vars)
	if err != nil {
		return vars, fmt.Errorf("error pargins env file: %v", err)
	}

	return vars, nil
}

func checkEnvVars(keys []string) []string {
	missingKeys := []string{}
	for _, key := range keys {
		_, exists := os.LookupEnv(key)
		if !exists {
			missingKeys = append(missingKeys, key)
		}
	}

	return missingKeys
}
