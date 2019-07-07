package main

import (
	htmlp "./htmlParser"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/urfave/cli.v1"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

var link string = ""
var rcpt string
var ttime string
var pages string
var fileName string
var port string

var KnownAds = make(map[string]string)
var UrlsToScan = make([]string, 0)

func main() {
	app := cli.NewApp()
	app.Name = "jfind"
	app.Usage = "Finds new car ads on avto.net"
	app.Author = "Domen Jesenovec"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "link, l",
			Usage:       "Finds car ads with predefined filter",
			Destination: &link,
		},
		cli.StringFlag{
			Name:        "rcpt, r",
			Usage:       "Recipient of email notifications",
			Destination: &rcpt,
		},
		cli.StringFlag{
			Name:        "ttime, t",
			Value:       "6",
			Usage:       "Interval on witch the program triggers",
			Destination: &ttime,
		},
		cli.StringFlag{
			Name:        "pages, p",
			Value:       "2",
			Usage:       "NUM Number of pages to search",
			Destination: &pages,
		},
		cli.StringFlag{
			Name:        "port, s",
			Value:       "8080",
			Usage:       "Run as a service",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "file, f",
			Usage:       "Read URLs form specified fileName",
			Destination: &fileName,
		},
	}

	app.Action = func(c *cli.Context) error {
		realTime, err := strconv.Atoi(ttime)

		if err != nil {
			return err
		}

		if fileName != "" {
			file, _ := os.OpenFile(fileName, os.O_RDONLY, 0660)
			defer file.Close()

			scanner := bufio.NewScanner(file)

			for ; scanner.Scan(); {
				UrlsToScan = append(UrlsToScan, scanner.Text())
			}

			if link != "" {
				UrlsToScan = append(UrlsToScan, link)
			}
		}

		if len(UrlsToScan) < 1 {
			fmt.Println("ERROR: Please specify the --link flag or --file flag")
			return nil
		}

		schedule(getData, time.Duration(realTime)*time.Hour)

		router := mux.NewRouter().StrictSlash(true)

		//status
		router.HandleFunc("/", GetStatus)
		corsObj := handlers.AllowedOrigins([]string{"*"})
		log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(corsObj)(router)))

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func getData() {
	file, _ := os.OpenFile("known.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for ; scanner.Scan(); {
		text := scanner.Text()
		KnownAds[text] = text
	}

	realPages, err := strconv.Atoi(pages)

	if err != nil {
		panic(err)
	}

	for j := 0; j < len(UrlsToScan); j++ {
		for i := 1; i <= realPages; i++ {
			fmt.Println("Getting: " + UrlsToScan[j] + strconv.Itoa(i))

			ads := htmlp.FindAds(UrlsToScan[j] + strconv.Itoa(i))

			body := ""

			for _, v := range ads {
				if KnownAds[v.Hash] == "" {
					KnownAds[v.Hash] = v.Hash

					_, err := file.WriteString(v.Hash + "\n")

					if err != nil {
						panic(err)
					}

					body += v.GetMailContents()
				}
			}

			if body != "" {
				sendMail(body, rcpt)
			}
		}
	}
}

func sendMail(body string, rcpt string) {
	from := "GMAIL"
	pass := "APP-PASSWORD"

	msg := "From: " + from + "\r\n" +
		"To: " + rcpt + "\r\n" +
		"Subject: Novi Avtomobili\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<body>" + body + "</body>"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{rcpt}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit whatever")
}

func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("works")
}
