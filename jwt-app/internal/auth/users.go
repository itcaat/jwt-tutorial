package auth

// Users - фиктивная база данных пользователей
var Users = map[string]struct {
	Username string
	Password string
	Role     string
}{
	"reader": {"reader", "password", "read"},
	"writer": {"writer", "password", "write"},
}
