# tplight
Control the TP-LINK LB130 IOT Lightbulb with Golang.\

## Example
```Go
package main

import (
	"github.com/cullenbass/tplight"
	"time"
	"fmt"
)


func main() {
	// create lightbulb object 
	bulb := tplight.NewBulb("192.168.1.128")

	// returns map[string]int, keys: onOff, hue, saturation, brightness
	fmt.Printf("%v\n", bulb.Info())
	time.Sleep(time.Second)

	// turn on bulb
	bulb.On()

	// set the hue, saturation, brightness
	bulb.SetHSB(0, 100, 100)

	// set the HSB, but with a fade into the new color specified in milliseconds
	bulb.SetHSBT(100, 100, 100, 5000)
	time.Sleep(6 * time.Second)

	// turn off bulb
	bulb.Off()
}
```