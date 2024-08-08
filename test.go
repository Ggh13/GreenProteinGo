package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // Импорт драйвера MySQL
)

func main() {
    // Строка подключения
    var addressDataBase = "user:password@tcp(147.45.163.58:3306)/test"

    // Открытие соединения с базой данных
    db, err := sql.Open("mysql", addressDataBase)
    if err != nil {
        fmt.Println("Error opening database:", err)
        return
    }
    defer db.Close()

    // Проверка соединения
    err = db.Ping()
    if err != nil {
        fmt.Println("Error connecting to the database:", err)
        return
    }

    // Выполнение SELECT запроса
    rows, err := db.Query("SELECT id FROM `persons`")
    if err != nil {
        fmt.Println("Error executing query:", err)
        return
    }
    defer rows.Close()

    // Обработка результатов
    for rows.Next() {

        var column2 int    // Замените на соответствующий тип данных

        if err := rows.Scan(&column2); err != nil {
            fmt.Println("Error scanning row:", err)
            return
        }

        fmt.Println(column2)
    }

    if err := rows.Err(); err != nil {
        fmt.Println("Error iterating rows:", err)
    }
}
