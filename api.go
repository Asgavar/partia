package partia

import (
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

	sql += " ORDER BY action ASC"

	rows, _ := DB.Query(sql, queryArgs...)

	var data []interface{}

	var action_id int
	var action_type string
	var project_id int
	var authority_id int

	var upvotes int
	var downvotes int

	for rows.Next() {
		rows.Scan(
			&action_id, &action_type, &project_id, &authority_id)

		DB.QueryRow(
			"SELECT count(*) FROM upvote WHERE action_id = $1", action_id).Scan(&upvotes)
		DB.QueryRow(
			"SELECT count(*) FROM downvote WHERE action_id = $1", action_id).Scan(&downvotes)

		rowInOutput := []interface{}{action_id, action_type, project_id, authority_id, upvotes, downvotes}
		data = append(data, rowInOutput)
	}

	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK", Data: data}
}

func Projects(rawJson string) PartiaOutput {
	timestamp := gjson.Get(rawJson, "projects.timestamp").String()
	member := int(gjson.Get(rawJson, "projects.member").Int())
	password := gjson.Get(rawJson, "projects.password").String()
	authority := int(gjson.Get(rawJson, "projects.authority").Int())

	if DoesMemberExist(member) && IsMemberLeader(member) {
		if !AreMemberCredsCorrect(member, password) {
			return PartiaError()
		}
	} else {
		return PartiaError()
	}

	var queryArgs []interface{}
	sql := "SELECT * FROM project"

	if authority != 0 {
		sql += " WHERE authority = $1"
		queryArgs = append(queryArgs, authority)
	}

	sql += " ORDER BY project ASC"

	var (
		project_id int
		authority_id int
		data []interface{}
	)

	rows, _ := DB.Query(sql, queryArgs...)

	for rows.Next() {
		rows.Scan(&project_id, &authority_id)
		rowInOutput := []interface{}{project_id, authority_id}
		data = append(data, rowInOutput)
	}

	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK", Data: data}
}

func Votes(rawJson string) PartiaOutput {
	timestamp := gjson.Get(rawJson, "votes.timestamp").String()
	member := int(gjson.Get(rawJson, "votes.member").Int())
	password := gjson.Get(rawJson, "votes.password").String()
	action := int(gjson.Get(rawJson, "votes.action").Int())
	project := int(gjson.Get(rawJson, "votes.project").Int())

	if DoesMemberExist(member) && IsMemberLeader(member) {
		if !AreMemberCredsCorrect(member, password) {
			return PartiaError()
		}
	} else {
		return PartiaError()
	}

	if action != 0 && project != 0 {
		return PartiaError()
	}

	var action_ids []interface{}

	if action != 0 || project != 0 {
		if action != 0 {
			action_ids = []interface{}{action}
		} else {
			actions_from_db, _ := DB.Query(
				"SELECT id FROM action WHERE project_id = $1", project)
			for actions_from_db.Next() {
				var new_action_id int
				actions_from_db.Scan(&new_action_id)
				action_ids = append(action_ids, new_action_id)
			}
		}
	}

	var data []interface{}

	all_users_from_db, _ := DB.Query("SELECT id FROM member ORDER BY id ASC")
	for all_users_from_db.Next() {
		var new_user_id int
		var tmp_upvotes int
		var tmp_downvotes int
		var upvotes_total int
		var downvotes_total int

		all_users_from_db.Scan(&new_user_id)

		if len(action_ids) == 0 {  // i.e. no filtering required
			DB.QueryRow(
				"SELECT count(*) FROM upvote WHERE member_id = $1", new_user_id).Scan(&upvotes_total)
			DB.QueryRow(
				"SELECT count(*) FROM downvote WHERE member_id = $1", new_user_id).Scan(&downvotes_total)
		}

		for _, action_id := range action_ids {
			DB.QueryRow(
				"SELECT count(*) FROM upvote WHERE member_id = $1 AND action_id = $2",
				new_user_id, action_id).Scan(&tmp_upvotes)
			DB.QueryRow(
				"SELECT count(*) FROM downvote WHERE member_id = $1 AND action_id = $2",
				new_user_id, action_id).Scan(&tmp_downvotes)

			upvotes_total += tmp_upvotes
			downvotes_total += tmp_downvotes
		}

		data = append(data, []interface{}{new_user_id, upvotes_total, downvotes_total})
	}

	UpdateMemberLastActive(member, timestamp)

	return PartiaOutput{Status: "OK", Data: data}
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
