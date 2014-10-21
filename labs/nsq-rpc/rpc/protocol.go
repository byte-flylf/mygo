package rpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	PROTOCOL_PRTOBUF = 1
	PROTOCOL_JSON    = 2
)

func decodeRequest(b []byte) (
	id uint32,
	topic string,
	addr string,
	msg interface{},
	err error) {
	buf := bytes.NewBuffer(b)

	if err = binary.Read(buf, binary.BigEndian, &id); err != nil {
		err = fmt.Errorf("read id, %q", err)
		return
	}

	var encoding byte
	if err = binary.Read(buf, binary.BigEndian, &encoding); err != nil {
		err = fmt.Errorf("read encoding, %q", err)
		return
	}

	var addrLen byte
	if err = binary.Read(buf, binary.BigEndian, &addrLen); err != nil {
		err = fmt.Errorf("read len addr, %q", err)
		return
	}
	addrSlice := make([]byte, int(addrLen))
	var n int
	n, err = buf.Read(addrSlice)
	if err != nil {
		err = fmt.Errorf("read addr, %q", err)
		return
	}
	if n != int(addrLen) {
		err = fmt.Errorf("read addr, except %d, get %d", addrLen, n)
		return
	}
	addr = string(addrSlice)

	var topicLen byte
	if err = binary.Read(buf, binary.BigEndian, &topicLen); err != nil {
		err = fmt.Errorf("read len topic, %q", err)
		return
	}
	topicSlice := make([]byte, int(topicLen))
	n, err = buf.Read(topicSlice)
	if err != nil {
		err = fmt.Errorf("read topic, %q", err)
		return
	}
	if n != int(topicLen) {
		err = fmt.Errorf("read topic, expect %d, get %d", topicLen, n)
		return
	}
	topic = string(topicSlice)

	bodyLen := len(b) - (4 + 1 + 1 + int(addrLen) + 1 + int(topicLen))
	bodySlice := make([]byte, bodyLen)
	n, err = buf.Read(bodySlice)
	if err != nil {
		err = fmt.Errorf("read body, %q", err)
		return
	}
	if n != int(bodyLen) {
		err = fmt.Errorf("read body, expect %d, get %d", bodyLen, n)
		return
	}
	body := string(bodySlice)

	enc := int(encoding)
	switch enc {
	case PROTOCOL_PRTOBUF:
		msg, err = parseProtobuf(body)
		break
	case PROTOCOL_JSON:
		break
	default:
		err = fmt.Errorf("invalid encoding type: %d", enc)
		return
	}

	return
}

func parseProtobuf(s string) (interface{}, error) {

	return
}
