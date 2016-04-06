package main

import (
	"fmt"

	"github.com/immesys/spawnpoint/spawnable"
	bw "gopkg.in/immesys/bw2bind.v2"
)

// The bwlifx service looks like
// baseuri/bwlifx/hsb-light.v1/slot/hsb
// This service also exposes some metadata
// as simple metadata POs. These can be configured
// in the params.yml under metadata
var btok string
var lid string

func main() {
	cl := bw.ConnectOrExit("")
	params := spawnable.GetParamsOrExit()
	eblob := params.GetEntityOrExit()
	cl.SetEntityOrExit(eblob)

	uri := params.MustString("svc_base_uri")
	btok = params.MustString("bearer_token")
	lid = params.MustString("light_id")

	command_uri := uri + "bwlifx/hsb-light.v1/slot/hsb"

	//Subscribe
	msgchan := cl.SubscribeOrExit(&bw.SubscribeParams{
		URI:       command_uri,
		AutoChain: true,
	})

	for m := range msgchan {
		dispatch(m)
	}
}

func dispatch(m *bw.SimpleMessage) {
	m.Dump()
	po := m.GetOnePODF(bw.PODFHSBLightMessage)
	if po == nil {
		return
	}
	var v map[string]interface{}
	po.(bw.MsgPackPayloadObject).ValueInto(&v)
	hue, hashue := v["hue"].(float64)
	sat, hassat := v["saturation"].(float64)
	bri, hasbri := v["brightness"].(float64)
	sta, hassta := v["state"].(bool)

	colorstr := ""
	pstr := "on"
	if hashue {
		clamp(&hue)
		colorstr += fmt.Sprintf("hue:%.3f ", hue*360)
	}
	if hassat {
		clamp(&sat)
		colorstr += fmt.Sprintf("saturation:%.3f ", sat)
	}
	if hasbri {
		clamp(&bri)
		colorstr += fmt.Sprintf("brightness:%.3f ", bri)
	}
	if hassta && !sta {
		pstr = "off"
	}

	if colorstr != "" {
		colorstr = ",\"color\":\"" + colorstr + "\""
	}
	msg := fmt.Sprintf("{\"power\":\"%s\",\"duration\":0.1%s}", pstr, colorstr)

	fmt.Println(spawnable.DoHttpPutStr(fmt.Sprintf("https://api.lifx.com/v1/lights/%s/state", lid),
		msg, []string{"Content-Type", `application/json`,
			"Authorization", "Bearer " + btok}))
}

func clamp(f *float64) {
	if *f > 1.0 {
		*f = 1.0
	}
	if *f < 0 {
		*f = 0
	}
}
