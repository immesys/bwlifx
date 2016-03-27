package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/immesys/spawnpoint/spawnable"
	//	bw "gopkg.in/immesys/bw2bind.v1"
	bw "github.com/immesys/bw2bind"
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
	eblob := spawnable.GetEntityOrExit(params)
	us := cl.SetEntityOrExit(eblob)

	uri, ok := params["svc_base_uri"]
	if !ok {
		fmt.Println("Could not get service base uri")
		os.Exit(1)
	}

	uris, ok := uri.(string)
	if !ok {
		fmt.Println("Bad base uri")
		os.Exit(1)
	}
	command_uri := uris + "bwlifx/hsb-light.v1/slot/hsb"

	btoki, ok := params["bearer_token"]
	if !ok {
		fmt.Println("Could not get bearer token")
		os.Exit(1)
	}
	btok, ok = btoki.(string)
	if !ok {
		fmt.Println("Bearer token invalid")
		os.Exit(1)
	}

	lidi, ok := params["light_id"]
	if !ok {
		fmt.Println("Could not get light identifier")
		os.Exit(1)
	}
	lid, ok = lidi.(string)
	if !ok {
		fmt.Println("Light ID invalid")
		os.Exit(1)
	}

	//Build a chain
	pac := cl.BuildAnyChainOrExit(command_uri, "C", us)

	//Subscribe
	msgchan, err := cl.Subscribe(&bw.SubscribeParams{
		URI:                command_uri,
		PrimaryAccessChain: pac.Hash,
		ElaboratePAC:       bw.ElaborateFull,
	})
	if err != nil {
		fmt.Println("Could not subscribe: ", err)
		os.Exit(1)
	}

	for m := range msgchan {
		dispatch(m)
	}
	fmt.Println("Terminating, chan over")
}

func dispatch(m *bw.SimpleMessage) {

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

	clamp := func(f *float64) {
		if *f > 1.0 {
			*f = 1.0
		}
		if *f < 0 {
			*f = 0
		}
	}
	colorstr := ""
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
	if !hassta {
		sta = true
	}
	pstr := "on"
	if !sta {
		pstr = "off"
	}
	msg := fmt.Sprintf("{\"power\":\"%s\",\"duration\":0.1", pstr)
	if colorstr != "" {
		msg += ",\"color\":\"" + colorstr + "\""
	}
	msg += "}"
	client := &http.Client{}
	bd := bytes.NewBufferString(msg)
	target := fmt.Sprintf("https://api.lifx.com/v1/lights/%s/state", lid)
	req, err := http.NewRequest("PUT", target, bd)
	if err != nil {
		fmt.Println("Got error: ", err)
		return
	}
	req.Header.Add("Content-Type", `application/json`)
	req.Header.Add("Authorization", "Bearer "+btok)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Got err2: ", err)
		return
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(contents))

}
