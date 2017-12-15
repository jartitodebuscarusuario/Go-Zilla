package main

import (
	"fmt"
	//"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"
)

type Inf struct {
	Uri string `json:"uri"`
}

type InfList struct{
	UriList []Inf `json:"uri-list"`
}

func main() {
	
    radl := `network publica (outbound = 'yes' and outports = '80,8080,3306')
			 network privada ()
			 system front
			 deploy front 1`
    
    clientAdd := &http.Client{}
    reqAdd, _ := http.NewRequest("POST", "http://104.155.107.21:8080/infrastructures/a6de1b3e-e188-11e7-ba70-42010af08f38", strings.NewReader(radl))
    reqAdd.Header.Add("Charset", "utf-8")
    reqAdd.Header.Add("Content-Type", "text/plain")
    reqAdd.Header.Add("Accept", "application/json")
    reqAdd.Header.Add("Authorization", "id = INDRAIM; type = InfrastructureManager; username = jartito; password = Papafrita9\\nid = INDRAVMRD; type = VMRC; host = http://104.155.107.21:8800/vmrc/vmrc; username = micafer; password = ttt5\\nid = GCE; type = GCE; project = gphp-preproduccion; username = icb-preproduccion@gphp-preproduccion.iam.gserviceaccount.com; password = -----BEGIN PRIVATE KEY-----\\\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCm5vdXkO3oP5cm\\\\nREM+UJVuHquta2MjCI8m//E8QFvSUpwSc43jPxIACy3C4eSYYTv8B79zT5K2x37D\\\\nAJOwuZmUsMQrs8gbbSM/WhuDqNNrnwJSQnXK1Saxf2mNcndnejS7KkR0Ikm4Rckl\\\\npTTC8QIDmVwjcCIWHKMsTU/qQLQn8mLAAKqUBPc+JJNnaYXq2r1XCsKAhebrEiyb\\\\n7t2qTPuIWscUD7T6BxtjNRA1I9QnS0HfQyuHM2Ee5M0RnzW91FXKRIUq0GhRcEe0\\\\nnSfsHE7wb6AruP03/BckRuJNOA9cQLGsbmBVA67oVJB9QJCYOw3s/0XJOvUqZyNb\\\\nY60lhEOPAgMBAAECggEADlHQ/9onPJ6+wkw8tAsoPvSdCHnRss7zBUa+n6Dqn97G\\\\n6uME1e3yhpQZnteltsMvk9M7BFLb722p/Q4UHdIBZhtNFsN2fuU8Dx0/22nFNehv\\\\nGQ1VsODdJdZyK8NvO4QQeSHKZ2e3s8/c+atfyFKdjcnUeL/lVMir4UsCzmzdvXLa\\\\nqyXh6Ij4pACFIvXLHIheBBkXFQl7ULjRVBVfKFk7u5LlNtHLBIP7LpA46kuDA3+t\\\\nwTSnl/gprBnjJ76MyHDfKONbnHh7ABEgZHQ+JRgt3o6I9KfS1gXU/6W02uybBtkp\\\\nNLCIWf/t/Ej+rOsmaqGJiPbxEbTlrhnXW4g8LB62sQKBgQDaZTEDf45NGvbZFuDG\\\\nht3OX6Y8Rr8Eg4bYUodk963F0t3XHz99pMY65afDh96OLPplg54Fa1x2wGDRDihh\\\\nqkKLah4mQZtqKJYo5f36BSSHGZFuCjZnhbq5KnmhTOeUYHoUNDc7CLyXSr5ln4JE\\\\ns4AuxCPbdcAc8wo3MxrCt3DUcQKBgQDDo/kuB9W2fAcudYWUgrDDBnUzzjKHRvZj\\\\nN5Ah/GnkXZjD1juiGz4NgEXZ0ov1K+sME8hCtdQYB66gIXc+mUXCvNpeSAgKzJj+\\\\nU4tLQLphGcymzAZiCyCQDuj+9LqnyxrRYQm+ggq07hb4Q9lqkSYo9iGz00mzmo3C\\\\n7emyzPKX/wKBgEVTwg9eOon3eUzImmnq/hY4/sg7nP+N0QxyhlBi32Lg4VMctEbq\\\\nO5MOvAax5tAzLvlyooMN5bg8sX8rg14dcipXcWKriO5WG/S3rbvkTggk8amAzGxo\\\\nYzHMbffqNclAJwCq4q12xIcyTuZrkCrG4HX4BXnxEx8dd6y2KFSPbt3BAoGAehCC\\\\n1h95TiRQbsJQl/p6wxPyaGJM0G6MKBdwzGOqxhtHx1iRWHFa5B2Wd3OQc2X1f1GQ\\\\nb173eA7C+5Ilzl7fUcN3E8AplGNXScdib49xOkhYkfFWQjHjHT7QTNLw6uQkVWMQ\\\\nK1cDyyOKHVhn/L+XaZM4L/SyVWcm7+p1F2QcMI0CgYEAwbMSqJHpQssS3OPecJlc\\\\njAOZ/svhUodjNDFDIao8gbgc9Cj3mVlcde1FR1qg6m6dBp+SL8AQwGrXhDkJ52O+\\\\n1D7VJ5w5DCpFa+OwsEHs8fLGdmUbGSn4HdCxUz+ukBLaZXfqXmdmYHKUuxQA/Igq\\\\nphA2iPF5ViMegkU3GArHT2w=\\\\n-----END PRIVATE KEY-----\\\\n")
    var VmAdded string
    if respAdd, err := clientAdd.Do(reqAdd); err != nil {
    	fmt.Println("Error: ", err)
    } else {
	    //fmt.Println(respAdd.Status)
	    decResp := json.NewDecoder(respAdd.Body)
	    
	    nvm := &InfList{}
	    decResp.Decode(nvm)
	    
	    UriVmAdded := nvm.UriList[0].Uri
	    ArrVmAdded := strings.Split(UriVmAdded, "/")
	    VmAdded = ArrVmAdded[len(ArrVmAdded)-1] 
    }
    
    fmt.Println("Added vm", VmAdded)
    
    clientDel := &http.Client{}
    reqDel, _ := http.NewRequest("DELETE", "http://104.155.107.21:8080/infrastructures/a6de1b3e-e188-11e7-ba70-42010af08f38/vms/" + VmAdded, nil)
    reqDel.Header.Add("Charset", "utf-8")
    //reqDel.Header.Add("Content-Type", "text/plain")
    //reqDel.Header.Add("Accept", "application/json")
    reqDel.Header.Add("Authorization", "id = INDRAIM; type = InfrastructureManager; username = jartito; password = Papafrita9\\nid = INDRAVMRD; type = VMRC; host = http://104.155.107.21:8800/vmrc/vmrc; username = micafer; password = ttt5\\nid = GCE; type = GCE; project = gphp-preproduccion; username = icb-preproduccion@gphp-preproduccion.iam.gserviceaccount.com; password = -----BEGIN PRIVATE KEY-----\\\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCm5vdXkO3oP5cm\\\\nREM+UJVuHquta2MjCI8m//E8QFvSUpwSc43jPxIACy3C4eSYYTv8B79zT5K2x37D\\\\nAJOwuZmUsMQrs8gbbSM/WhuDqNNrnwJSQnXK1Saxf2mNcndnejS7KkR0Ikm4Rckl\\\\npTTC8QIDmVwjcCIWHKMsTU/qQLQn8mLAAKqUBPc+JJNnaYXq2r1XCsKAhebrEiyb\\\\n7t2qTPuIWscUD7T6BxtjNRA1I9QnS0HfQyuHM2Ee5M0RnzW91FXKRIUq0GhRcEe0\\\\nnSfsHE7wb6AruP03/BckRuJNOA9cQLGsbmBVA67oVJB9QJCYOw3s/0XJOvUqZyNb\\\\nY60lhEOPAgMBAAECggEADlHQ/9onPJ6+wkw8tAsoPvSdCHnRss7zBUa+n6Dqn97G\\\\n6uME1e3yhpQZnteltsMvk9M7BFLb722p/Q4UHdIBZhtNFsN2fuU8Dx0/22nFNehv\\\\nGQ1VsODdJdZyK8NvO4QQeSHKZ2e3s8/c+atfyFKdjcnUeL/lVMir4UsCzmzdvXLa\\\\nqyXh6Ij4pACFIvXLHIheBBkXFQl7ULjRVBVfKFk7u5LlNtHLBIP7LpA46kuDA3+t\\\\nwTSnl/gprBnjJ76MyHDfKONbnHh7ABEgZHQ+JRgt3o6I9KfS1gXU/6W02uybBtkp\\\\nNLCIWf/t/Ej+rOsmaqGJiPbxEbTlrhnXW4g8LB62sQKBgQDaZTEDf45NGvbZFuDG\\\\nht3OX6Y8Rr8Eg4bYUodk963F0t3XHz99pMY65afDh96OLPplg54Fa1x2wGDRDihh\\\\nqkKLah4mQZtqKJYo5f36BSSHGZFuCjZnhbq5KnmhTOeUYHoUNDc7CLyXSr5ln4JE\\\\ns4AuxCPbdcAc8wo3MxrCt3DUcQKBgQDDo/kuB9W2fAcudYWUgrDDBnUzzjKHRvZj\\\\nN5Ah/GnkXZjD1juiGz4NgEXZ0ov1K+sME8hCtdQYB66gIXc+mUXCvNpeSAgKzJj+\\\\nU4tLQLphGcymzAZiCyCQDuj+9LqnyxrRYQm+ggq07hb4Q9lqkSYo9iGz00mzmo3C\\\\n7emyzPKX/wKBgEVTwg9eOon3eUzImmnq/hY4/sg7nP+N0QxyhlBi32Lg4VMctEbq\\\\nO5MOvAax5tAzLvlyooMN5bg8sX8rg14dcipXcWKriO5WG/S3rbvkTggk8amAzGxo\\\\nYzHMbffqNclAJwCq4q12xIcyTuZrkCrG4HX4BXnxEx8dd6y2KFSPbt3BAoGAehCC\\\\n1h95TiRQbsJQl/p6wxPyaGJM0G6MKBdwzGOqxhtHx1iRWHFa5B2Wd3OQc2X1f1GQ\\\\nb173eA7C+5Ilzl7fUcN3E8AplGNXScdib49xOkhYkfFWQjHjHT7QTNLw6uQkVWMQ\\\\nK1cDyyOKHVhn/L+XaZM4L/SyVWcm7+p1F2QcMI0CgYEAwbMSqJHpQssS3OPecJlc\\\\njAOZ/svhUodjNDFDIao8gbgc9Cj3mVlcde1FR1qg6m6dBp+SL8AQwGrXhDkJ52O+\\\\n1D7VJ5w5DCpFa+OwsEHs8fLGdmUbGSn4HdCxUz+ukBLaZXfqXmdmYHKUuxQA/Igq\\\\nphA2iPF5ViMegkU3GArHT2w=\\\\n-----END PRIVATE KEY-----\\\\n")
    
    if respDel, err := clientDel.Do(reqDel); err != nil {
    	fmt.Println("Error: ", err)
    } else {
	    //fmt.Println(respAdd.Status)
		fmt.Println("Response status:", respDel.Status)
		fmt.Println("Deleted vm", VmAdded)
    }
    
}
