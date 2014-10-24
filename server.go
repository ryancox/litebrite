package main

import (
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	"github.com/kidoman/embd/sensor/isl29125"
)

type ColorReading struct {
	Red   uint16
	Green uint16
	Blue  uint16
}

var isl *isl29125.ISL29125

func DataServer(ws *websocket.Conn) {
	for {
		r, err := isl.Reading()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got: %v\n", r)

		cr := &ColorReading{Red: r.Red, Green: r.Green, Blue: r.Blue}
		websocket.JSON.Send(ws, cr)
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}

	bus := embd.NewI2CBus(1)

	isl = isl29125.New(isl29125.DefaultConfig, bus)

	if err := isl.Init(); err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.Handle("/data", websocket.Handler(DataServer))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
