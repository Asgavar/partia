package partia

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PartiaOutput struct {
	Status string
	Data   []interface{}
	Debug  string
}

func Dispatch(rawJson string) PartiaOutput {
	var buf map[string]interface{}
	var commandName string

	json.Unmarshal([]byte(rawJson), &buf)
	fmt.Println(buf)
	for key := range buf { // actually, there's only one key there
		commandName = key
	}

	switch commandName {
	case "leader":
		return Leader(rawJson)
	case "support":
		return Support(rawJson)
	case "protest":
		return Protest(rawJson)
	case "upvote":
		return Upvote(rawJson)
	case "downvote":
		return Downvote(rawJson)
	case "actions":
		return Actions(rawJson)
	case "projects":
		return Projects(rawJson)
	case "votes":
		return Votes(rawJson)
	default:
		return PartiaError()
	}
}

func RenderOutputJson(output PartiaOutput) string {
	jsonedAsBytes, _ := json.Marshal(output)
	jsoned := string(jsonedAsBytes)

	jsoned = strings.ReplaceAll(jsoned, ",\"Data\":null", "")
	jsoned = strings.ReplaceAll(jsoned, ",\"Debug\":\"\"", "")

	jsoned = strings.ReplaceAll(jsoned, "Status", "status")
	jsoned = strings.ReplaceAll(jsoned, "Data", "data")
	jsoned = strings.ReplaceAll(jsoned, "Debug", "debug")

	return jsoned
}
