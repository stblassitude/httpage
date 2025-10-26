package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	os.Exit(run())
}

type result struct {
	Age    int64       `json:"age"`
	Status interface{} `json:"status"`
}

func printJSON(r result) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(r)
}

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <url>\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return 2
	}

	url := flag.Arg(0)

	// Create HTTP client with a sane timeout
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		// Treat inability to create/send the request as timeout per requirements
		printJSON(result{Age: 0, Status: 0})
		return 0
	}

	// Ask servers/proxies not to give us cached metadata if possible
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		// Cannot send request or receive response
		printJSON(result{Age: 0, Status: 0})
		return 0
	}
	defer resp.Body.Close()

	// Compute age if Last-Modified is present and valid; otherwise age is 0
	var seconds int64 = 0
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		if t, err := http.ParseTime(lm); err == nil {
			seconds = int64(time.Since(t) / time.Second)
		}
	}

	// Always report JSON with age and HTTP status code
	printJSON(result{Age: seconds, Status: resp.StatusCode})
	return 0
}

// Keep fprintfErr in case other tooling expects it, but it's no longer used for output
func fprintfErr(format string, a ...any) {
	// Ensure newline if format didn't include it
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}
