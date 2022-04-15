package listener

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"strconv"
)

type Tcp struct {
	Listener
}

func (t *Tcp) Create(l Endpoint, return_ch *chan string) {
	listen.Address = l.Address
	listen.Port = l.Port

	var address = ":" + strconv.Itoa(listen.Port)

	s, err := net.Listen("tcp4", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	go func() {
		//defer s.Close()

		for {
			c, err := s.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go handleConnection(c, return_ch)
		}
	}()
}

func handleConnection(c net.Conn, r *chan string) error {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	reader := bufio.NewReader(c) //.ReadString('\n')
	tp := textproto.NewReader(reader)
	for {

		line, err := tp.ReadLine()
		if len(line) == 0 && err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return err
		}

		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				break
			}
			return err
		}

		*r <- line
	}
	defer c.Close()

	return nil
}
