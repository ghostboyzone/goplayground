package main

import (
	"encoding/binary"
	// "encoding/gob"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

type Message struct {
	MsgType int    `json:"msg_type"`
	Data    string `json:"data"`
	Path    string `json:"path"`
}

func main() {
	log.Println("Test")

	ff_bytes, err := ioutil.ReadFile("test.go")

	// []byte

	// var tmp string

	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.BigEndian, ff_bytes)

	msg := Message{
		MsgType: 1,
		Data:    buf.String(),
		Path:    "xxxxx",
	}

	// enc := gob.NewEncoder(&tmp)
	// err = enc.Encode(ff_bytes)
	aa, err := json.Marshal(msg)
	log.Println("11", err, string(aa))
	var bb Message
	json.Unmarshal(aa, &bb)
	log.Println(bb.MsgType)
	return

	var s string
	err = json.Unmarshal(ff_bytes, &s)
	log.Println(s, err)
}
