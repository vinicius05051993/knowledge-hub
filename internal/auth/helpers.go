package auth

import "net/http"

func GetAuthContext(
	r *http.Request,
) *AuthContext {

	value := r.Context().Value(
		ContextKey,
	)

	if value == nil {
		return nil
	}

	authContext, ok := value.(*AuthContext)

	if !ok {
		return nil
	}

	return authContext
}

func GetNamespace(
	r *http.Request,
) string {

	authContext := GetAuthContext(r)

	if authContext == nil {
		return ""
	}

	return authContext.Namespace
}