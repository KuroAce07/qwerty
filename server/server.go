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
	senddata(data)
	*result = strconv.FormatFloat(data.A, 'f', 5, 64) + strconv.FormatFloat(data.B, 'f', 5, 64) + strconv.FormatFloat(data.C, 'f', 5, 64)
	return nil
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Berhasil terhubung ke broker")
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Terima message dari broker : ", string(msg.Payload()), " dengan topik ", msg.Topic())
}

var opts = mqtt.NewClientOptions()
var client mqtt.Client
var client_rpc *rpc.Client
var err error
var listener net.Listener

func main() {
	// opts := mqtt.NewClientOptions()
	// Deklarasikan opsi untuk koneksi dari pub/sub ke broker
	opts.AddBroker("tcp://127.0.0.1:1883")
	// Deklarasikan callback function untuk handling connection
	opts.OnConnect = connectHandler
	// Deklarasikan callback function untuk handle message masuk
	opts.SetDefaultPublishHandler(messageHandler)
	// Kirim permintaan koneksi MQTT ke broker
	client = mqtt.NewClient(opts)
	// Koneksikan dari client ke broker
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println("Terdapat error koneksi : ", token.Error())
	}

	//RPC CLIENT
	client_rpc, err = rpc.DialHTTP("tcp", "127.0.0.1:1234")
	handleError(err)

	// Inisiasi struct arith
	RS := &RumahSakit{}
	// Registrasikan struct dan method ke RPC
	rpc.Register(RS)
	// Deklarasikan bahwa kita menggunakan protokol HTTP sebagai mekanisme pengiriman pesan
	rpc.HandleHTTP()
	// Deklarasikan listerner HTTP dengan layer transport TCP dan Port 1234
	listener, err = net.Listen("tcp", ":"+os.Args[1])
	sub()
	http.Serve(listener, nil)
}

func senddata(data *Data) {
	var result string
	err = client_rpc.Call("RumahSakit.Send", data, &result)
	handleError(err)
}

func sub() {
	topic := "/sensor"
	// Variabel topik message
	// Subscribe dengan topik tertentu dan QoS Level 1
	client.Subscribe(topic, 1, nil)
	fmt.Println("Berhasil subscribe :", topic)
	// Menunggu susbcribe berhasil
}

func pub(data *Data) {
	topic := "/sensor"
	message := "data sensor: " + strconv.FormatFloat(data.A, 'f', 5, 64) + strconv.FormatFloat(data.B, 'f', 5, 64) + strconv.FormatFloat(data.C, 'f', 5, 64)
	fmt.Println("Publish data :", data.A, data.B, data.C, "ke broker :", topic)
	token := client.Publish(topic, 1, false, message)
	token.Wait()
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
	}
}
