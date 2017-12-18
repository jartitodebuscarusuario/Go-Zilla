package main

import (
    "time"
)

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
      	timer := time.NewTimer(time.Second * time.Duration(time2activesp))
		if VmAdded, ok := deployVm2Inf(infid); ok {
			wprint("Succesfully added vm", VmAdded)
		} else {
			wprint("Error adding vm, exiting")
			//os.Exit(1)
		}
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
      	for VmAdded, _ := range infmap.Data[infid].Data {
      		if ok := delVmFromInf(infid, VmAdded); ok {
				wprint("Successfully deleted vm", VmAdded)
			} else {
				wprint("Error deleting vm", VmAdded)
			}
			break
      	}
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
	for i := 0; i < int(defconf["numsamples"].(float64)); i++ {
	  emptyAlarm = append(emptyAlarm, false)	
	}
	infmap.RLock()
	if _, ok := infmap.Data[infid]; !ok { //Initialize non-existent infid with empty vmid 
	infmap.RUnlock()
	infmap.Lock()
	infmap.Data[infid] = &Infdata{ 
	  Data: map[string]*Vmdata{},
	  Conf: map[string]int{
	    "upcpu": int(defconf["upcpu"].(float64)),
		"downcpu": int(defconf["downcpu"].(float64)),
		"upmem": int(defconf["upmem"].(float64)),
		"downmem": int(defconf["downmem"].(float64)),
		"numalert": int(defconf["numalert"].(float64)),
		"numsamples": int(defconf["numsamples"].(float64)),
		"evalsp": evalsp[string(defconf["evalsp"].(string))],
		"activesp": int(defconf["activesp"].(float64)),
		"tactsp": int(defconf["tactsp"].(float64)),
		"maxvm": int(defconf["maxvm"].(float64)),
		"minvm": int(defconf["minvm"].(float64)),
	    "nvm": int(defconf["nvm"].(float64)),
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

func addData2InfidVmid (infid string, vmid string, data map[string]int) {

  //infmap.RLock()
  //infmap.Data[infid].RLock()
  //infmap.infidRLock(infid)
  //defer infmap.infidRUnlock(infid)
  //Add received data
  //infmap.Data[infid].Data[vmid].Lock()
  //infmap.Data[infid].Data[vmid].Data = data
  //infmap.Data[infid].Data[vmid].Unlock()
  infmap.AddVmidData(infid, vmid, data)
  //infmap.RLock()
  //infmap.Data[infid].RLock()
  infmap.infidRLock(infid)
  defer infmap.infidRUnlock(infid)
  
  if (infmap.Data[infid].Conf["activesp"] > 0) {
	  
	  //Find param to evaluate sp (cpu or mem)
	  var paramsp string
	  for key, val := range evalsp {
	  	if val == infmap.Data[infid].Conf["evalsp"] {
	  	  paramsp = key	
	  	}
	  }
	  //Calculate average paramsp in infid
	  sum := 0
	  count := 0
	  for vmid, val := range infmap.Data[infid].Data {
	  	infmap.Data[infid].Data[vmid].RLock()
	  	for key, val := range val.Data {
	  		if key == paramsp && val != -1 {
		  		sum =  sum + val
		  		count++
	  		}
	  	}
	  	infmap.Data[infid].Data[vmid].RUnlock()
	  }
	  //Only evaluate if have values in all machines in infid
	  wprint("Count: ", count, " nvm ", infmap.Data[infid].Conf["nvm"])
	  if count == infmap.Data[infid].Conf["nvm"] {
	  wprint("Evaluating in infid " + infid + " paramsp " + paramsp + " average ", sum/count)	
	  average := sum/count	
	  //Evaluate is received data over limits and modify Alarm slices
		  if average > infmap.Data[infid].Conf["up" + paramsp] {
		  	infmap.Data[infid].RUnlock()
		  	infmap.Data[infid].Lock()
		  	infmap.Data[infid].Alarm["up" + paramsp] = infmap.Data[infid].Alarm["up" + paramsp][1:]
		  	infmap.Data[infid].Alarm["up" + paramsp] = append(infmap.Data[infid].Alarm["up" + paramsp], true)
		  	infmap.Data[infid].Unlock()
		  	infmap.Data[infid].RLock()
		  } else if average < infmap.Data[infid].Conf["down" + paramsp] {
		  	infmap.Data[infid].RUnlock()
		  	infmap.Data[infid].Lock()
		  	infmap.Data[infid].Alarm["down" + paramsp] = infmap.Data[infid].Alarm["down" + paramsp][1:]
		  	infmap.Data[infid].Alarm["down" + paramsp] = append(infmap.Data[infid].Alarm["down" + paramsp], true)
		  	infmap.Data[infid].Unlock()
		  	infmap.Data[infid].RLock()
		  }
		  //Set values of paramsp to -1 in all machines in infid, wait for new values in all machines of infid
		  for vmid, val := range infmap.Data[infid].Data {
		  	infmap.Data[infid].Data[vmid].Lock()
		  	for key, _ := range val.Data {
		  		if key == paramsp {
			  		val.Data[key] = -1
		  		}
		  	}
		  	infmap.Data[infid].Data[vmid].Unlock()
		  }
	  }
  }
  //infmap.Data[infid].RUnlock()
  //infmap.RUnlock()
  //infmap.infidRUnlock(infid)
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
