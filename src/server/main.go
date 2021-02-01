package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"

	"./middleware"
	"./models"
	"./secrets"
	"./sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var activeAdmin models.Admin
var activeUser models.User
var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/calendar", calendarGetHandler).Methods("GET")
	r.HandleFunc("/add-address", middleware.UserAuthRequired(addAddressGetHandler)).Methods("GET")
	r.HandleFunc("/add-address", middleware.UserAuthRequired(addAddressPostHandler)).Methods("POST")
	r.HandleFunc("/list-addresses", middleware.UserAuthRequired(listAddressesGetHandler)).Methods("GET")
	r.HandleFunc("/list-addresses", middleware.UserAuthRequired(listAddressesPostHandler)).Methods("POST")
	r.HandleFunc("/delete-address", middleware.UserAuthRequired(deleteAddressPostHandler)).Methods("POST")
	r.HandleFunc("/edit-address", middleware.UserAuthRequired(editAddressPostHandler)).Methods("POST")
	r.HandleFunc("/update-address", middleware.UserAuthRequired(updateAddressPostHandler)).Methods("POST")
	r.HandleFunc("/logout", logoutGetHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	fmt.Println("Hi!")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	username, ok := session.Values["username"]
	fmt.Println(username, ok)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}
	http.Redirect(w, r, "/list-addresses", 302)
	return
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	fmt.Println(username, password)
	user, err := models.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		log.Fatal(err)

	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["username"] = username
	fmt.Println(username)
	session.Save(r, w)
	activeUser, err = models.GetUser(username)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("username: %s user id: %d\n", activeUser.Name, activeUser.ID)
	http.Redirect(w, r, "/", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal(err)

	}
	err = models.InsertUser(username, string(hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		fmt.Println("Register error occured!")
		return
	}
	http.Redirect(w, r, "/login", 302)

}

func calendarGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "calendar.html", nil)
}

func addAddressGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("username: %s user id: %d\n", activeUser.Name, activeUser.ID)
	templates.ExecuteTemplate(w, "add-address.html", struct {
		User models.User
	}{
		User: activeUser,
	})
}

func addAddressPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	il := r.PostForm.Get("il")
	ilce := r.PostForm.Get("ilce")
	mahalle := r.PostForm.Get("mahalle")
	sokak := r.PostForm.Get("sokak")
	err := models.InsertAddress(il, ilce, mahalle, sokak)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		fmt.Println("Register error occured!")
		return
	}
	http.Redirect(w, r, "/", 302)

}

func listAddressesGetHandler(w http.ResponseWriter, r *http.Request) {
	addresses := models.GetAddresses()
	templates.ExecuteTemplate(w, "list-addresses.html", struct {
		Addresses []models.Address
		User      models.User
	}{
		Addresses: addresses,
		User:      activeUser,
	})
}

func listAddressesPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	il := r.PostForm.Get("il")
	ilce := r.PostForm.Get("ilce")
	mahalle := r.PostForm.Get("mahalle")
	sokak := r.PostForm.Get("sokak")
	err := models.InsertAddress(il, ilce, mahalle, sokak)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		fmt.Println("Register error occured!")
		return
	}
	http.Redirect(w, r, "/", 302)

}

func deleteAddressPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	err := models.DeleteAddress(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		fmt.Println("Register error occured!")
		return
	}
	http.Redirect(w, r, "/list-addresses", 302)
}

func editAddressPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.PostForm.Get("id")
	address, _ := models.GetAddress(id)
	fmt.Println(address.Sokak)
	templates.ExecuteTemplate(w, "edit-address.html", struct {
		Address models.Address
		User    models.User
	}{
		Address: address,
		User:    activeUser,
	})
}

func updateAddressPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	il := r.PostForm.Get("il")
	id := r.PostForm.Get("id")
	ilce := r.PostForm.Get("ilce")
	mahalle := r.PostForm.Get("mahalle")
	sokak := r.PostForm.Get("sokak")
	err := models.UpdateAddress(id, il, ilce, mahalle, sokak)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error!"))
		fmt.Println("Register error occured!")
		return
	}
	http.Redirect(w, r, "/list-addresses", 302)
}
func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	session.Values["username"] = ""
	err := session.Save(r, w)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/login", 302)
}

func sendRecovery(to, id string) error {
	fmt.Println("sendRecovery running..")
	msg := []byte("To:  " + to + "\r\n" +
		"Subject: Password Recovery\r\n" +
		"\r\n" +
		"Follow this link to reset your password: localhost:8080/admin/reset-pass?id=" + id + "\r\n")
	err := sendMail(to, msg)
	return err
}

func sendMail(to string, msg []byte) error {
	from := secrets.GetSMTPMail()
	pass := secrets.GetSMTPPass()
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	/*
		msg := []byte("To:  " + to + "\r\n" +
			"Subject: discount Gophers!\r\n" +
			"\r\n" +
			"This is the email body.\r\n")
	*/
	err := smtp.SendMail("smtp.gmail.com:587",
		auth,
		from, []string{to}, msg)

	return err
}
