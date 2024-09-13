package main // v1 work

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
"strings"
"github.com/gorilla/sessions"
)


type User struct{
   Id, Name, Surname, Email, Password, Nickname string
   Is_that_authorized_user bool
}

type anthropometry_data_for_person struct{
  Height, Weight, Neck_size, Shoulder_size, Chest_size, Waist_size, Bicep_size, Forearm_size, Thigh_size, Quadriceps_size, Calf_size, Wrist_size, Ankle_size string
}
type Data_for_record_anthropometry struct{
  Authorized_user_data User
  Anthropometry_data anthropometry_data_for_person
}
type Data_for_personal_page struct{
  Authorized_user_data User
  Persona User
  Icon1 string
  Background_video string
  Sport_achive map[string]string
  Unique_types_of_exercises map [string] string

  Anthropometry anthropometry_data_for_person

  Error string
}
type Data_for_searching_personal_page struct{
  Authorized_user_data User
  Personas []Data_for_personal_page
  Name_of_searching, Surname_of_searching, Nickname_of_searching string
  Error string
}

type Data_for_create_training_programm struct{
  Authorized_user_data User
  Id_train_programm string
  Unique_types_of_exercises map [string] string
  Error string
}


type part_of_training_programm struct{

  Id_train_programm string
  Date string
  Id_train string

  Type_of_train string
  Type_of_train_translated string
  Weight string
  Count string
  Queue string
  Time_to_chill string

  Have_this_id_in_bd bool
  WeightDone int
  CountDone int


}
type Data_for_watch_training_programm struct{
  Authorized_user_data User

  Id_train_programm string
  Id_author_train_programm string
  Style_sport string
  Name_training_program string
  Current_data string

  Parts_training_programm map[string][]part_of_training_programm

  Unique_types_of_exercises map [string] string
  Error string
}

type Training_programma struct{
  Name_of_programma,Style_of_trainings, Description, Id  string

}
type Data_for_send_to_page_View_created_training_programms struct{
  Authorized_user_data User

  Training_programms []Training_programma
  Author User
  Is_this_Author bool
  Authors_Icon1  string
}



type Data_articles_page struct{
  Authorized_user_data User
  Themes []string
  Themes_and_name_articles  map [string] []string
}


type Data_referals struct{
  Authorized_user_data User
  Points int

}

//"root:@tcp(127.127.126.50)/test"
//"user:password@tcp(147.45.163.58:3306)/test"
//var authorized_user User
var sessionName = "name_session"
var adress_data_base = "root:@tcp(127.127.126.50)/test"

var store = sessions.NewCookieStore([]byte("super-secret-key"))




func populateUserFromSession(r *http.Request, sessionName string) (User, error) {
    session, err := store.Get(r, sessionName)
    if err != nil{
      panic(err)
    }

    user := User{
        Id:                    getSessionValueAsString(session, "curret_user_id"),
        Name:                  getSessionValueAsString(session, "name"),
        Surname:               getSessionValueAsString(session, "surname"),
        Email:                 getSessionValueAsString(session, "email"),
        Password:              getSessionValueAsString(session, "password"),
        Nickname:              getSessionValueAsString(session, "nickname"),
        Is_that_authorized_user: getSessionValueAsBool(session, "is_authorized"),
    }

    return user, nil
}



func saveUserToSession(r *http.Request, w http.ResponseWriter, sessionName string, user User) error{
    session, err := store.Get(r, sessionName)
    if err != nil {
        return err
    }

    session.Values["curret_user_id"] = user.Id
    session.Values["name"] = user.Name
    session.Values["surname"] = user.Surname
    session.Values["email"] = user.Email
    session.Values["password"] = user.Password
    session.Values["nickname"] = user.Nickname
    session.Values["is_authorized"] = user.Is_that_authorized_user



    fmt.Println("Данные о сохранении:  : : :")
    fmt.Println("------------")
    return session.Save(r, w)
}

func save_session_lang_user(r *http.Request, w http.ResponseWriter, sessionName string, lang string) error{
  session, err := store.Get(r, sessionName)
  if err != nil {
      return err
  }
  session.Values["lang"] = lang
  return session.Save(r, w)
}

func get_session_lang_user(r *http.Request, sessionName string) string{
  session, err := store.Get(r, sessionName)
  if err != nil{
    panic(err)
  }
  return  getSessionValueAsString(session, "lang")
}

func PrintSessionData(w http.ResponseWriter, r *http.Request) {
    // Получаем сессию
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Проходим по всем данным в сессии и выводим их
    fmt.Println(w, "Session Data:")
    fmt.Println(session.Values["curret_user_id"])
    fmt.Println(session.Values["name"] )
    fmt.Println(session.Values["surname"] )
    fmt.Println(session.Values["email"])
    fmt.Println(session.Values["password"])
    fmt.Println(session.Values["nickname"] )
    fmt.Println(session.Values["is_authorized"] )
}

func getSessionValueAsString(session *sessions.Session, key string) string {
    if val, ok := session.Values[key].(string); ok {
        return val
    }
    return ""
}

func getSessionValueAsBool(session *sessions.Session, key string) bool {
    if val, ok := session.Values[key].(bool); ok {
        return val
    }
    return false
}





func home_page(w http.ResponseWriter, r *http.Request){
      t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
      if err != nil{
        fmt.Fprintf(w, err.Error())

      }
      var data Data_for_personal_page

      authorized_user, err  := populateUserFromSession(r, sessionName)
      data.Authorized_user_data  = authorized_user
      t.ExecuteTemplate(w, "index", data)
}

func create_personal_account(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/create_personal_page.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
  authorized_user, err  := populateUserFromSession(r, sessionName)
    if r.Method == http.MethodPost {

      name := r.FormValue("name")
      surname := r.FormValue("surname")
      email := r.FormValue("email")
      password := r.FormValue("password")
      nickname := r.FormValue("nickname")

      fmt.Println(name, surname)
      db, err := sql.Open("mysql", adress_data_base)
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
      uploadsDir := "../personal_static/" + name+"_"+surname+"_"+ strconv.Itoa(idC)
      fmt.Println("WAY TO :  ",uploadsDir)
      err = os.Mkdir(uploadsDir, os.FileMode(0755))
      if err != nil{
        panic(err)
      }
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
      saveUserToSession(r, w, sessionName, authorized_user)
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
  db, err := sql.Open("mysql", adress_data_base)
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

func get_anthropometry_data_for_person_by_id(person_id string) (anthropometry_data_for_person) {
  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  var anthropometry anthropometry_data_for_person
  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
  var zapros = fmt.Sprintf("SELECT height, weight, neck_size, shoulder_size, chest_size, waist_size, bicep_size, forearm_size, thigh_size, quadriceps_size, calf_size, wrist_size, ankle_size FROM `anthropometria_sportsmen` WHERE person_id='%s'", person_id)
  res,err := db.Query(zapros)
  fmt.Println(zapros)
  fmt.Println(res)
  for res.Next(){
    err = res.Scan(&anthropometry.Height, &anthropometry.Weight, &anthropometry.Neck_size, &anthropometry.Shoulder_size, &anthropometry.Chest_size, &anthropometry.Waist_size, &anthropometry.Bicep_size, &anthropometry.Forearm_size, &anthropometry.Thigh_size, &anthropometry.Quadriceps_size, &anthropometry.Calf_size, &anthropometry.Wrist_size, &anthropometry.Ankle_size)
    if err != nil{
      panic(err)
    }
  }

  fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
  fmt.Println(anthropometry)
  return anthropometry
}


func authorization(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/input_personal_account.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    fmt.Fprintf(w, err.Error())

  }
  fmt.Println("Start authorized")
  PrintSessionData(w,r)
  authorized_user, err  := populateUserFromSession(r, sessionName)
  fmt.Println(authorized_user)
  if(authorized_user.Is_that_authorized_user){
    http.Redirect(w, r, "/user_page/" + authorized_user.Id, http.StatusSeeOther)
  }
    if r.Method == http.MethodPost {
      email := r.FormValue("email")
      password := r.FormValue("password")

      db, err := sql.Open("mysql", adress_data_base)
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
       fmt.Print("after input: ")
       fmt.Println(authorized_user)
       fmt.Println()
       err = saveUserToSession(r, w, sessionName, authorized_user)
       fmt.Println(err)
       PrintSessionData(w,r)

         http.Redirect(w, r, "/user_page/" + cur_user.Id, http.StatusSeeOther)
     }

   }
  t.ExecuteTemplate(w, "input_personal_account", nil)
}





func get_Unique_types_of_exercises_m(lang string ) (map [string] string){
  db, err := sql.Open("mysql", adress_data_base)
  if lang == ""{
    lang = "rus"
  }
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")

  //Установка данных
  lang_name := lang + "_name"
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 var zapros = fmt.Sprintf("SELECT DISTINCT eng_name, %s FROM `names_of_exercises`", lang_name)
 res,err := db.Query(zapros)
 fmt.Println(zapros)
 Unique_types_of_exercises := make(map [string] string)
 for res.Next(){
   var temp1, temp2 string
   _ = res.Scan(&temp1, &temp2)
   Unique_types_of_exercises[temp1] = temp2

 }
 return Unique_types_of_exercises
}







func personal_account(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/personal_page.html", "templates/header.html", "templates/footer.html")

  authorized_user, err  := populateUserFromSession(r, sessionName)

  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
//  fmt.Fprintf(w, "Category: %v\n", vars["id_user"])
  var current_user_id = vars["id_user"]

  var data Data_for_personal_page

  data.Persona = get_user_data_by_id(current_user_id)
  if(data.Persona.Id == authorized_user.Id){
    data.Persona.Is_that_authorized_user = true
  }


  data.Anthropometry = get_anthropometry_data_for_person_by_id(current_user_id)
  data.Unique_types_of_exercises = make(map [string] string)
  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))


  fmt.Println(current_user_id)


  data.Icon1 = "/personal_static/" + data.Persona.Name+"_"+data.Persona.Surname+"_"+ data.Persona.Id + "/icon_1.jpg"
  data.Background_video = "/personal_static/" + data.Persona.Name+"_"+data.Persona.Surname+"_"+ data.Persona.Id + "/background_video.mp4"

  db, err := sql.Open("mysql", adress_data_base)
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

     trans_type_training := data.Unique_types_of_exercises[type_training]
     data.Sport_achive[trans_type_training] =  strconv.Itoa(maximim_weight) + " на " + strconv.Itoa(count_for_maximim_weight) + " повторений"
   }
 }


 data.Authorized_user_data = authorized_user
  t.ExecuteTemplate(w, "personal_page", data)
}


func submit_achive(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodPost {
    authorized_user, err  := populateUserFromSession(r, sessionName)
    type_training := r.FormValue("options")
    type_of_sports_load := r.FormValue("options2")


    fmt.Println(type_training, type_of_sports_load)
    db, err := sql.Open("mysql", adress_data_base)
    if err != nil{
      panic(err)
    }

    defer db.Close()
    fmt.Printf("Подключено")
    //Установка данных
   //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
   result, err := db.Exec("insert into test.view_personal_achievements (id_user, type_training, type_of_sports_load) values (?, ?, ?)", authorized_user.Id, type_training, type_of_sports_load)
   fmt.Println(result)
   http.Redirect(w, r, "/user_page/" + authorized_user.Id, http.StatusSeeOther)
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
  data.Unique_types_of_exercises = make(map [string] string)

  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 var zapros = fmt.Sprintf("SELECT name_of_train FROM `trainings` WHERE id = (SELECT MAX(id) FROM trains WHERE id_person = %s) AND id_person = %s;", current_user_id)
 res,err := db.Query(zapros)
 var temp string
 for res.Next(){
   _ = res.Scan(&temp)
 }

  data.Unique_types_of_exercises = append(data.Unique_types_of_exercises, temp)
  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))
  if r.Method == http.MethodPost {
    authorized_user, err  := populateUserFromSession(r, sessionName)
    option_of_train := r.FormValue("option_of_train")
    weight := r.FormValue("weight")
    count := r.FormValue("count")

    fmt.Println(option_of_train, weight, count)
    db, err := sql.Open("mysql", adress_data_base)
    if err != nil{
      panic(err)
    }

    defer db.Close()
    fmt.Printf("Подключено")
    //Установка данных

    var time_now  = getCurrentDate()
    var seconds = SecondsSinceStartOfDay()
   //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
   result, err := db.Exec("insert into test.trainings (id_person, name_of_train, weight, count, date, second) values (?, ?, ?, ?, ?,? )", authorized_user.Id, option_of_train, weight, count, time_now, seconds)
   fmt.Println(result, "   RES345")
   http.Redirect(w, r, "/create_train", http.StatusSeeOther)
 }

 data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  t.ExecuteTemplate(w, "create_train", data)
}


func graphic(times []time.Time, values []float64) {
    // Определяем начальную и конечную даты
    start := times[0]
    end := times[len(times)-1]

    // Создаем мапу для хранения данных с датами
    data := make(map[time.Time]float64)
    for i, t := range times {
        data[t] = values[i]
    }

    // Создаем точки для графика
    var pts plotter.XYs
    var lastValue float64

    for t := start; !t.After(end); t = t.AddDate(0, 0, 1) { // Двигаемся по каждому дню
        if value, exists := data[t]; exists {
            lastValue = value
        }
        pts = append(pts, plotter.XY{
            X: float64(t.Unix()),
            Y: lastValue,
        })
    }

    // Обязательно включаем последний день
    pts = append(pts, plotter.XY{
        X: float64(end.Unix()),
        Y: values[len(values)-1],
    })

    // Создание нового графика
    p := plot.New()

    // Настройка оси X для отображения дат
    p.X.Label.Text = "Дата"
    p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}

    // Настройка пределов оси X
    p.X.Min = float64(start.Unix())
    p.X.Max = float64(end.Unix())

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
  data.Unique_types_of_exercises = make(map [string] string)
  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))


  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  var zapros string
  if r.Method == http.MethodPost {

    option_of_train := r.FormValue("option_of_train")
    zapros = fmt.Sprintf("SELECT date, MAX(weight) AS max_weight FROM trainings WHERE name_of_train = '%s' AND id_person = '%s' GROUP BY date ORDER BY date ASC", option_of_train, current_user_id)

  }else{
    zapros = fmt.Sprintf("SELECT date, MAX(weight) AS max_weight FROM trainings WHERE name_of_train = '%s' AND id_person = '%s' GROUP BY date ORDER BY date ASC",  "bench press", current_user_id)

  }

  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 res,err := db.Query(zapros)
 fmt.Println(zapros)

 var times []time.Time
 var weights []float64

 fmt.Println(res)

 for res.Next(){
   var dateS string
   var weight int
   err = res.Scan(&dateS, &weight)
   fmt.Println(dateS, weight)
   date, err := time.Parse("2006-01-02", dateS)
    if err != nil {
        fmt.Println("Ошибка при парсинге даты:", err)
        return
    }


   times = append(times, date)
   weights = append(weights, float64(weight))
 }
 fmt.Println(times, weights)
 graphic(times, weights)
 data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  t.ExecuteTemplate(w, "personal_statistic", data)
}


func verification_of_authorization(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/need_authorization.html", "templates/header.html", "templates/footer.html")
  path := r.URL.Path
  fmt.Println(path)
  authorized_user, _  := populateUserFromSession(r, sessionName)
  if(authorized_user.Is_that_authorized_user){
    if(path == "/create_train_programm_step_1"){
      fmt.Println("All good")
      create_train_programm_step_1(w,r)

    }

     if strings.Contains(path, "/watch_training_programm") {

       watch_training_programm(w,r)
     }
     if strings.Contains(path, "/view_created_training_programms") {
       view_created_training_programms(w,r)
     }
     if strings.Contains(path, "/view_current_training_programm") {
       view_current_training_programm(w,r)
     }



  }else{
    t.ExecuteTemplate(w, "need_authorization", nil)
  }
}
func search_people(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/search_people_page.html", "templates/header.html", "templates/footer.html")
  fmt.Println("---09090------")
  if err != nil{
    panic(err)
  }
  var data Data_for_searching_personal_page


  if r.Method == http.MethodPost {
    name := r.FormValue("firstName")
    surname := r.FormValue("lastName")
    nickname := r.FormValue("nickname")
    if(name == "" && surname == "" && nickname == "" ){
      data.Error = "Введите имя, фамилию или никнейм"
    }else{
      db, err := sql.Open("mysql", adress_data_base)
      if err != nil{
        panic(err)
      }
      var zapros string
      zapros = fmt.Sprintf("SELECT id, name, surname, nickname FROM `persons` WHERE ")
      defer db.Close()
      fmt.Printf("Подключено")
      if(name != ""){
        zapros = zapros + fmt.Sprintf("name='%s'", name)
      }
      if(surname != ""){
        zapros = zapros + fmt.Sprintf("surname='%s'", surname)
      }
      if(nickname != ""){
        zapros = zapros + fmt.Sprintf("nickname='%s'", nickname)
      }

      //Установка данных
     //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
     //var zapros = fmt.Sprintf("SELECT id, name, surname, nickname FROM `persons` WHERE id_user = '%s'", current_user_id)
     res,err := db.Query(zapros)
     fmt.Println(zapros)

     for res.Next(){
       var temp Data_for_personal_page
       err = res.Scan(&temp.Persona.Id, &temp.Persona.Name, &temp.Persona.Surname, &temp.Persona.Nickname)
       temp.Icon1 = "/personal_static/" + temp.Persona.Name+"_"+temp.Persona.Surname+"_"+ temp.Persona.Id + "/icon_1.jpg"
       data.Personas = append(data.Personas, temp)
     }
    }
  }else{
    db, err := sql.Open("mysql", adress_data_base)
    if err != nil{
      panic(err)
    }
    zapros := fmt.Sprintf("SELECT id, name, surname, nickname FROM `persons` ")
    res,err := db.Query(zapros)
    fmt.Println(zapros)

    for res.Next(){
      var temp Data_for_personal_page
      err = res.Scan(&temp.Persona.Id, &temp.Persona.Name, &temp.Persona.Surname, &temp.Persona.Nickname)
      temp.Icon1 = "/personal_static/" + temp.Persona.Name+"_"+temp.Persona.Surname+"_"+ temp.Persona.Id + "/icon_1.jpg"
      data.Personas = append(data.Personas, temp)
    }
  }

fmt.Println(data)
data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
t.ExecuteTemplate(w, "search_people_page", data)
}

func processing_create_train_programm_step_1(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodPost {
    fmt.Println("ПЕРВЫЙ !!! ____ ----____----____-____---____----____--__--_---__----_------_")
    name_of_train_programm := r.FormValue("name_of_train_programm")
    option_of_style_training := r.FormValue("option_of_style_training")
    Description := r.FormValue("Description")
    fmt.Println(name_of_train_programm, option_of_style_training)

    db, err := sql.Open("mysql", adress_data_base)
    if err != nil{
      panic(err)
    }
    authorized_user, err  := populateUserFromSession(r, sessionName)
    defer db.Close()
    fmt.Printf("Подключено")

    result, err := db.Exec("insert into test.Training_programms (id_author_of_training_program	, name_training_program, style_of_training, Description) values (?, ?, ?, ?)", authorized_user.Id, name_of_train_programm, option_of_style_training, Description)
    fmt.Println(result)

   var zapros = fmt.Sprintf("SELECT id FROM `Training_programms` WHERE id_author_of_training_program = '%s' AND name_training_program = '%s'", authorized_user.Id, name_of_train_programm)
   res,err := db.Query(zapros)
   fmt.Println(zapros)
   var id_programm string
   for res.Next(){
     err = res.Scan(&id_programm)
   }
   http.Redirect(w, r, "/view_created_training_programms/" + authorized_user.Id, http.StatusSeeOther)
  }
}
func create_train_programm_step_1(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/create_train_programm.html", "templates/header.html", "templates/footer.html")
  if r.Method == http.MethodPost {
        log.Println("Начало обработки POST-запроса")
        // Остальной код...
    } else {
        log.Println("Обработка GET-запроса")
    }
    t.ExecuteTemplate(w, "create_train_programm", nil)
}



func getCurrentDate() string {
    now := time.Now()
    year, month, day := now.Date()
    return fmt.Sprintf("%d-%02d-%02d", year, int(month), day)
}

func SecondsSinceStartOfDay() int {
    now := time.Now()                // Текущее время
    midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // Начало дня (00:00:00)
    seconds := int(now.Sub(midnight).Seconds()) // Количество секунд с начала дня
    return seconds
}

func watch_training_programm(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/watch_training_programm.html", "templates/header.html", "templates/footer.html")


  vars := mux.Vars(r)
  if err != nil{
    panic(err)
  }
  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }

  defer db.Close()







  var data Data_for_watch_training_programm
  data.Parts_training_programm = make(map[string][]part_of_training_programm)
  data.Id_train_programm = vars["id_training_programm"]


  var zapros = fmt.Sprintf("SELECT id_author_of_training_program FROM `Training_programms` WHERE id = %s", data.Id_train_programm)
  res,err := db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    err = res.Scan(&data.Id_author_train_programm)
  }


  data.Current_data = getCurrentDate()
  data.Unique_types_of_exercises = make(map [string] string)
  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))



  authorized_user, err  := populateUserFromSession(r, sessionName)

    if r.Method == http.MethodPost {
      weight := r.FormValue("weight")
      count := r.FormValue("count")
      option_of_train := r.FormValue("hiddenFieldOption")
      //
      id_of_train := r.FormValue("hiddenFieldID")
      date := getCurrentDate()
      fmt.Println(weight, count, option_of_train, id_of_train, date)

      result, _ := db.Exec("insert into test. trainings (`id_of_all_train`, `name_of_train`,	`weight`,	`count`,	`id_person`,	`date`) values (?, ?, ?, ?, ?, ?)",id_of_train ,option_of_train, weight, count, authorized_user.Id, date)
      fmt.Println(result)

    }

  zapros = fmt.Sprintf("SELECT e.id, n.rus_name AS name_of_train, n.eng_name AS name_of_train,  e.weight, e.count, number_o_execution_sequence, time_to_chill_before_next, e.date FROM Exercises_in_training_programs e JOIN names_of_exercises n ON e.name_of_train = n.eng_name WHERE e.id_of_programm = %s;", data.Id_train_programm)
  res,err = db.Query(zapros)
  fmt.Println(zapros)

  var part part_of_training_programm
  for res.Next(){
    err = res.Scan(&part.Id_train, &part.Type_of_train_translated, &part.Type_of_train, &part.Weight, &part.Count, &part.Queue, &part.Time_to_chill, &part.Date)
    fmt.Println(part)
    part.Id_train_programm = data.Id_train_programm

    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM trainings WHERE id_of_all_train = ? AND id_person = ?)"

    // Выполнение запроса
    err = db.QueryRow(query, part.Id_train, authorized_user.Id).Scan(&exists)
    if err != nil {
        panic(err)
    }

    // Вывод результата
    if exists {
        part.Have_this_id_in_bd = true
        var zapros2 = fmt.Sprintf("SELECT  weight, count FROM `trainings` WHERE id_of_all_train = '%s' AND id_person = '%s' ", part.Id_train,  authorized_user.Id)
        res2,_ := db.Query(zapros2)
        fmt.Println(zapros2)
        for res2.Next(){
            err = res2.Scan(&part.WeightDone, &part.CountDone)
            fmt.Println(part.WeightDone, part.CountDone)

        }
    } else {
        part.Have_this_id_in_bd = false
    }


    data.Parts_training_programm[part.Date] = append(data.Parts_training_programm[part.Date], part)
    fmt.Println(data.Parts_training_programm[part.Date])
  }
  fmt.Println(data)

//  fmt.Fprintf(w, "Category: %v\n", vars["id_user"])
  data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  t.ExecuteTemplate(w, "watch_training_programm", data)

}

func constacs(w http.ResponseWriter, r *http.Request){
      t, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")
      fmt.Println("!!!!!!!!")
      if err != nil {
          panic(err)
      }

      var data Data_for_personal_page

      authorized_user, err  := populateUserFromSession(r, sessionName)
      data.Authorized_user_data  = authorized_user

      t.ExecuteTemplate(w, "contacts", data)
}

func exit(w http.ResponseWriter, r *http.Request){
    authorized_user, _  := populateUserFromSession(r, sessionName)
    authorized_user.Id = "-1"
    authorized_user.Name = "-1"
    authorized_user.Surname = "-1"
    authorized_user.Is_that_authorized_user = false
    saveUserToSession(r, w, sessionName, authorized_user)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

 func view_created_training_programms(w http.ResponseWriter, r *http.Request){
   t, err := template.ParseFiles("templates/View_created_training_programms.html", "templates/header.html", "templates/footer.html")

   vars := mux.Vars(r)
   //w.WriteHeader(http.StatusOK)
  authorized_user, err  := populateUserFromSession(r, sessionName)
   author_id := vars["id_author"]
   var data Data_for_send_to_page_View_created_training_programms
   data.Author = get_user_data_by_id(author_id)
   if(authorized_user.Id == author_id){
     data.Is_this_Author = true
   }else{
     data.Is_this_Author = false
   }
   data.Authors_Icon1 = "/personal_static/" + data.Author.Name+"_"+data.Author.Surname+"_"+ data.Author.Id + "/icon_1.jpg"


   if err != nil {
       panic(err)
   }

   db, err := sql.Open("mysql", adress_data_base)
   if err != nil{
     panic(err)
   }

   defer db.Close()
   var zapros2 = fmt.Sprintf("SELECT  id, name_training_program, style_of_training, Description FROM `Training_programms` WHERE 	id_author_of_training_program = '%s' ", author_id)
   res2,_ := db.Query(zapros2)
   fmt.Println(zapros2)
   for res2.Next(){
      var programma Training_programma
      err = res2.Scan(&programma.Id, &programma.Name_of_programma, &programma.Style_of_trainings, &programma.Description)
      data.Training_programms = append(data.Training_programms, programma)


   }


   data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
   t.ExecuteTemplate(w, "view_created_training_programms", data)
 }



 func create_train_set_exercises(w http.ResponseWriter, r *http.Request){

   vars := mux.Vars(r)
   //w.WriteHeader(http.StatusOK)
   var data Data_for_create_training_programm
   data.Id_train_programm = vars["id_programm"]
   data.Unique_types_of_exercises = make(map [string] string)
   data.Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))

 //  fmt.Fprintf(w, "Category: %v\n", vars["id_user"])
   //var current_user_id =
   if r.Method == http.MethodPost {

     weight := r.FormValue("weight")
     count := r.FormValue("count")
     option_of_train := r.FormValue("option_of_train")
     date := r.FormValue("date")
     time_to_chill_before_next := r.FormValue("time_to_chill_before_next")
     queue := r.FormValue("queue")
     fmt.Println(date)

     db, err := sql.Open("mysql", adress_data_base)
     if err != nil{
       panic(err)
     }

     defer db.Close()


    result, err := db.Exec("insert into test.Exercises_in_training_programs (`id_of_programm`, `name_of_train`,	`weight`,	`count`,	`number_o_execution_sequence`,	`time_to_chill_before_next`, `date`) values (?, ?, ?, ?, ?, ?, ?)",data.Id_train_programm, option_of_train, weight, count, queue, time_to_chill_before_next, date)
    fmt.Println(result)
    http.Redirect(w, r, "/watch_training_programm/" + data.Id_train_programm, http.StatusSeeOther)

   }

 }



 func view_current_training_programm(w http.ResponseWriter, r *http.Request){ // удалить через неделю после 14.08.2024
       t, err := template.ParseFiles("templates/View_current_training_programm.html", "templates/header.html", "templates/footer.html")
       fmt.Println("!!!!!!!!")
       if err != nil {
           panic(err)
       }
       t.ExecuteTemplate(w, "View_current_training_programm", nil)
 }

func training_summary(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/training_summary.html", "templates/header.html", "templates/footer.html")
  vars := mux.Vars(r)
  //w.WriteHeader(http.StatusOK)
  Person := get_user_data_by_id(vars["id_person"])
  var data Data_for_watch_training_programm
  data.Parts_training_programm = make(map[string][]part_of_training_programm)
  Unique_types_of_exercises := make(map[string]string)
  Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))
  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }

  defer db.Close()

  var zapros = fmt.Sprintf("SELECT name_of_train, weight, count, date FROM `trainings` WHERE 	id_person = '%s' ", Person.Id)
  res,err := db.Query(zapros)
  fmt.Println(zapros)

  var part part_of_training_programm
  var date string
  for res.Next(){
    err = res.Scan(&part.Type_of_train, &part.Weight, &part.Count, &date)
    part.Type_of_train_translated = Unique_types_of_exercises[part.Type_of_train]
    data.Parts_training_programm[date] = append(data.Parts_training_programm[date], part)


  }

  fmt.Println("!!!!!!!!")
  if err != nil {
      panic(err)
  }
  data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  t.ExecuteTemplate(w, "training_summary", data)
}

 func handlerTTT(w http.ResponseWriter, r *http.Request) {
    // Извлечение заголовка Referer
    referer := r.Header.Get("Referer")

    // Вывод страницы-источника
    fmt.Fprintf(w, "Переадресация произошла с: %s", referer)
}



func record_anthropometric_changes(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/record_anthropometric_changes.html", "templates/header.html", "templates/footer.html")
  //vars := mux.Vars(r)
  //w.WriteHeader(http.StatusOK)
  //Person := get_user_data_by_id(vars["id_person"])
  var data Data_for_record_anthropometry

  data.Authorized_user_data, _ = populateUserFromSession(r, sessionName)

  if r.Method == http.MethodPost {

    data.Anthropometry_data.Height = r.FormValue("Height")
    data.Anthropometry_data.Weight = r.FormValue("Weight")
    data.Anthropometry_data.Neck_size = r.FormValue("Neck_size")
    data.Anthropometry_data.Shoulder_size = r.FormValue("Shoulder_size")
    data.Anthropometry_data.Chest_size = r.FormValue("Chest_size")
    data.Anthropometry_data.Waist_size = r.FormValue("Waist_size")
    data.Anthropometry_data.Bicep_size = r.FormValue("Bicep_size")
    data.Anthropometry_data.Forearm_size = r.FormValue("Forearm_size")
    data.Anthropometry_data.Thigh_size = r.FormValue("Thigh_size")
    data.Anthropometry_data.Quadriceps_size = r.FormValue("Quadriceps_size")
    data.Anthropometry_data.Calf_size = r.FormValue("Calf_size")
    data.Anthropometry_data.Wrist_size = r.FormValue("Wrist_size")
    data.Anthropometry_data.Ankle_size = r.FormValue("Ankle_size")


    fmt.Println(data)

    db, err := sql.Open("mysql", adress_data_base)
    if err != nil{
      panic(err)
    }

    defer db.Close()


   result, err := db.Exec("insert into test.anthropometria_sportsmen (	`person_id`, `height`, `weight`, `neck_size`, `shoulder_size`, `chest_size`, `waist_size`, `bicep_size`, `forearm_size`, `thigh_size`, `quadriceps_size`, `calf_size`, `wrist_size`, `ankle_size`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",data.Authorized_user_data.Id, data.Anthropometry_data.Height, data.Anthropometry_data.Weight, data.Anthropometry_data.Neck_size, data.Anthropometry_data.Shoulder_size, data.Anthropometry_data.Chest_size, data.Anthropometry_data.Waist_size, data.Anthropometry_data.Bicep_size, data.Anthropometry_data.Forearm_size, data.Anthropometry_data.Thigh_size, data.Anthropometry_data.Quadriceps_size, data.Anthropometry_data.Calf_size, data.Anthropometry_data.Wrist_size, data.Anthropometry_data.Ankle_size)
   fmt.Println(result)

   http.Redirect(w, r, "/watch_anthropometric/" + data.Authorized_user_data.Id, http.StatusSeeOther)
 }

  t.ExecuteTemplate(w, "record_anthropometric_changes", data)
}
func watch_anthropometric(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/watch_anthropometric.html", "templates/header.html", "templates/footer.html")
  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
  Person := get_user_data_by_id(vars["id_person"])
  var data Data_for_record_anthropometry
  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }
  data.Authorized_user_data, _ = populateUserFromSession(r, sessionName)
  defer db.Close()

  var zapros = fmt.Sprintf("SELECT height, weight, neck_size, shoulder_size, chest_size, waist_size, bicep_size, forearm_size, thigh_size, quadriceps_size, calf_size, wrist_size, ankle_size FROM `anthropometria_sportsmen` WHERE person_id = '%s' ", Person.Id)
  res,err := db.Query(zapros)
  fmt.Println(zapros)


  for res.Next(){
    _ = res.Scan(&data.Anthropometry_data.Height, &data.Anthropometry_data.Weight, &data.Anthropometry_data.Neck_size, &data.Anthropometry_data.Shoulder_size, &data.Anthropometry_data.Chest_size, &data.Anthropometry_data.Waist_size, &data.Anthropometry_data.Bicep_size, &data.Anthropometry_data.Forearm_size, &data.Anthropometry_data.Thigh_size, &data.Anthropometry_data.Quadriceps_size, &data.Anthropometry_data.Calf_size, &data.Anthropometry_data.Wrist_size, &data.Anthropometry_data.Ankle_size)
    fmt.Println(data)
  }
  t.ExecuteTemplate(w, "record_anthropometric_changes", data)

}

func create_name_of_theme_article(w http.ResponseWriter, r *http.Request){
  db, err := sql.Open("mysql", adress_data_base)

  vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)
  id_person := vars["id_person"]


  if err != nil{
    panic(err)
  }

  defer db.Close()

  if r.Method == http.MethodPost {

    name_new_article := r.FormValue("name_new_article")
    _, _ = db.Exec("insert into test.theme_article ( `id_author`	`name_of_theme_article`) values (?, ?)",id_person, name_new_article)
    http.Redirect(w, r, ("/articles_page/" + id_person), http.StatusSeeOther)
  }

}
func articles_page(w http.ResponseWriter, r *http.Request){
      t, err := template.ParseFiles("templates/articles_page.html", "templates/header.html", "templates/footer.html")
      if err != nil{
        fmt.Fprintf(w, err.Error())
      }
      var data Data_articles_page

      data.Themes_and_name_articles = make(map[string] []string)
      data.Themes = make([]string, 0)
      db, err := sql.Open("mysql", adress_data_base)
      if err != nil{
        panic(err)
      }
      data.Authorized_user_data, _ = populateUserFromSession(r, sessionName)

      defer db.Close()
      vars := mux.Vars(r)
      w.WriteHeader(http.StatusOK)
      id_person := vars["id_person"]
      fmt.Println(data)
      var zapros = fmt.Sprintf("SELECT name_of_theme_article FROM `theme_article` WHERE id_author = '%s' ", id_person)
      res,err := db.Query(zapros)
      fmt.Println(zapros)
      if err != nil{
        panic(err)
      }


      for res.Next(){
        var name_them string
        err = res.Scan(&name_them)
        if err != nil{
          panic(err)
        }
        data.Themes = append(data.Themes, name_them)

      }
      fmt.Println(data)
      t.ExecuteTemplate(w, "articles_page", data)
}

func referal_page(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/referal_page.html", "templates/header.html", "templates/footer.html")
  fmt.Println("!!!!!!!!")
  var Data_ref Data_referals
  if err != nil {
      panic(err)
  }
  authorized_user, err  := populateUserFromSession(r, sessionName)
  db, err := sql.Open("mysql", adress_data_base)
  if err != nil{
    panic(err)
  }

  defer db.Close()

  var zapros = fmt.Sprintf("SELECT points FROM `score_system` WHERE 	id_user = '%s' ", authorized_user.Id)
  res,err := db.Query(zapros)
  var points int
  for res.Next(){
    err = res.Scan(&points)
  }
  Data_ref.Points = points
  t.ExecuteTemplate(w, "referal_page", Data_ref)
}
 func main() {

   store.Options = &sessions.Options{
        Path:     "/",       // Путь для куков
        MaxAge:   86400 * 7, // Время жизни куков - 7 дней
        HttpOnly: true,      // Запрет доступа к кукам через JavaScript
        Secure:   false,     // Для HTTP должен быть false
    }


 http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

 fs := http.FileServer(http.Dir("../personal_static"))
 http.Handle("/personal_static/", http.StripPrefix("/personal_static/", fs))




 r := mux.NewRouter()







  r.HandleFunc("/", home_page)
  r.HandleFunc("/create_personal_page", create_personal_account)
  r.HandleFunc("/constacs", constacs)
  r.HandleFunc("/exit", exit)
  r.HandleFunc("/input_personal_account", authorization)
  r.HandleFunc("/submit_achive", submit_achive)
  r.HandleFunc("/user_page/{id_user}", personal_account)
  r.HandleFunc("/create_train", create_train)
  r.HandleFunc("/search_people", search_people)
  r.HandleFunc("/personal_statistic/{id_user}", personal_statistic)
  r.HandleFunc("/create_train_programm_step_1", verification_of_authorization)
  r.HandleFunc("/create_train_programm_step_2/{id_training_programm}", handlerTTT)

  r.HandleFunc("/watch_training_programm/{id_training_programm}", verification_of_authorization)
  r.HandleFunc("/watch_anthropometric/{id_person}", watch_anthropometric)

  r.HandleFunc("/view_created_training_programms/{id_author}", verification_of_authorization)
  r.HandleFunc("/training_summary/{id_person}", training_summary)
  r.HandleFunc("/record_anthropometric_changes/{id_person}", record_anthropometric_changes)
  //r.HandleFunc("/view_current_training_programm/{id_programm}", verification_of_authorization)
  r.HandleFunc("/articles_page/{id_person}", articles_page)
  r.HandleFunc("/create_name_of_theme_article/{id_person}", articles_page)



  r.HandleFunc("/create_train_set_exercises/{id_programm}", create_train_set_exercises)

  r.HandleFunc("/processing_create_train_programm_step_1", processing_create_train_programm_step_1)
  r.HandleFunc("/referal_page", referal_page)
 fmt.Println()
 http.Handle ("/", r)
 http.ListenAndServe(":8080", nil)
}
