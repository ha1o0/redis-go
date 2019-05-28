package main

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
)

//var zeroParamCommands = []string{"PING", "EXIT"}
//var oneParamCommands = []string{"GET", "EXISTS", "DEL"}
//var twoParamCommands = []string{"SET"}
//var atLeastThreeParamCommands = []string{"RPUSH"}
//var commands = map[string][]string{
//	"*1": zeroParamCommands,
//	"*2": oneParamCommands,
//	"*3": twoParamCommands,
//}
var commandsMap = map[string]string{
	"PING": "1",
	"GET": "2",
	"SET": "3",
	"EXISTS": "2",
	"DEL": "2",
	"RPUSH": ">=3",
	"RPOP": "2",
	"LPOP": "2",
	"LLEN": "2",
	"EXIT": "1",
}

var commandReflect = map[string]interface{}{
    "ping": ping,
    "get": get,
    "set": set,
    "exists": exists,
    "del": del,
    "rpush": rpush,
    "rpop": rpop,
    "lpop": lpop,
    "llen": llen,
	"exit": exit,
}

var valueMap = make(map[string]interface{})

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
		log.Fatal(err)
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
	paramNumberStr := reqArr[0]
	commandName := strings.ToUpper(reqArr[2])
	_, ok := commandsMap[commandName]
	if !ok {
		handleCommandError(1000, commandName, conn)
		return
	}
	paramNumber, _ := strconv.ParseInt(strings.TrimPrefix(paramNumberStr, "*"), 0, 64)
	paramRequireNumberStr := commandsMap[commandName]
	if strings.Contains(paramRequireNumberStr, ">=") {
		paramRequireNumber, _ := strconv.ParseInt(strings.TrimPrefix(paramRequireNumberStr, ">="), 0, 64)
		if paramNumber < paramRequireNumber {
			handleCommandError(1001, commandName, conn)
			return
		}
		Apply(commandReflect[strings.ToLower(commandName)], []interface{}{reqArr, conn, int(paramNumber)})
	} else {
		paramRequireNumber, _ := strconv.ParseInt(paramRequireNumberStr, 0, 64)
		if paramRequireNumber != paramNumber {
			handleCommandError(1001, commandName, conn)
			return
		}
		Apply(commandReflect[strings.ToLower(commandName)], []interface{}{reqArr, conn})
	}
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
        conn.Write([]byte("+\"" + result.(string) + "\"\r\n"))
    }
}

func exists(reqArr []string, conn net.Conn) {
	_, ok := valueMap[reqArr[4]]
	result := "0"
	if ok {
		result = "1"
	}
	conn.Write([]byte(":" + result + "\r\n"))
}

func del(reqArr []string, conn net.Conn) {
	_, ok := valueMap[reqArr[4]]
	result := "0"
	if ok {
		result = "1"
		delete(valueMap, reqArr[4])
	}
	conn.Write([]byte(":" + result + "\r\n"))
}

func rpush(reqArr []string, conn net.Conn, paramNumber int) {
	sliceTemp, ok := valueMap[reqArr[4]]
	valueNumber := paramNumber - 2
	if ok {
		for i := 1; i <= valueNumber; i++ {
			sliceTemp = append(sliceTemp.([]interface{}), reqArr[2 * i + 4])
		}
		valueMap[reqArr[4]] = sliceTemp
		conn.Write([]byte(":" + strconv.Itoa(len(sliceTemp.([]interface{}))) + "\r\n"))
	} else {
		newSliceTemp := []interface{}{}
		for i := 1; i <= valueNumber; i++ {
			fmt.Println(reqArr[2 * i + 4])
			newSliceTemp = append(newSliceTemp, reqArr[2 * i + 4])
			fmt.Println(newSliceTemp)
		}
		valueMap[reqArr[4]] = newSliceTemp
		conn.Write([]byte(":" + strconv.Itoa(valueNumber) + "\r\n"))
	}
}

func rpop(reqArr []string, conn net.Conn) {
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		sliceLength := len(sliceTemp.([]interface{}))
		lastElement := sliceTemp.([]interface{})[sliceLength - 1].(string)
		if sliceLength == 1 {
			delete(valueMap, reqArr[4])
		} else {
			valueMap[reqArr[4]] = sliceTemp.([]interface{})[:sliceLength-1]
		}
		conn.Write([]byte("+\"" + lastElement + "\"\r\n"))
	} else {
		conn.Write([]byte("+(nil)\r\n"))
	}
}

func lpop(reqArr []string, conn net.Conn) {
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		sliceLength := len(sliceTemp.([]interface{}))
		firstElement := sliceTemp.([]interface{})[0].(string)
		if sliceLength == 1 {
			delete(valueMap, reqArr[4])
		} else {
			valueMap[reqArr[4]] = sliceTemp.([]interface{})[1:sliceLength]
		}
		conn.Write([]byte("+\"" + firstElement + "\"\r\n"))
	} else {
		conn.Write([]byte("+(nil)\r\n"))
	}
}

func llen(reqArr []string, conn net.Conn) {
	sliceTemp, ok := valueMap[reqArr[4]]
	sliceLength := 0
	if ok {
		sliceLength = len(sliceTemp.([]interface{}))
	}
	conn.Write([]byte(":" + strconv.Itoa(sliceLength) + "\r\n"))
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

