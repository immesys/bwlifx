// This is a rudimentary test app you can use to send commands to a
// bwlifx service. Just modify the code to work for you lol
package main

import bw "gopkg.in/immesys/bw2bind.v2"

type hsbcmd struct {
	Hue        float64 `msgpack:"hue,omitempty"`
	Saturation float64 `msgpack:"saturation,omitempty"`
	Brightness float64 `msgpack:"brightness,omitempty"`
	State      bool    `msgpack:"state,omitempty"`
}

func main() {
	cl := bw.ConnectOrExit("")
	cl.SetEntityFileOrExit("thekey.key")

	cmd := hsbcmd{
		Hue:        0.7,
		Saturation: 0.5,
		Brightness: 0.7,
		State:      true,
	}

	po, _ := bw.CreateMsgPackPayloadObject(bw.PONumHSBLightMessage, &cmd)

	cl.PublishOrExit(&bw.PublishParams{
		URI:            "castle.bw2.io/michael/0/bwlifx/hsb-light.v1/slot/hsb",
		PayloadObjects: []bw.PayloadObject{po},
		AutoChain:      true,
	})

}
