package main

import ("fmt";"net/http";"html/template")
import( "database/sql"

        _ "github.com/go-sql-driver/mysql")
//import( "database/sql")

//	_ "github.com/go-sql-driver/mysql")

func index(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
 t.ExecuteTemplate(w, "index", nil)
}


func create(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
  t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request){
    fmt.Printf("Подключено1")
    title := r.FormValue("title")
    anons := r.FormValue("anons")
    full_text := r.FormValue("full_text")
    fmt.Printf("Подключено2")
    fmt.Println(title, anons, full_text, "123")

    db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
    if err != nil{
      panic(err)
    }

    defer db.Close()
    fmt.Printf("Подключено")
    //Установка данных
   //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
   result, err := db.Exec("insert into test.articles (title, anons, full_text) values (?, ?, ?)", title, anons, full_text)
   if err != nil {
      panic(err)
    }
    fmt.Println(result.LastInsertId())  // id добавленного объекта
    fmt.Println(result.RowsAffected())  // количество затронутых строк

   http.Redirect(w, r, "/", http.StatusSeeOther)






}


func handleR(){

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  http.HandleFunc("/create", create)
  http.HandleFunc("/save_article", save_article)
  http.HandleFunc("/", index)

  http.ListenAndServe(":8080", nil)
}

func main() {
  handleR()
}
