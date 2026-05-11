/*Testing */
package main 

import 
(
    "fmt"
    "net/http"
    "html/template"
)

type FolderPageData struct {
    FolderName string
    Files []string
    UploadSuccess bool 
}

func session(w http.ResponseWriter, r *http.Request){
    if r.URL.Path == "/Hello"{
        w.Header().Set("Content-Type", "text/html")
        tmplt, err := template.ParseFiles("static/hello.html")
        if err != nil{
            http.Error(w, "1Internal Server Error", http.StatusInternalServerError)
            return
        }

        err = tmplt.Execute(w, nil)
        if err != nil{
            http.Error(w, "2Internal Server Error", http.StatusInternalServerError)
            return 
        }
    }else if r.URL.Path == "/"{
        if r.Method == http.MethodGet {
        tmplt, err:= template.ParseFiles("static/login.html")
        if err != nil {
            http.Error(w, "1Internal Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        err = tmplt.Execute(w, nil)
        if err != nil {
            http.Error(w, "2Internal Error", http.StatusInternalServerError)
            return 
        }
    }else if r.Method == http.MethodPost {
        Login(w, r)
    }
    }else if r.URL.Path == "/AddUser"{
        if r.Method == http.MethodGet{
            if !RequireLogin(w, r){
                return 
            }
            tmplt, err := template.ParseFiles("static/addUser.html")
            if err != nil{
                http.Error(w, "1Internal Error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "text/html")
            err = tmplt.Execute(w, nil)
            if err != nil{
                http.Error(w, "2Internal Server Error", http.StatusInternalServerError)
                return
            }

        }else if r.Method == http.MethodPost{
            if !RequireLogin(w, r){
                return 
            }

            r.ParseForm()
            adminUsername := r.FormValue("adminUsername")
            adminPassword := r.FormValue("adminPassword")
            newUsername := r.FormValue("newUsername")
            newPassword := r.FormValue("newPassword")

            wasAdded := addUserToDB(adminUsername, adminPassword, newUsername, newPassword, false)
            if !wasAdded{
                fmt.Println("Failed to add user")
                return
            }
        }
    }else if r.URL.Path == "/Home"{
        if !RequireLogin(w, r) {
            return
        }

        tmplt, err := template.ParseFiles("static/home.html")

        if err != nil{
            http.Error(w, "1Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        err = tmplt.Execute(w, nil)
        if err != nil{
            http.Error(w, "2Internal Server Error", http.StatusInternalServerError)
            return
        }
    }else if r.URL.Path == "/Calc"{
        if  validUser := RequireLogin(w, r); !validUser {
            fmt.Println(validUser)
            return
        }

        if r.Method == http.MethodGet{
            tmplt, err := template.ParseFiles("static/calc.html")
            if err != nil{
                http.Error(w, "1Internal Server Error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            err = tmplt.Execute(w, nil)
            if err != nil{
                http.Error(w, "2Internal Server Error", http.StatusInternalServerError)
            }
        }else if r.Method == http.MethodPost{
            if !RequireLogin(w, r){
                return 
            }

            r.ParseForm()
            num1 := r.FormValue("num1")
            num2 := r.FormValue("num2")
            op := r.FormValue("operation")
            answer, err := Calc(num1, num2, op)
            if err != nil { answer = 0 }
            
            finalText := fmt.Sprintf("You entered %s and %s with operation %s = %f", num1, num2, op, answer)

            w.Header().Set("Content-Type", "text/html")
            tmplt, err := template.ParseFiles("static/calcResult.html")
            if err != nil{
                http.Error(w, "1Internal Error", http.StatusInternalServerError)
                return
            }

            err = tmplt.Execute(w, finalText)
            if err != nil{
                http.Error(w, "2Internal Erro", http.StatusInternalServerError)
                return 
            }
        }else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }else if r.URL.Path == "/Games"{
        if !RequireLogin(w, r) {
            return
        }

        if r.Method == http.MethodGet{
            tmplt, err := template.ParseFiles("static/games.html")
            if err != nil{
                http.Error(w, "1Games Internal Error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("content-Type", "text/html")
            err = tmplt.Execute(w, nil)
            if err != nil{
                http.Error(w, "2Games Internal Error", http.StatusInternalServerError)
                return 
            }
        }
    }else if r.URL.Path == "/Games/TicTacToe"{
        if !RequireLogin(w, r){
            return
        }

        if r.Method == http.MethodGet{
            tmplt, err := template.ParseFiles("static/TicTacToe.html")
            if err != nil{
                http.Error(w, "1TicTacToe Internal Error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-type", "text/html")
            err = tmplt.Execute(w, nil)
            if err != nil{
                //http.Error(w, "2TicTacToe Internal Error", http.StatusInternalServerError)
                fmt.Println("Failed")
                return 
            }
        }
    }else if r.URL.Path == "/Cloud" {
        if r.Method == http.MethodGet{
            tmplt, err := template.ParseFiles("static/storageSelection.html")
            if err != nil{
                http.Error(w, "1Cloud Menu Internal Error", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-type", "text/html")
            folders := getStorageFolderNames(r);
            err = tmplt.Execute(w, folders)
            if err != nil{
                http.Error(w, "2Cloud Menu Exectute Error", http.StatusInternalServerError)
                return
            }
        }else if r.Method == http.MethodPost{
            //get folder name clicked and pass that to function to display 
            r.ParseForm()
            //folder := r.FormValue("folder")
            return
        }
    }else if r.URL.Path == "/Cloud/Upload" {
        tmplt, err := template.ParseFiles("static/uploadPage.html")
        if err != nil {
            http.Error(w, "Failed opening upload page", http.StatusInternalServerError)
            return 
        }

        w.Header().Set("Content-type", "text/html")
        err = tmplt.Execute(w, nil)
        if err != nil {
            http.Error(w, "Failed to set Header", http.StatusInternalServerError)
            return
        }
    }
}


func main(){
    fmt.Println("We Ball")

    storageFolder := "Storage"
    err := InitStorage(storageFolder)
    if err != nil{
        panic(err)
    }

    http.Handle("/files/", FilesHandler(storageFolder))
    http.Handle("/Cloud/Browse", ListFilesHandler(storageFolder))
    http.Handle("/upload", UploadHandler(storageFolder))
    http.HandleFunc("/Cloud/download", DownloadHandler(storageFolder))
    http.HandleFunc("/Cloud/delete", DeleteFileHandler(storageFolder))

    http.HandleFunc("/", session)
    fmt.Println("Listening on Port: 8080")

    //Enable File server
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    if err := http.ListenAndServeTLS(":8080", "Cert/server.crt", "Cert/server.key",  nil); err != nil {
        fmt.Println("Server Failed", err)
    }
}
