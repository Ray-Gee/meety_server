package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	// "strings"
	"time"

	"github.com/Ryuichi-g/meety_server/env"
	"github.com/Ryuichi-g/meety_server/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
    log.Printf("DIALECT: %s", DIALECT)

    dbURI := fmt.Sprintf("user=postgres password=postgres dbname=%s sslmode=disable host=meety-db-1 port=5432", DBNAME) // docker
    // dbURI := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=5432", DBUSER, PASSWORD, DBNAME) // local
    db, err = gorm.Open(DIALECT, dbURI)
    if err != nil {
        log.Fatal("personControllers.go L:56", err)
    } else {
        fmt.Println("Successfully connected to database")
    }
        
    // db.AutoMigrate(&models.Person{})
    // db.AutoMigrate(&models.Book{})
}
func setupResponse(w *http.ResponseWriter, r *http.Request) {
    // (*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	// (*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    // (*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Headers", "*")
    (*w).Header().Set("Access-Control-Max-Age", "864000")
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
    // file, err := os.Create("sample.png")
    // if err != nil {
    //     ctx.JSON(http.StatusInternalServerError, ErrorDocument("internal sesrver error"))
    // }
    // defer file.Close()
    // //画像データのコピー
    // io.Copy(file, ctx.Request.Body)
    fmt.Printf("getpeopleの中")
    // setupResponse(&w, r)
    // if r.Method == "OPTIONS" {
    //     w.WriteHeader(http.StatusOK)
    //     return
    // }
    // if r.Method == "OPTIONS" {
    //     //ヘッダーにAuthorizationが含まれていた場合はpreflight成功
    //     s := r.Header.Get("Access-Control-Request-Headers")
    //     fmt.Printf("ヘッダー中身%s", s)
    //     if strings.Contains(s, "authorization") == true || strings.Contains(s, "Authorization") == true {
    //     w.WriteHeader(204)
    //     }
    //     w.WriteHeader(400)
    //     return
    // }
    fmt.Printf("OPTIONSの下")

	var people []models.Person
	db.Limit(12).Order("id desc").Find(&people)
    fmt.Printf("people: %#v", people)
	for i, v := range people {
		fmt.Printf("%+v\n", v.Name)
		birthday := v.Birthday
		age, _ := calcAge(birthday)
		fmt.Printf("%+v\n", age)
		people[i].Age = age
        const layout = "2006-01-02"
        fmt.Printf("birthday format: %v\n", birthday.Format(layout))
        people[i].BirthdayFormatted = birthday.Format(layout)
	}
	json.NewEncoder(w).Encode(&people)
    fmt.Printf("getpeople終わり")
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

func CreatePerson(w http.ResponseWriter, r *http.Request) {
    // defer db.Close()
    fmt.Printf("HERE ## creeatePerson: golang %#v", r)
    setupResponse(&w, r)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    fmt.Printf("createPerson, OPTIONSの下")

	var person models.Person

    // パース
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        fmt.Println("ParseMultipartForm error")
    }
    mf := r.MultipartForm

    // 通常のリクエスト
    for k, v := range mf.Value {
        fmt.Printf("MultipartForm %v : %v\n", k, v)
    }
    fmt.Printf("mf.Value name: %v\n", mf.Value["name"])
    fmt.Printf("mf.Value email: %v\n", mf.Value["email"])
    person.Name = mf.Value["name"][0]
    person.Email = mf.Value["email"][0]

    // ファイルのリクエスト
    var fileName string
    var f multipart.File
    for k, v := range mf.File {
        fmt.Printf("file => %v : %v\n", k, v)
        for _, vv := range v {
            f, _ = vv.Open()
            fmt.Printf("open: %v\n", f)

            // アンコメントするとバイト型でperson.Sourceに保存
            // img, err := jpeg.Decode(f)
            // fmt.Printf("img %v\n", img) // ... 132 133 133 134 134 133 133 132 126 126 126 126 126 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125 125] 400 200 YCbCrSubsampleRatio422 (0,0)-(400,602)}
            // if err != nil {
            //     log.Fatal(err)
            // }
            // buffer := new(bytes.Buffer)
            // if err := jpeg.Encode(buffer, img, nil); err != nil {
            //     log.Println("unable to encode image.")
            // }
            // imageBytes := buffer.Bytes()
            // fmt.Printf("imageBytes %+v\n", imageBytes)
            // person.Source = imageBytes  

            fileName = vv.Filename
			// fmt.Printf("%v : %v\n", k, vv.Filename) // image : loadjpeg.jpg
			// fmt.Printf("%v : %v\n", k, vv.Header) // image : map[Content-Disposition:[form-data; name="image"; filename="loadjpeg.jpg"] Content-Type:[image/jpeg]]
			// fmt.Printf("%v : %v\n", k, vv.Size) // image : 57682
		}
    }
    // uuid作成
    extension := filepath.Ext(fileName)
    uu, err := uuid.NewRandom()
    if err != nil {
            fmt.Println(err)
    }
    newPath := uu.String() + extension
    fmt.Printf("newPath%+v\n", newPath)

    // clientFilePath := filepath.Join("./images/", fileName)
    clientFilePath := filepath.Join("./images/", newPath)
    filePath := filepath.Join("../client/public/", clientFilePath)
    fmt.Printf("filePathNew: %+v\n", filePath)
    nf, err := os.Create(filePath)
    if err != nil {
        log.Printf("[os.Create] %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer nf.Close()
    fmt.Printf("nf: %+v\n", nf)
    _, err = io.Copy(nf, f)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
	person.ImagePath = clientFilePath
	json.NewDecoder(r.Body).Decode(&person)
    fmt.Printf("filePath %v", filePath)
    fmt.Printf("person %v", person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
    } else {
        json.NewEncoder(w).Encode(&person)
    }
    // w.Header().Set("Content-Type", "text/html")
    // w.Header().Set("location", "http://www.yahoo.co.jp/")
    // w.WriteHeader(http.StatusMovedPermanently)
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
