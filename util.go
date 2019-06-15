package partia

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func DoesMemberExist(id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM member where id = $1", id).Scan(&count)
	return count == 1
}

func HasNumberAlreadyBeenUsed(number int) bool {
	var count int
	DB.QueryRow(
		"SELECT count(*) FROM already_used_numbers where number = $1", number).Scan(&count)
	fmt.Println(count)
	return count == 1
}

func AreMemberCredsCorrect(id int, unhashedPassword string) bool {
	// TODO
	return true
}

func IsUserActive(id string) {

}

func CreateMember(id int, unhashedPassword string) {
	hashedPassword := hashPassword(unhashedPassword)
	_, err := DB.Query(
		"INSERT INTO member (id, password) "+
			"VALUES ($1, $2)", id, hashedPassword)
	fmt.Println(err)
	markNumberAsUsed(id)
}

func UpdateMemberLastActive(member_id int, timestamp string) {
	_, err := DB.Query(
		"INSERT INTO member_lastactive VALUES ($1, TO_TIMESTAMP($2)) " +
			"ON CONFLICT (member_id) DO UPDATE SET last_active = EXCLUDED.last_active",
		member_id, timestamp)
	fmt.Println(err)
}

func MarkMemberAsLeader(member_id int) {
	DB.Query(
		"INSERT INTO member_isleader VALUES ($1)", member_id)
}

func IsMemberActiveEnough(member_id int) bool {
	var is_he_or_she bool
	DB.QueryRow(
		"SELECT NOW() - last_active < '365 days' FROM member_lastactive WHERE member_id = $1",
		member_id).Scan(&is_he_or_she)
	fmt.Println(is_he_or_she)
	return is_he_or_she
}

func DoesProjectExist(project_id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM project where id = $1", project_id).Scan(&count)
	return count == 1
}

func DoesAuthorityExist(authority_id int) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM authority where id = $1", authority_id).Scan(&count)
	return count == 1
}

func CreateProjectAndMaybeAuthority(project_id, authority_id int) {
	_, err := DB.Query(
		"INSERT INTO authority VALUES ($1)" +
			"ON CONFLICT(id) DO NOTHING", authority_id)
	fmt.Println(err)
	_, err = DB.Query(
		"INSERT INTO project VALUES ($1, $2)",
		project_id, authority_id)
	fmt.Println(err)

	markNumberAsUsed(project_id)
	markNumberAsUsed(authority_id)
}

func CreateAction(action_id, proposed_by_member, project_id int, of_type string) {
	_, err := DB.Query(
		"INSERT INTO action VALUES ($1, $2, $3, $4)",
		action_id, proposed_by_member, project_id, of_type)
	fmt.Println(err)
}

func markNumberAsUsed(number int) {
	DB.Query(
		"INSERT INTO already_used_numbers VALUES ($1)", number)
}

func hashPassword(password string) []byte {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return hashed
}
