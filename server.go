package main

import (
    "fmt"
    "log"
    "net"
    "reflect"
    "strings"
)

var zeroParamCommands = []string{"PING", "EXIT"}
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
	"EXIT": "*1",
}

var commandReflect = map[string]interface{}{
    "ping": ping,
    "exit": exit,
    "get": get,
    "set": set,
}

type key interface{}
var valueMap = make(map[key]string)

func main() {
	fmt.Println("Start the tcp socket")
	startTcpServer()
}

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

// handle the tcp message from client
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
		fmt.Println("receive the client message：", reqArr)
		if reqArr[2] == "COMMAND" {
            break
        }
		handleCommands(reqArr, conn)
	}
}

func handleCommands(reqArr []string, conn net.Conn) {
	paramNumber := reqArr[0]
	commandName := strings.ToUpper(reqArr[2])
	_, ok := commandsMap[commandName]
	if !ok {
		handleCommandError(1000, commandName, conn)
		return
	}
	if commandsMap[commandName] != paramNumber {
		handleCommandError(1001, commandName, conn)
		return
	}
	//handleCommand(reqArr, commandName, conn)
    Apply(commandReflect[strings.ToLower(commandName)], []interface{}{reqArr, conn})
}

// handle the right command from client
func handleCommand(reqArr []string, commandName string, conn net.Conn) {
   switch commandName {
   case "PING":
       conn.Write([]byte("+PONG\r\n"))
   case "GET":
       result, ok := valueMap[reqArr[4]]
       if !ok {
           conn.Write([]byte("+(nil)\r\n"))
       } else {
           conn.Write([]byte("+\"" + result + "\"\r\n"))
       }
   case "SET":
       valueMap[reqArr[4]] = reqArr[6]
       conn.Write([]byte("+OK\r\n"))
   case "EXIT":
       conn.Close()
   default:
       conn.Write([]byte("+OTHER COMMAND\r\n"))
   }
   fmt.Println("this connect end")
}

func ping(reqArr []string, conn net.Conn) {
    conn.Write([]byte("+PONG\r\n"))
}

func exit(reqArr []string, conn net.Conn) {
    conn.Close()
}

func set(reqArr []string, conn net.Conn) {
    valueMap[reqArr[4]] = reqArr[6]
    conn.Write([]byte("+OK\r\n"))
}

func get(reqArr []string, conn net.Conn) {
    result, ok := valueMap[reqArr[4]]
    if !ok {
        conn.Write([]byte("+(nil)\r\n"))
    } else {
        conn.Write([]byte("+\"" + result + "\"\r\n"))
    }
}

func Apply(f interface{}, args []interface{})([]reflect.Value){
    fun := reflect.ValueOf(f)
    in := make([]reflect.Value, len(args))
    for k, param := range args{
        in[k] = reflect.ValueOf(param)
    }
    r := fun.Call(in)
    return r
}

// handle the error of the command from client
func handleCommandError(errorCode int, commandName string, conn net.Conn) {
	switch errorCode {
	case 1000:
		conn.Write([]byte("+(error) ERR unknown command '" + commandName + "' \r\n"))
	case 1001:
		conn.Write([]byte("+(error) ERR wrong number of arguments for " + commandName + " command\r\n"))
	default:

	}
	fmt.Println("this connect end")
}

// Judge if the element exists in array
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

