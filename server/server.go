package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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
	fmt.Println("method:", r.Method) //获取请求的方法
	fmt.Println(r)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("upfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		//
		bool_fileexist := checkFileIsExist("./test")
		fmt.Println("check file 1-----------------")
		fmt.Println("-------------------------bool_fileexist:", bool_fileexist)
		if bool_fileexist { //如果文件夹存在
			//f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
			f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			fmt.Println("文件夹存在")
		} else { //不存在文件夹时 先创建文件夹再上传
			err1 := os.Mkdir("./test", os.ModePerm) //创建文件夹
			if err1 != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("文件夹不存在")
			fmt.Println("文件夹创建成功！")
			f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}

	}
}