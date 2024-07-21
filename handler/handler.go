package handler


func home_page(w http.ResponseWriter, r *http.Request){
      t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

      t.ExecuteTemplate(w, "index", nil)
}
