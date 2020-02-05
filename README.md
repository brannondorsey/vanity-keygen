# Vanity Keygen

A vanity key generator for the P224, P256, P384, and P521 elliptic curves. This CLI uses Go's `crypto/elliptic` to spin up one concurrent goroutine on each processor core to search for public keys whose Base64 representations contain specific strings.

```
vanity-keygen --match-case Linux
...
[INFO] Generating p256 key pair using 12 threads
[INFO] String matching "Linux" in location: beginning
[INFO] Strict case matching: false
[INFO] Using Go's safe crypto/rand library for cryptographic random number generator
[INFO] Match found in 3.673709084s
[INFO] Public key:  BLInuxyBTk6D9f0BUhM4NGoezUyDh51ZcejXGB5ybeN3XSpgxPhnj65iek8qwKcUUd3QHhRqlE28qyycrzRey9k=
[INFO] Private key: XXXXXXXXXXXXXXX
```

Generated keys are in the following format:

```
Base64(concat(X, Y))
```

Where `x` and `y` are the points of a public key and `concat(x, y)` is the uncompressed form of the key specified in section 4.3.6 of ANSI X9.62. Finally the uncompressed version of the key is encoded in Base64 using for portability, and to make a wide variety of characters available for string matching.

## Download

Pre-compiled binaries are available for Linux (x64 and ARM), MacOS, and Windows can be downloaded for from the latest [release page](https://github.com/brannondorsey/vanity-keygen/releases/latest).

## Usage

In its simplest form, the tool can be invoked with `vanity-keygen <search-string>`. However, the flags listed below can be used to tune the behavior of the key search.

```
Usage: vanity-keygen [OPTIONS] <search-string> ...
  -t, --concurrency int         The number of concurrent goroutines to use to perform the search. Default 1 per CPU.
  -c, --curve string            The name of the curve to generate keys for. Accepted values include "p224", "p256", "p384", "p521". (default "p256")
  -h, --help                    Show this screen.
  -i, --insecure                Use the unvetted fastrand library for cryptographic randomness (https://gitlab.com/NebulousLabs/fastrand)
      --match-case              Enable strict case matching. This will dramatically increase the search time.
      --match-location string   The location of the search string in the generated key. Accepted values include "beginning", "end", "anywhere" (default "beginning")
      --verbose                 Print verbose output
Version: vanity-keygen version v0.1.0-snapshot+13c5df3 (built Sun Feb  2 21:09:00 UTC 2020)
```

## Recommendations

* Use *very* short values for `<search-string>`, as each letter increases the average search time by a 64x. The difficulty of the search explodes explodes very quickly and searches for longer strings can take well over a human lifetime.
* The `--insecure` flag can be used to increase the search speed by ~1.5-2x at the expense of sacrificing security by using the [`fastrand`](https://gitlab.com/NebulousLabs/fastrand) as the source of cryptographic random, instead of Go's `crypto/rand`. Fastrand uses `crypto/rand` to seed 32 bytes of strong random entropy, but then uses a deterministic non-guessable hashing algorithm to generate random bytes more quickly for the rest of the lifetime of the program. You can read more about the security concerns with this approach [here](https://gitlab.com/NebulousLabs/fastrand#security).
* The default `--match-case=false` flag can be used to control wether keys produced during the search match the exact case of the `<search-string>` provided as input (e.g. `Linux` vs `lInUX`). Matching exact case takes far longer and may result in searches that simply can't be completed in your lifetime.
* `--match-location` can be used to define the location your `<search-string>` will appear in the generated key. Accepted values include:
  * `beginning` (default): `BLInuxyBTk6D9f0BUhM4NGoezUyDh51ZcejXGB5ybeN3XSpgxPhnj65iek8qwKcUUd3QHhRqlE28qyycrzRey9k=`. Note the first character in a Base64 encoded uncompressed ANSI X9.62 public key will always be "B", so the vanity text will appear directly after it.
  * `anywhere`: `BLOfiARtTHJuzbZZqbWbxlsH8YGDLinuxySvytYhrxLDhZJqD3TI6Ga8vwvI4fra8GnR9N2iUiy/RBZn5W78d80=`. The search string can appear anywhere in the key. These searches are much faster!
  * `end`: `BRey9kyBTk6D9f0BUhM4NGoezUyDh51ZcejXGB5ybeN3XSpgxPhnj65iek8qwKcUUd3QHhRqlE28qyycrzLInux=` The search string appears at the end of the key. In my experience, these searches take the longest, and may never produce valid keys. (Note: this is an invalid public key).
* When using `--match-location=beginning` or `--match-location=end`, there is a chance that your `<search-string>` can't generate a valid key, and therefore your search will never terminate. This is because not all (x, y) pairs appear on an elliptic curve (in fact, most don't), and so its quite possible that no points exist on the curve whose Base64 representations begin or end with the `<search-string>` you've provided. This tool makes no attempt to determine the likelyhood your search string can even produce a valid key. The shorter your search string, the more likely you are to actually find a match.
* Use `--verbose` to see search speed printed to the screen periodically.
