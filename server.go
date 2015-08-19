package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"io/ioutil"
  "net/http"
  "path/filepath"
  //"path"
  //"reflect"
  "strings"
  "time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/plimble/ace"
	//"github.com/astaxie/beego/session"
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"
	"github.com/Unknwon/paginater"
)

type Item struct {
	Id          int
	Title       string
	Date        time.Time
	Description string
	Seller      string
	Price       string
	Image       string
}




var db *sqlx.DB = SetupDB()

func SetupDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=devmenon dbname=webshop sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}



func main() {

	defer db.Close()

	router := ace.Default()
	router.Use(ace.Logger())

	router.GET("/", Home)

	router.GET("/login", Login)

	router.POST("/login", PostLogin)

	router.GET("/signup", Signup)

	router.POST("/signup", PostSignup)

	router.GET("/logout", Logout)

	router.GET("/authfail", Authfail)

	router.GET("/user", User)

	router.GET("/item/:id", Itemid)

	router.GET("/public/css/:cssfile", fileLoadHandler)

	router.Run(":2500")
}

func Home(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	ShowItems(w, r)

}

func Login(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	SimplePage(w, r, "try")

}

func Authfail(c *ace.C) {

	log.Print("Authorization Failed")
}

func User(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	SimpleAuthenticatedPage(w, r, "user")
}

func Itemid(c *ace.C) {

	url := c.Param("id")
	str := strings.Trim(url, ":")
	fmt.Println(str)

	ShowItemid(c, str)

}

func Signup(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	SimplePage(w, r, "signup")
}

// end of page handlers
//action handlers begin

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {




	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func PostLogin(c *ace.C) {
	var w = c.Writer
	var req = c.Request



	var password_in_database string
	var email string

	username, password := req.FormValue("inputUsername"), req.FormValue("inputPassword")
	err := db.QueryRow("SELECT email, password FROM users WHERE username=$1", username).Scan(&email, &password_in_database)
	if err == sql.ErrNoRows {
		http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_in_database), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}




	http.Redirect(w, req, "/user", 302)
}

func PostSignup(c *ace.C) {
	var w = c.Writer
	var req = c.Request

	username := req.FormValue("reg_username")
	password := req.FormValue("reg_password")
	email := req.FormValue("reg_email")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, string(hashedPassword), email)
	if err != nil {
		log.Print(err)
	}





	http.Redirect(w, req, "/login", 302)
}

func Logout(c *ace.C) {
	var w = c.Writer
	var req = c.Request


	http.Redirect(w, req, "/", 302)
}

func ShowItems(w http.ResponseWriter, r *http.Request) {

	// Loop through rows using only one struct
	item := Item{}
	rows, err := db.Queryx("SELECT * FROM items")
	for rows.Next() {
		err = rows.StructScan(&item)
		if err != nil {
			log.Print("is this the problem?")
			log.Fatalln(err)
		}

//render := render.New(render.Options{
	//	IndentJSON: true,
	//})



 p := paginater.New(45, 10, 3, 3)
 Page := p

//render.HTML(w, http.StatusOK, "home", &item)
t, _ := template.ParseFiles("templates/home.tmpl")
	 t.Execute(w, Page)

}

}

func ShowItemid(c *ace.C, s string) {
	var w = c.Writer
	var r = c.Request
	fmt.Println(s)
	itemid := Item{}
	err := db.Get(&itemid, "SELECT * FROM items WHERE id=$1", s)
	if err != nil {
		log.Print("This must be the problem")
	}

	fmt.Printf("%#v\n", itemid.Id)
	if itemid.Id == 0 {
		fmt.Println("My Probz")

		http.Redirect(w, r, "/", 302)
	}

	render := render.New(render.Options{
		IndentJSON: true,
	})
	render.HTML(c.Writer, http.StatusOK, "item", &itemid)
}

func fileLoadHandler(c *ace.C) {
	url := c.Param("cssfile")
	str := strings.TrimLeft(url, ":")
	w := c.Writer

 baseDir, _ := filepath.Abs("/Users/devmenon/golang/src/webshop/public/css/")

 fileName := "/"+str
 fmt.Println(fileName)
 file, err := ioutil.ReadFile(fmt.Sprintf("%s%s", baseDir, fileName))

 // Here we set the header.. Without this the browser
 // won't use you css
 w.Header().Set("Content-Type", "text/css")
 fmt.Fprint(w, string(file))
 if err != nil {
  fmt.Println("Error reading file: ")
  fmt.Println(err)
 }
}
