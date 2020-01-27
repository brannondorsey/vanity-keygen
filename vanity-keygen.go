package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"flag"
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

func generateKeyStartingWith(curve elliptic.Curve, prefix string, matchCase bool, ch chan *keyPair, thread int) {
	prefixLower := strings.ToLower(prefix)
	prefixLen := len(prefix)
	for i := 1; true; i++ {
		private, x, y, _ := elliptic.GenerateKey(curve, rand.Reader)
		public := elliptic.Marshal(curve, x, y)
		publicString := base64.StdEncoding.EncodeToString(public)
		if strings.ToLower(publicString[1:prefixLen+1]) == prefixLower {
			if matchCase && publicString[1:prefixLen+1] != prefix {
				continue
			}
			ch <- &keyPair{Public: public, Private: private}
			return
		}
	}
}

type Arguments struct {
	NumThreads int
	MatchCase  bool
	CurveName  string
	Needle     string
}

func parseArgs() Arguments {
	var needle *string
	numThreads := flag.Int("threads", runtime.NumCPU(), "The number of threads to use to perform the search. Default 1 per CPU.")
	matchCase := flag.Bool("match-case", false, "Enable strict case matching. This will dramatically increase the search time.")
	curveName := flag.String("curve", "p256", "The name of the curve to generate keys for. Accepted values include \"p224\", \"p256\", \"p384\", \"p521\".")
	flag.Parse()
	if !(*curveName == "p224" || *curveName == "p256" || *curveName == "p384" || *curveName == "p521") {
		fmt.Printf("Error: Invalid curve name \"%s\"\n", *curveName)
		os.Exit(1)
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	} else {
		needle = &flag.Args()[0]
	}
	return Arguments{
		NumThreads: *numThreads,
		MatchCase:  *matchCase,
		CurveName:  *curveName,
		Needle:     *needle,
	}
}

func getCurve(name string) elliptic.Curve {
	switch name {
	case "p224":
		return elliptic.P224()
	case "p256":
		return elliptic.P256()
	case "p384":
		return elliptic.P384()
	default:
		return elliptic.P521()
	}
}

func main() {
	args := parseArgs()
	curve := getCurve(args.CurveName)
	start := time.Now()
	fmt.Printf("[INFO] Generating %s key pair that start with \"%s\" using %d threads\n", args.CurveName, args.Needle, args.NumThreads)
	ch := make(chan *keyPair)
	for thread := 1; thread < args.NumThreads+1; thread++ {
		go generateKeyStartingWith(curve, args.Needle, args.MatchCase, ch, thread)
	}
	key := <-ch
	fmt.Printf("[INFO] Match found in %s\n", time.Since(start))
	fmt.Printf("[INFO] Public key:  %s\n", base64.StdEncoding.EncodeToString(key.Public))
	fmt.Printf("[INFO] Private key: %s\n", base64.StdEncoding.EncodeToString(key.Private))
}
