package reverseproxy

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sidecarauth/utils"
	"sync"
	"time"
)

// ReverseProxy is a struct that implements the http.Handler interface
// and is used for forwarding requests to another HTTP service.
type ReverseProxy struct {
	target     *url.URL
	proxy      *httputil.ReverseProxy
	trustStore *x509.CertPool
	certCache  *sync.Map
}

func NewReverseProxy(targetURL string, trustStoreFile string) (*ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// Load the custom trust store
	trustStore, err := utils.LoadTrustStore(trustStoreFile)
	if err != nil {
		return nil, err
	}

	return &ReverseProxy{
		target:     target,
		proxy:      httputil.NewSingleHostReverseProxy(target),
		trustStore: trustStore,
		certCache:  &sync.Map{},
	}, nil
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("Incoming Request: %s %s", r.Method, r.URL.Path)
	log.Print("This is a test message")

	// Check the cache for the last verification time
	entry, found := rp.certCache.Load(rp.target.Host)
	if found {
		log.Print("Cache Found")

		cacheEntry := entry.(*utils.CacheEntry)
		// Check if the last verification was more than a certain duration ago
		if time.Since(cacheEntry.LastVerified) > 1*time.Hour {
			// Perform certificate verification
			err := rp.checkCertificate(rp.target)
			if err != nil {
				log.Printf("Error verifying certificate: %v", err)
			}
			// Update the last verification time in the cache
			rp.certCache.Store(rp.target.Host, &utils.CacheEntry{LastVerified: time.Now()})
		}
	} else {
		// Perform certificate verification for the first time
		log.Print("Cache not Found")

		//err := rp.checkCertificate(rp.target)

		/*if err != nil {
			log.Printf("Error verifying certificate: %v", err)
		}*/
		// Add the entry to the cache
		rp.certCache.Store(rp.target.Host, &utils.CacheEntry{LastVerified: time.Now()})
	}

	// Forward the request to the target service
	rp.proxy.ServeHTTP(w, r)
}

func (rp *ReverseProxy) checkCertificate(target *url.URL) error {
	log.Print("In cert verify funcs")

	conn, err := tls.Dial("tcp", target.Host, &tls.Config{
		InsecureSkipVerify: true,
		//RootCAs:            rp.trustStore,

	})
	log.Print("In cert verify funcx")

	if err != nil {
		return nil
	}
	defer conn.Close()

	// Additional certificate verification logic if needed
	return nil
}
