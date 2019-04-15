# tplight
Control Tp-LINK IOT Lightbulbs with Golang. Confirmed to work with LB130 and KL130 models.

## Example
```Go
package main

import (
	"time"
	"fmt"
	"github.com/cullenbass/tplight"
)


func main() {
	// create lightbulb object 
	bulb := tplight.NewBulb("192.168.1.128")

	// returns map[string]int, keys: onOff, hue, saturation, brightness
	info := bulb.Info()
	fmt.Printf("%v\n", info)

	// turn on bulb
	bulb.On()
	time.Sleep(time.Second)
	// set the hue, saturation, brightness
	bulb.SetHSB(0, 100, 100)
	time.Sleep(5 * time.Second)

	// set the HSB, but with a fade into the new color specified in milliseconds
	bulb.SetHSBT(100, 100, 100, 5000)
	time.Sleep(6 * time.Second)

	// turn off bulb
	bulb.Off()
}	
```
