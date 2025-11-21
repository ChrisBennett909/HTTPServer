package main

import (
  "net/http"
  "fmt"
)


func Login(w http.ResponseWriter, r *http.Request){
    if r.Method == http.MethodPost {
        r.ParseForm()
        username := r.FormValue("username")
        password := r.FormValue("password")
        folder := r.FormValue("folder")
         

        if userPresent, _:= checkUser(username, string(password)); userPresent{
             cookie := http.Cookie{
                   Name: "session",
                   Value: "loggedin", 
                   Path: "/",
                   HttpOnly: true,
               }
              http.SetCookie(w, &cookie)
              if folder == nil{
                http.Redirect(w, r, "/Home", http.StatusSeeOther)
              }else if folder != nil{
                //Redirect to folder cloud menu
              }
              return
          }
          http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
       }
}

func storageLogin(w http.ResponseWriter, r *http.Request, rightUser string, nextPage string) bool {
  //Check specific user 
  /*
  Get password
  display login page and take username and password
  if user is correct user or admin user return true
  if not return false
  */
  r.ParseForm()
  username := r.FormValue("username")
  password := r.FormValue("password")
  folder := r.FormValue("folder")

  if _, isAdmin := checkUser(rightUser, string(password)); isAdmin{
    return true
  }

  if username != rightUser{
    fmt.Println("Wronge User for folder")
    return false
  }

  if (checkUserPassword(username, password)){
    return true
  }

  return false
}

func checkSingleUser(rightUser string, username string, password string) bool {
  return false 
}

func RequireLogin(w http.ResponseWriter, r *http.Request) bool{
     cookie, err := r.Cookie("session")

     if err != nil || cookie.Value != "loggedin"{
         http.Redirect(w, r, "/", http.StatusSeeOther)
         return false
     }
     return true;
}
