# badgateway

It's a bad gateway. Use it to induce random 502 errors for testing HTTP client retry handling.

## Usage

- `go install .`
- `badgateway -port 12345 -fail 0.25 https://some.target.host`
- `ngrok http 12345`
- Point your client at the ngrok host instead of your original target host
- Watch the random failures
