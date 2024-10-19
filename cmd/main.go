package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	settings := flags()

	domains, err := readDomains(settings.DomainsFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read the domains file")
	}

	if len(domains) == 0 {
		log.Fatal().Msg("found no valid domains")
	}

	subdomains, err := readLinesUnique(settings.SubdomainsFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read the domains file")
	}

	if len(subdomains) == 0 {
		log.Fatal().Msg("found no valid subdomains")
	}

	log.Info().Msgf("Found %d domains, resolving %d subdomains each", len(domains), len(subdomains))

	resolvers := makeResolvers(settings.Nameservers, settings.Timeout)
	resolverChan := infiniterator(resolvers)

	var results []NameBlazerResult

	for host := range combiner(domains, subdomains) {
		resolver := <-resolverChan
		addrs, err := lookup(resolver, host, settings.IPv4, settings.IPv6, settings.Timeout)
		if err != nil {
			if !settings.Quiet {
				log.Warn().Err(err).Msgf("failed to resolve '%s' with resolver ", host)
			}
			continue
		}

		results = append(results, NameBlazerResult{
			Host:      host,
			Addresses: addrs,
		})
	}

	switch settings.OutputFormat {
	case FORMAT_JSON:
		writeJson(results, settings.OutputFile)
	case FORMAT_YAML:
		writeYaml(results, settings.OutputFile)
	case FORMAT_CSV:
		var lines [][]string
		for _, result := range results {
			var line []string
			line = append(line, result.Host)
			line = append(line, result.Addresses...)
			lines = append(lines, line)
		}
		writeCsv(lines, settings.OutputFile)
	}
}
