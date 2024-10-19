package main

import "time"

type format int

const (
	FORMAT_CSV  format = iota // CSV output format = sub.example.com,1.2.3.4
	FORMAT_JSON format = iota // JSON output format = {"host": "sub.example.com", "addresses": ["1.2.3.4"]}
	FORMAT_YAML format = iota // YAML output format = host: sub.example.com, addresses: [1.2.3.4]
)

type NameBlazerSettings struct {
	Workers        int
	Timeout        time.Duration
	DomainsFile    string
	SubdomainsFile string
	Nameservers    []string
	Quiet          bool
	IPv4           bool
	IPv6           bool
	OutputFile     string
	OutputFormat   format
}

type NameBlazerResult struct {
	Host      string   `json:"host" yaml:"host"`
	Addresses []string `json:"addresses" yaml:"addresses"`
}
