package storage

import (
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
