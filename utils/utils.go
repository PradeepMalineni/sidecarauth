package utils

import (
	"crypto/x509"
	"log"
	"os"
	"time"
)

// CacheEntry represents an entry in the certificate cache.
type CacheEntry struct {
	LastVerified time.Time
}

func LoadTrustStore(trustStoreFile string) (*x509.CertPool, error) {
	pemData, err := os.ReadFile(trustStoreFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(pemData)
	log.Print("trust pool", certPool)
	return certPool, nil
}
