package partia

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func Leader(rawJson string) PartiaOutput {
	timestamp := gjson.Get(rawJson, "leader.timestamp").String()
	password := gjson.Get(rawJson, "leader.password").String()
	member := int(gjson.Get(rawJson, "leader.member").Int())

	if HasNumberAlreadyBeenUsed(member) {
		return PartiaError()
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

func Upvote(rawJson string) PartiaOutput {
	return vote(rawJson, "up")
}

func Downvote(rawJson string) PartiaOutput {
	return vote(rawJson, "down")
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
	timestamp := gjson.Get(rawJson, whatDoWeDo+".timestamp").String()
	member := int(gjson.Get(rawJson, whatDoWeDo+".member").Int())
	password := gjson.Get(rawJson, whatDoWeDo+".password").String()
	action := int(gjson.Get(rawJson, whatDoWeDo+".action").Int())
	project := int(gjson.Get(rawJson, whatDoWeDo+".project").Int())
	authority := int(gjson.Get(rawJson, whatDoWeDo+".authority").Int())

	if DoesMemberExist(member) {
		if !(AreMemberCredsCorrect(member, password) && IsMemberActiveEnough(member, timestamp)) {
			return PartiaError()
		}
	} else {
		if !HasNumberAlreadyBeenUsed(member) {
			CreateMember(member, password)
		} else {
			return PartiaError()
		}
	}

	if !DoesProjectExist(project) {
		if HasNumberAlreadyBeenUsed(project) {
			return PartiaError()
		}
		if (!DoesAuthorityExist(authority)) && (authority == 0 || HasNumberAlreadyBeenUsed(authority)) {
			return PartiaError()
		}
		CreateProjectAndMaybeAuthority(project, authority)
	}

	if HasNumberAlreadyBeenUsed(action) {
		return PartiaError()
	}

	CreateAction(action, member, project, whatDoWeDo)
	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK"}
}

func vote(rawJson, upOrDown string) PartiaOutput {
	timestamp := gjson.Get(rawJson, upOrDown+"vote.timestamp").String()
	member := int(gjson.Get(rawJson, upOrDown+"vote.member").Int())
	password := gjson.Get(rawJson, upOrDown+"vote.password").String()
	action := int(gjson.Get(rawJson, upOrDown+"vote.action").Int())
	fmt.Println(password, "PASWORD")

	if DoesMemberExist(member) {
		if !(AreMemberCredsCorrect(member, password) && IsMemberActiveEnough(member, timestamp)) {
			return PartiaError()
		}
	} else {
		if !HasNumberAlreadyBeenUsed(member) {
			CreateMember(member, password)
		} else {
			return PartiaError()
		}
	}

	if HasUserAlreadyVotedForThisAction(member, action) || ! DoesActionExist(action) {
		return PartiaError()
	}

	InsertVote(member, action, upOrDown)
	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK"}
}

func PartiaError() PartiaOutput {
	return PartiaOutput{Status: "ERROR"}
}
