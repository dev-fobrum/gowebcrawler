package main

import (
	"flag"
	"gowebcrawler/crawler"
)

var (
	link string
)

func init() {
	flag.StringVar(&link, "url", "https://aprendagolang.com.br", "Url para iniciar as visitas")
}

func main() {
	flag.Parse()
	crawler.Start(link)
}
