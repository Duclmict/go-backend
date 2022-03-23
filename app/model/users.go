package model

import (
	"fmt"
	"time"
	"strconv"
	
	"golang.org/x/crypto/bcrypt"
	"github.com/Duclmict/go-backend/app/helper"
	"github.com/Duclmict/go-backend/app/service/log_service"
	clone "github.com/huandu/go-clone"
	"github.com/Duclmict/go-backend/config"
	"github.com/Duclmict/go-backend/app/service/mail_service"
)

const IdentityKey = "id"

var (
	UsersModelName string = "Users"
	UsersOrder string = "id asc"
	UsersSearch []SearchParams = []SearchParams {
		SearchParams {
			Name: "user_name",
			Field: "name",
			ID:	SearchLike,
		},
		SearchParams {
			Name: "email",
			ID:	SearchLike,
			Option: nil,
		},
		SearchParams {
			Name: "age",
			ID:	SearchMatch,
			Option: nil,
		},
		SearchParams {
			Name: "birthday",
			ID:	SearchDate,
			Sign: "=",
			Option: nil,
		},
	}
)

type Users struct {
	Name         string   
	Email        string   
	Age          int	 
	Birthday     time.Time
	Role		 int		
	Default      Default 	  `gorm:"embedded"`
}

func UsersCheckStore(resquestData map[string][]string) (map[string][]string, error) {

	return resquestData, nil
}

func UsersBeforeStore(resquestData map[string][]string) (map[string][]string, error) {
	// ignore Password and PasswordConfirmation field
	temp := clone.Clone(resquestData).(map[string][]string)
	delete(temp , "Password")
	delete(temp , "PasswordConfirmation")
	log_service.Debug("data: " + fmt.Sprint(temp))
	return temp, nil
}

func UsersAfterStore(object interface{}, resquestData map[string][]string) (error) {

	users, ok := object.(*Users)
	if !ok {
		log_service.Error("type not converted")
		return helper.ErrCanNotConvertType
	}

	password, err := helper.MapRDGetKeybyValue(resquestData, "Password")
	log_service.Debug("Password: " + fmt.Sprint(password))
	if err != nil {
		return err
	}
	
	err_v := createorUpdatePasswordUser(users, password)
	if err_v != nil {
		return err_v
	}

	err_s := setRole(users, RoleUsers)
	if err_s != nil {
		return err_s
	}

	// send mail
	err_m := mail_service.Send(config.MAIL_TEMPLATE["hello"], []*string {&users.Email} , "hello world", nil, nil)
	if err_m != nil {
		return err_m
	}
	
	return nil
}

func UsersCheckUpdate(id string, resquestData map[string][]string) (error) {

	// Check current version with lastest version
	currentVersion, err := helper.MapRDGetKeybyValue(resquestData, "Default.CurrentVersion")
	log_service.Debug("CurrentVersion request: " + fmt.Sprint(currentVersion))
	if err != nil {
		return err
	}

	cur_version, err_c := strconv.Atoi(currentVersion)
	if err_c != nil {
		return err
	}

	var user Users
	DB.First(&user, id)
	log_service.Debug("CurrentVersion get from database: " + fmt.Sprint(user.Default.CurrentVersion))

	if cur_version != int(user.Default.CurrentVersion) {
		return helper.ErrCurentVersionNotLastest
	}

	return nil
}

func UsersBeforeUpdate(id string, resquestData map[string][]string) (map[string][]string, error) {
	// ignore Password and PasswordConfirmation field
	temp := clone.Clone(resquestData).(map[string][]string)
	delete(temp , "Password")
	delete(temp , "PasswordConfirmation")
	log_service.Debug("data: " + fmt.Sprint(temp))
	return temp, nil
}

func UsersAfterUpdate(id string,req interface{},resquestData map[string][]string) (error) {

	// Check current version with lastest version
	currentVersion, err := helper.MapRDGetKeybyValue(resquestData, "Default.CurrentVersion")
	log_service.Debug("CurrentVersion: " + fmt.Sprint(currentVersion))
	if err != nil {
		return err
	}

	cur_version, err_c := strconv.Atoi(currentVersion)
	if err_c != nil {
		return err
	}

	updateCurrentVersion(id, req, cur_version + 1)

	return nil
}

func VerifyLogin(email string, password string) (res interface{}, err error) {

	var users Users
	var credentials Credentials

	// find user
	if err_u := DB.Where("email = ?", email).Where("is_deleted = ?", "0").First(&users).Error; err_u != nil {
		return nil, err_u
	}

	// find PasswordHash
	if err_c := DB.Where("owner_id = ?", users.Default.ID).Where("is_deleted = ?", "0").Where("type = ?", TypePassword).First(&credentials).Error; err_c != nil {
		return nil, err_c
	}

    // Comparing the password with the hash
    err_s := bcrypt.CompareHashAndPassword([]byte(credentials.Token), []byte(password))
    if err_s != nil {
		return nil, err_s
    }

	return users, nil
}