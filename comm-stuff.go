package partia

import (
	"encoding/json"
	"fmt"
)

type PartiaOutput struct {
	Status string
	Data []interface{}
	Debug string
}

func Dispatch(rawJson string) PartiaOutput {
	var buf map[string]interface{}
	var commandName string

	json.Unmarshal([]byte(rawJson), &buf)
	fmt.Println(buf)
	for key := range buf {  // actually, there's only one key there
		commandName = key
	}

	switch commandName {
	case "leader":
		return Leader(rawJson)
	case "support":
		return Support(rawJson)
	case "protest":
		return Protest(rawJson)
	default:
		// TODO
		return Leader(rawJson)
	}
}
