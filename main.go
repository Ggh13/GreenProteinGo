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

"sort"
//"log"
"strings"
"github.com/gorilla/sessions"
)


type User struct{
   Id, Name, Surname, Email, Password, Nickname string
   Is_that_authorized_user bool
   Is_it_admin bool
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

  Technique_exercise_link string
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

  Id_day_int_train_programm string
  Id_author_train_programm string
  Style_sport string
  Name_training_program string
  Current_data string

  Parts_training_programm map[string][]part_of_training_programm

  Unique_types_of_exercises map [string] string

  Temp_mass []string
  Error string
}

type Training_programma struct{
  Author User
  Authors_Icon1  string
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
  Link_adress string
  Points int

}

type Training_day_discrpt struct{

  Id, Name,Number_of_training_day string
}

type Data_for_view_training_days_of_current_train_programm struct{
  Days []Training_day_discrpt
  Authorized_user_data User
  Id_programm string
  Is_this_author bool
}


type Product_data struct{
  Id, Name, Price, Brand, Weight string
  Date string
  Features_had []Features_for_products
  Quantity_in_store, Quantity_in_cart int
  Icon_product string
  In_cart bool

}
type Data_for_store_main_page struct{
  Products []Product_data
  Features []Features_for_products
  Authorized_user_data User
}


type Features_for_products struct{
  Id, Name, Unit_of_measurement, Value string
}

//"root:@tcp(127.127.126.50)/test"
//"user:password@tcp(147.45.163.58:3306)/test"
//"user:password@tcp(147.45.163.58:3306)/test"

//http://buzhor13.ru
//"http://147.45.163.58:8080"
//http://localhost:8080
var adress_web = "http://localhost:8080"
//var authorized_user User
var sessionName = "name_session"
var adress_sql = "root:@tcp(127.127.126.50)/"
var adress_data_base_test = adress_sql + "test"
var adress_data_base_store = adress_sql + "/store"



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
        Is_it_admin:  getSessionValueAsBool(session, "adminAuth"),
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
    session.Values["adminAuth"] = user.Is_it_admin


    fmt.Println("Данные о сохранении:  : : :")
    fmt.Println("------------")
    return session.Save(r, w)
}

func returnToLastPage(r *http.Request, w http.ResponseWriter){
  referer := r.Referer()

    // Если Referer пустой, используем URL по умолчанию
    if referer == "" {
        referer = "buzhor13.ru" // Замените на ваш URL по умолчанию
    }

    // Выполняем редирект
    http.Redirect(w, r, referer, http.StatusFound)
}
func saveUserWhoInvite(r *http.Request, w http.ResponseWriter, sessionName string, id_person_who_invite string) error{
    session, err := store.Get(r, sessionName)
    if err != nil {
        return err
    }

    session.Values["who_invite_id"] = id_person_who_invite




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

func get_who_inveted_id(r *http.Request, sessionName string) string{
  session, err := store.Get(r, sessionName)
  if err != nil{
    panic(err)
  }
  return  getSessionValueAsString(session, "who_invite_id")
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


      who_invite_you := get_who_inveted_id(r, sessionName)
      fmt.Println("I INVITE YOU! :" + who_invite_you)
      if(who_invite_you != ""){
            db, err := sql.Open("mysql", adress_data_base_test)
            if err != nil{
              panic(err)
            }

            defer db.Close()
            fmt.Printf("Подключено")
            query := `UPDATE score_system SET points = points + ? WHERE id_user = ?`


            _, err = db.Exec(query, 25, who_invite_you)
      }








      fmt.Println(name, surname)
      db, err := sql.Open("mysql", adress_data_base_test)
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

    result, err = db.Exec("insert into test.score_system (id_user	, points) values (?, ?)", idC, 100)
    fmt.Println(result)
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

    	//image2, _, err := r.FormFile("background_video")
    //	if err != nil {
    //		http.Error(w, "Failed to get image2", http.StatusBadRequest)
    //		return
    //	}
    //	defer image2.Close()

    	// Создаем папку для сохранения файлов, если её нет
    	dir := uploadsDir


    	// Сохраняем первое изображение
    	err = saveFile("icon_1.jpg", image1, dir)
    	if err != nil {
    		http.Error(w, "Failed to save image1", http.StatusInternalServerError)
    		return
    	}

    	// Сохраняем второе изображение
    //	err = saveFile("background_video.mp4", image2, dir)
    //	if err != nil {
    //		http.Error(w, "Failed to save image2", http.StatusInternalServerError)
    //		return
    //	}

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
  db, err := sql.Open("mysql", adress_data_base_test)
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
  db, err := sql.Open("mysql", adress_data_base_test)
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

      db, err := sql.Open("mysql", adress_data_base_test)
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
  db, err := sql.Open("mysql", adress_data_base_test)
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

  db, err := sql.Open("mysql", adress_data_base_test)
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
    db, err := sql.Open("mysql", adress_data_base_test)
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

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  //Установка данных
  data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  fmt.Println("temp")
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 var zapros = fmt.Sprintf("SELECT name_of_train FROM `trainings` WHERE id = (SELECT MAX(id) FROM `trainings` WHERE id_person = '%s') AND id_person = '%s';", data.Authorized_user_data.Id, data.Authorized_user_data.Id)

 res,err := db.Query(zapros)
 var temp string
 fmt.Println("temp2")
 fmt.Println(res)
 for res.Next(){
   _ = res.Scan(&temp)
 }
 fmt.Println(temp)
  fmt.Println("temp3")
  var zapros2 = fmt.Sprintf("SELECT rus_name  FROM `names_of_exercises` WHERE eng_name= '%s';", temp)
  res2,err := db.Query(zapros2)
  var temp2 string
  fmt.Println(zapros2)
  for res2.Next(){
    _ = res2.Scan(&temp2)
  }
  fmt.Println("temp3")

  data.Unique_types_of_exercises = get_Unique_types_of_exercises_m(get_session_lang_user(r, sessionName))
  data.Unique_types_of_exercises["Last_ex_768"] = temp
  fmt.Println(data.Unique_types_of_exercises)
  if r.Method == http.MethodPost {
    authorized_user, err  := populateUserFromSession(r, sessionName)
    option_of_train := r.FormValue("option_of_train")
    weight := r.FormValue("weight")
    count := r.FormValue("count")

    fmt.Println(option_of_train, weight, count)
    db, err := sql.Open("mysql", adress_data_base_test)
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


  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Printf("Подключено")
  var zapros string
  if r.Method == http.MethodPost {

    option_of_train := r.FormValue("option_of_train")
    fmt.Println()
    fmt.Println("!!!!!")
    fmt.Println(option_of_train)
    zapros = fmt.Sprintf("SELECT date, MAX(weight) AS max_weight FROM `trainings` WHERE name_of_train = '%s' AND id_person = '%s' GROUP BY date ORDER BY date ASC", option_of_train, current_user_id)

  }else{
    fmt.Println("----------")
    zapros = fmt.Sprintf("SELECT date, MAX(weight) AS max_weight FROM `trainings` WHERE id_person = '%s' GROUP BY date ORDER BY date ASC", current_user_id)

  }

  //Установка данных
 //insert, err := db.Query(fmt.Sprintf("INSERT INTO test.articles (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
 res,err := db.Query(zapros)
 fmt.Println(zapros)

 var times []time.Time
 var weights []float64

 //fmt.Println(res)
 var dateS string
 dateS = "Empty"
 for res.Next(){

   var weight int
   err = res.Scan(&dateS, &weight)
   fmt.Println(dateS, weight)
   date, err := time.Parse("2006-01-02", dateS)
    if err != nil {
        fmt.Println("Ошибка при парсинге даты:", err)
        return
    }

    fmt.Println("12345678987654321")
   times = append(times, date)
   weights = append(weights, float64(weight))
 }
 fmt.Println(times, weights)
 if( dateS != "Empty"){
   fmt.Println("----------------12345678987654321-----------------")
   graphic(times, weights)
 }

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
     if strings.Contains(path, "/add_to_cart") {
       add_to_cart(w,r)
     }



  }else{
    t.ExecuteTemplate(w, "need_authorization", nil)
  }
}

func search_train_programm(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/search_train_programm.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    panic(err)
  }

//  authorized_user, err  := populateUserFromSession(r, sessionName)

  var data Data_for_send_to_page_View_created_training_programms
//  if(authorized_user.Id == author_id){
//    data.Is_this_Author = true
//  }else{
//    data.Is_this_Author = false
//  }
//  data.Authors_Icon1 = "/personal_static/" + data.Author.Name+"_"+data.Author.Surname+"_"+ data.Author.Id + "/icon_1.jpg"


  if err != nil {
      panic(err)
  }

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  var zapros2 = fmt.Sprintf("SELECT  id,id_author_of_training_program, name_training_program, style_of_training, Description FROM `Training_programms`")
  res2,_ := db.Query(zapros2)
  fmt.Println(zapros2)
  for res2.Next(){
     var programma Training_programma
     var id_author string
     err = res2.Scan(&programma.Id, &id_author, &programma.Name_of_programma, &programma.Style_of_trainings, &programma.Description)
     programma.Author = get_user_data_by_id(id_author)
     programma.Authors_Icon1 =  "/personal_static/" + programma.Author.Name+"_"+programma.Author.Surname+"_"+ programma.Author.Id + "/icon_1.jpg"
     data.Training_programms = append(data.Training_programms, programma)


  }







  fmt.Println(data.Training_programms)
  data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  t.ExecuteTemplate(w, "search_train_programm", data)
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
      db, err := sql.Open("mysql", adress_data_base_test)
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
    db, err := sql.Open("mysql", adress_data_base_test)
    if err != nil{
      panic(err)
    }else{
      fmt.Println("GOOOOD ONE ")
    }

    zapros := fmt.Sprintf("SELECT id, name, surname, nickname FROM `persons` ")
    res,err := db.Query(zapros)

    fmt.Println(zapros)
    if(err != nil){
      panic(err)
    }
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

    db, err := sql.Open("mysql", adress_data_base_test)
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




func get_maximum_in_current_ex(id_person string, name_of_ex string ) string{
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  var zapros = fmt.Sprintf("SELECT MAX(weight) AS max_weight FROM `trainings` WHERE name_of_train = '%s' AND id_person = '%s';", name_of_ex, id_person)
  res,err := db.Query(zapros)
  var max_w string
  for res.Next(){
    err = res.Scan(&max_w)
  }

  return max_w
}



func delet_current_train_from_training_day(w http.ResponseWriter, r *http.Request){

  vars := mux.Vars(r)
  id_ex := vars["id_ex"]

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  query := fmt.Sprintf("DELETE FROM Exercises_in_training_programs WHERE id = %s", id_ex)
	_, err = db.Exec(query)

  referer := r.Header.Get("Referer")
  http.Redirect(w, r, referer, http.StatusSeeOther)
}

func watch_training_programm(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/watch_training_programm.html", "templates/header.html", "templates/footer.html")


  vars := mux.Vars(r)
  if err != nil{
    panic(err)
  }
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()







  var data Data_for_watch_training_programm
  data.Parts_training_programm = make(map[string][]part_of_training_programm)
  data.Id_day_int_train_programm = vars["id_day"]


  var zapros = fmt.Sprintf("SELECT tp.id_author_of_training_program FROM training_days td JOIN Training_programms tp ON td.id_training_programm = tp.id WHERE td.id = %s;", data.Id_day_int_train_programm)
  res,err := db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    err = res.Scan(&data.Id_author_train_programm)
  }

  data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
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

      result, _ := db.Exec("insert into test.trainings (`id_of_all_train`, `name_of_train`,	`weight`,	`count`,	`id_person`,	`date`) values (?, ?, ?, ?, ?, ?)",id_of_train ,option_of_train, weight, count, authorized_user.Id, date)
      fmt.Println(result)

    }

  zapros = fmt.Sprintf("SELECT e.id, n.rus_name AS name_of_train, n.eng_name AS name_of_train,  e.weight, e.count, number_o_execution_sequence, time_to_chill_before_next, unit_of_measurement, technique_exercise FROM Exercises_in_training_programs e JOIN names_of_exercises n ON e.name_of_train = n.eng_name WHERE e.id_of_day_in_programm = %s ORDER BY e.number_o_execution_sequence ASC;", data.Id_day_int_train_programm)
  res,err = db.Query(zapros)
  fmt.Println(zapros)

  var part part_of_training_programm
  for res.Next(){
    var unit_of_measurement string
    err = res.Scan(&part.Id_train, &part.Type_of_train_translated, &part.Type_of_train, &part.Weight, &part.Count, &part.Queue, &part.Time_to_chill, &unit_of_measurement, &part.Technique_exercise_link)


    if(part.Technique_exercise_link == ""){

      var zaprosTech = fmt.Sprintf("SELECT `technique_exercise_defoult`	 FROM `names_of_exercises` WHERE 	`eng_name` = '%s'", part.Type_of_train)
      resTech,_ := db.Query(zaprosTech)
      fmt.Println(zaprosTech)
      var tech_defoult string
      for resTech.Next(){
          _ = resTech.Scan(&tech_defoult)
      }
      part.Technique_exercise_link = tech_defoult
    }


    if(unit_of_measurement == "per_cent"){
      fmt.Print("Watch ------ ")
      fmt.Println(strconv.Atoi(get_maximum_in_current_ex(data.Authorized_user_data.Id, part.Type_of_train)))

      var temp4, _ = strconv.Atoi(part.Weight)
      temp1 := float64(temp4)
      var temp3, _ = strconv.Atoi(get_maximum_in_current_ex(data.Authorized_user_data.Id, part.Type_of_train))
      temp2 := float64(temp3)

      part.Weight = strconv.Itoa(int(temp1 * (temp2/100))) + " (" + part.Weight + "%)"
      fmt.Print("second wantch4   ")
      fmt.Println(temp2)
      fmt.Print("3 wantch   ")
      fmt.Println(strconv.Itoa(int(temp1 * (temp2/100))))
    }
    fmt.Println(part)
    part.Id_train_programm = data.Id_day_int_train_programm

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


    data.Parts_training_programm["part.Date"] = append(data.Parts_training_programm["part.Date"], part)

    fmt.Println(data.Parts_training_programm["part.Date"])
  }
  fmt.Println("Сам смотри")
  fmt.Println(data)

//  fmt.Fprintf(w, "Category: %v\n", vars["id_user"])

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

   db, err := sql.Open("mysql", adress_data_base_test)
   if err != nil{
     panic(err)
   }

   defer db.Close()
   var zapros2 = fmt.Sprintf("SELECT 	id_author_of_training_program, id, name_training_program, style_of_training, Description FROM `Training_programms` WHERE 	id_author_of_training_program = '%s' ", author_id)
   res2,_ := db.Query(zapros2)
   fmt.Println(zapros2)
   for res2.Next(){
      var programma Training_programma
      err = res2.Scan(	&programma.Author.Id, &programma.Id, &programma.Name_of_programma, &programma.Style_of_trainings, &programma.Description)
      data.Training_programms = append(data.Training_programms, programma)
      fmt.Println("Author.Id")
      fmt.Println(programma.Author.Id)

   }


   data.Authorized_user_data, err = populateUserFromSession(r, sessionName)

   fmt.Println("Authorized_user_data.Id")
   fmt.Println(data.Authorized_user_data.Id)
   t.ExecuteTemplate(w, "view_created_training_programms", data)
 }

func refresh_curent_train_programm(w http.ResponseWriter, r *http.Request){

  t, err := template.ParseFiles("templates/refresh_train_programm.html", "templates/header.html", "templates/footer.html")
  fmt.Println("!!!!!!!!")
  if err != nil {
      panic(err)
  }

  var data Data_for_personal_page

  authorized_user, err  := populateUserFromSession(r, sessionName)
  data.Authorized_user_data  = authorized_user




  vars := mux.Vars(r)
  id_programm := vars["id_programm"]
  fmt.Println("Refresh!")
  fmt.Println(id_programm)



  t.ExecuteTemplate(w, "refresh_train_programm", data)

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

     db, err := sql.Open("mysql", adress_data_base_test)
     if err != nil{
       panic(err)
     }

     defer db.Close()

     weight := r.FormValue("weight")
     count := r.FormValue("count")
     option_of_train := r.FormValue("option_of_train")

     unit_of_measurement := r.FormValue("unit_of_measurement")
     time_to_chill_before_next := r.FormValue("time_to_chill_before_next")
     queue := r.FormValue("queue")
     technique_exercise := r.FormValue("technique_exercise")



     fmt.Println("Check this shit!! ------ ")
     fmt.Println(data.Id_train_programm)
    result, err := db.Exec("insert into test.Exercises_in_training_programs (`id_of_day_in_programm`, `name_of_train`,	`weight`,	`count`,	`number_o_execution_sequence`,	`time_to_chill_before_next`, `unit_of_measurement`, `technique_exercise`) values (?, ?, ?, ?, ?, ?, ?,?)",data.Id_train_programm, option_of_train, weight, count, queue, time_to_chill_before_next, unit_of_measurement, technique_exercise)
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
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()

  var zapros = fmt.Sprintf("SELECT name_of_train, weight, count, date FROM `trainings` WHERE 	id_person = '%s'  ORDER BY date DESC ", Person.Id)
  res,err := db.Query(zapros)
  fmt.Println(zapros)


  var part part_of_training_programm
  var date string
  for res.Next(){
    err = res.Scan(&part.Type_of_train, &part.Weight, &part.Count, &date)
    part.Type_of_train_translated = Unique_types_of_exercises[part.Type_of_train]
    data.Parts_training_programm[date] = append(data.Parts_training_programm[date], part)


  }



  data.Temp_mass = make([]string, 0, len(data.Parts_training_programm))
    for k := range data.Parts_training_programm {
        data.Temp_mass = append(data.Temp_mass, k)
    }
  sort.Sort(sort.Reverse(sort.StringSlice(data.Temp_mass)))

  sort.Sort(sort.Reverse(sort.StringSlice(data.Temp_mass)))
  fmt.Println("!!!!!!!!")
  if err != nil {
      panic(err)
  }
  data.Authorized_user_data, err = populateUserFromSession(r, sessionName)
  t.ExecuteTemplate(w, "training_summary", data)
}




func deletTrainProgramm(w http.ResponseWriter, r *http.Request){
  t, _ := template.ParseFiles("templates/delet_train_programm.html", "templates/header.html", "templates/footer.html")

  vars := mux.Vars(r)
  //w.WriteHeader(http.StatusOK)
  id_train_programm := vars["id_train_programm"]
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()

  var zapros = fmt.Sprintf("SELECT id_author_of_training_program, name_training_program FROM `Training_programms` WHERE id = '%s' ", id_train_programm)
  res,err := db.Query(zapros)
  fmt.Println(zapros)

  var data Training_programma
  for res.Next(){
    data.Id = id_train_programm
    _ = res.Scan(&data.Author.Id, &data.Name_of_programma)

  }

  t.ExecuteTemplate(w, "delet_train_programm", data)
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

    db, err := sql.Open("mysql", adress_data_base_test)
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
  db, err := sql.Open("mysql", adress_data_base_test)
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
  db, err := sql.Open("mysql", adress_data_base_test)

  vars := mux.Vars(r)

  id_person := vars["id_person"]


  if err != nil{
    panic(err)
  }

  defer db.Close()
  fmt.Println("1-art")
  if r.Method == http.MethodPost {
    name_new_article := r.FormValue("name_new_article")
    fmt.Println("2-art")
    fmt.Println(id_person , name_new_article)
    result, err := db.Exec("insert into test.theme_article ( `id_author`,	`name_of_theme_article`) values (?, ?)",id_person, name_new_article)
    fmt.Println("4-art")

    fmt.Println(result)
    if(err != nil){
      fmt.Println(err)
    }

  }
  fmt.Println("new articles saved")
  http.Redirect(w, r, "/articles_page/" + id_person, http.StatusSeeOther)
}
func articles_page(w http.ResponseWriter, r *http.Request){
      t, err := template.ParseFiles("templates/articles_page.html", "templates/header.html", "templates/footer.html")
      if err != nil{
        fmt.Fprintf(w, err.Error())
      }
      var data Data_articles_page

      data.Themes_and_name_articles = make(map[string] []string)
      data.Themes = make([]string, 0)
      db, err := sql.Open("mysql", adress_data_base_test)
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



  db, err := sql.Open("mysql", adress_data_base_test)
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
  Data_ref.Authorized_user_data = authorized_user
  Data_ref.Link_adress = adress_web + "/invite_link_page/" + authorized_user.Id
  t.ExecuteTemplate(w, "referal_page", Data_ref)

}








func invite_link_page(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)

  id_person := vars["id_person"]

  t, err := template.ParseFiles("templates/invite_link_page.html", "templates/header.html", "templates/footer.html")

  var Data_ref Data_referals
  if err != nil {
      panic(err)
  }
  saveUserWhoInvite(r, w, sessionName, id_person)
  authorized_user, err  := populateUserFromSession(r, sessionName)


  Data_ref.Authorized_user_data = authorized_user
  http.Redirect(w, r, "/", http.StatusSeeOther)
  t.ExecuteTemplate(w, "invite_link_page", Data_ref)

}



func create_train_day_in_train_programm(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  fmt.Println("^^^^^^^^^^^^^")
  id_programm := vars["id_programm"]
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }



  if r.Method == http.MethodPost {
    name_train_day := r.FormValue("name_train_day")
    number_train_day := r.FormValue("number_train_day")
    result, err := db.Exec("insert into test.training_days ( `id_training_programm`,	`number_of_training_day`, `name_of_training_day`) values (?, ?,?)",id_programm, number_train_day, name_train_day)

    fmt.Println(result)
    if(err != nil){
      fmt.Println(err)
    }

  }



  http.Redirect(w, r, "/view_training_days_of_current_train_programm/" + id_programm, http.StatusSeeOther)
}


func select_in_favorites(w http.ResponseWriter, r *http.Request){


  referer := r.Header.Get("Referer")
  fmt.Println(referer)
  vars := mux.Vars(r)

  id_programm := vars["id_programm"]

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  auth_user,_ := populateUserFromSession(r, sessionName)
  _, err = db.Exec("insert into test.selected_programs ( `id_selected_programm`,	`id_user_whos_select`) values (?, ?)",id_programm, auth_user.Id)

  http.Redirect(w, r, referer, http.StatusSeeOther)
}

func view_training_days_of_current_train_programm(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)

  id_programm := vars["id_programm"]

  t, err := template.ParseFiles("templates/view_training_days_of_current_train_programm.html", "templates/header.html", "templates/footer.html")

  var Data Data_for_view_training_days_of_current_train_programm
  if err != nil {
      panic(err)
  }

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()
  Data.Authorized_user_data,err = populateUserFromSession(r, sessionName)
  var zapros = fmt.Sprintf("SELECT id, name_of_training_day, number_of_training_day FROM `training_days` WHERE 	id_training_programm = '%s'ORDER BY number_of_training_day ASC", id_programm)
  res,err := db.Query(zapros)
  var day_temp Training_day_discrpt
  for res.Next(){
    err = res.Scan(&day_temp.Id, &day_temp.Name, &day_temp.Number_of_training_day)
    Data.Days = append(Data.Days, day_temp)
  }

  zapros = fmt.Sprintf("SELECT id_author_of_training_program FROM `Training_programms` WHERE 	id = '%s' ", id_programm)
  var id_author string
  res,err = db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    fmt.Print("watch this hell ")


    err = res.Scan(&id_author)

    fmt.Print(id_author)
    fmt.Print(Data.Authorized_user_data.Id)
    if(id_author == Data.Authorized_user_data.Id){
      Data.Is_this_author = true;
    }else
    {
      Data.Is_this_author = false;
    }

  }



  Data.Id_programm = id_programm
  fmt.Print("Now see:   ")
  fmt.Println(Data)
  t.ExecuteTemplate(w, "view_training_days_of_current_train_programm", Data)

}



/*
func view_training_days_of_current_train_programm(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)

  id_programm := vars["id_programm"]
  id_day := vars["id_day"]

  t, err := template.ParseFiles("templates/view_training_days_of_current_train_programm.html", "templates/header.html", "templates/footer.html")
  t.ExecuteTemplate(w, "view_training_days_of_current_train_programm", Data)
}


*/

func buy_all_cart(w http.ResponseWriter, r *http.Request){
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()

  var data Data_for_store_main_page
  data.Authorized_user_data, err  = populateUserFromSession(r, sessionName)

  zapros := fmt.Sprintf("SELECT points FROM `score_system` WHERE id_user = %s;", data.Authorized_user_data.Id)
  var score int
  res,err := db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    err = res.Scan(&score)
  }


  zapros = fmt.Sprintf("SELECT p.id, c.quantity, p.price FROM cart c JOIN products_in_store p ON c.id_product = p.id WHERE c.id_user = %s;", data.Authorized_user_data.Id)
  var need_score_to_buy_cart int
  res,err = db.Query(zapros)
  fmt.Println(zapros)

  var querty_varible string

  for res.Next(){
    var quantity, price, id int
    err = res.Scan(&id, &quantity, &price)
    need_score_to_buy_cart = quantity * price
    fmt.Println(quantity)
    fmt.Println(price)
    //querty_varible += fmt.Sprintf(" UPDATE `score_system` SET points = points - %s WHERE id_user = %s; UPDATE `products_in_store` SET quantity = quantity - %s WHERE id = %s; ", strconv.Itoa(price), data.Authorized_user_data.Id, strconv.Itoa(quantity), strconv.Itoa(id))
    querty_varible += fmt.Sprintf(" UPDATE `products_in_store` SET quantity = quantity - %s WHERE id = %s; ",strconv.Itoa(quantity), strconv.Itoa(id))
    querty_varible += fmt.Sprintf(" UPDATE `score_system` SET points = points - %s WHERE id_user = %s; ", strconv.Itoa(price), data.Authorized_user_data.Id)

  }

  fmt.Println(score)

  fmt.Println(need_score_to_buy_cart)
  if(score < need_score_to_buy_cart){
    t, err := template.ParseFiles("templates/not_enough_score_page.html", "templates/header.html", "templates/footer.html")
    if err != nil{
      panic(err)
    }
    t.ExecuteTemplate(w, "not_enough_score_page", data)
  }else{
    zapros = fmt.Sprintf("SELECT p.id, c.quantity, p.price FROM cart c JOIN products_in_store p ON c.id_product = p.id WHERE c.id_user = %s;", data.Authorized_user_data.Id)
    res,err = db.Query(zapros)
    fmt.Println(zapros)
    for res.Next(){
      var quantity, price, id int
      err = res.Scan(&id, &quantity, &price)
      need_score_to_buy_cart = quantity * price
      fmt.Println(quantity)
      fmt.Println(price)

      //querty_varible += fmt.Sprintf(" UPDATE `score_system` SET points = points - %s WHERE id_user = %s; UPDATE `products_in_store` SET quantity = quantity - %s WHERE id = %s; ", strconv.Itoa(price), data.Authorized_user_data.Id, strconv.Itoa(quantity), strconv.Itoa(id))
      querty_varible = fmt.Sprintf(" UPDATE `products_in_store` SET quantity = quantity - %s WHERE id = %s; ",strconv.Itoa(quantity), strconv.Itoa(id))
      _, err = db.Exec(querty_varible)
      querty_varible = fmt.Sprintf(" DELETE FROM `cart` WHERE id_product = %s; ", strconv.Itoa(id))
      _, err = db.Exec(querty_varible)
      fmt.Println(querty_varible)
      if(err != nil){
        panic(err)
      }
      querty_varible = fmt.Sprintf(" UPDATE `score_system` SET points = points - %s WHERE id_user = %s; ", strconv.Itoa(price), data.Authorized_user_data.Id)
      _, err = db.Exec(querty_varible)

      result, err := db.Exec("insert into test.bought_products (`id_in_storage`,	`quantity_bought`, `price`, `id_user`) values (?, ?,?, ?)",id, quantity, price, data.Authorized_user_data.Id)

      fmt.Println(result)
      if(err != nil){
        fmt.Println(err)
      }
    }

    //temp := fmt.Sprintf("%s", querty_varible)
    //fmt.Println(temp)



    if err != nil{
      panic(err)
    }
    returnToLastPage(r, w);
  }


}
func add_new_product_in_store(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/add_new_product_in_store.html", "templates/header.html", "templates/footer.html")

  vars := mux.Vars(r)

  id_product := vars["id_product"]

  var data Data_for_store_main_page
  data.Authorized_user_data, err  = populateUserFromSession(r, sessionName)

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()






  if r.Method == http.MethodPost {
    name := r.FormValue("name")
    brand := r.FormValue("brand")
    weight := r.FormValue("weight")
    price := r.FormValue("price")
    quantity := r.FormValue("quantity")
    id_inputed := r.FormValue("id")








    if(id_inputed == ""){
      result, err := db.Exec("insert into test.products_in_store ( `name`,	`brand`, `weight`, `price`, `quantity`) values (?, ?,?,?,?)",name, brand, weight, price, quantity)

      fmt.Println(result)
      if(err != nil){
        fmt.Println(err)
      }
      zapros := fmt.Sprintf("SELECT id FROM `products_in_store` ORDER BY id DESC LIMIT 1;")
      var id_prod string
      res,err := db.Query(zapros)
      fmt.Println(zapros)
      for res.Next(){
        err = res.Scan(&id_prod)

      }
      uploadsDir := "../store_data/" + id_prod
      fmt.Println("WAY TO :  ",uploadsDir)
      err = os.Mkdir(uploadsDir, os.FileMode(0755))
      if err != nil{
        panic(err)
      }
      err = r.ParseMultipartForm(10 << 20) // Размер максимального загружаемого файла 10MB
      if err != nil {
        http.Error(w, "Failed to parse form", http.StatusInternalServerError)
        return
      }

      // Получаем изображения из формы по разным именам
      image1, _, err := r.FormFile("icon_product")
      if err != nil {
        http.Error(w, "Failed to get image1", http.StatusBadRequest)
        return
      }
      defer image1.Close()

      dir := uploadsDir

      err = saveFile("icon_product.jpg", image1, dir)
      if err != nil {
        http.Error(w, "Failed to save image1", http.StatusInternalServerError)
        return
      }
    }else{

      new_feature := r.FormValue("dropdown2")
      if(new_feature != ""){
        zapros := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM `Features_added_in_products` WHERE id_feature = '%s');", new_feature)

        res,_ := db.Query(zapros)
        var res_check bool
        err = res.Scan(&res_check)
        fmt.Println("Result of statment:")
        fmt.Println(res_check)
        if(!res_check){
          value_of_feature := r.FormValue("value_of_feature")
          _, _ = db.Exec("insert into test.Features_added_in_products (id_product_in_store, id_feature, value) values (?, ?,?)",id_inputed, new_feature, value_of_feature)

        }
      }




      fmt.Println("UPDATE-2")
      q := fmt.Sprintf("UPDATE test.products_in_store SET name = '%s',	brand = '%s', weight = %s, price = %s, quantity = %s WHERE id = %s",name, brand, weight, price, quantity, id_inputed)
      fmt.Println(q)
      _, _ = db.Exec(q)



      uploadsDir := "../store_data/" + id_inputed
      err = os.Remove(uploadsDir+"/icon_product.jpg")

      fmt.Println("AYYYY")

      fmt.Println("WAY TO :  ",uploadsDir)
      err = os.Mkdir(uploadsDir, os.FileMode(0755))

      err = r.ParseMultipartForm(10 << 20) // Размер максимального загружаемого файла 10MB
      if err != nil {
        http.Error(w, "Failed to parse form", http.StatusInternalServerError)
        return
      }

      // Получаем изображения из формы по разным именам
      image1, _, err := r.FormFile("icon_product")
      if err != nil {
        http.Error(w, "Failed to get image1", http.StatusBadRequest)
        return
      }
      defer image1.Close()

      dir := uploadsDir


      err = saveFile("icon_product.jpg", image1, dir)
      if err != nil {
        http.Error(w, "Failed to save image1", http.StatusInternalServerError)
        return
      }


    }
  }


  var product_data Product_data

  if(id_product != ""){
    zapros := fmt.Sprintf("SELECT id, name, brand, weight, price, quantity FROM `products_in_store` WHERE id = %s", id_product)

    res,_ := db.Query(zapros)
    fmt.Println(zapros)
    for res.Next(){
      _ = res.Scan(&product_data.Id, &product_data.Name, &product_data.Brand, &product_data.Weight, &product_data.Price, &product_data.Quantity_in_store)

    }


    zapros = fmt.Sprintf("SELECT id, name_of_features FROM `features_for_products`;")
    var featur Features_for_products
    res,err = db.Query(zapros)
    fmt.Println(zapros)
    for res.Next(){
      err = res.Scan(&featur.Id, &featur.Name)
      data.Features = append(data.Features, featur)
    }

    fmt.Println(data.Features)

  }
  data.Products = append(data.Products, product_data)
  fmt.Println("BMW")
  fmt.Println(data.Products)

  if err != nil{
    panic(err)
  }






  t.ExecuteTemplate(w, "add_new_product_in_store", data)
}

func add_to_cart(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)

  id_product := vars["id_product"]
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()

  User_data, err := populateUserFromSession(r, sessionName)
  result, err := db.Exec("insert into test.cart ( `id_user`,	`id_product`, `quantity`) values (?, ?,?)",User_data.Id, id_product, 1)

  fmt.Println(result)
  if(err != nil){
    fmt.Println(err)
  }
  returnToLastPage(r, w);
}



func delet_product_from_cart(id_product string){

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()


  result, err := db.Exec(fmt.Sprintf("DELETE FROM cart WHERE id_product = %s;", id_product))

  fmt.Println(result)
  if(err != nil){
    fmt.Println(err)
  }

}

func watch_product_page(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/product_page.html", "templates/header.html", "templates/footer.html")
  vars := mux.Vars(r)

  id_product := vars["id_product"]

  var data Data_for_store_main_page
  data.Authorized_user_data, err  = populateUserFromSession(r, sessionName)

  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()

  var product Product_data
  zapros := fmt.Sprintf("SELECT id, name, brand, weight, price, quantity FROM `products_in_store` WHERE id = '%s';" , id_product)
  res,_ := db.Query(zapros)
  fmt.Println(zapros)

  for res.Next(){

    err = res.Scan(&product.Id, &product.Name, &product.Brand, &product.Weight, &product.Price, &product.Quantity_in_store)
    product.Icon_product = "/store_data/" + product.Id + "/icon_product.jpg"

    zapros2 := fmt.Sprintf("SELECT quantity FROM cart WHERE id_product = '%s' AND id_user = '%s' ;", product.Id, data.Authorized_user_data.Id)
    res2,_ := db.Query(zapros2)
    fmt.Println(zapros2)
    for res2.Next(){
      err =  res2.Scan(&product.Quantity_in_cart)
    }

    if(product.Quantity_in_cart >= 1){
      product.In_cart=true
    }else{
      product.In_cart=false
      delet_product_from_cart(product.Id)
    }

    zapros3 := fmt.Sprintf("SELECT f.value, s.name_of_features, s.unit_of_measurement FROM Features_added_in_products f JOIN  features_for_products s ON f.id_feature = s.id WHERE f.id_product_in_store = '%s';", product.Id)
    res3,_ := db.Query(zapros3)
    fmt.Println(zapros3)
    var feature_in_product Features_for_products
    for res3.Next(){
      err =  res3.Scan(&feature_in_product.Value, &feature_in_product.Name, &feature_in_product.Unit_of_measurement)
      product.Features_had = append(product.Features_had, feature_in_product)
    }



  }
  data.Products = append(data.Products, product)
  fmt.Println("check rhis shiiiiit")
  fmt.Println(data.Products)
  fmt.Println(data.Authorized_user_data.Is_it_admin)
  t.ExecuteTemplate(w, "product_page", data)
}



func change_quantity_product_in_cart(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)

  id_product := vars["id_product"]
  plus_or_minus := vars["plus_or_minus"]
  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()


  fmt.Println(plus_or_minus)
  fmt.Println(id_product)
  temp := fmt.Sprintf("UPDATE cart SET quantity = quantity %s WHERE id_product = '%s'; ",plus_or_minus, id_product)
  fmt.Println(temp)
  result, err := db.Exec(temp)

  fmt.Println(result)
  if(err != nil){
    fmt.Println(err)
  }
  returnToLastPage(r, w);
}


func history_of_bought(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/history_of_bought.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    panic(err)
  }
  var data Data_for_store_main_page
  data.Authorized_user_data, err  = populateUserFromSession(r, sessionName)


  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }

  defer db.Close()

  zapros := fmt.Sprintf("SELECT id_in_storage, quantity_bought, date FROM `bought_products` WHERE id_user = '%s' ORDER BY date DESC", data.Authorized_user_data.Id)

  res,err := db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    var product_Id string
    var product Product_data
    err = res.Scan(&product_Id, &product.Quantity_in_cart, &product.Date)
    zapros2 := fmt.Sprintf("SELECT id, name, brand, weight, price, quantity FROM `products_in_store` WHERE id = '%s';" , product_Id)
    res2,_ := db.Query(zapros2)
    fmt.Println(zapros2)

    for res2.Next(){

      err = res2.Scan(&product.Id, &product.Name, &product.Brand, &product.Weight, &product.Price, &product.Quantity_in_store)
      product.Icon_product = "/store_data/" + product.Id + "/icon_product.jpg"

    }
    if(product.Quantity_in_cart > 0){
      data.Products = append(data.Products, product)
    }else{
      delet_product_from_cart(product.Id)
    }






  }



  t.ExecuteTemplate(w, "history_of_bought", data)
}


func cart_main(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/cart_main.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    panic(err)
  }
  var data Data_for_store_main_page
  data.Authorized_user_data, err  = populateUserFromSession(r, sessionName)


  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()

  zapros := fmt.Sprintf("SELECT id_product, quantity FROM `cart` WHERE id_user = '%s'", data.Authorized_user_data.Id)

  res,err := db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    var product_Id string
    var product Product_data
    err = res.Scan(&product_Id, &product.Quantity_in_cart)
    zapros2 := fmt.Sprintf("SELECT id, name, brand, weight, price, quantity FROM `products_in_store` WHERE id = '%s';" , product_Id)
    res2,_ := db.Query(zapros2)
    fmt.Println(zapros2)

    for res2.Next(){

      err = res2.Scan(&product.Id, &product.Name, &product.Brand, &product.Weight, &product.Price, &product.Quantity_in_store)
      product.Icon_product = "/store_data/" + product.Id + "/icon_product.jpg"

    }
    if(product.Quantity_in_cart > 0){
      data.Products = append(data.Products, product)
    }else{
      delet_product_from_cart(product.Id)
    }






  }



  t.ExecuteTemplate(w, "cart_main", data)
}


func store_page(w http.ResponseWriter, r *http.Request){
  t, err := template.ParseFiles("templates/store_main.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    panic(err)
  }
  var data Data_for_store_main_page
  data.Authorized_user_data, err  = populateUserFromSession(r, sessionName)
  sort_type := "new"
  if r.Method == http.MethodPost {
    sort_type = r.FormValue("dropdown2")
  }
  fmt.Println(sort_type)
  if(sort_type == "new"){
    sort_type = "ORDER BY id DESC"
  }else if(sort_type == "cheep"){
    sort_type = "ORDER BY price ASC"
  }else if(sort_type == "expensive"){
    sort_type = "ORDER BY price DESC"
  }




  db, err := sql.Open("mysql", adress_data_base_test)
  if err != nil{
    panic(err)
  }
  defer db.Close()
  fmt.Println(sort_type)
  zapros := fmt.Sprintf("SELECT id, name, brand, weight, price, quantity FROM `products_in_store` %s", sort_type)

  res,err := db.Query(zapros)
  fmt.Println(zapros)
  for res.Next(){
    var product Product_data
    err = res.Scan(&product.Id, &product.Name, &product.Brand, &product.Weight, &product.Price, &product.Quantity_in_store)
    zapros2 := fmt.Sprintf("SELECT quantity FROM cart WHERE id_product = '%s' AND id_user = '%s' ;", product.Id, data.Authorized_user_data.Id)
    res2,_ := db.Query(zapros2)
    //fmt.Println(zapros2)
    for res2.Next(){
      err =  res2.Scan(&product.Quantity_in_cart)
    }
    if(product.Quantity_in_cart >= 1){
      product.In_cart=true
    }else{
      product.In_cart=false
      delet_product_from_cart(product.Id)
    }


    fmt.Println(product.In_cart)



    product.Icon_product = "/store_data/" + product.Id + "/icon_product.jpg"
    data.Products = append(data.Products, product)
  }


  t.ExecuteTemplate(w, "store_main", data)
}

func admin_main(w http.ResponseWriter, r *http.Request){
 // t, _ := template.ParseFiles("templates/view_choosen_programm.html", "templates/header.html", "templates/footer.html")
  t, err := template.ParseFiles("templates/admin_main.html", "templates/header.html", "templates/footer.html")
  if err != nil{
    panic(err)
  }
  var admin_User User
  admin_User.Is_it_admin = true
  saveUserToSession(r, w, sessionName, admin_User)

  t.ExecuteTemplate(w, "admin_main", nil)
}



 func view_choosen_programm(w http.ResponseWriter, r *http.Request){
  // t, _ := template.ParseFiles("templates/view_choosen_programm.html", "templates/header.html", "templates/footer.html")
   t, err := template.ParseFiles("templates/search_train_programm.html", "templates/header.html", "templates/footer.html")
   if err != nil{
     panic(err)
   }

   //  authorized_user, err  := populateUserFromSession(r, sessionName)

   var data Data_for_send_to_page_View_created_training_programms
   //  if(authorized_user.Id == author_id){
   //    data.Is_this_Author = true
   //  }else{
   //    data.Is_this_Author = false
   //  }
   //  data.Authors_Icon1 = "/personal_static/" + data.Author.Name+"_"+data.Author.Surname+"_"+ data.Author.Id + "/icon_1.jpg"
   data.Authorized_user_data, err = populateUserFromSession(r, sessionName)

   if err != nil {
       panic(err)
   }

   db, err := sql.Open("mysql", adress_data_base_test)
   if err != nil{
     panic(err)
   }

   defer db.Close()





   var zapros3 = fmt.Sprintf("SELECT id_selected_programm  FROM `selected_programs` WHERE id_user_whos_select = %s", data.Authorized_user_data.Id)
   res3,_ := db.Query(zapros3)
   fmt.Println(zapros3)
   for res3.Next(){
       var id_choosen_temp string
       err = res3.Scan(&id_choosen_temp)
       var zapros2 = fmt.Sprintf("SELECT  id, id_author_of_training_program, name_training_program, style_of_training, Description FROM `Training_programms` WHERE id = %s", id_choosen_temp)
       res2,_ := db.Query(zapros2)
       fmt.Println(zapros2)
       for res2.Next(){
          var programma Training_programma
          var id_author string
          err = res2.Scan(&programma.Id, &id_author, &programma.Name_of_programma, &programma.Style_of_trainings, &programma.Description)
          programma.Author = get_user_data_by_id(id_author)
          programma.Authors_Icon1 =  "/personal_static/" + programma.Author.Name+"_"+programma.Author.Surname+"_"+ programma.Author.Id + "/icon_1.jpg"
          data.Training_programms = append(data.Training_programms, programma)


       }
     }







   fmt.Println(data.Training_programms)

   t.ExecuteTemplate(w, "search_train_programm", data)

   //t.ExecuteTemplate(w, "view_choosen_programm", data)
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

 fs = http.FileServer(http.Dir("../store_data"))
 http.Handle("/store_data/", http.StripPrefix("/store_data/", fs))


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
  r.HandleFunc("/search_train_programm", search_train_programm)

  r.HandleFunc("/personal_statistic/{id_user}", personal_statistic)
  r.HandleFunc("/create_train_programm_step_1", verification_of_authorization)
//  r.HandleFunc("/create_train_programm_step_2/{id_training_programm}", handlerTTT)

  r.HandleFunc("/watch_training_programm/{id_day}", verification_of_authorization)
  r.HandleFunc("/watch_anthropometric/{id_person}", watch_anthropometric)

  r.HandleFunc("/view_created_training_programms/{id_author}", verification_of_authorization)
  r.HandleFunc("/training_summary/{id_person}", training_summary)
  r.HandleFunc("/record_anthropometric_changes/{id_person}", record_anthropometric_changes)
  //r.HandleFunc("/view_current_training_programm/{id_programm}", verification_of_authorization)
  r.HandleFunc("/articles_page/{id_person}", articles_page)
  r.HandleFunc("/create_name_of_theme_article/{id_person}", create_name_of_theme_article)



  r.HandleFunc("/create_train_set_exercises/{id_programm}", create_train_set_exercises)

  r.HandleFunc("/refresh_curent_train_programm/{id_programm}", refresh_curent_train_programm)

  r.HandleFunc("/processing_create_train_programm_step_1", processing_create_train_programm_step_1)
  r.HandleFunc("/referal_page", referal_page)
  r.HandleFunc("/invite_link_page/{id_person}", invite_link_page)




  //
  r.HandleFunc("/view_training_days_of_current_train_programm/{id_programm}", view_training_days_of_current_train_programm)
  r.HandleFunc("/create_train_day_in_train_programm/{id_programm}", create_train_day_in_train_programm)
  //r.HandleFunc("/view_day_of_current_train_programm/{id_programm}/{id_day}", view_day_of_current_train_programm)

  r.HandleFunc("/select_in_favorites/{id_programm}", select_in_favorites)

  r.HandleFunc("/view_choosen_programm", view_choosen_programm)
  r.HandleFunc("/delete_current_train_from_training_day/{id_ex}", delet_current_train_from_training_day)


  r.HandleFunc("/store_page", store_page)
  r.HandleFunc("/cart_main", cart_main)
  r.HandleFunc("/add_new_product_in_store", add_new_product_in_store)
  r.HandleFunc("/add_new_product_in_store/{id_product}", add_new_product_in_store)
  r.HandleFunc("/add_to_cart/{id_product}", verification_of_authorization )
  r.HandleFunc("/watch_product_page/{id_product}", watch_product_page)


  r.HandleFunc("/buy_all_cart", buy_all_cart)
  r.HandleFunc("/history_of_bought", history_of_bought)

  r.HandleFunc("/change_quantity_product_in_cart/{id_product}/{plus_or_minus}", change_quantity_product_in_cart )

  r.HandleFunc("/admin", admin_main)

 fmt.Println()
 http.Handle ("/", r)
 http.ListenAndServe(":8080", nil)
}
