// This is a rudimentary test app you can use to send commands to a
// bwlifx service. Just modify the code to work for you lol
package main

import (
	"fmt"

	bw "gopkg.in/immesys/bw2bind.v1"
)

type hsbcmd struct {
	Hue        float64 `msgpack:"hue,omitempty"`
	Saturation float64 `msgpack:"saturation,omitempty"`
	Brightness float64 `msgpack:"brightness,omitempty"`
	State      bool    `msgpack:"state,omitempty"`
}

func main() {
	cl := bw.ConnectOrExit("")
	cl.OverrideAutoChainTo(true)

	// UPDATE THIS TO WORK FOR YOU
	cl.SetEntityFileOrExit("/home/immesys/.ssh/michael.key")

	cmd := hsbcmd{
		Hue:        0.3,
		Saturation: 0.5,
		Brightness: 0.5,
		State:      true,
	}

	po, _ := bw.CreateMsgPackPayloadObject(bw.PONumHSBLightMessage, &cmd)

	err := cl.Publish(&bw.PublishParams{
		URI:            "castle.bw2.io/michael/0/bwlifx/hsb-light.v1/slot/hsb",
		PayloadObjects: []bw.PayloadObject{po},
	})
	fmt.Println("Published, err was: ", err)

}
