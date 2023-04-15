package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
)

const (
	SERIAL_CODE        string = `//*[@id="e_9474"]`
	EMAIL              string = `//*[@id="e_9477"]`
	NAME               string = `//*[@id="e_9490"]`
	FURIGANA           string = `//*[@id="e_9479"]`
	AGE                string = `//*[@id="e_9480"]`
	POSTAL_CODE1       string = `//*[@id="e_9482[zip1]"]`
	POSTAL_CODE2       string = `//*[@id="e_9482[zip2]"]`
	ADRESS             string = `//*[@id="e_9483"]`
	TEL1               string = `//*[@id="e_9484[tel1]"]`
	TEL2               string = `//*[@id="e_9484[tel2]"]`
	TEL3               string = `//*[@id="e_9484[tel3]"]`
	SUBMIT             string = `//*[@id="__send"]`
	SCREENSHOT_QUALITY int    = 60 // スクショのクオリティ
	CONFIRM_BUTTON     string = `//*[@id="__commit"]`
)

var buf []byte

func main() {

	e := echo.New()

	// HTMLファイルの読み込み
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	// フォームの送信処理
	e.POST("/submit", func(c echo.Context) error {

		serialCodes := c.FormValue("serialCodes")
		service := c.FormValue("service")

		// ご希望のコース・希望内容であれこれ
		var SERVICE string
		var content string
		var CONTENT string
		switch service {
		case "2":
			SERVICE = `//*[@id="service_2"]`
			content = c.FormValue("e_9487")
			CONTENT = `//*[@id="e_9487"]`
		case "3":
			SERVICE = `//*[@id="service_3"]`
			content = c.FormValue("e_9751")
			CONTENT = `//*[@id="e_9751"]`
		case "4":
			SERVICE = `//*[@id="service_4"]`
			content = c.FormValue("e_9753")
			CONTENT = `//*[@id="e_9753"]`
		case "5":
			SERVICE = `//*[@id="service_5"]`
			content = c.FormValue("e_9755")
			CONTENT = `//*[@id="e_9755"]`
		}

		email := c.FormValue("email")
		name := c.FormValue("name")
		furigana := c.FormValue("furigana")
		age := c.FormValue("age")

		gender := c.FormValue("gender")
		if gender == "1" {
			gender = `//*[@id="tmp_0"]`
		} else {
			gender = `//*[@id="tmp_1"]`
		}

		postal_code1 := c.FormValue("postal_code1")
		postal_code2 := c.FormValue("postal_code2")
		address := c.FormValue("address")
		tel1 := c.FormValue("tel1")
		tel2 := c.FormValue("tel2")
		tel3 := c.FormValue("tel3")

		// Chromeの設定
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-extensions", false),
			chromedp.UserAgent("Chrome"),
		)

		// Chromeの起動
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel = chromedp.NewContext(ctx)
		defer cancel()

		// シリアルコードを配列に格納
		lineBreak := detectLineBreak(serialCodes)
		codeArray := strings.Split(serialCodes, lineBreak)

		for _, serial := range codeArray {
			if strings.TrimSpace(serial) != "" {
				// ChromeでWebページにアクセスし、シリアルコードを入力
				fmt.Println("シリアル： ", serial)
				err := chromedp.Run(ctx,
					chromedp.Navigate("https://krs.bz/kingrecords/m/9ks32kr"),
					chromedp.WaitVisible(SERIAL_CODE),
					chromedp.SendKeys(SERIAL_CODE, serial, chromedp.BySearch),
					chromedp.Click(SERVICE, chromedp.BySearch),
					chromedp.SetValue(CONTENT, content, chromedp.BySearch),
					chromedp.SendKeys(EMAIL, email, chromedp.BySearch),
					chromedp.SendKeys(NAME, name, chromedp.BySearch),
					chromedp.SendKeys(FURIGANA, furigana, chromedp.BySearch),
					chromedp.SendKeys(AGE, age, chromedp.BySearch),
					chromedp.Click(gender, chromedp.BySearch),
					chromedp.SendKeys(POSTAL_CODE1, postal_code1, chromedp.BySearch),
					chromedp.SendKeys(POSTAL_CODE2, postal_code2, chromedp.BySearch),
					chromedp.SendKeys(ADRESS, address, chromedp.BySearch),
					chromedp.SendKeys(TEL1, tel1, chromedp.BySearch),
					chromedp.SendKeys(TEL2, tel2, chromedp.BySearch),
					chromedp.SendKeys(TEL3, tel3, chromedp.BySearch),
					chromedp.Click(`//*[@id="e_9485[value][1]"]`),
					chromedp.Click(SUBMIT),
					chromedp.WaitVisible(CONFIRM_BUTTON),
					chromedp.FullScreenshot(&buf, SCREENSHOT_QUALITY),
					chromedp.Sleep(1*time.Second),
					chromedp.Click(CONFIRM_BUTTON),
					chromedp.WaitVisible(`/html/body/button`),
				)

				if err := ioutil.WriteFile("screenshots/serial_"+serial+".jpg", buf, 0o644); err != nil {
					log.Fatal(err)
				}

				if err != nil {
					log.Println("完了！")
				}

			}
		}

		return c.String(http.StatusOK, "完了！")
	})

	// サーバーの起動
	e.Logger.Fatal(e.Start(":1323"))
}

// 改行コードを判定する関数
func detectLineBreak(data string) string {
	if strings.Contains(data, "\r\n") {
		return "\r\n"
	} else if strings.Contains(data, "\r") {
		return "\r"
	} else {
		return "\n"
	}
}
