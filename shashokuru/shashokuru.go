package shashokuru

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	URL                    = "https://shashokuru.jp"
	LOGIN_PATH             = "/login"
	PRODUCT_PATH           = "/product"
	CART_PATH              = "/cart"
	CART_CONFIRMATION_PATH = CART_PATH + "/confirmation"
)

type Shashokuru struct {
	url    string
	client *http.Client
	Bento  *BentoService
}

func NewClient() *Shashokuru {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	bento := NewBentoService(client)

	return &Shashokuru{
		url:    URL,
		client: client,
		Bento:  bento,
	}
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
