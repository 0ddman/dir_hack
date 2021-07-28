package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/axgle/mahonia"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	url string
	ext string
	all bool
)

var xpath=[]string{
	"/html/body/table/tbody/tr/td/a/@href",
	"/html/body/ul/li/a/@href",
	"/html/body/pre/a/@href",
}
var namelist = []string{
	"hack",
	"shell",
	"0",
	"1",
	"2",
	"3",
	"connect",
	"adminer",
	"upload",
	"db",
	"x",
}

var extname = []string{
	"zip",
	"rar",
	"7z",
	"tar",
	"gz",
	"xz",
	"bz2",
	"inc",
	"bak",
	"mdb",
	"sql",
	"db",
}

func main() {
	flag.StringVar(&url, "url", "", "dir url")
	flag.StringVar(&ext, "ext", "", "file ext")
	flag.BoolVar(&all, "all", true, "true:get all file / false:auto list")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `example :
	dirhack: ./dir -url https://www.t00ls.net/uploads/ -ext php  -all=true/-all=false 
	
`)
		flag.PrintDefaults()
	}
	flag.Parse()
	if url != "" && ext != "" {
		getBodystr(url, ext, all)
	} else {
		flag.Usage()
	}

}

func getBodystr(url string, extInput string, all bool) {
	txt,n:=getText(url)
	for key, _ := range txt {
		if n==2{
			if key >4 {
				str1 := htmlquery.InnerText(txt[key])
				if strings.HasSuffix(str1, "/") {
					url_new := url + str1
					getBodystr(url_new, extInput, all)
				} else {
					fileCheck(url, str1, extInput, all)
				}
			}
		}else{
			if key > 0 {
				str1 := htmlquery.InnerText(txt[key])
				if strings.HasSuffix(str1, "/") {
					url_new := url + str1
					getBodystr(url_new, extInput, all)
				} else {
					fileCheck(url, str1, extInput, all)
				}
			}
		}
	}
}

func getText(url string)([] *html.Node, int){
	r, err := http.Get(url)
	if err != nil {
		fmt.Println("Network error,check your url and try again!")
		os.Exit(-1)
	}
	defer func() { _ = r.Body.Close() }()
	body, _ := ioutil.ReadAll(r.Body)
	bodystr := mahonia.NewDecoder("gbk").ConvertString(string(body))
	root, err := htmlquery.Parse(strings.NewReader(bodystr))
	if err != nil {
		fmt.Println("parse error!")
	}
	var txt []*html.Node
	var n int
	for num,value:= range xpath{
		n=num
		txt = htmlquery.Find(root, value)
		if txt !=nil{
			break
		}
	}
	return txt,n
}

func fileCheck(url string, fileFullname string, extInput string, all bool) {
	a := strings.Split(fileFullname, ".")
	extLower := strings.ToLower(a[len(a)-1])
	nameLower := strings.ToLower(a[0])
	if all {
		if extLower == extInput {
			fmt.Println(url + fileFullname)
		}
		for _, value := range extname {
			if extLower == value {
				fmt.Println(url + fileFullname)
			}
		}
	} else {
		for _, value := range namelist {
			if nameLower == value && extLower == extInput {
				fmt.Println(url + fileFullname)
			}
		}
		for _, value := range extname {
			if extLower == value {
				fmt.Println(url + fileFullname)
			}
		}
	}
}
