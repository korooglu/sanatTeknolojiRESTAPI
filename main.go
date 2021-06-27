package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	helper "go-multitenancy/helpers"
)

func main() {
	uName, pwd, pwdConfirm,ticket := "", "", "",""
	db, err := sql.Open("mysql", "root:1234567@/demodb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	createStatement := " 'users'('ID' int(11) NOT NULL AUTO_İNCREMENT, 'username' varchar(45) NOT NULL,'password' varchar(45) NOT NULL)"

	_, err = db.Exec("CREATE TABLE IF NOT EXİST" + createStatement)
	if err != nil {
		panic(err.Error())
	}

	var tickets []string
	mux := http.NewServeMux()

	//Signup
	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		uName = r.FormValue("username")

		pwd = r.FormValue("password")
		pwdConfirm = r.FormValue("confirm")

		uNameCheck := helper.IsEmpty(uName)

		pwdCheck := helper.IsEmpty(pwd)
		pwdConfirmCheck := helper.IsEmpty(pwdConfirm)

		if uNameCheck || pwdCheck || pwdConfirmCheck {
			fmt.Fprintf(w, "Boş alan bulunuyor.")
			return
		}

		if pwd == pwdConfirm {
			sorgu1:="INSERT INTO users(username, password) VALUES(%s,%s)"
			_, err := db.Exec(sorgu1,(uName,pwd))
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, "Kayıt Olma Başarılı")
		} else {
			fmt.Fprintf(w, "Şifreler Aynı Olmalı.")
		}

	})

	//Login
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		uName = r.FormValue("username")
		pwd = r.FormValue("password")
		uNameCheck := helper.IsEmpty(uName)
		pwdCheck := helper.IsEmpty(pwd)

		if uNameCheck || pwdCheck {
			fmt.Fprintf(w, "Boş Alanlar var.")
			return
		}
		
		sorgu:="SELECT password FROM users WHERE username=%s"
		dbPwd := db.Exec(sorgu,(uName,))

		if  pwd==dbPwd{
			fmt.Fprintln(w, "Giriş Başarılı.")

		} else {
			fmt.Fprintln(w, "Giriş Başarısız.")
		}
	})
	//pwdUpdate
	mux.HandleFunc("/pwdUpdate",func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		lastPwd:=r.FormValue("lastpwd")
		updatedPwd:=r.FormValue("updatedpwd")
		lastPwdCheck:=helper.IsEmpty(lastPwd)
		updatedPwdCheck:=helper.IsEmpty(updatedPwd)
		if lastPwdCheck || updatedPwdCheck{
			fmt.Fprintf(w,"Boş Alan Bulunuyor.")
		}else{
			sorgu2:="UPDATE users
			SET password = %s
			WHERE password = %s"
			_, err := db.Exec(sorgu2,(updatedPwd,laslastPwd))
			if err != nil {
				panic(err.Error())
			}
		}
		

	})

	//Ticket

	//Create Ticket
	mux.HandleFunc("/createTicket", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ticket = r.FormValue("ticket")
		
		ticketCheck := helper.IsEmpty(ticket)
		if ticketCheck {
			fmt.Fprintf(w, "Sorununuzu bizimle paylaşın.")
			return
		} else {
			createStatement=" 'ticket'('content' varchar(45) NOT NULL)"
			_, err = db.Exec("CREATE TABLE IF NOT EXİST" + createStatement)
			if err != nil {
			panic(err.Error())
			}
			sorgu3:="INSERT INTO ticket(content) VALUES(%s)"
			_, err := db.Exec(sorgu3,(ticket))
			if err != nil {
				panic(err.Error())
			}
			
			fmt.Fprintln(w, "Ticket Oluşturuldu, Sizinle iletişime geçilecektir.")
		}
		//List Tİckets
		var content string
		mux.HandleFunc("listTickets", func(w http.ResponseWriter, r *http.Request) {
			rows,err :=db.Query("SELECT * FROM ticket")
			if err != nil {
				panic(err.Error())
				}

			for rows.Next() {
				err = rows.Scan(&content)
				if err != nil {
					panic(err.Error())
					}
			}
			fmt.Fprintf(w,content)
		})

	})
	http.ListenAndServe(":8080", mux)

}
