package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

func Convert(amount float64, from string, to string) (*http.Response, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://www.xe.com/currencyconverter/convert/?Amount=%f&From=%s&To=%s", amount, strings.ToUpper(from), strings.ToUpper(to))
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authority", "www.xe.com")
	req.Header.Add("authority", "www.xe.com")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("accept-language", "en-US,en;q=0.9,uk-UA;q=0.8,uk;q=0.7,ru-RU;q=0.6,ru;q=0.5")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp, err := client.Do(req)
	return resp, err
}

func ExtractConversion(resp *http.Response) string {
	const selector = "p[class*=result__BigRate]"
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	conversionResult := doc.Find(selector).First()
	return strings.Split(conversionResult.Text(), " ")[0]
}

func main() {
	if len(os.Args[1:]) != 3 {
		panic("Need amount, from, to as arguments.")
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		panic(err)
	}

	from, to := os.Args[2], os.Args[3]

	resp, err := Convert(amount, from, to)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("HTTP request not OK")
	}

	conversion := ExtractConversion(resp)
	fmt.Printf("%s %s", conversion, to)
}
