package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
)

type MyArgs struct {
	A, B, C float64
}

func main() {

	// Inisiasi koneksi HTTP dari client ke server
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:"+os.Args[1])
	handleError(err)
	// Argument yang akan dikirimkan dari client ke server
	suhu, err := strconv.ParseFloat(os.Args[2], 64)
	kelembapan, err := strconv.ParseFloat(os.Args[3], 64)
	karbon, err := strconv.ParseFloat(os.Args[4], 64)
	clientargs := &MyArgs{suhu, kelembapan, karbon}
	// Pointer untuk menampung hasil eksekusi dari server
	var result string
	err = client.Call("RumahSakit.Send", clientargs, &result)
	handleError(err)
	fmt.Println("Hasil eksekusi RPC server : ", result)

}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}
