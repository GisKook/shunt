package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/giskook/gotcp"
	"github.com/giskook/shunt"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// read configuration
	shunt.Config, _ = shunt.ReadConfig("./conf.json")

	port := shunt.Config.ServerConfig.BindPort

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+fmt.Sprint(port))
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &shunt.Callback{}, &shunt.TrackerProtocol{})

	// starts service
	srv.Start(listener, time.Second)
	log.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
