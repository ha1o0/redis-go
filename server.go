package main

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	"EXPIRE": "3",
	"SETEX": "4",
	"SETNX": "3",
	"EXISTS": "2",
	"DEL": "2",
	"RPUSH": ">=3",
	"RPOP": "2",
	"LPOP": "2",
	"LLEN": "2",
	"LINDEX": "3",
	"LRANGE": "4",
	"LTRIM": "4",
	"EXIT": "1",
}

var commandReflect = map[string]interface{}{
    "ping": ping,
    "get": get,
    "set": set,
    "expire": expire,
    "setex": setex,
    "setnx": setnx,
    "exists": exists,
    "del": del,
    "rpush": rpush,
    "rpop": rpop,
    "lpop": lpop,
    "llen": llen,
    "lindex": lindex,
    "lrange": lrange,
    "ltrim": ltrim,
	"exit": exit,
}

var valueMap = make(map[string]interface{})

func main() {
	fmt.Println("Start the tcp socket")
	fmt.Println(len("aaa"))
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
	paramNumber, err := strconv.ParseInt(strings.TrimPrefix(paramNumberStr, "*"), 0, 64)
	if err != nil {
		/// to do
		handleCommandError(0, commandName, conn)
		return
	}
	paramRequireNumberStr := commandsMap[commandName]
	if strings.Contains(paramRequireNumberStr, ">=") {
		paramRequireNumber, err := strconv.ParseInt(strings.TrimPrefix(paramRequireNumberStr, ">="), 0, 64)
		if err != nil {
			/// to do
			handleCommandError(0, commandName, conn)
			return
		}
		if paramNumber < paramRequireNumber {
			handleCommandError(1001, commandName, conn)
			return
		}
		Apply(commandReflect[strings.ToLower(commandName)], []interface{}{reqArr, conn, int(paramNumber)})

	} else {
		paramRequireNumber, err := strconv.ParseInt(paramRequireNumberStr, 0, 64)
		if err != nil {
			/// to do
			handleCommandError(0, commandName, conn)
			return
		}
		if paramRequireNumber != paramNumber {
			handleCommandError(1001, commandName, conn)
			return
		}
		Apply(commandReflect[strings.ToLower(commandName)], []interface{}{reqArr, conn})
	}
}

func ping(reqArr []string, conn net.Conn) {
	response(conn, 2, "PONG")
}

func exit(reqArr []string, conn net.Conn) {
    conn.Close()
}

func set(reqArr []string, conn net.Conn) {
    valueMap[reqArr[4]] = reqArr[6]
	response(conn, 2, "OK")
}

func get(reqArr []string, conn net.Conn) {
    result, ok := valueMap[reqArr[4]]
    if !ok {
		response(conn, 2, "(nil)")
    } else {
		if !checkIfString(result) {
			handleCommandError(1002, "get", conn)
			return
		}
		response(conn, 1, result.(string))
    }
}

func expire(reqArr []string, conn net.Conn) {
	keyString := reqArr[4]
	_, ok := valueMap[keyString]
	result := 0
	if !ok {
		response(conn, 0, result)
		return
	}
	expireSeconds, err := strconv.ParseInt(reqArr[6], 0, 64)
	if err != nil {
		handleCommandError(1003, "expire", conn)
		return
	}
	result = 1
	response(conn, 0, result)
	// start a new thread to exec the timer
	go setExpireTimer(keyString, int(expireSeconds))
}

func setex(reqArr []string, conn net.Conn) {
	keyString := reqArr[4]
	expireSeconds, err := strconv.ParseInt(reqArr[6], 0, 64)
	if err != nil {
		handleCommandError(1003, "expire", conn)
		return
	}
	valueMap[keyString] = reqArr[8]
	response(conn, 2, "OK")
	// start a new thread to exec the timer
	go setExpireTimer(keyString, int(expireSeconds))
}

func setnx(reqArr []string, conn net.Conn) {
	result := 0
	keyString := reqArr[4]
	_, ok := valueMap[keyString]
	if !ok {
		valueMap[keyString] = reqArr[6]
		result = 1
	}
	response(conn, 0, result)
}

func exists(reqArr []string, conn net.Conn) {
	_, ok := valueMap[reqArr[4]]
	result := 0
	if ok {
		result = 1
	}
	response(conn, 0, result)
}

func del(reqArr []string, conn net.Conn) {
	_, ok := valueMap[reqArr[4]]
	result := 0
	if ok {
		result = 1
		delete(valueMap, reqArr[4])
	}
	response(conn, 0, result)
}

func rpush(reqArr []string, conn net.Conn, paramNumber int) {
	sliceTemp, ok := valueMap[reqArr[4]]
	valueNumber := paramNumber - 2
	if ok {
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, "rpush", conn)
			return
		}
		for i := 1; i <= valueNumber; i++ {
			sliceTemp = append(sliceTemp.([]interface{}), reqArr[2 * i + 4])
		}
		valueMap[reqArr[4]] = sliceTemp
		response(conn, 0, len(sliceTemp.([]interface{})))
	} else {
		newSliceTemp := []interface{}{}
		for i := 1; i <= valueNumber; i++ {
			fmt.Println(reqArr[2 * i + 4])
			newSliceTemp = append(newSliceTemp, reqArr[2 * i + 4])
			fmt.Println(newSliceTemp)
		}
		valueMap[reqArr[4]] = newSliceTemp
		response(conn, 0, valueNumber)
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
		response(conn, 1, lastElement)
	} else {
		response(conn, 2, "(nil)")
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
		response(conn, 1, firstElement)
	} else {
		response(conn, 2, "(nil)")
	}
}

func llen(reqArr []string, conn net.Conn) {
	sliceTemp, ok := valueMap[reqArr[4]]
	sliceLength := 0
	if ok {
		sliceLength = len(sliceTemp.([]interface{}))
	}
	response(conn, 0, sliceLength)
}

func lindex(reqArr []string, conn net.Conn) {
	index, err := strconv.ParseInt(reqArr[6], 0, 64)
	if err != nil {
		handleCommandError(1003, "lindex", conn)
		return
	}
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		fmt.Println(reflect.TypeOf(sliceTemp).String())
		// judge if target value is slice
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, "lindex", conn)
			return
		}
		positionIndex := int(index)
		result := ""
		targetValue := sliceTemp.([]interface{})
		if positionIndex < 0 {
			positionIndex += len(targetValue)
		}
		if positionIndex >= 0 && positionIndex < len(targetValue) {
			result = targetValue[positionIndex].(string)
			response(conn, 1, result)
		} else {
			response(conn, 2, "(nil)")
		}
	} else {
		response(conn, 2, "(nil)")
	}
}

func lrange(reqArr []string, conn net.Conn) {
	startIndex, startErr := strconv.ParseInt(reqArr[6], 0, 64)
	endIndex, endErr := strconv.ParseInt(reqArr[8], 0, 64)
	if startErr != nil || endErr != nil {
		handleCommandError(1003, "lrange", conn)
		return
	}
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		fmt.Println(reflect.TypeOf(sliceTemp).String())
		// judge if target value is slice
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, "lrange", conn)
			return
		}
		startPositionIndex := int(startIndex)
		endPositionIndex := int(endIndex)
		result := ""
		targetSlice := sliceTemp.([]interface{})
		targetSliceLength := len(targetSlice)
		startPositionIndex, endPositionIndex = getPositiveIndex(startPositionIndex, endPositionIndex, targetSliceLength)
		if startPositionIndex > endPositionIndex || endPositionIndex < 0 {
			response(conn, 2, "(empty list or set)")
			return
		}
		j := 1
		for i := startPositionIndex; i <= endPositionIndex; i++ {
			commonResult := strconv.Itoa(j) + ") \"" + targetSlice[i].(string) + "\""
			if i == endPositionIndex {
				result += commonResult
			} else {
				result += commonResult + "\n"
			}
			j++
		}
		response(conn, 2, result)
	} else {
		response(conn, 2, "(empty list or set)")
	}
}

func ltrim(reqArr []string, conn net.Conn) {
	startIndex, startErr := strconv.ParseInt(reqArr[6], 0, 64)
	endIndex, endErr := strconv.ParseInt(reqArr[8], 0, 64)
	if startErr != nil || endErr != nil {
		handleCommandError(1003, "ltrim", conn)
		return
	}
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		fmt.Println(reflect.TypeOf(sliceTemp).String())
		// judge if target value is slice
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, "ltrim", conn)
			return
		}
		startPositionIndex := int(startIndex)
		endPositionIndex := int(endIndex)
		targetSlice := sliceTemp.([]interface{})
		targetSliceLength := len(targetSlice)
		startPositionIndex, endPositionIndex = getPositiveIndex(startPositionIndex, endPositionIndex, targetSliceLength)
		positionSub := endPositionIndex - startPositionIndex
		if positionSub < 0 || endPositionIndex < 0 {
			response(conn, 2, "(empty list or set)")
			return
		}
		valueMap[reqArr[4]] = targetSlice[startPositionIndex:endPositionIndex + 1]
		response(conn, 2, "OK")
	} else {
		response(conn, 2, "(empty list or set)")
	}
}

func setExpireTimer(keyString string, expireSeconds int) {
	expireTimer := time.NewTimer(time.Duration(int(expireSeconds)) * time.Second)
	select {
	case <-expireTimer.C:
		_, ok := valueMap[keyString]
		if !ok {
			return
		}
		delete(valueMap, keyString)
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

func response(conn net.Conn, responseType int, message interface{}) {
	switch responseType {
	// number
	case 0:
		conn.Write([]byte(":" + strconv.Itoa(message.(int)) + "\r\n"))
	// "string"
	case 1:
		conn.Write([]byte("+" + stringWithQuotation(message.(string)) + "\r\n"))
	// string
	case 2:
		conn.Write([]byte("+" + message.(string) + "\r\n"))
	}
}

// handle the error of the command from client
func handleCommandError(errorCode int, commandName string, conn net.Conn) {
	switch errorCode {
	case 1000:
		conn.Write([]byte("+(error) ERR unknown command '" + commandName + "' \r\n"))
	case 1001:
		conn.Write([]byte("+(error) ERR wrong number of arguments for " + commandName + " command\r\n"))
	case 1002:
		conn.Write([]byte("+(error) WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
	case 1003:
		conn.Write([]byte("+(error) ERR value is not an integer or out of range\r\n"))
	case 1004:

	default:
		conn.Write([]byte("+(error) unknown error\r\n"))
	}
	fmt.Println("this connect end")
}

func checkIfString(target interface{}) bool {
	return reflect.TypeOf(target).String() == "string"
}

func checkIfSlice(target interface{}) bool {
	return strings.Contains(reflect.TypeOf(target).String(), "[]")
}

func getPositiveIndex(startPositionIndex int, endPositionIndex int, targetSliceLength int) (int, int) {
	if startPositionIndex < 0 {
		startPositionIndex += targetSliceLength
		if startPositionIndex < 0 {
			startPositionIndex = 0
		}
	}
	if endPositionIndex < 0 {
		endPositionIndex += targetSliceLength
	}
	if endPositionIndex >= targetSliceLength {
		endPositionIndex = targetSliceLength - 1
	}
	return startPositionIndex, endPositionIndex
}

func stringWithQuotation(str string) string {
	return "\"" + str + "\""
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

