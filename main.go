package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	myfirego "github.com/khaiql/firego"
	"github.com/zabawaba99/firego"
)

const (
	url        = "https://benchmark-7f1d0.firebaseio.com"
	secret     = "3Bvs506qQgA0qyHg04r3xyTvcwrGpouR0ufeLiVR"
	maxKey     = 30
	maxRequest = 1000
)

func main() {

	doneChan := make(chan bool)
	go func() {
		for i := 0; i < maxRequest; i++ {
			index := fmt.Sprintf("%d", i)
			go func() {
				_url := fmt.Sprintf("%s/myfirego/%s", url, index)
				fmt.Println(_url)
				myFirego := myfirego.New(_url, nil)
				myFirego.Auth(secret)
				myFirego.SetParams("print", "silence")
				myfirego.TimeoutDuration = 1 * time.Minute
				fmt.Printf("My Firego: %s\n", index)
				message := produceMessage("my_firego", index)
				if err := myFirego.Set(message); err != nil {
					fmt.Printf("MyFirego - Error: %s. URL: %s\n", err.Error(), myFirego.String())
				}
			}()
		}
	}()

	go func() {
		for i := 0; i < maxRequest; i++ {
			index := fmt.Sprintf("%d", i)
			go func() {
				firebaseTimeout := 3 * time.Minute
				httpClient := &http.Client{
					Transport: &http.Transport{
						DisableKeepAlives: true,
						Dial: func(network, address string) (net.Conn, error) {
							c, err := net.DialTimeout(network, address, firebaseTimeout)
							return c, err
						},
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				}
				firego := firego.New(fmt.Sprintf("%s/firego/%s", url, index), httpClient)
				firego.Auth(secret)
				fmt.Printf("Firego: %s\n", index)
				message := produceMessage("firego", index)
				if err := firego.Set(message); err != nil {
					fmt.Printf("Firego - Error: %s. URL: %s\n", err.Error(), firego.String())
				}
			}()
		}
	}()

	<-doneChan
}

func produceMessage(identifer string, count string) map[string]string {
	var (
		data       = make(map[string]string)
		key, value string
	)
	for i := 0; i < maxKey; i++ {
		key = fmt.Sprintf("KEY_%s_%s_%d", identifer, count, i)
		value = fmt.Sprintf("VALUE_%s_%d", count, i)
		data[key] = value
	}

	return data
}
