package main

import (
	"fmt"
	"io/ioutil"
	"meety/server/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func add(a, b int) int {
// 	return a + b
// }


func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world!"))
}

func TestGetPeople(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	controllers.GetPeople(w, r)
	// helloHandler(w, r)

	resp := w.Result()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read test response: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("got = %d, want = 200", resp.StatusCode)
	}

	if string(body) == "" {
		t.Errorf("got = %s, 中身がからです。", body)
	}
}

func ItoaByFmt(n int) {
	var s []string
	for i := 0; i < n; i++ {
		s = append(s, fmt.Sprint(i))
	}
}
func BenchmarkItoaByFmt(b *testing.B)      { ItoaByFmt(b.N) }
// func TestAdd(t *testing.T) {
// 	type args struct {
// 		a int
// 		b int
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want int
// 	}{
// 		{
// 			name: "normal",
// 			args: args{a: 1, b: 2},
// 			want: 3,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := add(tt.args.a, tt.args.b); got != tt.want {
// 				t.Errorf("add() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestAdd(t *testing.T) {
// 	if add(1, 2) != 3 {
// 		t.Errorf("add() = %v, want %v", add(1, 2), 3)
// 	}
// }