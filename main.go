package main

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/google/uuid"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"strings"
)

const (
	DIR = "imgs"
)

var q, _ = queue.New(
			2,                                           // Number of consumer threads
			&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
		)

func resolvHtml()  {
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	c.OnHTML(".postContent", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("href"))
		e.ForEach("img", func(i int, element *colly.HTMLElement) {
			//e.Request.Visit(element.Attr("src"))
			q.AddURL(element.Attr("src"))
		})
	})
	c.OnResponse(func(resp *colly.Response) {
		if strings.Contains(resp.Headers.Get("Content-Type"), "image/jpeg") {
			download(resp.Body)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	q.Run(c)
}

func download(img []byte) {
	guid := uuid.New()
	file_name := guid.String() + ".jpg"
	out, _ := os.Create(file_name)
	io.Copy(out, bytes.NewReader(img))
}

func main() {
	app := &cli.App{
		Name:    "clicmd",
		Usage:   "free cli cmd",
		Version: "1.0.1",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "s",
				Value: 1,
				Usage: "起始页",
			},
			&cli.IntFlag{Name: "e",
				Value: 100,
				Usage: "截止页",
			},
		},
	}
	app.Action = func(context *cli.Context) error {
		if args := context.Args(); len(args) > 0 {
			return fmt.Errorf("invalid command: %q", args.Get(0))
		}
		start := context.Int("s")
		log.Println("起始页:", start)
		end := context.Int("e")
		log.Println("截止页:", end)

		for i := start; i <= end; i++ {
			url := fmt.Sprintf("https://www.meizitu.com/a/%d.html", i)
			q.AddURL(url)
		}
		resolvHtml()
		return nil

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
