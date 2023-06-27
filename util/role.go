package util

const (
	ADMIN = 1
	USER  = 2
)

func IsSupportedRole(role int) bool {
	switch role {
	case ADMIN, USER:
		return true
	}
	return false
}

func User() int {
	return USER
}

func Admin() int {
	return ADMIN
}
