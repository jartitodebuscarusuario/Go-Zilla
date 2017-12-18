package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "path/filepath"
    "github.com/shirou/gopsutil/mem"
    //"github.com/shirou/gopsutil/load"
    "github.com/shirou/gopsutil/cpu"
)

func main() {
	
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }
    fmt.Println(dir)
    
    v, _ := mem.VirtualMemory()

    // almost every return value is a struct
    fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

    // convert to JSON. String() is also implemented
    fmt.Println(v)
    
    //v, _ := load.Avg()
    
    // almost every return value is a struct
    //fmt.Printf("Load1: %v, Load5:%v, Load15:%v\n", v.Load1, v.Load5, v.Load15)

    // convert to JSON. String() is also implemented
    //fmt.Println(v)
    
    c, _ := cpu.Percent(15 * time.Second, false)
    
    fmt.Printf("Percent(15, false): %v\n", c)
    
    // convert to JSON. String() is also implemented
    fmt.Println(c[0])    
    
}
