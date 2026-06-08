package auth

type contextKey string

const ContextKey contextKey = "auth"

type AuthContext struct {
	Namespace   string
	APIKeyID    int64
	Permissions []string
}