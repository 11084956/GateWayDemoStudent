package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const Msg_Header = "12345678"

func main() {
	//类比接收缓冲区 net.Conn
	bytesBuffer := bytes.NewBuffer([]byte{})

	//发送
	if err := Encode(bytesBuffer, "hello, world 0!!!"); err != nil {
		panic(err)
	}
	if err := Encode(bytesBuffer, "hello world 1!!!"); err != nil {
		panic(err)
	}

	//读取
	for true {
		if bt, err := Decode(bytesBuffer); err == nil {
			fmt.Println(string(bt))
			continue
		}

		break
	}
}

func Encode(byteBuffer io.Writer, content string) error {
	//msg_header+content_len+content
	//8			 + 4		 + content

	clen := int32(len([]byte(content)))
	if err := binary.Write(byteBuffer, binary.BigEndian, clen); err != nil {
		return err
	}

	if err := binary.Write(byteBuffer, binary.BigEndian, []byte(content)); err != nil {
		return err
	}

	return nil
}

func Decode(byteBuffer io.Reader) (bodyBuf []byte, err error) {
	MagicBuf := make([]byte, len(Msg_Header))
	if _, err = io.ReadFull(byteBuffer, MagicBuf); err != nil {
		return nil, err
	}

	if string(MagicBuf) != Msg_Header {
		return nil, errors.New("msg_header error")
	}

	lengthBuf := make([]byte, 4)
	if _, err = io.ReadFull(byteBuffer, lengthBuf); err != nil {
		return nil, err
	}

	return bodyBuf, err
}
