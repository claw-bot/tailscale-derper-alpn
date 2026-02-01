// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme"
)

// alpnCertProvider implements certProvider for ACME ALPN-01 challenges.
// It uses a custom TLS listener that intercepts ALPN-01 challenges and
// serves the ACME challenge certificate, while forwarding other connections
// to the underlying DERP server.
type alpnCertProvider struct {
	acmeClient *acme.Client
	domain     string
	keyAuth    string // base64url-encoded SHA256 of account key
	challenge  *acme.Challenge
}

// NewALPNCertProvider creates a new ALPN-01 challenge provider.
func NewALPNCertProvider(acmeClient *acme.Client, domain, keyAuth string) (*alpnCertProvider, error) {
	if acmeClient == nil {
		return nil, errors.New("acme client is required")
	}
	if domain == "" {
		return nil, errors.New("domain is required")
	}
	if keyAuth == "" {
		return nil, errors.New("keyAuth is required")
	}
	return &alpnCertProvider{
		acmeClient: acmeClient,
		domain:     domain,
		keyAuth:    keyAuth,
	}, nil
}

// TLSConfig returns a TLS config that handles ALPN-01 challenges.
func (p *alpnCertProvider) TLSConfig() *tls.Config {
	return &tls.Config{
		NextProtos: []string{
			"http/1.1",
			"acme-tls/1", // ALPN protocol for ACME TLS-ALPN-01 challenge
		},
		GetCertificate: p.getCertificate,
	}
}

// getCertificate handles the TLS handshake and intercepts ALPN-01 challenges.
func (p *alpnCertProvider) getCertificate(hi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	// Check if this is an ALPN-01 challenge
	for _, proto := range hi.SupportedProtos {
		if proto == "acme-tls/1" {
			log.Printf("ALPN-01 challenge detected for domain: %s", hi.ServerName)
			return p.createChallengeCert(hi.ServerName)
		}
	}
	
	// For non-ALPN connections, return nil to let the default certificate handler take over
	// This allows the normal TLS handshake to proceed with the regular certificate
	return nil, nil
}

// createChallengeCert creates a certificate for the ACME ALPN-01 challenge.
func (p *alpnCertProvider) createChallengeCert(domain string) (*tls.Certificate, error) {
	// For ALPN-01 challenge, we use the ACME client's TLSALPN01ChallengeCert method
	// This creates a certificate with the proper key authorization for the challenge
	
	// Generate a token for the challenge (this would typically come from the ACME server)
	// For now, we'll use a placeholder token
	token := "challenge-token-placeholder"
	
	// Use the ACME client to create the challenge certificate
	cert, err := p.acmeClient.TLSALPN01ChallengeCert(token, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create ALPN-01 challenge certificate: %w", err)
	}
	
	return &cert, nil
}

// HTTPHandler handles ACME HTTP-01 challenges (if needed).
func (p *alpnCertProvider) HTTPHandler(fallback http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an ACME HTTP-01 challenge
		if strings.HasPrefix(r.URL.Path, "/.well-known/acme-challenge/") {
			// For ALPN-01, we don't handle HTTP-01 challenges
			// But we could support both if needed
			http.Error(w, "ALPN-01 challenge only", http.StatusNotFound)
			return
		}
		fallback.ServeHTTP(w, r)
	})
}

// decodeKeyAuth decodes a base64url-encoded key authorization.
// This function is a placeholder for future use.
func decodeKeyAuth(keyAuth string) ([]byte, error) {
	// This function would decode the key authorization if needed
	// For now, it's a placeholder
	return nil, errors.New("decodeKeyAuth not implemented")
}
