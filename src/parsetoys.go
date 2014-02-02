package main

import (
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"flag"
	"encoding/json"
)

type GlobalSet struct {
	Donors []Donor `json:"Donors"`
}

type Donor struct {
	MainUrl string `json:"MainUrl"`
	CurrentLink string `json:"CurrentLink"`
	ProductName []string `json:"ProductName"`
	NewProduct []string  `json:"NewProduct"`
	HitProduct []string  `json:"HitProduct"`
}

type Tovar struct {
	NameProduct string
	Action_acia string
	Action_new string
	Price string
	Art string
	Pol string
	Age string
	Descrip string
	Pic []string
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
		if(isNotInBase(donor)){
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

func getDataFromDOM(s *gq.Selection, arr []string) string{
	var dt string
	if(arr[0] == "text"){
		dt = s.Text()
	}else{
		dt, _ = s.Attr(arr[0])
	}
	return encode_string(dt)
}

func getAndSaveTovarIfExist(doc *gq.Document, donor Donor){
	var curTovar Tovar
	doc.Find(donor.ProductName[1]).Each(func(i int, s *gq.Selection) {
		curTovar.NameProduct = getDataFromDOM(s, donor.ProductName)
	})
	doc.Find(donor.NewProduct[1]).Each(func(i int, s *gq.Selection) {
		if(getDataFromDOM(s, donor.NewProduct) == donor.NewProduct[2]){
			curTovar.Action_new = "1"
		}
	})
	doc.Find(donor.HitProduct[1]).Each(func(i int, s *gq.Selection) {
		if(getDataFromDOM(s, donor.HitProduct) == donor.HitProduct[2]){
			curTovar.Action_acia = "1"
		}
	})
	fmt.Printf("ТОВАР: %+v\n", curTovar.Action_acia)
	fmt.Printf("ТОВАР2: %+v\n", curTovar.Action_new)
}

/*
func get_full_info_product_2(id int, link string) {
	var action_acia, action_new, price, art, pol, age, desc string = "none", "none", "none", "none", "none", "none", "none"
	var pic string = "pic;"
	doc, err := gq.NewDocument(SITE + link)
	error_log(err)
	db, err := sql.Open("mysql", DB)
	error_log(err)
	tx, err := db.Begin()
	error_log(err)
	doc.Find("div.product>a>img").Each(func(i int, s *gq.Selection) {
		val, _ := s.Attr("src")
		if val == "/tpl/icon/act.png" {
			action_acia = "акция"
		}
		if val == "/tpl/icon/new.png" {
			action_new = "новое"
		}
		fmt.Printf("%s %s %s \n ", val, action_acia, action_new)
	})
	doc.Find("div.product div.fll.price span").Each(func(i int, s *gq.Selection) {
		val := encode_string(s.Text())
		price = val
		fmt.Printf("%s \n", price)
	})
	doc.Find("div.product div.details div.row>div.row").Each(func(i int, s *gq.Selection) {
		var temp string = "no"
		s.Find("strong").Each(func(i int, sel *gq.Selection) {
			temp = encode_string(sel.Text())
		})
		if temp == "Артикул:" {
			s.Find("p").Each(func(i int, sel *gq.Selection) {
				art = encode_string(sel.Text())
			})
		}
		if temp == "Пол:" {
			s.Find("p").Each(func(i int, sel *gq.Selection) {
				pol = encode_string(sel.Text())
			})
		}
		if temp == "Возраст:" {
			s.Find("p").Each(func(i int, sel *gq.Selection) {
				age = encode_string(sel.Text())
			})
		}
		fmt.Printf("%s \n", temp)
		fmt.Printf("%s \n", art)
		fmt.Printf("%s \n", pol)
		fmt.Printf("%s \n", age)
	})
	doc.Find("div.product div.info>p").Each(func(i int, s *gq.Selection) {
		desc = encode_string(s.Text())
		fmt.Printf("%s \n", desc)
	})
	doc.Find("div.product div.row.gallery  img").Each(func(i int, s *gq.Selection) {
		val, _ := s.Attr("src")
		pic = pic + val + ";"
		fmt.Printf("%s \n", pic)
	})
	//fmt.Println("INSERT INTO product_full SET art='"+wsc(art)+"', pol='"+wsc(pol)+"',age='"+wsc(age)+"', opis='"+wsc(desc)+"',price='"+wsc(price)+"', acia='"+wsc(action_acia)+"',newtov='"+wsc(action_new)+"', pic='"+wsc(pic)+"'")
	_, err = db.Exec("INSERT INTO product_full SET art='"+wsc(art)+"', pol='"+wsc(pol)+"',age='"+wsc(age)+"', opis='"+wsc(desc)+".',price='"+wsc(price)+"', acia='"+wsc(action_acia)+"',newtov='"+wsc(action_new)+"', pic='"+wsc(pic)+"', parent_id=?", id)
	error_log(err)
	_, err = db.Exec("UPDATE product_list SET status = 1 WHERE id =?", id)
	error_log(err)

	tx.Commit()
	defer db.Close()
}
*/

func init(){
	flag.StringVar(&settingsFile, "setfile", "/data/donors.json", "Путь до файла настроек")
	flag.StringVar(&mongoURL, "mdb","127.0.0.1:27017", "Путь до монги")
}

func main() {
	var gs GlobalSet
	flag.Parse()
	dropOldLinks()
	text, err := ioutil.ReadFile(settingsFile)
	err = json.Unmarshal(text, &gs)
	if err != nil{
		log.Fatal(err)
	}
	for _, v := range gs.Donors{
		crawler(v)
		fmt.Printf("%+v\n", v)
	}

fmt.Printf("%+v\n", gs)
}
