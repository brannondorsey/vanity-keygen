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

func getFindKeyFunc(curve elliptic.Curve, needle string, matchCase bool, matchLocation string, ch chan *keyPair, thread int) func() {

	needleLen := len(needle)
	if !matchCase {
		needle = strings.ToLower(needle)
	}

	return func() {
		for i := 1; true; i++ {
			private, x, y, _ := elliptic.GenerateKey(curve, rand.Reader)
			public := elliptic.Marshal(curve, x, y)
			haystack := base64.StdEncoding.EncodeToString(public)
			if !matchCase {
				haystack = strings.ToLower(haystack)
			}
			switch matchLocation {
			case "beginning":
				haystack = haystack[1 : needleLen+1]
				if needle != haystack {
					continue
				}
			case "end":
				hackstackLen := len(haystack)
				haystack = haystack[hackstackLen-needleLen-1 : hackstackLen-1]
				if needle != haystack {
					continue
				}
			default:
				if !strings.Contains(haystack, needle) {
					continue
				}
			}
			ch <- &keyPair{Public: public, Private: private}
			return
		}
	}
}

type Arguments struct {
	NumThreads    int
	MatchCase     bool
	MatchLocation string
	CurveName     string
	Needle        string
}

func parseArgs() Arguments {
	var needle *string
	numThreads := flag.Int("threads", runtime.NumCPU(), "The number of threads to use to perform the search. Default 1 per CPU.")
	curveName := flag.String("curve", "p256", "The name of the curve to generate keys for. Accepted values include \"p224\", \"p256\", \"p384\", \"p521\".")
	matchCase := flag.Bool("match-case", false, "Enable strict case matching. This will dramatically increase the search time.")
	matchLocation := flag.String("match-location", "beginning", "The location of the search string in the generated key. Accepted values include \"beginning\", \"end\", \"anywhere\"")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <search-string> ...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if !(*curveName == "p224" || *curveName == "p256" || *curveName == "p384" || *curveName == "p521") {
		fmt.Printf("[ERROR] Invalid curve name \"%s\"\n", *curveName)
		os.Exit(1)
	}
	if !(*matchLocation == "beginning" || *matchLocation == "end" || *matchLocation == "anywhere") {
		fmt.Printf("[ERROR] Invalid match location \"%s\\n", *matchLocation)
		os.Exit(1)
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	} else {
		needle = &flag.Args()[0]
	}
	return Arguments{
		NumThreads:    *numThreads,
		MatchCase:     *matchCase,
		MatchLocation: *matchLocation,
		CurveName:     *curveName,
		Needle:        *needle,
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
	fmt.Printf("[INFO] Generating %s key pair using %d threads\n", args.CurveName, args.NumThreads)
	fmt.Printf("[INFO] String matching \"%s\" in location: %s\n", args.Needle, args.MatchLocation)
	fmt.Printf("[INFO] Strict case matching: %t\n", args.MatchCase)
	ch := make(chan *keyPair)
	for thread := 1; thread < args.NumThreads+1; thread++ {
		go getFindKeyFunc(curve, args.Needle, args.MatchCase, args.MatchLocation, ch, thread)()
	}
	key := <-ch
	fmt.Printf("[INFO] Match found in %s\n", time.Since(start))
	fmt.Printf("[INFO] Public key:  %s\n", base64.StdEncoding.EncodeToString(key.Public))
	fmt.Printf("[INFO] Private key: %s\n", base64.StdEncoding.EncodeToString(key.Private))
}
