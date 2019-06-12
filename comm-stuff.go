package partia

import (
	"encoding/json"
	"fmt"
)

type PartiaOutput struct {
	status string
	data []interface{}
	debug string
}

func Dispatch(rawJson string){
	var buf map[string]interface{}
	var commandName string

	json.Unmarshal([]byte(rawJson), &buf)
	fmt.Println(buf)
	for key := range buf {  // actually, there's only one key there
		commandName = key
	}

	switch commandName {
	case "leader":
		Leader(rawJson)
	default:
		// TODO
	}
}
