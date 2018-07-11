package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var IMG_DIRECTORY string = "img/"

func uploadTemplate(message string) string {
	return `<!DOCTYPE html>
               <html>
                   <head></head>
                   <body>
                       <form enctype="multipart/form-data" action="/upload" method="POST">
                           <h3>File Upload</h3>
                           <input type="file" placeholder="uploadfile" name="uploadfile"><br>
                           <input type="submit" value="Upload">
                           <div>` + message + `</div>
                       </form>
                   <body>
               </html>`
}

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		fmt.Fprintf(w, uploadTemplate(""))
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if nil != err {
		logger.Error(err)
		fmt.Fprintf(w, uploadTemplate(err.Error()))
		return
	}

	defer file.Close()

	f, err := os.OpenFile(IMG_DIRECTORY+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(err)
		fmt.Fprintf(w, uploadTemplate(err.Error()))
		return
	}
	defer f.Close()
	io.Copy(f, file)

	fmt.Fprintf(w, uploadTemplate(`{"status":"ok"}`))
}
