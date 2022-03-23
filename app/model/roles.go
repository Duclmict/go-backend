package model

type Roles struct {
	Name         string
	Description  string
	Default  Default `gorm:"embedded"`
}

const (
	RoleSupperAdmin 	= 0
	RoleAdmin     		= 1
	RoleManager     	= 2
	RoleUsers     		= 3
)

func setRole(user *Users, role int) error{
	user.Role = role
	DB.Save(&user)

	return nil
}