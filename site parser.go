package main

import (
	"fmt"
	"golang.org/x/net/html"
	_ "golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

type site struct {
	name string
	url  string
}

func main() {
	var s site
	s.url = "https://www.citilink.ru/search/?text=смартфоны"
	s.name = "citilink"

	fmt.Printf("Это программа парсер, которая парсит сайт \"%s\" и выводит все смартфоны\n", s.url)

	response, err := http.Get(s.url)

	//if there was an error, report it and exit
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	////check response status code
	//if response.StatusCode != http.StatusOK {
	//	log.Fatalf("Response status code was %d\n", response.StatusCode)
	//}

	//check response content type
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html"){
		log.Fatalf("Response content type was %s (not text/html)\n",contentType)
	}

	//loop until we find the title element and its content
	//or encounter an error (which includes the end of the stream)
	htmlTokens := html.NewTokenizer(response.Body)
	for {
		tokenType := htmlTokens.Next()
		//if it's an error token, we either reached
		//the end of the file, or the HTML was malformed
		if tokenType==html.ErrorToken{
			err:=htmlTokens.Err()
			if err==io.EOF{
				break
			}
			log.Fatalf("error tokenizing html: %v",htmlTokens.Err())
		}
		//Todo figure out how to get product fields
		////the token has text content withing tag
		//if tokenType==html.TextToken{
		//	token:=htmlTokens.Token()
		//	if strings.HasPrefix(token.Data,"<div class=\"product_data__gtm-js product_data__pageevents-js"){
		//		fmt.Println(token.Data)
		//	}
		//}
	}

}
