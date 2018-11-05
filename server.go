package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var zeroParamCommands = []string{"PING"}
var oneParamCommands = []string{"GET"}
var twoParamCommands = []string{"SET"}
var commands = map[string][]string{
	"*1": zeroParamCommands,
	"*2": oneParamCommands,
	"*3": twoParamCommands,
}
var commandsMap = map[string]string{
	"PING": "*1",
	"GET": "*2",
	"SET": "*3",
}

func main() {
	fmt.Println("Start the tcp socket")
	//http.HandleFunc("/", handler)	//	each	request	calls	handler
	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
	startTcpServer()

}

//	handler	echoes	the	Path	component	of	the	request	URL	r.
//func handler(w http.ResponseWriter, r *http.Request) {
//	//fmt.Fprintf(w,	"URL.Path	=	%q\n",	r)
//	fmt.Printf("URL.Path	=%q\n", r.Method)
//}

func startTcpServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:6378")
	chkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	chkError(err)
	for {
		conn, err := listener.AcceptTCP()
		chkError(err)
		go handleStuff(conn)
	}
}

func chkError(err error) {
	if err != nil {
		log.Fatal(err);
	}
}

func handleStuff(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		fmt.Println("req", n)
		chkError(err)
		rAddr := conn.RemoteAddr()
		fmt.Println(rAddr.String())
		req := string(buf[:n])
		reqArr := strings.Split(req, "\r\n")
		reqArr = reqArr[0:len(reqArr)-1]
		fmt.Println("接收到客户端的消息：", reqArr)
		handleCommands(reqArr, conn)
	}
}

func handleCommands(reqArr []string, conn net.Conn) {
	paramNumber := reqArr[0]
	commandName := strings.ToUpper(reqArr[2])
	_, ok := commandsMap[commandName]
	if !ok {
		handleCommandError(1200, commandName, conn)
		return
	}
	if commandsMap[commandName] != paramNumber {
		handleCommandError(1201, commandName, conn)
		return
	}
	handleCommand(commandName, conn)
}

func handleCommandError(errorCode int, commandName string, conn net.Conn) {
	switch errorCode {
	case 1200:
		conn.Write([]byte("+(error) ERR unknown command '" + commandName + "' \r\n"))
	case 1201:
		conn.Write([]byte("+(error) ERR wrong number of arguments for " + commandName + " command\r\n"))
	default:

	}
	fmt.Println("this connect end")
}

func handleCommand(commandName string, conn net.Conn) {
	switch commandName {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	default:
		conn.Write([]byte("+OTHER COMMAND\r\n"))
	}
	fmt.Println("this connect end")
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

