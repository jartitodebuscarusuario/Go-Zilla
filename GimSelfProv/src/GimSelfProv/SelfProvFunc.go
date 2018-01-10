package main

import (
    "time"
)

func evaluatesp (infid string) {
  infmap.RLock()
  infmap.Data[infid].RLock()
  if (infmap.Data[infid].Conf["activesp"].(int) > 0) {	
	  //Create emptyAlarm slice
	  var emptyAlarm []bool
	  for i := 0; i < infmap.Data[infid].Conf["numsamples"].(int); i++ {
		emptyAlarm = append(emptyAlarm, false)
	  }
	  //Find param to evaluate sp (cpu or mem)
	  /*var paramsp string
	  for key, val := range evalsp {
	    if val == infmap.Data[infid].Conf["evalsp"] {
		  	  paramsp = key
	    }
	  }*/
	  paramsp := infmap.Data[infid].Conf["evalsp"].(string)
	  //Evaluate cpu slice of alarms (count true values in array of alarms)
	  //alarm := infmap.Data[infid].Alarm["up" +  paramsp]
	  wprint("Going to evaluate alarm slices, up =>", infmap.Data[infid].Alarm["up" +  paramsp], ", down =>", infmap.Data[infid].Alarm["down" +  paramsp])
	  count := 0
	  if infmap.Data[infid].Conf["nvm"].(int) < infmap.Data[infid].Conf["maxvm"].(int) {
		  for _, val := range infmap.Data[infid].Alarm["up" +  paramsp] {
		    if val == true {	
		      count++
		      wprint("Detected alarm up", paramsp, ", count", count)
		      if count >= infmap.Data[infid].Conf["numalert"].(int) && infmap.Data[infid].Conf["nvm"].(int) < infmap.Data[infid].Conf["maxvm"].(int) {
		      	infmap.Data[infid].RUnlock()
		      	infmap.Data[infid].Lock()
		      	//Empty alarm slices
		      	infmap.Data[infid].Alarm["up" +  paramsp] = emptyAlarm
		      	infmap.Data[infid].Alarm["down" +  paramsp] = emptyAlarm
		      	//Don't evaluate sp, give time to launch new machine
		      	infmap.Data[infid].Conf["activesp"] = 0
		      	time2activesp := infmap.Data[infid].Conf["tactsp"]
		      	timer := time.NewTimer(time.Second * time.Duration(time2activesp.(int)))
				if VmAdded, ok := deployVm2Inf(infid); ok {
					wprint("Succesfully added vm", VmAdded, "paramsp", paramsp)
					//infmap.Data[infid].Conf["nvm"].(int)++
					infmap.Data[infid].Conf["nvm"] = infmap.Data[infid].Conf["nvm"].(int) + 1
					//Reset monitorized values in vm's map
					for vmid, val := range infmap.Data[infid].Data {
			  	      infmap.Data[infid].Data[vmid].Lock()
			  	      for key, _ := range val.Data {
			  		    if key == paramsp {
				  		  val.Data[key] = -1
			  		    }
			  	      }
			  	      infmap.Data[infid].Data[vmid].Unlock()
			        }			
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
		      	}(time2activesp.(int), infid)
		      	wprint("Triggered up" + paramsp + " sp for infid " + infid)
		      	infmap.Data[infid].Unlock()
		      	infmap.Data[infid].RLock()
		      }
		      //break
		    }
		  }
	  } else {
	  	infmap.Data[infid].Alarm["up" +  paramsp] = emptyAlarm
	  } 	  
	  //alarm = infmap.Data[infid].Alarm["down" +  paramsp]
	  count = 0
	  if infmap.Data[infid].Conf["nvm"].(int) > infmap.Data[infid].Conf["minvm"].(int) {
		  for _, val := range infmap.Data[infid].Alarm["down" +  paramsp] {
		    if val == true {	
		      count++
		      wprint("Detected alarm down", paramsp, ", count", count)
		      if count >= infmap.Data[infid].Conf["numalert"].(int) && infmap.Data[infid].Conf["nvm"].(int) > infmap.Data[infid].Conf["minvm"].(int) {
		      	infmap.Data[infid].RUnlock()
		      	infmap.Data[infid].Lock()
		      	//Empty alarm slices
		      	infmap.Data[infid].Alarm["down" +  paramsp] = emptyAlarm
		      	infmap.Data[infid].Alarm["up" +  paramsp] = emptyAlarm
		      	//Don't evaluate sp, give time to delete machine
		      	infmap.Data[infid].Conf["activesp"] = 0
		      	time2activesp := infmap.Data[infid].Conf["tactsp"]
		      	timer := time.NewTimer(time.Second * time.Duration(time2activesp.(int)))
		      	for VmDeleted, _ := range infmap.Data[infid].Data {
		      		if ok := delVmFromInf(infid, VmDeleted); ok {
						wprint("Successfully deleted vm", VmDeleted, "paramsp", paramsp)
						//infmap.Data[infid].Conf["nvm"]--
						infmap.Data[infid].Conf["nvm"] = infmap.Data[infid].Conf["nvm"].(int) - 1
						//Delete machine from infid map
						delete(infmap.Data[infid].Data, VmDeleted)
						//Reset monitorized values in vm's map
						for vmid, val := range infmap.Data[infid].Data {
				  	      infmap.Data[infid].Data[vmid].Lock()
				  	      for key, _ := range val.Data {
				  		    if key == paramsp {
					  		  val.Data[key] = -1
				  		    }
				  	      }
				  	      infmap.Data[infid].Data[vmid].Unlock()
				        }
					} else {
						wprint("Error deleting vm", VmDeleted)
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
		      	}(time2activesp.(int), infid)
		      	wprint("Triggered down" + paramsp + " sp for infid " + infid)
		      	infmap.Data[infid].Unlock()
		      	infmap.Data[infid].RLock()
		      }
		      //break
		    }
		  }
	   } else {
	   	 infmap.Data[infid].Alarm["down" +  paramsp] = emptyAlarm
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
	  Conf: map[string]interface{}{
	    "upcpu": int(defconf["upcpu"].(float64)),
		"downcpu": int(defconf["downcpu"].(float64)),
		"upmem": int(defconf["upmem"].(float64)),
		"downmem": int(defconf["downmem"].(float64)),
		"numalert": int(defconf["numalert"].(float64)),
		"numsamples": int(defconf["numsamples"].(float64)),
		//"evalsp": evalsp[string(defconf["evalsp"].(string))],
		"evalsp": string(defconf["evalsp"].(string)),
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
  
  if (infmap.Data[infid].Conf["activesp"].(int) > 0) {
	  
	  //Find param to evaluate sp (cpu or mem)
	  var paramsp string
	  /* for key, val := range evalsp {
	  	if val == infmap.Data[infid].Conf["evalsp"] {
	  	  paramsp = key	
	  	}
	  } */
	  paramsp = infmap.Data[infid].Conf["evalsp"].(string)
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
		  if average > infmap.Data[infid].Conf["up" + paramsp].(int) {
		  	infmap.Data[infid].RUnlock()
		  	infmap.Data[infid].Lock()
		  	infmap.Data[infid].Alarm["up" + paramsp] = infmap.Data[infid].Alarm["up" + paramsp][1:]
		  	infmap.Data[infid].Alarm["up" + paramsp] = append(infmap.Data[infid].Alarm["up" + paramsp], true)
		  	infmap.Data[infid].Unlock()
		  	infmap.Data[infid].RLock()
		  } else if average < infmap.Data[infid].Conf["down" + paramsp].(int) {
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