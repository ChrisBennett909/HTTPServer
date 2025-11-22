package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "html/template"
    "log"
    "strings"
)


func InitStorage(folder string) error {
    return os.MkdirAll(folder, os.ModePerm)
}

func UploadHandler(storageFolder string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            fmt.Fprintf(w, "Only POST requests are allowed")
            return
        }

        r.ParseMultipartForm(10<<20) //10 MB max per file 

        folder := r.URL.Query().Get("folder")
        if folder == ""{
            http.Error(w, "No folder specified", http.StatusBadRequest)
            return
        }

        file, handler, err := r.FormFile("file")
        if err != nil{
            http.Error(w, "Error reading file: "+err.Error(), http.StatusBadRequest)
            return
        }
        defer file.Close()

        dstPath := filepath.Join(storageFolder, folder, handler.Filename)
        dst, err := os.Create(dstPath)
        if err != nil {
            http.Error(w, "Unable to create file: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer dst.Close()

        _, err = io.Copy(dst, file)
        if err != nil {
            http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
            return
        }

        //fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
        http.Redirect(w, r, "/Cloud/Browse?folder="+folder+"&success=true", http.StatusSeeOther)
    }
}

func FilesHandler(storageFolder string) http.Handler {
    return http.StripPrefix("/files/", http.FileServer(http.Dir(storageFolder)))
}

func ListFilesHandler(storageFolder string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        folder := r.URL.Query().Get("folder")
        if folder == ""{
            http.Error(w, "No folder specified", http.StatusBadRequest)
            return
        }
        folder = filepath.Clean(folder)
        if folder == "." || strings.Contains(folder, ".."){
            http.Error(w, "Invalid folder", http.StatusBadRequest)
            return
        }

        files, err := os.ReadDir(filepath.Join(storageFolder, folder))
        if err!= nil {
            http.Error(w, "Unable to read cloud files", http.StatusInternalServerError)
            return
        }

        var fileNames []string
        for _, f:= range files {
            fileNames = append(fileNames, f.Name())
        }

        uploadSuccess := r.URL.Query().Get("success") == "true"
        

        data := FolderPageData{
            FolderName: folder,
            Files: fileNames,
            UploadSuccess: uploadSuccess,
        }

        tmplt, err := template.ParseFiles("static/fileBrowserPage.html")
        if err != nil {
            http.Error(w, "Template Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-type", "text/html")
        err = tmplt.Execute(w, data)
        if err != nil{
            log.Println("Template execution error: ", err)
        }
    }
}

func DownloadHandler(storageFolder string) http.HandlerFunc {
    return func(w http.ResponseWriter, r* http.Request){
        filename := r.URL.Query().Get("file")
        if filename == ""{
            http.Error(w, "Missing File", http.StatusBadRequest)
            return
        }

        filename = filepath.Clean(filename)
        if strings.Contains(filename, ".."){
            http.Error(w, "Invalid file path", http.StatusBadRequest)
            return
        }

        filePath := filepath.Join(storageFolder, filename)
        file, err := os.Open(filePath)
        if err != nil{
            http.Error(w, "File not found", http.StatusNotFound)
            return
        }
        defer file.Close()

        w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filename))
        w.Header().Set("Content-Type", "application/octet-stream")

        _, err = io.Copy(w, file)
        if err != nil{
            http.Error(w, "Error sending file", http.StatusInternalServerError)
        }
    }
}

func getStorageFolderNames(r *http.Request) []string {
    folderPath := "Storage"
    entries, err := os.ReadDir(folderPath)

    if err != nil {
        log.Println("Error reading storage folder: ", err)
        return []string{}
    }

    var folders []string
    for _, entry := range entries{
        if entry.IsDir() {
            folders = append(folders, entry.Name())
        }
    }

    cookie, err := r.Cookie("session")
    if err != nil{
        fmt.Println("User has no cookie")
        return []string{}
    }

    token := cookie.Value
    if checkAdmin(token){
        return folders
    }else{
        username := getUsername(token)
        if contains(folders, username){
            return []string{username} 
        }
    }

    return []string{}
}

func contains(slice []string, str string) bool{
    for _, s := range slice{
        if s == str{
            return true
        }
    }
    return false 
}
