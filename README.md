# go-shashokuru

```go
package main

import (
	"github.com/nukosuke/go-shashokuru/shashokuru"
	"os"
	"time"
	"fmt"
)

const URL_DATE_FORMAT = "20060102"

func main() {
	client := shashokuru.NewClient()
	client.Login("your_shashokuru_email@example.com", "your_shashokuru_password")

	// :bento: メニューを取得する日付
	date, _ := time.Parse(URL_DATE_FORMAT, "20180122")

	// :bento: メニュー取得
	bentoList, err := client.Bento.GetListOnDate(date)
	if err != nil {
		fmt.Println(":bento: チャレンジ失敗")
		os.Exit(1)
	}

	// :bento: 予約
	// リストの最初の弁当を1つ予約
	err = client.Bento.Reserve(bentoList[0], 1)
	if err != nil {
 		fmt.Println(":bento: チャレンジ失敗")
		os.Exit(1)
	}

	fmt.Println(":bento: 獲得した。これで明日も生きていける。")
}
```
