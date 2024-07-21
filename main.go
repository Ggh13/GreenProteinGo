package main

import ("fmt";"net/http";"html/template")
import( "github.com/gorilla/mux"
    "database/sql"
_ "github.com/go-sql-driver/mysql"
"os"
"strconv"
"path/filepath"

 "mime/multipart"
 "io"
//"os"
//"log"


)


type User struct{
   Id, Name, Surname, Email, Password, Nickname string
   Authorized bool
}

type Data_for_personal_page struct{
  Persona User
  Icon1 string
  Background_video string
}
var authorized_user User
func home_page(w http.ResponseWriter, r *http.Request){
      t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
      if err != nil{
        fmt.Fprintf(w, err.Error())

      }
      t.ExecuteTemplate(w, "index", nil)
}

func create_personal_account(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/create_personal_page.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
    if r.Method == http.MethodPost {

      name := r.FormValue("name")
      surname := r.FormValue("surname")
      email := r.FormValue("email")
      password := r.FormValue("password")
      nickname := r.FormValue("nickname")

      fmt.Println(name, surname)
      db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
      if err != nil{
        panic(err)
      }

      defer db.Close()
      fmt.Printf("Подключено")
      //Установка данных
     //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
     result, err := db.Exec("insert into test.persons (name, surname, email, password, nickname) values (?, ?, ?, ? ,?)", name, surname, email, password, nickname)
     fmt.Println(result)
     var zapros = fmt.Sprintf("SELECT id FROM `persons` WHERE email = '%s'", email)
     res,err := db.Query(zapros)
     fmt.Println(zapros)
     var idC int
     for res.Next(){
       err = res.Scan(&idC)
     }
  //   var lastInsertID int64
  //  if err := result.Scan(&lastInsertID); err != nil {
  //      log.Fatal(err)
  //  }
    fmt.Printf("Последний вставленный ID: %d\n", idC)


      // file, fileHeader, err := r.FormFile("image")
    //  if err != nil {
    //      http.Error(w, "Не удалось получить файл из формы", http.StatusInternalServerError)
    //      return
    //  }
    //  defer file.Close()

      // Путь для сохранения изображения (папка "uploads" в текущей директории)
      uploadsDir := "./static/personal_static/" + name+"_"+surname+"_"+ strconv.Itoa(idC)
      fmt.Println("WAY TO :  ",uploadsDir)
      os.Mkdir(uploadsDir, os.FileMode(0522))
      // Сохраняем изображение







      err = r.ParseMultipartForm(10 << 20) // Размер максимального загружаемого файла 10MB
    	if err != nil {
    		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
    		return
    	}

    	// Получаем изображения из формы по разным именам
    	image1, _, err := r.FormFile("icon1")
    	if err != nil {
    		http.Error(w, "Failed to get image1", http.StatusBadRequest)
    		return
    	}
    	defer image1.Close()

    	image2, _, err := r.FormFile("background_video")
    	if err != nil {
    		http.Error(w, "Failed to get image2", http.StatusBadRequest)
    		return
    	}
    	defer image2.Close()

    	// Создаем папку для сохранения файлов, если её нет
    	dir := uploadsDir


    	// Сохраняем первое изображение
    	err = saveFile("icon_1.jpg", image1, dir)
    	if err != nil {
    		http.Error(w, "Failed to save image1", http.StatusInternalServerError)
    		return
    	}

    	// Сохраняем второе изображение
    	err = saveFile("background_video.mp4", image2, dir)
    	if err != nil {
    		http.Error(w, "Failed to save image2", http.StatusInternalServerError)
    		return
    	}

    //	fmt.Fprintf(w, "Images uploaded successfully")
      authorized_user.Name = name
      authorized_user.Surname = surname
      authorized_user.Email = email
      authorized_user.Password = password
      authorized_user.Id = strconv.Itoa(idC)
      authorized_user.Authorized = true
      http.Redirect(w, r, "/", http.StatusSeeOther)



    }

  t.ExecuteTemplate(w, "create_personal_page", nil)

}



func saveFile(fileName string, file multipart.File, dir string) error {
	// Создаем полный путь для сохранения файла
	filePath := filepath.Join(dir, fileName)

	// Открываем файл для записи
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Копируем содержимое файла в новый файл
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}




func authorization(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/input_personal_account.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
    if r.Method == http.MethodPost {
      email := r.FormValue("email")
      password := r.FormValue("password")

      db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
      if err != nil{
        panic(err)
      }

      defer db.Close()
      fmt.Printf("Подключено")
      //Установка данных
     //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
     var zapros = fmt.Sprintf("SELECT id, name, surname, email, password, nickname FROM `persons` WHERE email = '%s' AND password='%s'", email, password)
     res,err := db.Query(zapros)
     fmt.Println(zapros)
     var cur_user User
     for res.Next(){

       err = res.Scan(&cur_user.Id, &cur_user.Name, &cur_user.Surname, &cur_user.Email, &cur_user.Password, &cur_user.Nickname)
     }
     if (cur_user.Id != ""){
       authorized_user = cur_user
         http.Redirect(w, r, "/personal_account", http.StatusSeeOther)
     }

   }
   t.ExecuteTemplate(w, "input_personal_account", nil)
}

func personal_account(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/personal_page.html", "templates/header.html", "templates/footer.html")
  var data Data_for_personal_page
  data.Persona = authorized_user
  data.Icon1 = "./static/personal_static/" + authorized_user.Name+"_"+authorized_user.Surname+"_"+ authorized_user.Id + "/icon_1.jpg"
  data.Background_video = "./static/personal_static/" + authorized_user.Name+"_"+authorized_user.Surname+"_"+ authorized_user.Id + "/background_video.mp4"
  t.ExecuteTemplate(w, "personal_page", data)
}

 func main() {
 http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
 r := mux.NewRouter()

 r.HandleFunc("/", home_page)
r.HandleFunc("/create_personal_page", create_personal_account)
r.HandleFunc("/input_personal_account", authorization)
r.HandleFunc("/personal_account", personal_account)

 fmt.Println()
 http.Handle ("/", r)
 http.ListenAndServe(":8080", nil)
}
