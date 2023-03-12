//lint:file-ignore ST1006 reason: ...

package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const tcp string = "tcp";

type Message struct {
	from string
	payload []byte
}

type Server struct {
	listenAddr string
	listener   net.Listener
	quitChannel chan struct{}
	messageChannel chan Message
}

func NewServer(listenAddr string) *Server {

	return &Server{
		listenAddr: listenAddr,
		quitChannel: make(chan struct{}),
		messageChannel: make(chan Message, 10),
	};

}

func (self *Server) Start() error {

	listener, err := net.Listen(tcp, self.listenAddr);

	if err != nil {
		return err;
	}

	defer listener.Close();

	self.listener = listener;

	go self.AcceptLoop();
	
	<-self.quitChannel;
	close(self.messageChannel);

	return nil;
}

func (self *Server) AcceptLoop(){
	
	for {

		connection, err := self.listener.Accept(); 

		if err != nil {
			fmt.Println("accept error: ", err);
			continue;
		}

		fmt.Println("new connection to the server, from ", 
		connection.RemoteAddr().String())
		go self.readLoop(connection);
	}
}



func (self *Server) readLoop(connection net.Conn){

	defer connection.Close();

	buffer := make([]byte, 2048);

	for {

		n, err := connection.Read(buffer);

		if err != nil {
			if err == io.EOF {
				break;
			}
			fmt.Println("read error: ", err);
			continue
		}

		message := buffer[:n];

		self.messageChannel <- Message{
			from : connection.LocalAddr().String(),
			payload: message,
			};

		
	}
}

func main() {
	server := NewServer(":3000");

	go func(){
		for msg := range server.messageChannel {
			fmt.Printf("received message from connection {%s}: {%s} \n",
			 msg.from, 
			 string(msg.payload))
		}
	}();

	log.Fatal(server.Start());
}