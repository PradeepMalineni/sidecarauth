// main.go
package main

import (
	"fmt"
	"sidecarauth/auth"
	"sidecarauth/cache"
	"sidecarauth/listener"
)

func main() {
	// Initialize and use the auth sidecar proxy
	//authProxy := auth.NewAuthProxy()
	//authResult := authProxy.Authenticate("username", "password")
	fmt.Println("Authentication Result:", "hello2")
	fmt.Println("Authentication Result:", auth.X)

	// Initialize and use the cache module
	cacheInstance := cache.NewCache()
	cacheInstance.Set("key", "value")
	cachedValue, exists := cacheInstance.Get("key")
	if exists {
		fmt.Println("Cached Value:", cachedValue)
	} else {
		fmt.Println("Key not found in cache.")
	}

	//Start the Listner Function

	fmt.Println("Starting the Listner")
	listener.Listner()
}
