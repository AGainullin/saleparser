package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// SalePosition is data type, containing the necessary information about clothing
type SalePosition struct {
	URL      string
	Cost     string
	SiteName string
}

// NewSales return an array of clothing items from our url.txt file
func NewSales(fileName string) []SalePosition {

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	salesLine := strings.Split(string(file), "\r\n")
	sales := make([]SalePosition, 0)
	regularMaskForSiteName := regexp.MustCompile(`\w+\.\w+`)

	for i := 0; i < len(salesLine); i += 2 {
		siteName := strings.Split(regularMaskForSiteName.FindString(salesLine[i]), ".")
		sales = append(sales, SalePosition{URL: salesLine[i], Cost: salesLine[i+1], SiteName: siteName[0]})
	}

	return sales
}

func deleteSpace(str string) string {
	replace := regexp.MustCompile(`[[:space:]]`)
	return replace.ReplaceAllString(str, "")
}

func exit() {
	fmt.Println("Press 'q' to quit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exit := scanner.Text()
		if exit == "q" {
			break
		} else {
			fmt.Println("Press 'q' to quit")
		}
	}
}

func main() {

	sales := NewSales("url.txt")

	for i := 0; i < len(sales); i++ {

		resp, err := http.Get(sales[i].URL)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var costRegularMask *regexp.Regexp
		var realCostStr []string

		switch sales[i].SiteName {
		case "zolla":
			costRegularMask = regexp.MustCompile(`class="price-cur">\d\s?\d+`)
		case "ostin":
			costRegularMask = regexp.MustCompile(`class="o-product__price"[\w\d:;#"=\- ]+>\n\s+\d+`)
		}

		realCostStr = strings.Split(costRegularMask.FindString(string([]byte(body))), ">")
		realCostInt, err := strconv.Atoi(deleteSpace(realCostStr[1]))
		if err != nil {
			log.Fatal(err)
		}
		saleCostInt, err := strconv.Atoi(sales[i].Cost)
		if err != nil {
			log.Fatal(err)
		}

		if saleCostInt > realCostInt {
			fmt.Printf("На сайте \"%v\" cтало дешевле!!!\n----- Старая цена: %v\n----- Новая цена: %v\nСсылка: %v\n",
				sales[i].SiteName, saleCostInt, realCostInt, sales[i].URL)
		}

	}
	exit()
}
