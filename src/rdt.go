package main

import (
	"net"
	"fmt"
	"io"
	"runtime"
	"regexp"
)

type Socket struct {
	addr string
	ln   net.Listener
	conn net.Conn
	err  error
}

func NewSocket(addr string) *Socket {
	return &Socket{
		addr:addr,
		ln:nil,
		conn:nil,
		err:nil,
	}
}

func EstablishSocket(A Socket, B Socket) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Socket A started at:", A.addr)
	fmt.Println("Socket B started at:", B.addr)

	A.ln, A.err = net.Listen("tcp", A.addr)
	HandleError(A.err, "A listen")

	for {
		B.conn, B.err = net.Dial("tcp", B.addr)
		HandleError(B.err, "B dial")

		A.conn, A.err = A.ln.Accept()
		HandleError(A.err, "A accept")

		go HandleTransmission(A, B, "A to B")
		go HandleTransmission(B, A, "B to A")
	}
}

func HandleTransmission(A, B Socket, tag string) {
	for {
		buf := make([]byte, 2048)
		i, err := A.conn.Read(buf)
		if HandleError(err, tag, "Read") < 0 {
			return
		}
		_, err = B.conn.Write(buf[:i])
		if HandleError(err, tag, "Write") < 0 {
			return
		}
	}
}

func HandleError(err error, tag... string) int {
	if err != nil {
		errTimedOut, _ := regexp.MatchString("operation timed out", err.Error())
		if errTimedOut {
			return 0
		} else if err != io.EOF {
			fmt.Println(tag, "INFO:", err)
		}
		return -1
	}
	return 0
}