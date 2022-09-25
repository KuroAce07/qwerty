package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

type MyArgs struct {
	A, B, C float64
}

type RumahSakit struct {
}

func (t *RumahSakit) Send(args *MyArgs, result *string) error {
	fmt.Println("Menerima data: ")
	fmt.Printf("\tSuhu: %s", strconv.FormatFloat(args.A, 'f', 2, 64))
	fmt.Printf("\tKelembapan: %s", strconv.FormatFloat(args.B, 'f', 2, 64))
	fmt.Printf("\tCo2: %s\n", strconv.FormatFloat(args.C, 'f', 2, 64))
	*result = "data terkirim"
	return nil
}

func main() {

	// Inisiasi struct arith
	arith := &RumahSakit{}
	// Registrasikan struct dan method ke RPC
	rpc.Register(arith)
	// Deklarasikan bahwa kita menggunakan protokol HTTP sebagai mekanisme pengiriman pesan
	rpc.HandleHTTP()
	// Deklarasikan listerner HTTP dengan layer transport TCP dan Port 1234
	listener, err := net.Listen("tcp", ":1234")
	handleError(err)
	// Jalankan server HTTP
	http.Serve(listener, nil)

}
func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}
