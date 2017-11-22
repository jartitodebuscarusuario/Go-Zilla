package main

import (
	"fmt"
	"encoding/json"
	"net"
	"os"
	"runtime"
	"reflect"
	"sync"
	"github.com/jeffail/tunny"
	"github.com/xiaonanln/keylock"
)

type data struct {
	Infid string
	Vmid string
	Cpu int
	Mem int
}

type infdata struct {
	Infid   string
	Vmid    string
	Mondata string
}

type infvm struct {
	Infid   string
}

type infalarm struct {
	Infid string
}

type infconf struct {
	Infid   string
	Confvar string
}

var noprint bool

var params []string

var mapinf sync.Map
var mapalarm sync.Map
var mapconf sync.Map

var hashinf = make(map[infdata]int)

var hashinfvm = make(map[string]map[string]bool)

var hashalarm = make(map[infalarm][]bool)

var hashconf = make(map[infconf]int)

func getField(v *data, field string) int {
    r := reflect.ValueOf(v)
    f := reflect.Indirect(r).FieldByName(field)
    return int(f.Int())
}

func evalsp(infid string) string {
	
	result := ""
	if (hashconf[infconf{infid,"nvm"}] == hashconf[infconf{infid,"infnvm"}]) {
		wprint("We have all vmids!")
		wprint(hashconf[infconf{infid,"nvm"}])
		wprint(hashconf[infconf{infid,"infnvm"}])
		
	} else {
		wprint("We have to wait a bit more ...")
		wprint(hashconf[infconf{infid,"nvm"}])
		wprint(hashconf[infconf{infid,"infnvm"}])
	}
	return result
	
}

func initinf(infid string, confvars map[string]int) {

	for confvar := range confvars {
		hashconf[infconf{infid, confvar}] = confvars[confvar]
		//wprint("key["+confvar+"] value [", confvars[confvar], "]")
	}
	mapconf.Store(infid, confvars)
	//tmpconf, _ := mapconf.Load(infid)
	//wprint(confvar, ": ", tmpconf.(map[string]int)[confvar])
	hashalarm[infalarm{infid}] = make([]bool, confvars["numsamples"])
	mapalarm.Store(infid, make([]bool, confvars["numsamples"]))
	/* tmpalarm, _ := mapalarm.Load(infid)
	wprint(infid, ": ", tmpalarm.([]bool))
	hashinfvm[infvm{infid}] = make([]string, hashconf[infconf{infid, "infnvm"}])
	wprint(hashalarm[infalarm{infid}])
	mapinf.Store(infid, hashconf)
	fmt.Println(infid + "hashconf ...")
	fmt.Println(mapinf.Load(infid))
	fmt.Println(infid + "hashconf!") */

}

//func addvmid2inf(infid string, vmid string, cpu int, mem int) {
func addvmid2inf(message *data) {

	for _, element := range params {
		hashinf[infdata{message.Infid, message.Vmid, element}] = getField(message,element)
	    wprint(hashinf[infdata{message.Infid, message.Vmid, element}])
	}
	//hashinfvm[infvm{message.Infid}] = 
	if _, ok := hashinfvm[message.Infid][message.Vmid]; ok {
		wprint("Vmid ", message.Vmid, " already exists in ", message.Infid)
	} else {
		hashinfvm[message.Infid] = make(map[string]bool)
		hashinfvm[message.Infid][message.Vmid] = true
		wprint("Added ", message.Vmid, "to inf ", message.Infid)
	}
	hashconf[infconf{message.Infid,"nvm"}]++

}

func check(e error) {

	if e != nil {

		panic(e)

	}

}

func checkError(err error) {

	if err != nil {

		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())

		os.Exit(1)

	}

}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}


func main() {

	noprint = false
	confvars := map[string]int{}
	confvars["upcpu"] = 70
	confvars["downcpu"] = 30
	confvars["upmem"] = 70
	confvars["downmem"] = 30
	confvars["numalert"] = 3
	confvars["numsamples"] = 5
	confvars["evalsp"] = 0
	confvars["activesp"] = 1
	confvars["tactsp"] = 10
	confvars["maxvm"] = 2
	confvars["minvm"] = 1
	confvars["infnvm"] = 1
	confvars["nvm"] = 0
	params = []string{"Cpu", "Mem"}

	f, err := os.Create("D:\\agracia\\Documents\\g\\perl\\GoReport_server.txt")

	check(err)

	defer f.Close()
	
	type A struct {
	    B string
	    C string
	}
	
	testArray := []A{A{}}
	testArray[0].B = "test1"
	fmt.Println(testArray[0].B)
	testMap := make(map[string]A)
	testMap["key"] = A{}
	testMap["key"] = A{B: "test2", C:"test3"}
	fmt.Println(testMap["key"].B)
	
	type vminflock struct {
		sync.RWMutex
		vmid map[string]map[string]int
	}
	
	type inflock struct {
		sync.RWMutex
		infid map[string]*vminflock
	}
	
	//var vmiddatainit = map[string]int{"cpu": 5}
	
	//var vmidinit = map[string]map[string]int{"0": {"cpu": 5}}
	
	//var lockvminf = vminflock{vmid: map[string]map[string]int{"0": {"cpu": 5}}}
	
	//var infdatainit = map[string]*vminflock{"uuid": &lockvminf}
	//var infdatainit = map[string]*vminflock{"uuid": &vminflock{vmid: map[string]map[string]int{"0": {"cpu": 5}}}}
	
	//var lockinf = inflock{infid: infdatainit}
	var infid string = "uuid"
	var vmid string = "0"
	var cpu int = 5
	var paramcpu string = "cpu"
	var lockinf = inflock{
		infid: map[string]*vminflock{
			infid: &vminflock{
				vmid: map[string]map[string]int{
					vmid: {
						paramcpu: cpu,
					},
				},
			},
		},
	}
	
	lockinf.Lock()
	lockinf.Unlock()
	
	lockinf.infid["uuid"].Lock()
	lockinf.infid["uuid"].Unlock()
	
	fmt.Println("cpu: ",lockinf.infid["uuid"].vmid["0"]["cpu"])
	lockinf.infid["uuid"].vmid["1"] = map[string]int{"mem": 50}
	fmt.Println("mem: ",lockinf.infid["uuid"].vmid["1"]["mem"])
	
	var maplock *keylock.KeyRWLock = keylock.NewKeyRWLock()
	maplock.RLock(infid)
	maplock.RUnlock(infid)
	
	var vmids = lockinf.infid["uuid"].vmid
	fmt.Println("vmids:",vmids,len(vmids))
	
	/*var maplock = struct{
		sync.RWMutex
		m map[string]innermaplock
	}{m: map[string]innermaplock{}}*/
	
	/*var inflock := &inflock{
		infid: map[string]vminflock{
				   
			   }{}
	}
	
	maplock["some_key"].lock()
	
	/*maplock.Lock()
	maplock.m["some_key"] = map[string]map[string]int{}
	maplock.m["some_key"]["some_other_key"] = map[string]int{}
	maplock.m["some_key"]["some_other_key"]["some_other_other_key"]++
	maplock.Unlock()
	
	maplock.RLock()
	n := maplock.m["some_key"]["some_other_key"]["some_other_other_key"]
	maplock.RUnlock()
	fmt.Println("some_key:", n) */

	numCPUs := runtime.NumCPU()

	runtime.GOMAXPROCS(numCPUs + 1) // numCPUs hot threads + one for async tasks.

	pool, _ := tunny.CreatePool(numCPUs, func(object interface{}) interface{} {
		conn, _ := object.(net.Conn)
		defer conn.Close()
		request := make([]byte, 1024)

		var message data

		// Do something that takes a lot of work

		for {

			read_len, err := conn.Read(request)
			if err != nil {
				fmt.Println(err)
				break
			}

			if read_len == 0 {

				break // connection already closed by client

			} else {

				json.Unmarshal(request[:read_len], &message)

				if _, ok := hashconf[infconf{message.Infid, "upcpu"}]; ok {
					//fmt.Println("Already configured inf " + message.Infid);
				} else {
					initinf(message.Infid, confvars)
					wprint("Configured inf " + message.Infid)
				}

				if (hashconf[infconf{message.Infid, "activesp"}] == 1) {

					//addvmid2inf(message.Infid, message.Vmid, message.Cpu, message.Mem)
					addvmid2inf(&message)
					evalsp(message.Infid)

					f.Write(request[:read_len])
				}
			}

		}

		return true

	}).Open()

	defer pool.Close()

	service := "127.0.0.1:8888"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)

	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)

	checkError(err)

	for {

		conn, err := listener.Accept()

		defer conn.Close()

		if err != nil {

			continue

		}

		// result, _ := pool.SendWork(conn)

		pool.SendWork(conn)

		// go handleClient(conn)

	}

}

