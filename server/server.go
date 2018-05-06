package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"strings"
)

const FILE_DIR  = "/data/static/img/"

const ALI_STATIC_URL  = "http://static.xiaohuo.me/img/"

//检查目录是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Print(filename + " not exist")
		exist = false
	}
	return exist
}

func main() {
	http.HandleFunc("/xupload", upload)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	log.Println("success!")
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {

	type Ret struct {
		Status	string
		Message	string
	}

	//fmt.Println("method:", r.Method) //获取请求的方法
	//fmt.Println(r)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("upfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		bool_fileexist := checkFileIsExist(FILE_DIR)
		if bool_fileexist { //如果文件夹存在
			f, err := os.OpenFile(FILE_DIR+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		} else { //不存在文件夹时 先创建文件夹再上传
			err1 := os.Mkdir(FILE_DIR, os.ModePerm) //创建文件夹
			if err1 != nil {
				fmt.Println(err)
				return
			}

			f, err := os.OpenFile(FILE_DIR+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}

		ret := Ret{"200", ALI_STATIC_URL + handler.Filename}
		result, err := json.Marshal(ret)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintln(w,strings.ToLower(string(result)))
	}
}