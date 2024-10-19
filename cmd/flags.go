package main

import (
	"flag"
	"log"
	"strings"
	"time"
)

type stringlist []string

func (l *stringlist) String() string {
	return strings.Join(*l, ", ")
}

func (l *stringlist) Set(val string) error {
	*l = append(*l, val)
	return nil
}

func flags() *NameBlazerSettings {
	workers := flag.Int("workers", 10, "Number of concurrent workers")
	flag.IntVar(workers, "w", 10, "Number of concurrent workers (alias)")

	timeout := flag.Duration("timeout", 500*time.Millisecond, "DNS Timeout")
	flag.DurationVar(timeout, "t", 500*time.Millisecond, "DNS timeout (alias)")

	domainsFile := flag.String("domains", "", "File containing the base domain names")
	flag.StringVar(domainsFile, "d", "", "File containing the base domain names (alias)")

	subdomainsFile := flag.String("subdomains", "", "File containing the subdomains to append to each domain")
	flag.StringVar(subdomainsFile, "s", "", "File containing the subdomains to append to each domain (alias)")

	outputFile := flag.String("output", "", "File to write the output (.json, .yaml/.yml, .csv, or .txt)")
	flag.StringVar(outputFile, "o", "", "File to write the output (.json, .yaml/.yml, .csv, or .txt) (alias)")

	quiet := flag.Bool("quiet", false, "Suppress warning messages while resolving")
	flag.BoolVar(quiet, "q", false, "Suppress warning messages while resolving (alias)")

	ipv4 := flag.Bool("ipv4", false, "Resolve IPv4")
	flag.BoolVar(ipv4, "v4", false, "Resolve IPv4 (alias)")

	ipv6 := flag.Bool("ipv6", false, "Resolve IPv6")
	flag.BoolVar(ipv6, "v6", false, "Resolve IPv6 (alias)")

	var nameservers stringlist

	flag.Var(&nameservers, "nameservers", "DNS Nameserver IP")
	flag.Var(&nameservers, "ns", "DNS Nameserver IP")

	flag.Parse()

	useIPv4 := *ipv4
	useIPv6 := *ipv6

	// if nether is specified, use both IPv4 (A) and IPv6 (AAAA)
	if !useIPv4 && !useIPv6 {
		useIPv4 = true
		useIPv6 = true
	}

	if *domainsFile == "" {
		log.Fatal("-domains/-d file is not provided")
	}

	if *subdomainsFile == "" {
		log.Fatal("-subdomains/-s file is not provided")
	}

	if *outputFile == "" {
		log.Fatal("-output/-o file is not provided")
	}

	outputFileLower := strings.ToLower(*outputFile)
	var outputFormat format

	// detect the proper format to use
	if strings.HasSuffix(outputFileLower, ".json") {
		outputFormat = FORMAT_JSON
	} else if strings.HasSuffix(outputFileLower, ".yaml") || strings.HasSuffix(outputFileLower, ".yml") {
		outputFormat = FORMAT_YAML
	} else {
		outputFormat = FORMAT_CSV
	}

	if len(nameservers) == 0 {
		nameservers = append(nameservers, "1.1.1.1", "1.0.0.1", "8.8.8.8", "8.8.4.4")
	}

	return &NameBlazerSettings{
		Workers:        *workers,
		Timeout:        *timeout,
		DomainsFile:    *domainsFile,
		SubdomainsFile: *subdomainsFile,
		Nameservers:    nameservers,
		Quiet:          *quiet,
		IPv4:           useIPv4,
		IPv6:           useIPv6,
		OutputFile:     *outputFile,
		OutputFormat:   outputFormat,
	}
}
