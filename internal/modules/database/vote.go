package database

func voteIntToBool(vote int32) bool {
	if vote == 1 {
		return true
	}
	return false
}
