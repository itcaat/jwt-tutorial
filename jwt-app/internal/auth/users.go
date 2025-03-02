package auth

// Users - фиктивная база данных пользователей
var Users = map[string]struct {
	Password string
	Role     string
}{
	"reader": {"password", "read"},
	"writer": {"password", "write"},
}
