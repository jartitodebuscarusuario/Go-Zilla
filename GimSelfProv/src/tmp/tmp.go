package main

import (
	"fmt"
	//"os"
    //"log"
    //"bytes"
    //"strconv"
    //"io/ioutil"
    //"net/http"
    "encoding/xml"
    //"github.com/rogpeppe/go-charset/charset"
	//_ "github.com/rogpeppe/go-charset/data"
)

func main(){

//	type IMResponse struct { 
//    	IMR []interface{} 
//    }
	
//	type IMResponse struct { 
//    	ArrayResult string `xml:"params>param>value>array>data>"`
//    	BoolResult string `xml:"value>boolean"`
//    	ArrayInf string `xml:"value>array>data"`
//    	StringInf string `xml:"value>string>"`
//    }
    	
    
    var vIMR []interface{} 
    vIMR = append(vIMR, true)
    vIMR = append(vIMR, []string{ "e4454b32-d9c2-11e7-ba70-42010af08f38", "d547b12a-d9d5-11e7-ba70-42010af08f38" })
//    vIMResponse := &IMResponse {
//    	IMR: vIMR,
//    }
//	vIMResponse := &IMResponse {
//		IMR: vIMR,
//	}
	vIMResponse := make([]interface{}, 2)
	vIMResponse[0] = true
	vIMResponse[1] = []string{ "e4454b32-d9c2-11e7-ba70-42010af08f38", "d547b12a-d9d5-11e7-ba70-42010af08f38" }
    //reply = vIMResponse
    if enc, err := xml.Marshal(vIMResponse); err != nil {
  		fmt.Printf("error: %v\n", err)
  	} else {
		fmt.Printf("no error: %s\n", enc)
  	}	
}