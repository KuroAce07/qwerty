package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Data struct {
	A, B, C float64
}

type RumahSakit struct {
}

func (t *RumahSakit) Send(data *Data, result *string) error {
	pub(data)
	*result = strconv.FormatFloat(data.A, 'f', 5, 64) + strconv.FormatFloat(data.B, 'f', 5, 64) + strconv.FormatFloat(data.C, 'f', 5, 64)
	return nil
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Berhasil terhubung ke broker")
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Terima message dari broker : ", string(msg.Payload()), " dengan topik ", msg.Topic())
	//fmt.Println("Terima message")
}

var opts = mqtt.NewClientOptions()
var client_mqtt mqtt.Client

//var token mqtt.Token

func main() {
	//mqtt
	// opts := mqtt.NewClientOptions()
	// Deklarasikan opsi untuk koneksi dari pub/sub ke broker
	opts.AddBroker("tcp://127.0.0.1:1883")
	// Deklarasikan callback function untuk handling connection
	opts.OnConnect = connectHandler
	// Deklarasikan callback function untuk handle message masuk
	opts.SetDefaultPublishHandler(messageHandler)
	// Kirim permintaan koneksi MQTT ke broker
	client := mqtt.NewClient(opts)
	// Koneksikan dari client ke broker
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println("Terdapat error koneksi : ", token.Error())
	}

	// Inisiasi struct arith
	RS := &RumahSakit{}
	// Registrasikan struct dan method ke RPC
	rpc.Register(RS)
	// Deklarasikan bahwa kita menggunakan protokol HTTP sebagai mekanisme pengiriman pesan
	rpc.HandleHTTP()
	// Deklarasikan listerner HTTP dengan layer transport TCP dan Port 1234
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	handleError(err)
	// Jalankan server HTTP
	sub()
	http.Serve(listener, nil)
	// Tangkap input dari user untuk menentukan program dijalankan sebagai publisher atau susbcriber
}

func sub() {
	topic := "sensor"
	// Variabel topik message
	// Subscribe dengan topik tertentu dan QoS Level 1
	client_mqtt.Subscribe(topic, 1, messageHandler)
	fmt.Println("Berhasil subscribe :", topic)
	// Menunggu susbcribe berhasil
	//token.Wait()
}

func pub(data *Data) {
	topic := "/sensor"
	message := "data sensor: " + strconv.FormatFloat(data.A, 'f', 5, 64) + strconv.FormatFloat(data.B, 'f', 5, 64) + strconv.FormatFloat(data.C, 'f', 5, 64)
	fmt.Println("Publish data :", data.A, data.B, data.C, "ke broker :", topic)
	token := client_mqtt.Publish(topic, 1, false, message)
	token.Wait()
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}
