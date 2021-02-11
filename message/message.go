package message

import (
	"encoding/gob"
	"fmt"
	"kitty_mon_client/reading"
	"kitty_mon_client/util"
	"net"
	"time"
)

type Message struct {
	Id     uint64
	Type   string
	Param  string
	Param2 string
	Token  string
	//	reqSeq int64
	//	respSeq int64
	Reading   reading.Reading // (think NoteChange) Will include a ref to it's device
	CreatedAt time.Time
}

func SendMsg(encoder *gob.Encoder, msg Message) {
	encoder.Encode(msg)
	PrintMsg(msg, false)
	//time.Sleep(10)
}

func RcxMsg(decoder *gob.Decoder, msg *Message) {
	//time.Sleep(10)
	decoder.Decode(&msg)
	PrintMsg(*msg, true)
}

func PrintHangupMsg(conn net.Conn) {
	fmt.Printf("Closing connection: %+v\n----------------------------------------------\n", conn)
}

func PrintMsg(msg Message, rcx bool) {
	util.Pl("\n----------------------------------------------")
	if rcx {
		print("Received: ")
	} else {
		print("Sent: ")
	}
	util.Pl("Msg Type:", msg.Type, " Msg Param:", util.Short_sha(msg.Param))
	msg.Reading.Print()
}
