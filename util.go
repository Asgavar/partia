package partia

func DoesMemberExist(id string) bool {
	var count int
	DB.QueryRow("SELECT count(*) FROM member where id = ?", id).Scan(&count)
	return count == 1
}

func AreMemberCredsCorrect(id, unhashedPassword string) {

}

func IsUserActive(id string) {

}

func CreateUser(id, unhashedPassword string) {

}
