package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type keyPair struct {
	Public  []byte
	Private []byte
}

func generateKeyStartingWith(curve elliptic.Curve, prefix string, ch chan *keyPair, thread int) {
	prefixLower := strings.ToLower(prefix)
	prefixLen := len(prefix)
	for i := 1; true; i++ {
		private, x, y, _ := elliptic.GenerateKey(curve, rand.Reader)
		public := elliptic.Marshal(curve, x, y)
		publicString := base64.StdEncoding.EncodeToString(public)
		if strings.ToLower(publicString[1:prefixLen+1]) == prefixLower {
			ch <- &keyPair{Public: public, Private: private}
			return
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <needle>\n", os.Args[0])
		os.Exit(1)
	}
	needle := os.Args[1]
	curve := elliptic.P256()
	numCpus := runtime.NumCPU()
	start := time.Now()
	fmt.Printf("[INFO] Generating P256 key pair that start with \"%s\" using %d threads\n", needle, numCpus)
	ch := make(chan *keyPair)
	for thread := 1; thread < numCpus+1; thread++ {
		go generateKeyStartingWith(curve, needle, ch, thread)
	}
	key := <-ch
	fmt.Printf("[INFO] Match found in %s\n", time.Since(start))
	fmt.Printf("[INFO]  Public key:  %s\n", base64.StdEncoding.EncodeToString(key.Public))
	fmt.Printf("[INFO]  Private key: %s\n", base64.StdEncoding.EncodeToString(key.Private))
}
