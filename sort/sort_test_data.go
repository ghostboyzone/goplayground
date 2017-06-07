package main

import (
	"math/rand"
	"log"
	"reflect"
	"unsafe"
	"os"
	"strconv"
	"time"
	"os/exec"
)

func B2S(buf []byte) string {
    return *(*string)(unsafe.Pointer(&buf))
}
 
func S2B(s *string) []byte {
    return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(s))))
}

const MAX_NUM = 1000
const MAX_COUNT = 1200000
func main() {
	cmdObj1 := exec.Command("/bin/sh", "-c", "rm test_data.in")
	cmdObj1.Output()
	cmdObj2 := exec.Command("/bin/sh", "-c", "rm test_data.out*")
	cmdObj2.Output()

	get_num := MAX_COUNT
	var result []int
	f, err := os.OpenFile("test_data.in", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().Unix())
	for {
		tmp_num := rand.Intn(MAX_NUM)
		tmp_s := strconv.Itoa(tmp_num)
		// f.Write(S2B(&tmp_s))
		_, err := f.WriteString(tmp_s + "\n")
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, tmp_num)
		get_num--;
		if get_num <= 0 {
			break
		}
	}
	log.Println("Get ", len(result), "nums" )
}