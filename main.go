package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup
const BytesToRead = 30000

func readStuff(bc chan int64, canRead chan bool, file *os.File, final *string) {
	buf := make([]byte, BytesToRead)
	<-canRead
	startByte := <- bc
	b, err := file.ReadAt(buf, startByte)
	if err != nil {
		log.Fatal("ReadAt failed: ", err)
	}

	*final += string(buf)
	bc <- startByte + int64(b)
	wg.Done()
	canRead <- true
}

func main() {
	var final string
	file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal("Failed to open:", err)
	}
	defer file.Close()

	startByte := make(chan int64, 1)
	canRead := make(chan bool, 1)
	startByte <- 0
	canRead <- true

	wg.Add(2)
	go readStuff(startByte, canRead, file, &final)
	go readStuff(startByte, canRead, file, &final)

	wg.Wait()
	fmt.Println(len(final))
}
