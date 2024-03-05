package mongodb

type UserRoles string

const (
	Admin UserRoles = "ADMIN"
	User  UserRoles = "USER"
)
