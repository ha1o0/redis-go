package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type command struct {
	ParamNumber string
	Function interface{}
}

var commandsMap = map[string]command {
	"PING": {"1", ping},
	"GET": {"2", get},
	"SET": {"3", set},
	"EXPIRE": {"3", expire},
	"SETEX": {"4", setex},
	"SETNX": {"3", setnx},
	"EXISTS": {"2", exists},
	"DEL": {"2", del},
	"RPUSH": {">=3", rpush},
	"RPOP": {"2", rpop},
	"LPOP": {"2", lpop},
	"LLEN": {"2", llen},
	"LINDEX": {"3", lindex},
	"LRANGE": {"4", lrange},
	"LTRIM": {"4", ltrim},
	"SAVE": {"1", save},
	"RESGRDB": {"1", resgrdb},
	"HSET": {"4", hset},
	"HGETALL": {"2", hgetall},
	"HGET": {"3", hget},
	"HLEN": {"2", hlen},
	"HMSET": {">=4%", hmset},//%偶数#奇数
}

//var commandsMap = map[string]string{
//	"PING": "1",
//	"GET": "2",
//	"SET": "3",
//	"EXPIRE": "3",
//	"SETEX": "4",
//	"SETNX": "3",
//	"EXISTS": "2",
//	"DEL": "2",
//	"RPUSH": ">=3",
//	"RPOP": "2",
//	"LPOP": "2",
//	"LLEN": "2",
//	"LINDEX": "3",
//	"LRANGE": "4",
//	"LTRIM": "4",
//	"SAVE": "1",
//	"RESGRDB": "1",
//	"HSET": "4",
//	"HGETALL": "2",
//	"HGET": "3",
//	"HLEN": "2",
//	"HMSET": ">=4%", //%偶数#奇数
//}

//var commandReflect = map[string]interface{}{
//    "ping": ping,
//    "get": get,
//    "set": set,
//    "expire": expire,
//    "setex": setex,
//    "setnx": setnx,
//    "exists": exists,
//    "del": del,
//    "rpush": rpush,
//    "rpop": rpop,
//    "lpop": lpop,
//    "llen": llen,
//    "lindex": lindex,
//    "lrange": lrange,
//    "ltrim": ltrim,
//    "save": save,
//    "resgrdb": resgrdb,
//    "hset": hset,
//    "hgetall": hgetall,
//    "hget": hget,
//    "hlen": hlen,
//    "hmset": hmset,
//}

var valueMap = make(map[string]interface{})
const originDumpFileName = "./dump.json"
const socketAP = "127.0.0.1:6378"
const saveInterval = 60

func main() {
	fmt.Println(time.Now(),":Server initialized")
	checkIfMap(valueMap)
	resgrdb()
	go saveCron()
	startTcpServer()
}

func startTcpServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socketAP)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.AcceptTCP()
		checkError(err)
		go handleStuff(conn)
	}
}

func checkError(err error) {
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
		fmt.Println(time.Now(),":req", n)
		if n == 0 || err != nil {
			saveGrdb()
			break
		}
		checkError(err)
		rAddr := conn.RemoteAddr()
		fmt.Println(rAddr.String())
		req := string(buf[:n])
		reqArr := strings.Split(req, "\r\n")
		reqArr = reqArr[0:len(reqArr)-1]
		fmt.Println(time.Now(),":receive the client message：", reqArr)
		handleCommands(reqArr, conn)
	}
}

func handleCommands(reqArr []string, conn net.Conn) {
	paramNumberStr := reqArr[0]
	commandName := strings.ToUpper(reqArr[2])
	commandStruct, ok := commandsMap[commandName]
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
	paramRequireNumberStr := commandStruct.ParamNumber
	if strings.Contains(paramRequireNumberStr, ">=") {
		numberStr := strings.TrimPrefix(paramRequireNumberStr, ">=")
		if strings.Contains(paramRequireNumberStr, "%") {
			numberStr = strings.Trim(numberStr, "%")
			paramRequireNumber, err := strconv.ParseInt(numberStr, 0, 64)
			if err != nil {
				/// to do
				handleCommandError(0, commandName, conn)
				return
			}
			if paramNumber < paramRequireNumber || paramNumber % 2 != 0 {
				handleCommandError(1001, commandName, conn)
				return
			}
		} else {
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
		}
		Apply(commandStruct.Function, []interface{}{reqArr, conn, int(paramNumber)})
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
		Apply(commandStruct.Function, []interface{}{reqArr, conn})
	}
}

func saveCron() {
	intervalDuration := time.Duration(saveInterval) * time.Second
	expireTimer := time.NewTimer(intervalDuration)
	select {
	case <-expireTimer.C:
		go saveGrdb()
		expireTimer.Reset(intervalDuration)
	}
}

func ping(_ []string, conn net.Conn) {
	//jsonString := "{\"a\":\"hello\",\"b\":\"123\",\"books\":[\"1\",\"a\",\"9\",\"4\"]}"
	//_, result := json2Map(jsonString)
	//valueMap = result
	response(conn, 2, "PONG")
}

func set(reqArr []string, conn net.Conn) {
	fmt.Println(reqArr[6])
    valueMap[reqArr[4]] = reqArr[6]
	response(conn, 2, "OK")
}

func get(reqArr []string, conn net.Conn) {
    result, ok := valueMap[reqArr[4]]
    if !ok {
		response(conn, 2, "(nil)")
    } else {
		if !checkIfString(result) {
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
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
		handleCommandError(1003, strings.ToUpper(reqArr[2]), conn)
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
		handleCommandError(1003, strings.ToUpper(reqArr[2]), conn)
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
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
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
		handleCommandError(1003, strings.ToUpper(reqArr[2]), conn)
		return
	}
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		fmt.Println(reflect.TypeOf(sliceTemp).String())
		// judge if target value is slice
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
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
		handleCommandError(1003, strings.ToUpper(reqArr[2]), conn)
		return
	}
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		fmt.Println(reflect.TypeOf(sliceTemp).String())
		// judge if target value is slice
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
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
		handleCommandError(1003, strings.ToUpper(reqArr[2]), conn)
		return
	}
	sliceTemp, ok := valueMap[reqArr[4]]
	if ok {
		fmt.Println(reflect.TypeOf(sliceTemp).String())
		// judge if target value is slice
		if !checkIfSlice(sliceTemp) {
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
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

func hset(reqArr []string, conn net.Conn) {
	key := reqArr[4]
	keyMapKey := reqArr[6]
	keyMapValue := reqArr[8]
	valueTemp, ok := valueMap[key]
	result := 0
	if ok {
		if !checkIfMap(valueTemp) {
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
			return
		}
		_, okk := valueTemp.(map[string]interface{})[keyMapKey]
		if !okk {
			result = 1
		}
	} else {
		valueTemp = make(map[string]interface{})
		result = 1
	}
	valueTemp.(map[string]interface{})[keyMapKey] = keyMapValue
	valueMap[key] = valueTemp
	response(conn, 0, result)
}

func hgetall(reqArr []string, conn net.Conn) {
	key := reqArr[4]
	valueTemp, ok := valueMap[key]
	if !ok {
		response(conn, 2, "(empty list or set)")
		return
	}
	if !checkIfMap(valueTemp) {
		handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
		return
	}
	result := ""
	j := 1
	for mapKey, mapKeyValue := range valueTemp.(map[string]interface{}) {
		commonResult := strconv.Itoa(j) + ") \"" + mapKey + "\"\n"
		commonResult += strconv.Itoa(j + 1) + ") \"" + mapKeyValue.(string) + "\"\n"
		result += commonResult
		j += 2
	}
	response(conn, 2, result[0: len(result) - 1])
}

func hget(reqArr []string, conn net.Conn) {
	key := reqArr[4]
	keyMapKey := reqArr[6]
	valueTemp, ok := valueMap[key]
	if !ok || !checkIfMap(valueTemp) {
		handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
		return
	}
	keyMapValue, okk := valueTemp.(map[string]interface{})[keyMapKey]
	if !okk {
		response(conn, 2, "(nil)")
		return
	}
	response(conn, 1, keyMapValue)
}

func hlen(reqArr []string, conn net.Conn) {
	key := reqArr[4]
	valueTemp, ok := valueMap[key]
	count := 0
	if !ok {
		response(conn, 0, count)
		return
	}
	if !checkIfMap(valueTemp) {
		handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
		return
	}
	for range valueTemp.(map[string]interface{}) {
		count++
	}
	response(conn, 0, count)
}

func hmset(reqArr []string, conn net.Conn, paramNumber int) {
	key := reqArr[4]
	valueTemp, ok := valueMap[key]
	result := 0
	if ok {
		if !checkIfMap(valueTemp) {
			handleCommandError(1002, strings.ToUpper(reqArr[2]), conn)
			return
		}
	} else {
		valueTemp = make(map[string]interface{})
	}
	for i := 0; i < (paramNumber - 2) / 2; i++ {
		keyMapKey := reqArr[6 + 4 * i]
		keyMapValue := reqArr[8 + 4* i]
		_, okk := valueTemp.(map[string]interface{})[keyMapKey]
		if !okk {
			result++
		}
		valueTemp.(map[string]interface{})[keyMapKey] = keyMapValue
	}
	valueMap[key] = valueTemp
	response(conn, 0, result)
}

func save(_ []string, conn net.Conn) {
	response(conn, 2, "OK")
	go saveGrdb()
}

func saveGrdb() {
	dumpJsonExist := checkFileIsExist(originDumpFileName)
	code, storeString := map2Json(valueMap)
	if code != 0 {
		fmt.Println(time.Now(),":error occurs when map to json")
		return
	}
	targetFileName := originDumpFileName
	if dumpJsonExist {
		targetFileName = originDumpFileName + ".temp"
	}
	err := ioutil.WriteFile(targetFileName, []byte(storeString), 0666) //写入文件(字节数组)
	checkError(err)
	if dumpJsonExist {
		err = os.Remove(originDumpFileName)
		checkError(err)
		err = os.Rename(targetFileName, originDumpFileName)
		checkError(err)
	}
}

func resgrdb() {
	dumpJsonExist := checkFileIsExist(originDumpFileName)
	if !dumpJsonExist {
		return
	}
	content, err := ioutil.ReadFile(originDumpFileName)
	if err != nil {
		fmt.Println(time.Now(),":ioutil ReadFile error: ", err)
		return
	}
	//fmt.Println(time.Now(),":content: ", string(content))
	code, result := json2Map(string(content))
	if code != 0 {
		fmt.Println(time.Now(),":error occurs when json to map")
		return
	}
	valueMap = result
	fmt.Println(time.Now(),":DB loaded from disk")
}

func Apply(f interface{}, args []interface{}) {
    fun := reflect.ValueOf(f)
    in := make([]reflect.Value, len(args))
    for k, param := range args{
        in[k] = reflect.ValueOf(param)
    }
    _ = fun.Call(in)
}

func response(conn net.Conn, responseType int, message interface{}) {
	switch responseType {
	// number
	case 0:
		_, _ = conn.Write([]byte(":" + strconv.Itoa(message.(int)) + "\r\n"))
	// "string"
	case 1:
		_, _ = conn.Write([]byte("+" + stringWithQuotation(message.(string)) + "\r\n"))
	// string
	case 2:
		_, _ = conn.Write([]byte("+" + message.(string) + "\r\n"))
	}
}

// handle the error of the command from client
func handleCommandError(errorCode int, commandName string, conn net.Conn) {
	switch errorCode {
	case 1000:
		_, _ = conn.Write([]byte("+(error) ERR unknown command '" + commandName + "' \r\n"))
	case 1001:
		_, _ = conn.Write([]byte("+(error) ERR wrong number of arguments for " + commandName + " command\r\n"))
	case 1002:
		_, _ = conn.Write([]byte("+(error) WRONGTYPE Operation against a key holding the wrong kind of value\r\n"))
	case 1003:
		_, _ = conn.Write([]byte("+(error) ERR value is not an integer or out of range\r\n"))

	default:
		_, _ = conn.Write([]byte("+(error) unknown error\r\n"))
	}
	//fmt.Println(time.Now(),":this connect end")
}

func checkIfMap(target interface{}) bool {
	return strings.Contains(reflect.TypeOf(target).String(), "map")
}

func checkIfString(target interface{}) bool {
	return reflect.TypeOf(target).String() == "string"
}

func checkIfSlice(target interface{}) bool {
	typeString := reflect.TypeOf(target).String()
	return strings.Contains(typeString, "[]") && !strings.Contains(typeString, "map")
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

// Store and Restore between json and map
func map2Json(mapTarget map[string]interface{}) (int, string) {
	jsonResult, err := json.Marshal(mapTarget)
	if err != nil {
		return 1, ""
	}
	fmt.Println(string(jsonResult))
	return 0, string(jsonResult)
	//for _, v := range mapTarget {
	//	vType := reflect.TypeOf(v).String()
	//
	//	if strings.Contains(vType, "map") {
	//
	//	}
	//}
	//return 0, ""
}

func json2Map(jsonString string) (int, map[string]interface{}) {
	mapResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &mapResult)
	if err != nil {
		return 1, make(map[string]interface{})
	}
	return 0, mapResult
}

/**
 * check if the dump json exists
 */
func checkFileIsExist(filename string) bool {
	exist := true
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		exist = false
	}
	return exist
}