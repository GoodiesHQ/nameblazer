package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/publicsuffix"
	"gopkg.in/yaml.v2"
)

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" { // ignore empty lines
			lines = append(lines, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func readLinesUnique(filename string) ([]string, error) {
	lines, err := readLines(filename)
	if err != nil {
		return nil, err
	}

	return uniq(lines), nil
}

func readDomains(filename string) ([]string, error) {
	lines, err := readLinesUnique(filename)
	if err != nil {
		return nil, err
	}

	var baseDomains []string

	for _, line := range lines {
		etld, err := publicsuffix.EffectiveTLDPlusOne(line)
		if err != nil {
			log.Warn().Msgf("unable to parse the eTLD from '%s'", line)
			continue
		}
		baseDomains = append(baseDomains, etld)
	}

	return uniq(baseDomains), nil
}

func writeJson(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: formats JSON with indentation.
	return encoder.Encode(data)
}

func writeYaml(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(data)
}

func writeCsv(data [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	return writer.WriteAll(data)
}
