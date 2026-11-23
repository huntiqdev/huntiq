package handlers

import (
	"net"
	"testing"
)

func TestDisallowedIPs(t *testing.T) {
	disallowed := []string{
		"127.0.0.1",
		"127.0.0.2",
		"10.0.0.1",
		"10.254.0.1",
		"172.16.0.1",
		"172.31.255.255",
		"192.168.1.1",
		"169.254.169.254", // AWS/GCP metadata IP
		"0.0.0.0",
		"::1",
		"fe80::1",
	}

	for _, ipStr := range disallowed {
		ip := net.ParseIP(ipStr)
		if !isDisallowedIP(ip) {
			t.Errorf("Expected IP %s to be disallowed, but was allowed", ipStr)
		}
	}

	allowed := []string{
		"8.8.8.8",
		"1.1.1.1",
		"93.184.216.34", // example.com
	}

	for _, ipStr := range allowed {
		ip := net.ParseIP(ipStr)
		if isDisallowedIP(ip) {
			t.Errorf("Expected public IP %s to be allowed, but was disallowed", ipStr)
		}
	}
}

func TestValidateHostForSSRF(t *testing.T) {
	badHosts := []string{
		"localhost",
		"127.0.0.1",
		"169.254.169.254",
		"10.0.0.5",
		"192.168.1.50",
	}

	for _, host := range badHosts {
		if err := validateHostForSSRF(host); err == nil {
			t.Errorf("Expected host '%s' to be rejected by SSRF validator, but passed", host)
		}
	}
}
