
package main

import (
    "fmt"
    "database/sql"
    _"github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
    "time"
    "log"
)


func checkUser(name, pass string) (bool, bool){
    db, err := sql.Open("sqlite3", "./DB/credDB.db")
    if err != nil {
        fmt.Println("Failed to open database", err)
        return false, false
    }
    defer db.Close()

    query := "SELECT password, isAdmin FROM users WHERE username = ?"
    var storedPassword string
    var isAdmin bool
    qerr := db.QueryRow(query, name).Scan(&storedPassword, &isAdmin)

    if qerr == sql.ErrNoRows{
        fmt.Println("User does not exist")
        return false, false
    }else if qerr != nil{
        fmt.Println("Failed to query database")
        return false, false 
    }else if qerr == nil{
        err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(pass))
        if err != nil{
            return false, false
        }else {
            return true, isAdmin
        }
    }
    return false, false 
}


func addUserToDB(adminUsername, adminPassword, name, pass string, admin bool) bool {
    const maxAttemps = 5
    const retryDelay = 500 * time.Millisecond
    var command string = "INSERT INTO users (username, password, isAdmin) VALUES (?, ?, ?)"

    db, err := sql.Open("sqlite3", "./DB/credDB.db")
    if err != nil{
        fmt.Println("Failed to open DB")
        return false;
    }
    defer db.Close()

    //Check if Admin 
    query := "SELECT password, isAdmin FROM users WHERE username = ?"
    var storedPassword string
    var isAdmin bool 

    qerr := db.QueryRow(query, adminUsername).Scan(&storedPassword, &isAdmin)
    if qerr == sql.ErrNoRows{
        fmt.Println("User Does no Exist")
        return false
    }else if qerr != nil{
        fmt.Println("Faield to query data")
        return false
    }else{
        if !isAdmin{
            fmt.Println("Only Admins can add users")
            return false
        }

        err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(adminPassword))
        if err != nil{
            fmt.Println("failed to compare passwords")
            return false 
        }else{
            hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
            if err != nil{
                fmt.Println("Failed to hash password")
                return false
            }

            for attempt :=1; attempt <= maxAttemps; attempt++{

                _, err = db.Exec(command, name, string(hashedPass), admin)
                if err != nil {
                    log.Println("Failed to Add User", attempt)
                    time.Sleep(retryDelay)
                    continue
                }

                return false
            }
        }
    }

    return true;
}

func checkUserPassword(username string, password string) bool {
    return false
}
func checkAdmin(user string) bool {
    db, err := sql.Open("sqlite3", "./DB/credDB.db")
    if err != nil{
        fmt.Println("Failed to open database: ", err)
        return false;
    }
    defer db.Close();

    var isAdmin bool
    queryString := "SELECT isAdmin FROM users WHERE username == ?"
    query:= db.QueryRow(queryString, user).Scan(&isAdmin)

    if query == sql.ErrNoRows{
        fmt.Println("User does not exist")
        return false 
    }else if query != nil{
        fmt.Println("Failed to query data")
        return false
    }else{
        if isAdmin{
            fmt.Println("Admin True")
            return true
        }
    }

    return false 
}
