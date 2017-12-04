package main

import (
	"flag"
    "fmt"
    "net"
    "os"
    "sync"
    //"sync/atomic"
    "time"
    //"strconv"
    //"math/rand"
    //"strconv"
    "encoding/json"
    //"encoding/gob"
    //"bytes"
    //"github.com/satori/go.uuid"
)

//Type to store monitored data from remote instances
type Tmessage struct {
	Infid string
	Vmid string
	Data map[string]int
}

//Type to store monitored data from remote instances in global memory
type Vmdata struct {
	sync.RWMutex
	Data map[string]int
}

//Type to store remote instances in global memmory
type Infdata struct {
	sync.RWMutex
	Data map[string]*Vmdata
	Conf map[string]int
	Alarm map[string][]bool
}

//Type to store remote infraestructures
type Infmap struct {
	sync.RWMutex
	Data map[string]*Infdata
}

//Global vars
var noprint bool
//var infmap Infmap
var infmap *Infmap
var evalsp map[string]int
var nvm int
var encoder string
//var timeCounter int32
//var timeElapsed time.Duration

func main() {

    var CONN_HOST string
    var CONN_TYPE string
    var CONN_PORT string
    
    flag.StringVar(&CONN_HOST, "host", "::1", "Hostname to bind (default: agracia7)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Hostname to bind (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp6", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", true, "Supress output (default: false)")
    flag.IntVar(&nvm, "nvm", 1, "Machines per inf (default: 1)")
    flag.StringVar(&encoder, "encoder", "json", "Message encoding type (default: json)")
    
    flag.Parse()
    
    // Open file to write received data
	//f, err = os.Create("D:\\agracia\\Documents\\g\\perl\\GoReport_server.txt")
	
	//Initialize infmap
	infmap = &Infmap{ Data: map[string]*Infdata{} }
	//infmap.Data = map[string]*Infdata{}
	
	//Initialize evasp (how we evaluate sp)
	evalsp = map[string]int{
		"cpu": 0,
		"mem": 1,
	}
	
	//timeCounter = 0;
	
    // Listen for incoming connections.
    tcpAddr, err := net.ResolveTCPAddr(CONN_TYPE, "[" + CONN_HOST + "]:" + CONN_PORT)
    //l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    l, err := net.ListenTCP(CONN_TYPE, tcpAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        //atomic.AddInt32(&timeCounter,1)
        defer conn.Close()
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  	
  //defer conn.Close()
  // Make a buffer to hold incoming data.
  //dec := gob.NewDecoder(conn)
  dec := json.NewDecoder(conn)
	
  for dec.More() {
  message := &Tmessage{}
  dec.Decode(message)
  // Write received data to disk
  //f.Write(request[:read_len])
  
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received " + message.Infid + " " + message.Vmid))

  //TEST
  start := time.Now()
  
  checkInfid(message.Infid)
  checkVmid(message.Infid, message.Vmid)
  addData2InfidVmid(message.Infid, message.Vmid, message.Data)
  evaluatesp(message.Infid)
  wprint("Added data", message.Data, "to vmid", message.Vmid, "in infid", message.Infid)
  
  wprint("timeCounter: ", time.Now().Sub(start))
  //TEST  
  
  }
  // Close the connection when you're done with it.
  conn.Close()
  
  //go func() {
//  start := time.Now()
//  
//  checkInfid(message.Infid)
//  checkVmid(message.Infid, message.Vmid)
//  addData2InfidVmid(message.Infid, message.Vmid, message.Data)
//  evaluatesp(message.Infid)
//  wprint("Added data", message.Data, "to vmid", message.Vmid, "in infid", message.Infid)
//  
//  wprint("timeCounter: ", time.Now().Sub(start))
  
//  infmap.RLock()
//  defer infmap.RUnlock()
//  for key, value := range infmap.Data {
//  	fmt.Println("key: ",key," =>")
//    for key, value := range value.Data {
//      fmt.Println("key: ",key," =>")
//      for key, value := range value.Data {
//	    fmt.Println("Key:", key, "Value:", value)
//      }
//    }
//  }
//  fmt.Println("====================================")
  
}

//func readInf (infid string) (infdata map[string]*Vmdata) {
//	infmap.RLock()
//	infmap.Data[infid].RLock()
//	infdata = infmap.Data[infid].Data
//	infmap.Data[infid].RUnlock()
//	infmap.RUnlock()
//	return infdata
//}

func addData2InfidVmid (infid string, vmid string, data map[string]int) {

  //infmap.RLock()
  //infmap.Data[infid].RLock()
  //infmap.infidRLock(infid)
  //defer infmap.infidRUnlock(infid)
  //Add received data
  //infmap.Data[infid].Data[vmid].Lock()
  //infmap.Data[infid].Data[vmid].Data = data
  //infmap.Data[infid].Data[vmid].Unlock()
  //infmap.AddVmidData(infid, vmid, data)
  //infmap.RLock()
  //infmap.Data[infid].RLock()
  //infmap.infidRLock(infid)
  //defer infmap.infidRUnlock(infid)
  
  infmap.RLock()
  infmap.Data[infid].RLock()
  infmap.Data[infid].Data[vmid].Lock()
  infmap.Data[infid].Data[vmid].Data = data
  infmap.Data[infid].Data[vmid].Unlock()
  infmap.Data[infid].RUnlock()
  infmap.RUnlock()
  
  infmap.RLock()
  infdata := infmap.Data[infid]
  infmap.RUnlock()
  
  infdata.RLock()
  infconf := infdata.Conf
  vmidata := infdata.Data
  infdata.RUnlock()
  
  if (infconf["activesp"] > 0) {
	  
	  //Find param to evaluate sp (cpu or mem)
	  var paramsp string
	  for key, val := range evalsp {
	  	if val == infconf["evalsp"] {
	  	  paramsp = key	
	  	}
	  }
	  //Calculate average paramsp in infid
	  sum := 0
	  count := 0
	  tmpvmidata := vmidata
	  for vmid, val := range tmpvmidata {
	  	tmpvmidata[vmid].RLock()
	  	for key, val := range val.Data {
	  		if key == paramsp && val != -1 {
		  		sum =  sum + val
		  		count++
	  		}
	  	}
	  	tmpvmidata[vmid].RUnlock()
	  }
	  //Only evaluate if have values in all machines in infid
	  wprint("Count: ", count, " nvm ", infconf["nvm"])
	  if count == infconf["nvm"] {
	  wprint("Evaluating in infid " + infid + " paramsp " + paramsp + " average ", sum/count)	
	  average := sum/count	
	  //Evaluate is received data over limits and modify Alarm slices
		  if average > infconf["up" + paramsp] {
		  	infdata.Alarm["up" + paramsp] = infdata.Alarm["up" + paramsp][1:]
		  	infdata.Alarm["up" + paramsp] = append(infdata.Alarm["up" + paramsp], true)
		  } else if average < infdata.Conf["down" + paramsp] {
		  	infdata.Alarm["down" + paramsp] = infdata.Alarm["down" + paramsp][1:]
		  	infdata.Alarm["down" + paramsp] = append(infdata.Alarm["down" + paramsp], true)
		  }
		  //Set values of paramsp to -1 in all machines in infid, wait for new values in all machines of infid
		  for vmid, val := range tmpvmidata {
		  	tmpvmidata[vmid].Lock()
		  	for key, _ := range val.Data {
		  		if key == paramsp {
			  		val.Data[key] = -1
		  		}
		  	}
		  	tmpvmidata[vmid].Unlock()
		  }
		  infmap.RLock()
		  infmap.Data[infid].Lock()
		  infmap.Data[infid].Alarm = infdata.Alarm
		  infmap.Data[infid].Conf = infdata.Conf
		  infmap.Data[infid].Data = tmpvmidata
		  infmap.Data[infid].Unlock()
		  infmap.RUnlock()
	  }
  }
  //infmap.Data[infid].RUnlock()
  //infmap.RUnlock()
  //infmap.infidRUnlock(infid)
}

func (mapinf *Infmap)infidRLock(idinf string) {
	mapinf.RLock()
	mapinf.Data[idinf].RLock()
}

func (mapinf *Infmap)infidRUnlock(idinf string) {
	mapinf.Data[idinf].RUnlock()
	mapinf.RUnlock()
}

func (mapinf *Infmap)AddVmidData(idinf string, idvm string, vmdata map[string]int) {
	mapinf.RLock()
	defer mapinf.RUnlock()
	mapinf.Data[idinf].RLock()
	defer mapinf.Data[idinf].RUnlock()
	mapinf.Data[idinf].Data[idvm].Lock()
	defer mapinf.Data[idinf].Data[idvm].Unlock()
	mapinf.Data[idinf].Data[idvm].Data = vmdata	
}

func evaluatesp (infid string) {
  infmap.RLock()
  infmap.Data[infid].RLock()	
  //Create emptyAlarm slice
  var emptyAlarm []bool
  for i := 0; i < 5; i++ {
	emptyAlarm = append(emptyAlarm, false)
  }
  //Find param to evaluate sp (cpu or mem)
  var paramsp string
  for key, val := range evalsp {
    if val == infmap.Data[infid].Conf["evalsp"] {
	  	  paramsp = key
    }
  }
  //Evaluate cpu slice of alarms (count true values in array of alarms)
  alarm := infmap.Data[infid].Alarm["up" +  paramsp]
  count := 0
  for _, val := range alarm {
    if val == true {
      count++
      if count >= infmap.Data[infid].Conf["numalert"] && infmap.Data[infid].Conf["nvm"] < infmap.Data[infid].Conf["maxvm"] {
      	infmap.Data[infid].RUnlock()
      	infmap.Data[infid].Lock()
      	infmap.Data[infid].Alarm["up" +  paramsp] = emptyAlarm
      	infmap.Data[infid].Conf["activesp"] = 0
      	time2activesp := infmap.Data[infid].Conf["tactsp"]
      	//timer := time.NewTimer(time.Second * time.Duration(time2activesp))
      	go func (time int, infid string, timer *time.Timer) {
      		<-timer.C
	        // If main() finishes before the 60 second timer, we won't get here
	        infmap.RLock()
	        infmap.Data[infid].Lock()
	        infmap.Data[infid].Conf["activesp"] = 1
	        infmap.Data[infid].Unlock()
	        infmap.RUnlock()
	        wprint("Congratulations! Your ", time2activesp, " second timer for infid " + infid + " finished.")
      	}(time2activesp, infid, time.NewTimer(time.Second * time.Duration(time2activesp)))
      	wprint("Triggered up" + paramsp + " sp for infid " + infid)
      	infmap.Data[infid].Unlock()
      	infmap.Data[infid].RLock()
      }
    }
  }
  alarm = infmap.Data[infid].Alarm["down" +  paramsp]
  count = 0   
  for _, val := range alarm {
    if val == true {
      count++
      if count >= infmap.Data[infid].Conf["numalert"] && infmap.Data[infid].Conf["nvm"] > infmap.Data[infid].Conf["minvm"] {
      	infmap.Data[infid].RUnlock()
      	infmap.Data[infid].Lock()
      	infmap.Data[infid].Alarm["down" +  paramsp] = emptyAlarm
      	infmap.Data[infid].Conf["activesp"] = 0
      	time2activesp := infmap.Data[infid].Conf["tactsp"]
      	timer := time.NewTimer(time.Second * time.Duration(time2activesp))
      	go func (time int, infid string) {
      		<-timer.C
	        // If main() finishes before the 60 second timer, we won't get here
	        infmap.RLock()
	        infmap.Data[infid].Lock()
	        infmap.Data[infid].Conf["activesp"] = 1
	        infmap.Data[infid].Unlock()
	        infmap.RUnlock()
	        wprint("Congratulations! Your ", time2activesp, " second timer for infid " + infid + " finished.")
      	}(time2activesp, infid)
      	wprint("Triggered down" + paramsp + " sp for infid " + infid)
      	infmap.Data[infid].Unlock()
      	infmap.Data[infid].RLock()
      }
    }
  }
  infmap.Data[infid].RUnlock()
  infmap.RUnlock()
}

func checkInfid(infid string) {
	var emptyAlarm []bool
	for i := 0; i < 5; i++ {
	  emptyAlarm = append(emptyAlarm, false)	
	}
	infmap.RLock()
	if _, ok := infmap.Data[infid]; !ok { //Initialize non-existent infid with empty vmid 
	infmap.RUnlock()
	infmap.Lock()
	infmap.Data[infid] = &Infdata{ 
	  Data: map[string]*Vmdata{},
	  Conf: map[string]int{
	    "upcpu": 70,
		"downcpu": 30,
		"upmem": 80,
		"downmem": 20,
		"numalert": 3,
		"numsamples": 5,
		"evalsp": evalsp["cpu"],
		"activesp": 1,
		"tactsp": 10,
		"maxvm": 4,
		"minvm": 1,
	    "nvm": nvm,
	  },
	  Alarm: map[string][]bool {
	  	"upcpu": emptyAlarm,
	  	"downcpu": emptyAlarm,
	  	"upmem": emptyAlarm,
	  	"downmem": emptyAlarm,
	  },
	}
	infmap.Unlock()
	infmap.RLock()
  }
  infmap.RUnlock()
}

func checkVmid (infid string, vmid string) {
	infmap.RLock()
	infmap.Data[infid].RLock()
	if _, ok := infmap.Data[infid].Data[vmid]; !ok { //Initialize non-existent vmid in infid with empty map[string]int
	  infmap.Data[infid].RUnlock()
	  infmap.Data[infid].Lock()
	  infmap.Data[infid].Data[vmid] = &Vmdata{
	    Data: map[string]int{},
	  }
	  infmap.Data[infid].Unlock()
	  infmap.Data[infid].RLock()
    }
	infmap.Data[infid].RUnlock()
	infmap.RUnlock()
}

//func (infmap Infmap) vmidRLock(infid string, vmid string) {
//	infmap.RLock()
//	infmap.Data[infid].RLock()
//	infmap.Data[infid].Data[vmid].RLock()
//}
//
//func (infmap Infmap) vmidRUnlock(infid string, vmid string) {
//	infmap.Data[infid].Data[vmid].RUnlock()
//	infmap.Data[infid].RUnlock()
//	infmap.RUnlock()
//}
//
//func (infmap Infmap) vmidLock(infid string, vmid string) {
//	infmap.RLock()
//	infmap.Data[infid].RLock()
//	infmap.Data[infid].Data[vmid].Lock()
//}
//
//func (infmap Infmap) vmidUnlock(infid string, vmid string) {
//	infmap.Data[infid].Data[vmid].Unlock()
//	infmap.Data[infid].RUnlock()
//	infmap.RUnlock()
//}
//
//func (infmap Infmap) infidRLock (infid string) {
//	infmap.RLock()
//	infmap.Data[infid].RLock()
//}
//
//func (infmap Infmap) infidRUnlock(infid string) {
//	infmap.Data[infid].RUnlock()
//	infmap.RUnlock()
//}
//
//func (infmap Infmap) infidLock(infid string) {
//	infmap.RLock()
//	infmap.Data[infid].Lock()
//}
//
//func (infmap Infmap) infidUnlock(infid string) {
//	infmap.Data[infid].Unlock()
//	infmap.RUnlock()
//}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}