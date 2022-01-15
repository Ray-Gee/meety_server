package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// "io/ioutil"

	// "math/rand"
	"strconv"

	// "sync"
	"time"

	"log"
	"meety/server/env"
	"net/http"
	"os"

	// "time"

	// "github.com/gorilla/mux"
	// "gorm.io/gorm"
	"meety/server/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// type Person struct {
// 	gorm.Model

// 	Name string
// 	Email string `gorm:"typevarchar(100);unique_index"`
// 	Books []Book
// }

// type Book struct {
// 	gorm.Model

// 	Title string
// 	Author string
// 	CallNumber int `gorm:"unique_index"`
// 	PersonID int
// }

var db *gorm.DB
var err error

func init() {
    
    env.LoadEnv()
    DIALECT := os.Getenv("DIALECT")
    DBUSER := os.Getenv("DBUSER")
    DBNAME := os.Getenv("DBNAME")
    PASSWORD := os.Getenv("PASSWORD")
    log.Printf("user=%s password=%s dbname=%s sslmode=disable\n", DBUSER, PASSWORD, DBNAME)

    dbURI := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DBUSER, PASSWORD, DBNAME)
    db, err = gorm.Open(DIALECT, dbURI)
    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println("Successfully connected to database")
    }
        
    db.AutoMigrate(&models.Person{})
    db.AutoMigrate(&models.Book{})
}
func setupResponse(w *http.ResponseWriter, r *http.Request) {

	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    // (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    (*w).Header().Set("Access-Control-Allow-Headers", "*")
    (*w).Header().Set( "Access-Control-Allow-Credentials", "true" )
}

func calcAge(t time.Time) (int, error) {
    // 現在日時を数値のみでフォーマット (YYYYMMDD)
    dateFormatOnlyNumber := "20060102" // YYYYMMDD

    now := time.Now().Format(dateFormatOnlyNumber)
    birthday := t.Format(dateFormatOnlyNumber)

    // 日付文字列をそのまま数値化
    nowInt, err := strconv.Atoi(now)
    if err != nil {
        return 0, err
    }
    birthdayInt, err := strconv.Atoi(birthday)
    if err != nil {
        return 0, err
    }

    // (今日の日付 - 誕生日) / 10000 = 年齢
    age := (nowInt - birthdayInt) / 10000
    return age, nil
}

// func DateValidation(fl validator.FieldLevel) bool {
//     _, err := time.Parse("2006-01-02", fl.Field().String())
//     if err != nil {
//         return false
//     }
//     return true
// }

func GetPeople(w http.ResponseWriter, r *http.Request) {

    setupResponse(&w, r)
	var people []models.Person
	db.Limit(10).Order("id desc").Find(&people)
	for i, v := range people {
		fmt.Printf("%+v\n", v.Name)
		birthday := v.Birthday
		age, _ := calcAge(birthday)
		fmt.Printf("%+v\n", age)
		people[i].Age = age
	}
	json.NewEncoder(w).Encode(&people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person models.Person
	var books []models.Book

	db.First(&person, params["id"])
	db.Model(&person).Related(&books)

	person.Books = books
	
	json.NewEncoder(w).Encode(&person)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    var person models.Person
    // 生年月日形式
    // {"birthday": "2022-02-02T04:10:03Z"}
    db.First(&person, params["id"])
    json.NewDecoder(r.Body).Decode(&person)
    updatedPerson := db.Save(&person)
    err = updatedPerson.Error
    if err != nil {
        json.NewEncoder(w).Encode(err)
    } else {
        json.NewEncoder(w).Encode(&person)
    }
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
    // defer db.Close()

    setupResponse(&w, r)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    fmt.Printf("before formfile: %v\n", r)
    file, h, _ := r.FormFile("image")
    bs, err := ioutil.ReadAll(file)    // ファイルの中身を読む
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    s := string(bs)
    fmt.Printf("s: %v\n", s)
    fmt.Printf("h: %v\n", h.Filename)
    filePath := filepath.Join("./images/", h.Filename)
    // fmt.Printf("filePath: %v\n", filePath)
    // cur, _ := os.Getwd()
	// fmt.Println("pwd", cur)
    nf, err := os.Create(filePath)
    fmt.Printf("nf: %v\n", nf)
    if err != nil {
        log.Printf("[os.Create] %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer nf.Close()
    _, err = nf.Write(bs)    // 作成したファイルに、読み込んだファイルと同じ内容を書き込む
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // fmt.Printf("s: %v\n", s)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer file.Close()
    
	var person models.Person
    
	json.NewDecoder(r.Body).Decode(&person)
	person.ImagePath = filePath
    fmt.Printf("filePath %v", filePath)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
	
}

func DeletePerson(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	var person models.Person

	db.First(&person, params["id"])
	db.Delete(&person)

	json.NewEncoder(w).Encode(&person)
}



func GetBooks(w http.ResponseWriter, r *http.Request)  {
	var books []models.Book

	db.Find(&books)
	json.NewEncoder(w).Encode(&books)
}

func GetBook(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	var book models.Book

	db.First(&book, params["id"])
	json.NewEncoder(w).Encode(&book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
	
}

func DeleteBook(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	var book models.Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}

//成功version
// db.Model(&person).Where("id=?", i).Update("birthday", time.Unix(unixtime, 0).Format(date_format))