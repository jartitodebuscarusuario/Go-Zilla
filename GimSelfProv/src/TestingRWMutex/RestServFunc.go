package main

import (
    "fmt"
    "net/http"
    "encoding/json"    
)

func loadconf(conf map[string]interface{}, dconf map[string]interface{}) (result map[string]int, intresult map[string]interface{}) {
	result = make(map[string]int)
	intresult = make(map[string]interface{})
	if _, ok := conf["defaultconf"]; !ok {
		for key, _ := range dconf {
			switch param := key; param {
			  	case "evalsp":
				  	if _, ok := conf[key]; !ok {
				  		intresult[key] = string(dconf[key].(string))
				  	} else {
				  		intresult[key] = string(conf[key].(string))
				  	}
				  	//intresult[key] = result[key]
				  	var paramsp int
					for key, val := range evalsp {
					  	if key == string(result["evalsp"]) {
					  	paramsp = val
					  	}
					}
					result[key] = paramsp
			  	default:
				  	if _, ok := conf[key]; !ok {
				  		result[key] = int(dconf[key].(float64))
				  	} else {
				  		result[key] = int(conf[key].(float64))
				  	}
				  	intresult[key] = result[key]
		  	}
		}
	} else {
		for key, _ := range dconf {
			switch param := key; param {
			  	case "evalsp":
				  	var paramsp int
					for key, val := range evalsp {
					  	if key == string(dconf["evalsp"].(string)) {
					  	paramsp = val
					  	}
					}
					result[key] = paramsp
					intresult[key] = dconf[key]
			  	default:
				  	result[key] = int(dconf[key].(float64))
				  	intresult[key] = result[key]	  	
		  	}
		}
	}
	return result, intresult
}

func regInfid(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()
    var t Tcommand
    err := decoder.Decode(&t)
    if err != nil {
        w.Header().Set("Status", "400 Bad Request")
		fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
    }
    wprint("Command/Value: ", t.Command, t.Value)
    //if r.Method == "PUT" {
    switch method := r.Method; method {	
	    //if t.Command == "regInfid" {
	    case "PUT":
	        switch command := t.Command; command {
	        	case "regInfid":	
			    	infmap.RLock()
			    	if _, ok := infmap.Data[t.Value]; !ok {
				    	infmap.RUnlock()
				    	newconf, sendconf := loadconf(t.Data, defconf)
				    	var emptyAlarm []bool
						for i := 0; i < newconf["numsamples"]; i++ {
						  emptyAlarm = append(emptyAlarm, false)	
						}
						
				    	infmap.Lock()
						infmap.Data[t.Value] = &Infdata{ 
						  Data: map[string]*Vmdata{},
						  Conf: newconf,
						  Alarm: map[string][]bool {
						  	"upcpu": emptyAlarm,
						  	"downcpu": emptyAlarm,
						  	"upmem": emptyAlarm,
						  	"downmem": emptyAlarm,
						  },
						}
				    	
						infmap.Unlock()
						infmap.RLock()
						rmessage := &Tresponse{
					      Result: "OK",
					      Infid:t.Value,
					      Data:sendconf,
						}	
						w.Header().Set("Status", "201 Created")
						json.NewEncoder(w).Encode(rmessage)
						wprint("Registered inf: " + t.Value)
			    	} else {
			    		w.Header().Set("Status", "409 Conflict")
				    	fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
			    	}
			    	infmap.RUnlock()
			    default:
				    w.Header().Set("Status", "409 Conflict")
			    	fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
	        }		
	    default:
	    	w.Header().Set("Status", "405 Method Not Allowed")
	    	fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
    }	    
}