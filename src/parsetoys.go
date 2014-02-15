package main

import (
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"flag"
	"encoding/json"
)

type GlobalSet struct {
	Donors []Donor `json:"Donors"`
}

type Donor struct {
	Code                  string `json:"Code"`
	MainUrl               string `json:"MainUrl"`
	CurrentLink           string `json:"CurrentLink"`
	ProductName           []string `json:"ProductName"`
	NewProduct            []string  `json:"NewProduct"`
	HitProduct            []string  `json:"HitProduct"`
	Price                 []string  `json:"Price"`
	Articul               []string  `json:"Articul"`
	Sex                   []string  `json:"Sex"`
	Age                   []string  `json:"Age"`
	Description           []string  `json:"Description"`
	Pictures              []string `json:"Pictures"`
}

type Tovar struct {
	Site        string
	Link        string
	NameProduct string
	Action_acia string
	Action_new  string
	Price       string
	Art         string
	Sex         string
	Age         string
	Descrip     string
	Pic         []string
}

var (
	settingsFile string
	mongoURL string
)

func crawler(donor Donor) {
	doc, err := gq.NewDocument(donor.MainUrl + donor.CurrentLink)
	error_log(err)
	doc.Find("a").Each(func(i int, s *gq.Selection) {
		donor.CurrentLink, _ = s.Attr("href")
		fmt.Println(donor.CurrentLink)
		if (!inBaseDonor(donor)) {
			addVisitedLinks(donor)
			getAndSaveTovarIfExist(doc, donor)
			if len(donor.CurrentLink) > 0 {
				if donor.CurrentLink[0] == '/' {
					crawler(donor)
				}
			}
		}
	})
}

func getDataFromDOM(s *gq.Selection, arr []string, code string) string {
	var dt string
	if (arr[0] == "text") {
		dt = s.Text()
	}else {
		dt, _ = s.Attr(arr[0])
	}
	return encode_string(dt, code)
}

func getAndSaveTovarIfExist(doc *gq.Document, donor Donor) {
	var curTovar Tovar
	doc.Find(donor.ProductName[1]).Each(func(i int, s *gq.Selection) {
		curTovar.NameProduct = getDataFromDOM(s, donor.ProductName, donor.Code)
	})
	doc.Find(donor.NewProduct[1]).Each(func(i int, s *gq.Selection) {
		if (getDataFromDOM(s, donor.NewProduct, donor.Code) == donor.NewProduct[2]) {
			curTovar.Action_new = "1"
		}else {
			curTovar.Action_new = "0"
		}
	})
	doc.Find(donor.HitProduct[1]).Each(func(i int, s *gq.Selection) {
		if (getDataFromDOM(s, donor.HitProduct, donor.Code) == donor.HitProduct[2]) {
			curTovar.Action_acia = "1"
		}else {
			curTovar.Action_acia = "0"
		}
	})
	doc.Find(donor.Price[1]).Each(func(i int, s *gq.Selection) {
		curTovar.Price = getDataFromDOM(s, donor.Price, donor.Code)
	})
	doc.Find(donor.Articul[1]).Each(func(i int, s *gq.Selection) {
		curTovar.Art = getDataFromDOM(s, donor.Articul, donor.Code)
	})
	doc.Find(donor.Sex[1]).Each(func(i int, s *gq.Selection) {
		curTovar.Sex = getDataFromDOM(s, donor.Sex, donor.Code)
	})
	doc.Find(donor.Age[1]).Each(func(i int, s *gq.Selection) {
		curTovar.Age = getDataFromDOM(s, donor.Age, donor.Code)
	})
	doc.Find(donor.Description[1]).Each(func(i int, s *gq.Selection) {
		curTovar.Descrip = getDataFromDOM(s, donor.Description, donor.Code)
	})
	doc.Find(donor.Pictures[1]).Each(func(i int, s *gq.Selection) {
		curTovar.Pic = append(curTovar.Pic, donor.MainUrl+getDataFromDOM(s, donor.Pictures, donor.Code))
	})

	fmt.Printf("Название: %+v\n", curTovar.NameProduct)
	fmt.Printf("Новый: %+v\n", curTovar.Action_acia)
	fmt.Printf("Хит: %+v\n", curTovar.Action_new)
	fmt.Printf("Цена: %+v\n", curTovar.Price)
	fmt.Printf("Артикул: %+v\n", curTovar.Art)
	fmt.Printf("Пол: %+v\n", curTovar.Sex)
	fmt.Printf("Возраст: %+v\n", curTovar.Age)
	fmt.Printf("Описание: %+v\n", curTovar.Descrip)
	fmt.Printf("Картинки: %+v\n", curTovar.Pic)
	if (curTovar.Art != "") {
		curTovar.Site = donor.MainUrl
		curTovar.Link = donor.CurrentLink
		saveTovar(curTovar)
	}
}

func init() {
	flag.StringVar(&settingsFile, "setfile", "/data/donors.json", "Путь до файла настроек")
	flag.StringVar(&mongoURL, "mdb", "127.0.0.1:27017", "Путь до монги")
}

func main() {
	var gs GlobalSet
	flag.Parse()
	dropOldLinks()
	text, err := ioutil.ReadFile(settingsFile)
	err = json.Unmarshal(text, &gs)
	error_log(err)
	for _, v := range gs.Donors {
		crawler(v)
	}
}
