package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strings"
)

func deployVm2Inf(infid string) (vmid string, ok bool){
	ok = true
	clientAdd := &http.Client{}
    reqAdd, _ := http.NewRequest("POST", defconf["imhost"].(string) + "/infrastructures/" + infid, strings.NewReader(defconf["radl"].(string)))
    reqAdd.Header.Add("Charset", "utf-8")
    reqAdd.Header.Add("Content-Type", "text/plain")
    reqAdd.Header.Add("Accept", "application/json")
    reqAdd.Header.Add("Authorization", auth2Header((defconf["auth"]).([]interface{})))
    var VmAdded string
    if respAdd, err := clientAdd.Do(reqAdd); err != nil {
    	fmt.Println("Error: ", err)
    	ok = false
    } else {
    	//TODO: return error if respAdd.Status not 200
	    fmt.Println(respAdd.Status + "(" + defconf["imhost"].(string) + "/" + infid + ")")
	    decResp := json.NewDecoder(respAdd.Body)
	    
	    nvm := &InfList{}
	    decResp.Decode(nvm)
	    
	    UriVmAdded := nvm.UriList[0].Uri
	    ArrVmAdded := strings.Split(UriVmAdded, "/")
	    VmAdded = ArrVmAdded[len(ArrVmAdded)-1] 
    }
    
    fmt.Println("Added vm", VmAdded)
    return VmAdded, ok
}

func delVmFromInf(infid string, vmid string) (ok bool) {
	ok = true
	clientDel := &http.Client{}
    reqDel, _ := http.NewRequest("DELETE", defconf["imhost"].(string) + "/infrastructures/" + infid + "/vms/" + vmid, nil)
    reqDel.Header.Add("Charset", "utf-8")
    reqDel.Header.Add("Authorization", auth2Header((defconf["auth"]).([]interface{})))
    
    if respDel, err := clientDel.Do(reqDel); err != nil {
    	fmt.Println("Error: ", err)
    	ok = false
    } else {
	    //fmt.Println(respAdd.Status)
	    //TODO: return error if respDel.Status not 200
		fmt.Println("Response status:", respDel.Status)
		fmt.Println("Deleted vm", vmid)
    }
	return
}

func auth2Header(conf []interface{}) (authHeader string) {
	for key, value := range conf {
		//mapAuth := value.(map[string]interface{})
        for key, value := range value.(map[string]interface{}) {
        	valueConf := value.(string)
        	valueConf = strings.Replace(valueConf, "\n", "\\\\n", -1)
        	authHeader = authHeader + key + " = " + valueConf + "; "
        }
        if key != len(conf)-1 {
        	authHeader = strings.TrimRight(authHeader, "; ")
	        authHeader = authHeader + "\\n"
        }
    }
	authHeader = strings.TrimRight(authHeader, "; ")
	fmt.Println("Auth header:", authHeader)
	return
}

