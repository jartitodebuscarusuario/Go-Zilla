package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"net"
	"encoding/json"
	"io/ioutil"
)

//Global vars
var noprint bool
var wsocketPrint int32
var LogFile *os.File
var LogChan chan string
var errLogFile error
//var infmap Infmap
var infmap *Infmap
//var evalsp map[string]int
var defconf map[string]interface{}
var nvm int
var encoder string
//var timeCounter int32
//var timeElapsed time.Duration

//read config file in json format
func readConf(file string) map[string]interface{} {
	conf := make(map[string]interface{})
	dat, err := ioutil.ReadFile(file)
	var fconf interface{}
	err = json.Unmarshal(dat, &fconf)
		if err != nil {
		  log.Println("error:", err)
	}
	conf = fconf.(map[string]interface{})
	return conf
}

func openLogFile(logfile string) (filelog *os.File, errorlog error) {
	f, err := os.OpenFile(logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    //fmt.Printf("error opening file: %v\n", err)
	    return nil, err
	} else {
		return f, nil
	}	
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  	
  //defer conn.Close()
  // Make a buffer to hold incoming data.
  //dec := gob.NewDecoder(conn)
  //var dec interface{}
  dec := json.NewDecoder(conn)
  
  //decType := dec.(type)
	
  //for dec.(*json.Decoder).More() {
  for dec.More() {
  message := &Tmessage{}
  //dec.(*json.Decoder).Decode(message)
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

func wprint(param ...interface{}) {
	if !noprint {
		log.Println(param)
		if wsocketPrint != 0 {
		  LogChan <- fmt.Sprintf("%v", param)
		  log.Println("Send to channel:", param)
		}
	}
}

