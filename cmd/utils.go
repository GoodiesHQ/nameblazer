package main

import (
	"context"
	"net"
	"time"
)

func uniq(values []string) []string {
	seen := make(map[string]struct{})

	var unique []string

	for _, value := range values {
		if _, found := seen[value]; !found {
			seen[value] = struct{}{}
			unique = append(unique, value)
		}
	}

	return unique
}

func combiner(domains []string, subdomains []string) <-chan string {
	ch := make(chan string)
	go func() {
		for _, domain := range domains {
			for _, subdomain := range subdomains {
				combined := subdomain + "." + domain
				ch <- combined
			}
		}
		close(ch)
	}()
	return ch
}

func infiniterator[T any](values []T) <-chan T {
	ch := make(chan T)
	go func() {
		for {
			for _, value := range values {
				ch <- value
			}
		}
	}()
	return ch
}

func makeResolvers(addresses []string, timeout time.Duration) []*net.Resolver {
	var resolvers []*net.Resolver
	for _, address := range addresses {
		dialer := &net.Dialer{
			Timeout: timeout,
		}
		resolvers = append(resolvers, &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return dialer.DialContext(ctx, "udp", address+":53")
			},
		})
	}
	return resolvers
}

func lookup(resolver *net.Resolver, host string, ipv4, ipv6 bool, timeout time.Duration) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	network := "ip"

	if ipv4 && !ipv6 {
		network = "ip4"
	}

	if !ipv4 && ipv6 {
		network = "ip6"
	}

	ips, err := resolver.LookupIP(ctx, network, host)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, ip := range ips {
		addrs = append(addrs, ip.String())
	}

	return addrs, nil
}
