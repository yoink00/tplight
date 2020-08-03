package tplight

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
)

type Bulb interface {
	SetHSB(hue int, saturation int, brightness int) error
	On() error
	Off() error
	Info() (*map[string]int, error)
}

type ldata struct {
	address string
}

func NewBulb(host string) Bulb {
	return &ldata{address: host}
}

func (b ldata) SetHSB(hue, saturation, brightness int) error {
	message := []byte("{\"smartlife.iot.smartbulb.lightingservice\":" +
		"{\"transition_light_state\":" +
		"{\"ignore_default\":1," +
		"\"on_off\":1," +
		"\"transition_period\": 0," +
		"\"hue\":" + strconv.Itoa(hue) + "," +
		"\"saturation\":" + strconv.Itoa(saturation) + "," +
		"\"brightness\":" + strconv.Itoa(brightness) + "," +
		"\"color_temp\" : 0" +
		"}}}")
	_, err := send(b.address, message)
	return err
}

func (b ldata) SetHSBT(hue, saturation, brightness, transition_period int) error {
	message := []byte("{\"smartlife.iot.smartbulb.lightingservice\":" +
		"{\"transition_light_state\":" +
		"{\"ignore_default\":1," +
		"\"on_off\":1," +
		"\"transition_period\":" + strconv.Itoa(transition_period) + "," +
		"\"hue\":" + strconv.Itoa(hue) + "," +
		"\"saturation\":" + strconv.Itoa(saturation) + "," +
		"\"brightness\":" + strconv.Itoa(brightness) + "," +
		"\"color_temp\" : 0" +
		"}}}")
	_, err := send(b.address, message)
	return err
}

func (b ldata) On() error {
	message := []byte("{\"smartlife.iot.smartbulb.lightingservice\":" +
		"{\"transition_light_state\":" +
		"{\"ignore_default\":1," +
		"\"on_off\":1" +
		"}}}")
	_, err := send(b.address, message)
	return err
}

func (b ldata) Off() error {
	message := []byte("{\"smartlife.iot.smartbulb.lightingservice\":" +
		"{\"transition_light_state\":" +
		"{\"ignore_default\":1," +
		"\"on_off\":0," +
		"\"transition_period\":2" +
		"}}}")
	_, err := send(b.address, message)
	return err
}

func (b ldata) Info() (*map[string]int, error) {
	info := make(map[string]int)
	parsed, err := send(b.address, []byte("{\"system\" : {\"get_sysinfo\": {}}}")[:])
	if err != nil {
		return nil, err
	}
	data := parsed["system"].(map[string]interface{})["get_sysinfo"].(map[string]interface{})["light_state"].(map[string]interface{})
	info["onOff"] = int(data["on_off"].(float64))

	if info["onOff"] != 1 {
		info["hue"] = int(data["dft_on_state"].(map[string]interface{})["hue"].(float64))
		info["saturation"] = int(data["dft_on_state"].(map[string]interface{})["saturation"].(float64))
		info["brightness"] = int(data["dft_on_state"].(map[string]interface{})["brightness"].(float64))
	} else {
		info["hue"] = int(data["hue"].(float64))
		info["saturation"] = int(data["saturation"].(float64))
		info["brightness"] = int(data["brightness"].(float64))
	}

	return &info, nil

}

func encrypt(data []byte) (output []byte) {
	key := byte(0xAB)
	for i := 0; i < len(data); i++ {
		c := data[i]
		output = append(output, c^key)
		key = output[i]
	}
	return output
}

func decrypt(data []byte) (output []byte) {
	key := byte(0xAB)
	for i := 0; i < len(data); i++ {
		c := data[i]
		output = append(output, c^key)
		key = c
	}
	return output
}

func send(hostname string, message []byte) (map[string]interface{}, error) {
	data := encrypt(message)
	port := 9999
	conn, err := net.Dial("udp4", hostname+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	rData := make([]byte, 1500)
	rLen, err := bufio.NewReader(conn).Read(rData)
	if err != nil {
		return nil, err
	}
	dData := decrypt(rData[:rLen])
	var parsed map[string]interface{}
	err = json.Unmarshal(dData, &parsed)
	if err != nil {
		panic(err)
	}
	return parsed, nil
}
