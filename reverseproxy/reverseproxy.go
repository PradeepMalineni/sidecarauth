package reverseproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// ReverseProxy is a struct that implements the http.Handler interface
// and is used for forwarding requests to another HTTP service.
type ReverseProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	//trustStore *x509.CertPool
	certCache *sync.Map
}

func NewReverseProxy(targetURL string) (*ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// Load the custom trust store
	/*trustStore, err := utils.LoadTrustStore(trustStoreFile)
	if err != nil {s
		return nil, err
	}*/

	return &ReverseProxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
		//trustStore: trustStore,
		certCache: &sync.Map{},
	}, nil
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("Incoming Request: %s %s", r.Method, r.URL.Path)
	log.Print("This is a test message")

	// Check the cache for the last verification time

	rp.proxy.ServeHTTP(w, r)
}

/*
func (rp *ReverseProxy) checkCertificate(target *url.URL) error {
	log.Print("In cert verify funcs")

	conn, err := tls.Dial("tcp", target.Host, &tls.Config{
		InsecureSkipVerify: true,
		//RootCAs:            rp.trustStore,
	})
	if err != nil {
		return nil
	}
	defer conn.Close()

	// Additional certificate verification logic if needed
	return nil
}
*/
