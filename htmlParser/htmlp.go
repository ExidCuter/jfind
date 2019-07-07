package htmlParser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"strings"
)

type CarAd struct {
	Title            string
	FirstReg         string
	Km               string
	EngineType       string
	TransmissionType string
	Price            string
	Link             string
	ImgLink          string
	Hash             string
}

func getAds(url string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res.Status)

	return string(resBody)
}

func FindAds(url string) []CarAd {
	ads := make([]CarAd, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(getAds(url)))
	if err == nil {
		doc.Find(".ResultsAd").Each(func(i int, s *goquery.Selection) {
			carAd := CarAd{}

			s.Find(".ResultsAdDataTop").Each(func(i int, ss *goquery.Selection) {
				a := ss.Find("a")

				carAd.Title = a.Text()
				carAd.Link, _ = a.Attr("href")

				carAd.Link = "https://avto.net" + strings.Replace(carAd.Link, "..", "", 1)

				ss.Find("li").Each(func(i int, sss *goquery.Selection) {
					switch i {
					case 0:
						carAd.FirstReg = sss.Text()
						break
					case 1:
						carAd.Km = sss.Text()
						break
					case 2:
						carAd.EngineType = sss.Text()
						break
					case 3:
						carAd.TransmissionType = sss.Text()
						break
					}
				})
			})

			src, _:=s.Find(".ResultsAdPhotoTop").Find("img").Attr("src")

			carAd.ImgLink = src

			carAd.Price = s.Find(".ResultsAdPrice").Text()
			carAd.Price = strings.Replace(carAd.Price, "\n", "", -1)
			carAd.Price = strings.Replace(carAd.Price, "\t", "", -1)
			carAd.Price = strings.Replace(carAd.Price, " ", "", -1)

			carAd.Hash = fmt.Sprint(hash(carAd.Title + carAd.Price))

			ads = append(ads, carAd)
		})
	}

	return ads
}

func hash(text string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

func (carAd CarAd) GetMailContents() string {
	return "<table style=\"width: 100%;\">" + "<tbody>"+
		"<tr><td colspan=\"3\">"+ "<a href=\"" + carAd.Link + "\"> <h3>" + carAd.Title + "</h3></a>" + "</td></tr>" +
		"<tr>" +
			"<td>" +
				"<img src=\""+carAd.ImgLink+"\" />" +
			"</td>" +
			"<td>" +
				"<ul>" +
					"<li>" + carAd.FirstReg + "</li>" +
					"<li>" + carAd.Km + "</li>" +
					"<li>" + carAd.EngineType + "</li>" +
					"<li>" + carAd.TransmissionType + "</li>" +
				"</ul>" +
			"</td>" +
			"<td>" +
				carAd.Price +
			"</td>" +
		"</tr>"+ "</tbody>" +
		"</table>"
}
