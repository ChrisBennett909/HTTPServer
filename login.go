package main

import (
  "net/http"
  "crypto/rand"
  "encoding/hex"
  "fmt"
)


func Login(w http.ResponseWriter, r *http.Request){
    if r.Method == http.MethodPost {
        r.ParseForm()
        username := r.FormValue("username")
        password := r.FormValue("password")
         

        if userPresent, _:= checkUser(username, string(password)); userPresent{
          token := generateToken()
          storeSessionToken(username, token)

          cookie := http.Cookie{
                     Name: "session",
                     Value: token,
                     Path: "/",
                     HttpOnly: false,
          }
          http.SetCookie(w, &cookie)

          fmt.Println("User Present going ot home page")
          http.Redirect(w, r, "/Home", http.StatusSeeOther)
          return
        }
        http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
    }
}


func RequireLogin(w http.ResponseWriter, r *http.Request) bool{
     cookie, err := r.Cookie("session")
     if err != nil{
       fmt.Println("User has no Cookie")
       http.Redirect(w, r, "/", http.StatusSeeOther)
       return false
     }

     token := cookie.Value
     if !checkSessionToken(token){
       http.Redirect(w, r, "/", http.StatusSeeOther)
       return false
     }

     return true
}

func generateToken() string{
  b := make([]byte, 32)
  _, err := rand.Read(b)
  if err!= nil{
    fmt.Println("Failed to generate token")
    return ""
  }

  return hex.EncodeToString(b)
}
