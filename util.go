package partia

import (
	"golang.org/x/crypto/bcrypt"
)

func DoesMemberExist(id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM member WHERE id = $1", id).Scan(&count)
	return count == 1
}

func HasNumberAlreadyBeenUsed(number int) bool {
	var count int
	DB.QueryRow(
		"SELECT count(*) FROM already_used_numbers WHERE number = $1", number).Scan(&count)
	return count == 1
}

func AreMemberCredsCorrect(id int, unhashedPassword string) bool {
	var hashedPassword string
	DB.QueryRow(
		"SELECT password FROM MEMBER WHERE id = $1", id).Scan(&hashedPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(unhashedPassword))
	return err == nil
}

func IsMemberLeader(id int) bool {
	var count int
	DB.QueryRow(
		"SELECT count(*) FROM member_isleader WHERE member_id = $1", id).Scan(&count)
	return count > 0
}

func CreateMember(id int, unhashedPassword string) {
	hashedPassword := hashPassword(unhashedPassword)
	DB.Query(
		"INSERT INTO member (id, password) "+
			"VALUES ($1, $2)", id, hashedPassword)
	markNumberAsUsed(id)
}

func UpdateMemberLastActive(member_id int, timestamp string) {
	DB.Query(
		"INSERT INTO member_lastactive VALUES ($1, TO_TIMESTAMP($2)) "+
			"ON CONFLICT (member_id) DO UPDATE SET last_active = EXCLUDED.last_active",
		member_id, timestamp)
}

func MarkMemberAsLeader(member_id int) {
	DB.Query(
		"INSERT INTO member_isleader VALUES ($1)", member_id)
}

func IsMemberActiveEnough(member_id int, timestamp string) bool {
	var is_he_or_she bool
	DB.QueryRow(
		"SELECT TO_TIMESTAMP($1) - last_active < '365 days' FROM member_lastactive WHERE member_id = $2",
		timestamp, member_id).Scan(&is_he_or_she)
	return is_he_or_she
}

func DoesProjectExist(project_id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM project WHERE id = $1", project_id).Scan(&count)
	return count == 1
}

func DoesAuthorityExist(authority_id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM authority WHERE id = $1", authority_id).Scan(&count)
	return count == 1
}

func DoesActionExist(action_id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM action WHERE id = $1", action_id).Scan(&count)
	return count == 1
}

func CreateProjectAndMaybeAuthority(project_id, authority_id int) {
	DB.Query(
		"INSERT INTO authority VALUES ($1)"+
			"ON CONFLICT(id) DO NOTHING", authority_id)
	DB.Query(
		"INSERT INTO project VALUES ($1, $2)",
		project_id, authority_id)

	markNumberAsUsed(project_id)
	markNumberAsUsed(authority_id)
}

func CreateAction(action_id, proposed_by_member, project_id int, of_type string) {
	DB.Query(
		"INSERT INTO action VALUES ($1, $2, $3, $4)",
		action_id, proposed_by_member, project_id, of_type)
	markNumberAsUsed(action_id)
}

func HasUserAlreadyVotedForThisAction(member_id, action_id int) bool {
	var count int
	DB.QueryRow(
		"SELECT count(*) FROM vote WHERE member_id = $1 AND action_id = $2",
		member_id, action_id).Scan(&count)
	return count == 1
}

func InsertVote(member_id, action_id int, up_or_down string) {
	DB.Query(
		"INSERT INTO "+up_or_down+"vote (member_id, action_id) "+
			"VALUES ($1, $2)", member_id, action_id)
}

func markNumberAsUsed(number int) {
	DB.Query(
		"INSERT INTO already_used_numbers VALUES ($1)", number)
}

func hashPassword(password string) []byte {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return hashed
}
