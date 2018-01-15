package shashokuru

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	URL          = "https://shashokuru.jp"
	LOGIN_PATH   = "/login"
	PRODUCT_PATH = "/product"
)

type Shashokuru struct {
	url    string
	client *http.Client
}

func NewClient() *Shashokuru {
	jar, _ := cookiejar.New(nil)
	return &Shashokuru{url: URL, client: &http.Client{Jar: jar}}
}

func (this *Shashokuru) Login(email string, password string) error {
	csrfToken, err := this.getCsrfToken()
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("email", email)
	values.Add("password", password)
	values.Add("_token", csrfToken)

	req, err := http.NewRequest("POST", URL+LOGIN_PATH, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://shashokuru.jp/")

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Login Error: %s", resp.Status)
	}
	return nil
}

func (this *Shashokuru) PrintBentoList() error {
	req, err := http.NewRequest("GET", URL+PRODUCT_PATH, nil)
	if err != nil {
		return err
	}

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	doc.Find(".wrapper").Each(func(_ int, selection *goquery.Selection) {
		fmt.Println("title: ", selection.Find(".title").Text())
		fmt.Println("price: ", selection.Find(".price").Text())
	})
	return nil
}

//////////////////// private methods ////////////////////

func (this *Shashokuru) getCsrfToken() (string, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}

	token, exists := doc.Find("meta[name=csrf-token]").Attr("content")
	if exists == false {
		return "", fmt.Errorf("Error: Not found CSRF token")
	}
	return token, nil
}
