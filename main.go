package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"

	models "otten/models"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		data := makeRequest()
		response := buildResponse(data)

		return c.JSON(200, response)
	})

	e.Logger.Fatal(e.Start(":3000"))
}

func makeRequest() models.Data {
	url := "https://gist.githubusercontent.com/nubors/eecf5b8dc838d4e6cc9de9f7b5db236f/raw/d34e1823906d3ab36ccc2e687fcafedf3eacfac9/jne-awb.html"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("code: %d message:%s\n", res.StatusCode, res.Status)
	}

	histories := parsingHtml(res.Body)
	data := buildData(histories)

	return data
}

func parsingHtml(body io.ReadCloser) []models.History {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	var histories []models.History
	doc.Find(".main-content").Children().Each(func(i int, sel *goquery.Selection) {
		if i == 3 {
			sel.Find("tbody").Children().Each(func(j int, el *goquery.Selection) {
				createdAt := el.Children().Text()[0:16]
				description := el.Children().Text()[16:]
				year, _ := strconv.Atoi(createdAt[6:10])
				month, _ := strconv.Atoi(createdAt[3:5])
				day, _ := strconv.Atoi(createdAt[0:2])
				hour, _ := strconv.Atoi(createdAt[11:13])
				minute, _ := strconv.Atoi(createdAt[14:])

				formattedMonth := monthFormat(month)

				t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Now().Local().Location())
				h := models.History{
					CreatedAt:   t.Format(models.TIME_FORMAT_1),
					Description: description,
					Formatted: models.Format{
						CreatedAt: t.Format("02 " + formattedMonth + " 2006, 15:04 WIB"),
					},
				}
				if j == 0 {
					histories = append(histories, h)
				} else {
					var temp []models.History
					temp = append(temp, h)
					temp = append(temp, histories...)
					histories = make([]models.History, 0)
					histories = append(histories, temp...)
				}
			})
		}
	})
	return histories
}

func buildData(histories []models.History) models.Data {
	split1 := strings.Split(histories[0].Description, "DELIVERED TO [")[1]
	split2 := strings.Split(split1, "  ")[0]

	data := models.Data{
		ReceivedBy: split2,
		Histories:  histories,
	}

	return data
}

func buildResponse(data models.Data) map[string]interface{} {
	response := map[string]interface{}{
		"status": models.Status{
			Code:    "060101",
			Message: "Delivery tracking detail fetched successfully",
		},
		"data": data,
	}

	return response
}

func monthFormat(month int) string {
	var formattedMonth string

	switch month {
	case 1:
		formattedMonth = "Januari"
	case 2:
		formattedMonth = "Februari"
	case 3:
		formattedMonth = "Maret"
	case 4:
		formattedMonth = "April"
	case 5:
		formattedMonth = "Mei"
	case 6:
		formattedMonth = "Juni"
	case 7:
		formattedMonth = "Juli"
	case 8:
		formattedMonth = "Agustus"
	case 9:
		formattedMonth = "September"
	case 10:
		formattedMonth = "Oktober"
	case 11:
		formattedMonth = "November"
	case 12:
		formattedMonth = "Desember"
	}

	return formattedMonth
}
