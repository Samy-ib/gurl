package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Noooste/azuretls-client"
)

var browser_to_ja3 = map[string]string{
	"chrome":  "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,65037-10-65281-43-17513-5-0-18-35-16-13-11-51-27-23-45,29-23-24,0",
	"firefox": "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-16-5-34-51-43-13-45-28-65037-41,29-23-24-25-256-257,0",
}

var browser_to_http2 = map[string]string{
	"chrome":  "1:65536,2:0,4:6291456,6:262144|15663105|0|m,a,s,p",
	"firefox": "1:65536,4:131072,5:16384|12517377|3:0:0:201,5:0:0:101,7:0:0:1,9:0:7:1,11:0:3:1,13:0:0:241|m,p,a,s",
}

var browser_to_headers = map[string]azuretls.OrderedHeaders{
	"firefox": {
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "fr-FR,en-US;q=0.7,en;q=0.3"},
		{"Dnt", "1"},
		{"Sec-Fetch-Dest", "document"},
		{"Sec-Fetch-Mode", "navigate"},
		{"Sec-Fetch-Site", "none"},
		{"Sec-Fetch-User", "?1"},
		{"Upgrade-Insecure-Requests", "1"},
		{"User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0"},
	},
	"chrome": {
		{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"Accept-Language", "en-US,en;q=0.9"},
		{"Host", "httpbin.org"},
		{"Sec-Ch-Ua", "\"Not(A,Brand\";v=\"24\", \"Chromium\";v=\"122\""},
		{"Sec-Ch-Ua-Mobile", "?0"},
		{"Sec-Ch-Ua-Platform", "\"Linux\""},
		{"Sec-Fetch-Dest", "document"},
		{"Sec-Fetch-Mode", "navigate"},
		{"Sec-Fetch-Site", "none"},
		{"Sec-Fetch-User", "?1"},
		{"Upgrade-Insecure-Requests", "1"},
		{"User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"},
	},
}

var browser_to_headers_order = map[string][]string{
	"firefox": {"Accept", "Accept-Encoding", "Dnt", "Host", "Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-User", "Upgrade-Insecure-Requests", "User-Agent"},
	"chrome":  {"Accept", "Accept-Encoding", "Accept-Language", "Host", "Sec-Ch-Ua", "Sec-Ch-Ua-Mobile", "Sec-Ch-Ua-Platform", "Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-Site", "Upgrade-Insecure-Requests"},
}

func main() {

	f_browser := flag.String("browser", "chrome", "Browser to mimic. Will have same headers as the browser + TLS and HTTP2 fingerprints. Accepted values: chrome, firefox.")
	// f_tls := flag.String("tls", "", "Mimic the browser TLS, will overwrite the --browser TLS. Accepted values: chrome, firefox.")
	// f_ja3 := flag.String("ja3", "", "Mimic the browser TLS, will overwrite the --browser TLS. Accepted values: chrome, firefox.")
	// f_http2 := flag.String("http2", "", "Mimic the browser TLS, will overwrite the --browser TLS. Accepted values: chrome, firefox.")
	flag.Parse()

	session := azuretls.NewSession()
	defer session.Close()

	// response, err := session.Get(os.Args[1])

	switch *f_browser {
	case "chrome":
		if err := session.ApplyJa3(browser_to_ja3["chrome"], azuretls.Chrome); err != nil {
			panic(err)
		}
		if err := session.ApplyHTTP2(browser_to_http2["chrome"]); err != nil {
			panic(err)
		}
		session.OrderedHeaders = browser_to_headers["chrome"]
		session.HeaderOrder = browser_to_headers_order["chrome"]
	case "firefox":
		if err := session.ApplyJa3(browser_to_ja3["firefox"], azuretls.Firefox); err != nil {
			panic(err)
		}
		if err := session.ApplyHTTP2(browser_to_http2["firefox"]); err != nil {
			panic(err)
		}
		session.OrderedHeaders = browser_to_headers["firefox"]
		session.HeaderOrder = browser_to_headers_order["firefox"]
	}
	response, err := session.Get(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(response.StatusCode, string(response.Body))
}
