package partia

type FunctionName int

const (
	open FunctionName = iota
	leader
	support
	protest
	upvote
	downvote
	actions
	projects
	votes
	trolls
)

type FunctionCall struct {
	functionName FunctionName
	timestamp int
	member int
	password string
	action int
	project int
	authority int
	typeAsInSupportOrProtest string
}

func ParseFromJsonMap(rawMap map[string]interface{}) FunctionCall {

}
