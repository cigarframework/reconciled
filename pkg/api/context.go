package api

import "context"

var contextUser = struct{}{}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, contextUser, user)
}

func ContextUser(ctx context.Context) *User {
	in := ctx.Value(contextUser)
	if in == nil {
		return nil
	}
	user := in.(*User)
	return user
}

func WithToken(ctx context.Context, token string) context.Context {
	return WithUser(ctx, &User{Token: token})
}
