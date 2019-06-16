package partia

import (
	"fmt"
	"strings"

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

func Actions(rawJson string) PartiaOutput {
	timestamp := gjson.Get(rawJson, "actions.timestamp").String()
	member := int(gjson.Get(rawJson, "actions.member").Int())
	password := gjson.Get(rawJson, "actions.password").String()
	actionType := gjson.Get(rawJson, "actions.type").String()
	project := int(gjson.Get(rawJson, "actions.project").Int())
	authority := int(gjson.Get(rawJson, "actions.authority").Int())

	if DoesMemberExist(member) && IsMemberLeader(member) {
		if !AreMemberCredsCorrect(member, password) {
			return PartiaError()
		}
	} else {
		return PartiaError()
	}
	if project != 0 && authority != 0 {
		return PartiaError()
	}

	sql := "SELECT action.id, of_type, project_id, authority FROM action JOIN project ON project_id = project.id"
	var queryArgs []interface{}

	if actionType != "" {
		sql += " WHERE of_type = $1"
		queryArgs = append(queryArgs, actionType)
	}
	if project != 0 {
		if strings.Contains(sql, "WHERE") {
			sql += " AND project_id = $2"
		} else {
			sql += " WHERE project_id = $2"
		}
		queryArgs = append(queryArgs, project)
	}
	if authority != 0 {
		if strings.Contains(sql, "WHERE") {
			sql += " AND authority = $3"
		} else {
			sql += " WHERE authority = $3"
		}
		queryArgs = append(queryArgs, authority)
	}

	rows, err := DB.Query(sql, queryArgs...)
	fmt.Println("ERR -> ", err)
	fmt.Println("ROWS -> ", rows)

	var data []interface{}

	var action_id int
	var action_type string
	var project_id int
	var authority_id int

	var upvotes int
	var downvotes int

	for rows.Next() {
		err := rows.Scan(
			&action_id, &action_type, &project_id, &authority_id)
		fmt.Println(err)
		fmt.Println(
			action_id, action_type, project_id, authority_id)

		err = DB.QueryRow(
			"SELECT count(*) FROM upvote WHERE action_id = $1", action_id).Scan(&upvotes)
		fmt.Println(err)
		fmt.Println(upvotes, "UPVOTES")
		err = DB.QueryRow(
			"SELECT count(*) FROM downvote WHERE action_id = $1", action_id).Scan(&downvotes)
		fmt.Println(err)
		fmt.Println(downvotes, "DOWNVOTES")

		rowInOutput := []interface{}{action_id, action_type, project_id, authority_id, upvotes, downvotes}
		data = append(data, rowInOutput)
	}

	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK", Data: data}
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
	fmt.Println(password, "PASSWORD")

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

	if HasUserAlreadyVotedForThisAction(member, action) || !DoesActionExist(action) {
		return PartiaError()
	}

	InsertVote(member, action, upOrDown)
	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK"}
}

func PartiaError() PartiaOutput {
	return PartiaOutput{Status: "ERROR"}
}
