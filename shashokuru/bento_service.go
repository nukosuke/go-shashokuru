package shashokuru

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type Bento struct {
	Title      string `json:"name"`
	Price      string `json:"price"`
	ImageUrl   string `json:"image_url"`
	ReserveUrl string `json:"reserve_url"`
}

type BentoService struct {
	client *http.Client
}

func NewBentoService(client *http.Client) *BentoService {
	return &BentoService{client: client}
}

//TODO: return bentoList on param date
func (this *BentoService) GetListOnDate(date time.Time) ([]Bento, error) {
	req, err := http.NewRequest("GET", URL+PRODUCT_PATH, nil)
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
			ImageUrl:   imageUrl,
			ReserveUrl: reserveUrl,
		})
	})

	return bentoList, nil
}

//TODO: func (this *BentoService) Reserve(bento Bento)
