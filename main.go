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
 "time"
"log"
"gonum.org/v1/plot"
"gonum.org/v1/plot/plotter"
"gonum.org/v1/plot/plotutil"
"gonum.org/v1/plot/vg"
"gonum.org/v1/plot/vg/draw"
//"log"


)


type User struct{
   Id, Name, Surname, Email, Password, Nickname string
   Is_that_authorized_user bool
}

type Data_for_personal_page struct{
  Persona User
  Icon1 string
  Background_video string
  Sport_achive map[string]string
  Unique_types_of_exercises []string
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
      authorized_user.Is_that_authorized_user = true
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

func get_user_data_by_id(id_user_cur string) User{
  db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
  var zapros = fmt.Sprintf("SELECT id, name, surname, email, password, nickname FROM `persons` WHERE id = '%s'", id_user_cur)
  res,err := db.Query(zapros)
  fmt.Println(zapros)

  var cur_user User
  for res.Next(){
    err = res.Scan(&cur_user.Id, &cur_user.Name, &cur_user.Surname, &cur_user.Email, &cur_user.Password, &cur_user.Nickname)
  }
  return cur_user

}


func authorization(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/input_personal_account.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
  if(authorized_user.Is_that_authorized_user){
    http.Redirect(w, r, "/user_page/" + authorized_user.Id, http.StatusSeeOther)
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
       authorized_user.Is_that_authorized_user = true
         http.Redirect(w, r, "/user_page/" + cur_user.Id, http.StatusSeeOther)
     }

   }
   t.ExecuteTemplate(w, "input_personal_account", nil)
}





func get_Unique_types_of_exercises_m() []string{
  db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  var Unique_types_of_exercises []string
  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 var zapros = fmt.Sprintf("SELECT DISTINCT name_of_train FROM `trainings`")
 res,err := db.Query(zapros)
 fmt.Println(zapros)
 for res.Next(){
   var temp string
   _ = res.Scan(&temp)
   Unique_types_of_exercises = append(Unique_types_of_exercises, temp)

 }
 return Unique_types_of_exercises
}







func personal_account(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/personal_page.html", "templates/header.html", "templates/footer.html")

  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
//  fmt.Fprintf(w, "Category: %v\n", vars["id_user"])
  var current_user_id = vars["id_user"]

  var data Data_for_personal_page

  data.Persona = get_user_data_by_id(current_user_id)
  if(data.Persona.Id == authorized_user.Id){
    data.Persona.Is_that_authorized_user = true
  }






  fmt.Println(current_user_id)


  data.Icon1 = "./static/personal_static/" + data.Persona.Name+"_"+data.Persona.Surname+"_"+ data.Persona.Id + "/icon_1.jpg"
  data.Background_video = "./static/personal_static/" + data.Persona.Name+"_"+data.Persona.Surname+"_"+ data.Persona.Id + "/background_video.mp4"

  db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 var zapros = fmt.Sprintf("SELECT type_training, type_of_sports_load FROM `view_personal_achievements` WHERE id_user = '%s'", current_user_id)
 res,err := db.Query(zapros)
 fmt.Println(zapros)
data.Sport_achive = make(map[string]string)
 for res.Next(){
   var type_training, type_of_sports_load string
   err = res.Scan(&type_training, &type_of_sports_load)
   fmt.Println("9999---------------------9999 ",type_training, type_of_sports_load)
   if(type_of_sports_load == "absolute_power"){
     var zapros2 = fmt.Sprintf("SELECT weight, count FROM `trainings` WHERE name_of_train = '%s' AND id_person = '%s'", type_training, current_user_id)
     res2,err := db.Query(zapros2)
     if err != nil{
       panic(err)
     }
     fmt.Println(zapros2)
     fmt.Println("-----------------------")
     fmt.Println(res2)
     fmt.Println("-----------------------")
     var maximim_weight, count_for_maximim_weight int
     maximim_weight = -1
     count_for_maximim_weight = -1
     for res2.Next(){
       var weight, count int
       err = res2.Scan(&weight, &count)
       fmt.Println(weight, count, "        -----------------------------333333")
       if(maximim_weight < weight){
         maximim_weight = weight
         count_for_maximim_weight = count
       }
       if(maximim_weight == weight){
         count_for_maximim_weight = max(count, count_for_maximim_weight)
       }
     }
     data.Sport_achive[type_training] = type_training + " " + strconv.Itoa(maximim_weight) + " for " + strconv.Itoa(count_for_maximim_weight) + " times"
   }
 }
 data.Unique_types_of_exercises = get_Unique_types_of_exercises_m()


  t.ExecuteTemplate(w, "personal_page", data)
}


func submit_achive(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodPost {

    type_training := r.FormValue("options")
    type_of_sports_load := r.FormValue("options2")


    fmt.Println(type_training, type_of_sports_load)
    db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
    if err != nil{
      panic(err)
    }

    defer db.Close()
    fmt.Printf("Подключено")
    //Установка данных
   //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
   result, err := db.Exec("insert into test.view_personal_achievements (id_user, type_training, type_of_sports_load) values (?, ?, ?)", authorized_user.Id, type_training, type_of_sports_load)
   fmt.Println(result)
   http.Redirect(w, r, "/", http.StatusSeeOther)
 }
}

func getTimeInfo() [4]int {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	dayOfWeek := int(now.Weekday())
	secondsOfDay := now.Hour()*3600 + now.Minute()*60 + now.Second()

	return [4]int{year, month, dayOfWeek, secondsOfDay}
}

func create_train(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/create_train.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
  var data Data_for_personal_page
  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m()
  if r.Method == http.MethodPost {

    option_of_train := r.FormValue("option_of_train")
    weight := r.FormValue("weight")
    count := r.FormValue("count")

    fmt.Println(option_of_train, weight, count)
    db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
    if err != nil{
      panic(err)
    }

    defer db.Close()
    fmt.Printf("Подключено")
    //Установка данных

    var time_now  = getTimeInfo()
   //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
   result, err := db.Exec("insert into test.trainings (id_person, name_of_train, weight, count, year, month, day, seconds) values (?, ?, ?, ?, ? ,? ,? , ?)", authorized_user.Id, option_of_train, weight, count, time_now[0], time_now[1], time_now[2], time_now[3])
   fmt.Println(result, "   RES345")
   http.Redirect(w, r, "/", http.StatusSeeOther)
 }


  t.ExecuteTemplate(w, "create_train", data)
}


func graphic(times []time.Time, values []float64) {
    // Создание точек для графика
    pts := make(plotter.XYs, len(times))
    for i, t := range times {
        pts[i].X = float64(t.Unix())
        pts[i].Y = values[i]
    }

    // Создание нового графика
    p := plot.New()


    // Настройка оси X для отображения дат
    p.X.Label.Text = "Дата"
    p.X.Tick.Marker = plot.TimeTicks{Format: "02.01.2006"}

    // Добавление точек на график
    line, points, err := plotter.NewLinePoints(pts)
    if err != nil {
        log.Fatal(err)
    }
    line.Color = plotutil.Color(0)
    points.Shape = draw.CircleGlyph{}
    points.Color = plotutil.Color(1)
    p.Add(line, points)

    // Сохранение графика в файл
    if err := p.Save(10*vg.Inch, 6*vg.Inch, "static/plot.jpg"); err != nil {
        log.Fatal(err)
    }
    fmt.Println(times, values)
    println("График успешно сохранен в plot.png")
}


func personal_statistic(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/personal_statistic.html", "templates/header.html", "templates/footer.html")
  var data Data_for_personal_page
  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
//  fmt.Fprintf(w, "Category: %v\n", vars["id_user"])
  var current_user_id = vars["id_user"]
  data.Persona.Id = current_user_id
  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m()


  db, err := sql.Open("mysql", "root:@tcp(127.127.126.50:3306)/test")
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  var zapros string
  if r.Method == http.MethodPost {

    option_of_train := r.FormValue("option_of_train")
    zapros = fmt.Sprintf("WITH ranked_weights AS (SELECT weight, year, month, day, ROW_NUMBER() OVER (PARTITION BY year, month, day ORDER BY weight DESC) AS row_num FROM trainings WHERE name_of_train = '%s') SELECT weight, year, month, day FROM ranked_weights WHERE row_num = 1;", option_of_train)

  }else{
    zapros = fmt.Sprintf("WITH ranked_weights AS (SELECT weight, year, month, day, ROW_NUMBER() OVER (PARTITION BY year, month, day ORDER BY weight DESC) AS row_num FROM trainings) SELECT weight, year, month, day FROM ranked_weights WHERE row_num = 1")

  }
  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
  res,err := db.Query(zapros)
 fmt.Println(zapros)

 var times []time.Time
 var weights []float64


 for res.Next(){
   var year, month, day, weight int
   err = res.Scan(&weight, &year, &month, &day)
   times = append(times, time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
   weights = append(weights, float64(weight))
 }
 graphic(times, weights)

  t.ExecuteTemplate(w, "personal_statistic", data)
}





 func main() {
 http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
 r := mux.NewRouter()

 r.HandleFunc("/", home_page)
r.HandleFunc("/create_personal_page", create_personal_account)
r.HandleFunc("/input_personal_account", authorization)
r.HandleFunc("/submit_achive", submit_achive)
r.HandleFunc("/user_page/{id_user}", personal_account)
r.HandleFunc("/create_train", create_train)
r.HandleFunc("/personal_statistic/{id_user}", personal_statistic)

 fmt.Println()
 http.Handle ("/", r)
 http.ListenAndServe(":8080", nil)
}
