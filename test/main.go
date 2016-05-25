// This is a rudimentary test app you can use to send commands to a
// bwlifx service. Just modify the code to work for you lol
package main

import (
	"fmt"

	bw "github.com/immesys/bw2bind"
)

type hsbcmd struct {
	Hue        float64 `msgpack:"hue,omitempty"`
	Saturation float64 `msgpack:"saturation,omitempty"`
	Brightness float64 `msgpack:"brightness,omitempty"`
	State      bool    `msgpack:"state,omitempty"`
}

const BaseURI = "410.dev/lighting/s.lifx/0"

func main() {
	cl := bw.ConnectOrExit("")
	cl.SetEntityFromEnvironOrExit()

	cmd := hsbcmd{
		Hue:        0.7,
		Saturation: 0.5,
		Brightness: 0.7,
		State:      true,
	}

	po, _ := bw.CreateMsgPackPayloadObject(bw.PONumHSBLightMessage, &cmd)

	cl.PublishOrExit(&bw.PublishParams{
		URI:            BaseURI + "/i.hsblight/slot/hsb",
		PayloadObjects: []bw.PayloadObject{po},
		AutoChain:      true,
	})

	fmt.Println("Published")

}
