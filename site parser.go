package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

type Smartphone struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	smartphones, err := getTopTenSmartphones()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, smartphone := range smartphones {
		fmt.Println(fmt.Sprintf("%s \tURL: %s", smartphone.Name, smartphone.URL))
	}
}

func getTopTenSmartphones() ([]*Smartphone, error) {
	query := `
		LET doc = DOCUMENT("https://www.citilink.ru/catalog/smartfony")

		FOR el IN ELEMENTS(doc, ".ProductCardHorizontal__header-block")
			LIMIT 10
			LET url = ELEMENT(el, "a")
			LET name = ELEMENT(el, "a")
			LET price = ELEMENT(el,".ProductCardHorizontal__price_current-price")
		
			RETURN{
				name: TRIM(name.attributes.title),
				url: "https://www.citilink.ru" + url.attributes.href
			}
	`
	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime does not know about any HTML drivers
	// all HTML manipulations are done via functions from standard library
	// that assume that at least one driver is available
	ctx = drivers.WithContext(ctx, cdp.NewDriver())
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())

	out, err := program.Run(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]*Smartphone, 0, 10)

	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
