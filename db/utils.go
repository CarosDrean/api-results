package db

import (
	"fmt"
	"log"
	"sync"
)

var m = sync.Mutex{}

func checkError(err error, operation string, ctx string, description string){
	if err != nil {
		m.Lock()
		log.Println(fmt.Sprintf("Error: %s en %s, decription: %s", operation, ctx, description))
		log.Println(err.Error())
		m.Unlock()
	}
}
