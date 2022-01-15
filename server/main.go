package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"meety/server/controllers"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the server")

// var (
// 	person = &Person{Name: "Jack", Email: "jack@email.com"}
// 	books = []Book{
// 		{Title: "The Rules of Thinking", Author: "Richard Templar", CallNumber: 1234, PersonID: 1},
// 		{Title: "Book 2", Author: "Author 2", CallNumber: 12342, PersonID: 1},

// 	}
// )


func main() {
	// controllers.StartWebServer()

	router := mux.NewRouter()

	router.HandleFunc("/byte", byteSample).Methods("GET")

	router.HandleFunc("/person/{id}", controllers.GetPerson).Methods("GET")
	router.HandleFunc("/people", controllers.GetPeople).Methods("GET")
	router.HandleFunc("/create/person", controllers.CreatePerson).Methods("OPTIONS", "POST")
	// router.HandleFunc("/create/person", controllers.CreatePerson).Methods("POST")
	router.HandleFunc("/update/person/{id}", controllers.UpdatePerson).Methods("PUT")
	router.HandleFunc("/delete/person/{id}", controllers.DeletePerson).Methods("DELETE")
	
	router.HandleFunc("/book/{id}", controllers.GetBook).Methods("GET")
	router.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/create/book", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/delete/book/{id}", controllers.DeleteBook).Methods("DELETE")
	
	router.HandleFunc("/", RedirectHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8080", router))

	// パス変数で正規表現を使用
	// 	r.HandleFunc("/hello/{name}/{age:[0-9]+}", RegexHandler)

	// クエリ文字列の取得
	// r.HandleFunc("/hi/", QueryStringHandler)
 	
	// 静的ファイルの提供 	
	// $PROROOT/assets/about.html が http://localhost:8080/assets/about.html でアクセスできる
	// 	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))


}


func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/people", http.StatusFound)
}

func byteSample(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("バイトサンプル\n"))
}

// func RegexHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	fmt.Fprintf(w, "%s is %s years old\n", vars["name"], vars["age"])
// }

// func QueryStringHandler(w http.ResponseWriter, r *http.Request) {
// 	q := r.URL.Query()
// 	fmt.Fprintf(w, "%s Loves Gorilla\n", q.Get("name"))
// }

// const MyContextKey = 1

// func UseContext(handler http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		context.Set(r, MyContextKey, "Call SomeMiddleware")
// 		handler.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(fn)
// }

// func SomeHandler1(w http.ResponseWriter, r *http.Request) {
// 	contextVal := context.Get(r, MyContextKey)
// 	fmt.Fprintf(w, "%s Call SomeHandler1", contextVal)
// }

// func SomeHandler2(w http.ResponseWriter, r *http.Request) {
// 	contextVal := context.Get(r, MyContextKey)
// 	fmt.Fprintf(w, "%s Call SomeHandler2", contextVal)
// }
