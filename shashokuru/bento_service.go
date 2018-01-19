package shashokuru

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const URL_DATE_FORMAT = "20060102"

type Bento struct {
	Title      string `json:"name"`
	Price      string `json:"price"`
	Store      string `json:"store"`
	ImageUrl   string `json:"image_url"`
	ReserveUrl string `json:"reserve_url"`
}

type BentoService struct {
	client *http.Client
}

func NewBentoService(client *http.Client) *BentoService {
	return &BentoService{client: client}
}

func (this *BentoService) GetListOnDate(date time.Time) ([]Bento, error) {
	datePath := "/" + date.Format(URL_DATE_FORMAT)

	req, err := http.NewRequest("GET", URL+PRODUCT_PATH+datePath, nil)
	if err != nil {
		return []Bento{}, err
	}

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return []Bento{}, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return []Bento{}, err
	}

	bentoList := []Bento{}
	doc.Find(".wrapper").Each(func(_ int, selection *goquery.Selection) {
		imageUrl, _ := selection.Find("img").Attr("src")
		reserveUrl, _ := selection.Find("a.btn-a").Attr("href")

		bentoList = append(bentoList, Bento{
			Title:      selection.Find(".title").Text(),
			Price:      selection.Find(".price").Text(),
			Store:      selection.Find(".store").Text(),
			ImageUrl:   imageUrl,
			ReserveUrl: reserveUrl,
		})
	})

	return bentoList, nil
}

func (this *BentoService) Reserve(bento Bento, quantity int) error {
	nextFormValues, err := this.openDetail(bento, quantity)
	if err != nil {
		return err
	}

	nextFormValues, err = this.moveDetailToCart(nextFormValues, bento)
	if err != nil {
		return err
	}

	nextFormValues, err = this.moveCartToConfirmation(nextFormValues)
	if err != nil {
		return err
	}

	return this.moveConfirmationToComplete(nextFormValues)
}

//////////////////// private methods ////////////////////

// :bento: 詳細ページへ遷移
func (this *BentoService) openDetail(bento Bento, quantity int) (url.Values, error) {
	req, err := http.NewRequest("GET", bento.ReserveUrl, nil)
	if err != nil {
		return url.Values{}, err
	}

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return url.Values{}, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return url.Values{}, err
	}

	token, tokenExist := doc.Find("input[name=_token]").Attr("value")
	storeId, storeIdExist := doc.Find("input[name=store_id]").Attr("value")
	if !tokenExist || !storeIdExist {
		return url.Values{}, fmt.Errorf("Error: Cannot find CSRF token or store_id")
	}

	nextFormValues := url.Values{}
	nextFormValues.Add("_token", token)
	nextFormValues.Add("buy_quantity", strconv.Itoa(quantity))
	nextFormValues.Add("store_id", storeId)
	return nextFormValues, nil
}

// カート画面へ遷移
func (this *BentoService) moveDetailToCart(values url.Values, bento Bento) (url.Values, error) {
	req, err := http.NewRequest("POST", bento.ReserveUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return url.Values{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", URL)

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return url.Values{}, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return url.Values{}, err
	}

	token, tokenExist := doc.Find("input[name=_token]").Attr("value")
	paymentMethod, pmExist := doc.Find("input[name=payment_method]").Attr("value")
	if !tokenExist || !pmExist {
		return url.Values{}, fmt.Errorf("Error: Cannot find CSRF token or payment_method")
	}

	nextFormValues := url.Values{}
	nextFormValues.Add("_token", token)
	nextFormValues.Add("payment_method", paymentMethod)
	return nextFormValues, nil
}

// 予約確認画面へ遷移
func (this *BentoService) moveCartToConfirmation(values url.Values) (url.Values, error) {
	req, err := http.NewRequest("POST", URL+CART_PATH, strings.NewReader(values.Encode()))
	if err != nil {
		return url.Values{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", URL)

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return url.Values{}, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return url.Values{}, err
	}

	token, tokenExist := doc.Find("input[name=_token]").Attr("value")
	if !tokenExist {
		return url.Values{}, fmt.Errorf("Error: Cannot find CSRF token")
	}

	nextFormValues := url.Values{}
	nextFormValues.Add("_token", token)
	return nextFormValues, nil
}

// 予約完了画面に遷移
func (this *BentoService) moveConfirmationToComplete(values url.Values) error {
	req, err := http.NewRequest("POST", URL+CART_CONFIRMATION_PATH, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", URL)

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
