package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
	"fmt"
	"bytes"
	"mime/multipart"
	"io"
	"net/http"
	"io/ioutil"
)

func main() {
	app := cli.NewApp()

	app.Name = "xiaohuo"
	app.Usage = "a cli for me!"
	app.Version = "0.0.1"
	app.Author = "xiaohuo"
	app.Email = "xiaohuokevin@gmail.com"

	app.Commands = []cli.Command {
		{
			Name: "upload",
			Aliases: []string {"up"},
			Usage: "Upload file to ali server",
			Action: func(c *cli.Context)  error {
				upFile := c.Args().Get(0);

				if upFile == "" {
					return cli.NewExitError("ERROR:Please select the file to upload", 1)
				}
				fileInfo, err := os.Stat(upFile)

				if err != nil {
					return cli.NewExitError("ERROR:Getting file information failed", 1)
				}

				if fileInfo.IsDir() {
					return cli.NewExitError("ERROR:Folder uploads are not supported for the time being", 1)
				}

				if os.IsNotExist(err) {
					return cli.NewExitError("ERROR:file does not exist", 1)
				}

				fmt.Printf("starting to upload file %s\n", fileInfo.Name());
				fmt.Printf("fileSize is %dB\n", fileInfo.Size());

				targetUrl := "http://39.106.17.253:9090/xupload"
				status_code, res, err := postFile(upFile, targetUrl)

				fmt.Printf("upload result.... \n");
				fmt.Printf("http status code is %s\n", status_code);
				fmt.Printf("Result: file url is %s\n", res);

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func postFile(filename string, targetUrl string) (string, string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("upfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return "", "",err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return "", "", err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", "", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	fmt.Println(bodyBuf)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	//fmt.Println(resp.Status)
	//fmt.Println(string(resp_body))
	return resp.Status, string(resp_body),nil
}
