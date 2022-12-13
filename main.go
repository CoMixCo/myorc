package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	fmt.Println("running...")
	http.HandleFunc("/", HelloHandler)
	http.HandleFunc("/parse", ParseHandler)
	http.ListenAndServe(":8000", nil)
}

func ParsePic(pic_file string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(pic_file)
	return client.Text()
}

// 判断是否是url
func ISUrl(link string) bool {
	if len(link) == 0 {
		return false
	}
	reg := regexp.MustCompile(`^(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`)
	return reg.MatchString(link)
}

// 保存图片
func SavePic(src string, file_name string) error {
	v, err := http.Get(src)
	if err != nil {
		fmt.Printf("Http get [%v] failed! %v", src, err)
		return err
	}
	defer v.Body.Close()
	content, read_err := io.ReadAll(v.Body)
	if read_err != nil {
		fmt.Printf("Read http response failed! %v", read_err)
		return read_err
	}
	write_err := os.WriteFile(file_name, content, 0666)
	if write_err != nil {
		fmt.Printf("Save to file failed! %v", write_err)
		return write_err
	}
	return nil
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello ocr")
}

func ParseHandler(w http.ResponseWriter, r *http.Request) {
	src := r.PostFormValue("src")
	order_no := r.PostFormValue("order_no")
	if !ISUrl(src) || len(order_no) == 0 {
		fmt.Fprintf(w, "params error")
		return
	}
	run_file, _ := os.Executable()
	file_name := fmt.Sprintf("%s/%s.jpg", filepath.Dir(run_file), order_no)
	if err := SavePic(src, file_name); err != nil {
		fmt.Fprintf(w, "save pic error: %v", err)
		return
	}
	text, parse_err := ParsePic(file_name)
	if parse_err != nil {
		fmt.Fprintf(w, "parse pic error: %v", parse_err)
		return
	}
	fmt.Fprintf(w, text)
}
