// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"crypto/tls"
	"testing"

	"golang.org/x/crypto/acme"
)

func TestALPNCertProvider(t *testing.T) {
	// Test creating ALPN cert provider
	acmeClient := &acme.Client{
		DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
	}
	
	provider, err := NewALPNCertProvider(acmeClient, "example.com", "test-key-auth")
	if err != nil {
		t.Fatalf("Failed to create ALPN cert provider: %v", err)
	}
	
	// Test TLS config
	config := provider.TLSConfig()
	if config == nil {
		t.Fatal("TLSConfig returned nil")
	}
	
	// Check that acme-tls/1 is in NextProtos
	found := false
	for _, proto := range config.NextProtos {
		if proto == "acme-tls/1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("acme-tls/1 not found in NextProtos")
	}
	
	// Test GetCertificate
	if config.GetCertificate == nil {
		t.Error("GetCertificate is nil")
	}
}

func TestALPNCertProviderErrors(t *testing.T) {
	// Test nil ACME client
	_, err := NewALPNCertProvider(nil, "example.com", "test-key-auth")
	if err == nil {
		t.Error("Expected error for nil ACME client")
	}
	
	// Test empty domain
	acmeClient := &acme.Client{
		DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
	}
	_, err = NewALPNCertProvider(acmeClient, "", "test-key-auth")
	if err == nil {
		t.Error("Expected error for empty domain")
	}
	
	// Test empty keyAuth
	_, err = NewALPNCertProvider(acmeClient, "example.com", "")
	if err == nil {
		t.Error("Expected error for empty keyAuth")
	}
}

func TestNewALPNCertProviderFromFlags(t *testing.T) {
	// Test creating provider from flags
	provider, err := NewALPNCertProviderFromFlags("/tmp/certs", "example.com", "", "", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to create ALPN cert provider from flags: %v", err)
	}
	
	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestNewALPNCertProviderFromFlagsErrors(t *testing.T) {
	// Test missing hostname
	_, err := NewALPNCertProviderFromFlags("/tmp/certs", "", "", "", "test@example.com")
	if err == nil {
		t.Error("Expected error for missing hostname")
	}
	
	// Test missing email
	_, err = NewALPNCertProviderFromFlags("/tmp/certs", "example.com", "", "", "")
	if err == nil {
		t.Error("Expected error for missing email")
	}
}

func TestGetCertificate(t *testing.T) {
	acmeClient := &acme.Client{
		DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
	}
	
	provider, err := NewALPNCertProvider(acmeClient, "example.com", "test-key-auth")
	if err != nil {
		t.Fatalf("Failed to create ALPN cert provider: %v", err)
	}
	
	// Test with ALPN challenge
	hi := &tls.ClientHelloInfo{
		ServerName: "example.com",
		SupportedProtos: []string{"acme-tls/1"},
	}
	
	// This will fail because we don't have a real ACME account,
	// but it should at least call the method
	_, err = provider.getCertificate(hi)
	// We expect an error because we don't have a real ACME account
	if err == nil {
		t.Error("Expected error when creating challenge certificate without real ACME account")
	}
	
	// Test without ALPN challenge
	hi2 := &tls.ClientHelloInfo{
		ServerName: "example.com",
		SupportedProtos: []string{"http/1.1"},
	}
	
	cert, err := provider.getCertificate(hi2)
	if err != nil {
		t.Errorf("Expected no error for non-ALPN connection: %v", err)
	}
	if cert != nil {
		t.Error("Expected nil certificate for non-ALPN connection")
	}
}
