package partia

import (
	"github.com/tidwall/gjson"
)

func Leader(rawJson string) PartiaOutput {
	timestamp := gjson.Get(rawJson, "leader.timestamp").String()
	password := gjson.Get(rawJson, "leader.password").String()
	member := int(gjson.Get(rawJson, "leader.member").Int())

	if HasNumberAlreadyBeenUsed(member) {
		return PartiaOutput{Status: "ERROR"}
	}

	CreateMember(member, password)
	UpdateMemberLastActive(member, timestamp)
	MarkMemberAsLeader(member)

	return PartiaOutput{Status: "OK"}
}

func Support(rawJson string) PartiaOutput {
	return support_or_protest(rawJson, "support")
}

func Protest(rawJson string) PartiaOutput {
	return support_or_protest(rawJson, "protest")
}

func Upvote(rawJson string) {

}

func Downvote(rawJson string) {

}

func Actions(rawJson string) {

}

func Projects(rawJson string) {

}

func Votes(rawJson string) {

}

func Trolls(rawJson string) {

}

func support_or_protest(rawJson, whatDoWeDo string) PartiaOutput {
	timestamp := gjson.Get(rawJson, whatDoWeDo + ".timestamp").String()
	member := int(gjson.Get(rawJson, whatDoWeDo + ".member").Int())
	password := gjson.Get(rawJson, whatDoWeDo + ".password").String()
	action := int(gjson.Get(rawJson, whatDoWeDo + ".action").Int())
	project := int(gjson.Get(rawJson, whatDoWeDo + ".project").Int())
	authority := int(gjson.Get(rawJson, whatDoWeDo + ".authority").Int())

	if DoesMemberExist(member) {
		if (! (AreMemberCredsCorrect(member, password) && IsMemberActiveEnough(member))) {
			return PartiaOutput{Status: "ERROR"}
		}
	} else {
		if ! HasNumberAlreadyBeenUsed(member) {
			CreateMember(member, password)
		} else {
			return PartiaOutput{Status: "ERROR"}
		}
	}

	if ! DoesProjectExist(project) {
		if HasNumberAlreadyBeenUsed(project) {
			return PartiaOutput{Status: "ERROR"}
		}
		if (! DoesAuthorityExist(authority)) && (authority == 0 || HasNumberAlreadyBeenUsed(authority)) {
			return PartiaOutput{Status: "ERROR"}
		}
		CreateProjectAndMaybeAuthority(project, authority)
	}

	if HasNumberAlreadyBeenUsed(action) {
		return PartiaOutput{Status: "Error"}
	}

	CreateAction(action, member, project, whatDoWeDo)
	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{}
}
