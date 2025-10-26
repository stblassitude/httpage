# httpage

A tiny command-line tool that performs an HTTP HEAD request to a given URL and prints a small JSON object with:
- age: The age in seconds since the server’s Last-Modified header (0 if missing or invalid)
- status: Either the HTTP status code (number) or 0 if the request could not be sent or completed in time

Example output:

{"age": 123, "status": 200}
{"age": 0, "status": 0}

This is handy for monitoring freshness of cached content or endpoints that expose Last-Modified.


## Installation

- From source (requires Go):

  go install github.com/your/repo/httpage@latest

- From GitHub Releases:
  1. Download a binary for your platform from the Releases page.
  2. Make it executable (Linux/macOS): chmod +x httpage
  3. Place it in your PATH (e.g., /usr/local/bin).

Note: This repository ships a GitHub Actions workflow that builds release binaries for Linux, macOS, and Windows on amd64 and arm64 whenever a git tag is pushed.


## Usage

httpage <url>

- On success, prints one line of JSON: {"age":<seconds>,"status":<http status code>}
- If the request fails to be sent or completed (DNS/TCP/TLS error, timeout, etc.), prints {"age":0,"status":"timeout"}
- The program uses an internal request timeout of 15 seconds.

Examples:

# Basic
httpage https://example.com/

# With curl to compare headers
curl -I https://example.com/

Exit codes:
- 0 on normal JSON output (including timeout case)
- 2 on usage error (wrong number of arguments)


## Using with Telegraf [[inputs.exec]]

You can collect the JSON output using Telegraf’s exec input and parse it as JSON. Because the "status" field can be either a number (HTTP status code) or a string ("timeout"), it’s best to force Telegraf to treat status as a string using json_string_fields.

Minimal example (Linux/macOS):

```toml
[[inputs.exec]]
  commands = ["/usr/local/bin/httpage https://example.com/"]
  timeout = "20s"
  data_format = "json"
  json_string_fields = ["status"]
  name_override = "httpage"
  tags = { url = "https://example.com/" }
```

Notes:
- data_format = "json" tells Telegraf to parse the stdout as JSON.
- json_string_fields = ["status"] ensures the status is always a string (e.g., "200" or "timeout"). If you prefer to keep numeric status codes numeric, remove that line; be aware that a string value on timeout will lead to a field-type change.
- name_override = "httpage" sets the measurement name.
- tags = { url = "..." } lets you distinguish multiple monitored URLs.

Windows example:

```toml
[[inputs.exec]]
  commands = ["C:\\Program Files\\httpage\\httpage.exe https://example.com/"]
  timeout = "20s"
  data_format = "json"
  json_string_fields = ["status"]
  name_override = "httpage"
  tags = { url = "https://example.com/" }
```

Collecting multiple URLs:
- Add additional commands entries, one per URL, each with a different URL in tags, or
- Wrap multiple invocations in a small shell script and call that script from [[inputs.exec]]. Ensure each invocation prints exactly one JSON object per line.

Troubleshooting:
- If you see field-type conflicts for status, set json_string_fields = ["status"].
- If you get timeouts in Telegraf, increase the plugin timeout (timeout = "20s"). The httpage internal timeout is 15s; Telegraf’s timeout should be larger than that.


## License

This project is licensed under the Apache License, Version 2.0. See the LICENSE file for details.
