package model

import (
	"time"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"
	"github.com/Duclmict/go-backend/app/service/log_service"
)

const (  // iota is reset to 0
	TypePassword = 0
	TypeResetPassword = 1
)

type Credentials struct {
	OwnerID      uuid.UUID 
	Type		 uint
	Token  		 string
	ExpiresAt    time.Time
	Default  Default `gorm:"embedded"`
}

func generatePasswordHash(password string) (res string, err error){
	hashedpwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
        return "", err
    }
	return string(hashedpwd), nil
}

func createorUpdatePasswordUser(user *Users, password string) error{
	
	var cre Credentials

	cre.OwnerID = user.Default.ID
	cre.Type = TypePassword

	token,err_t := generatePasswordHash(password)
	if err_t != nil {
		log_service.Error("MODEL:[Credentials]" + fmt.Sprint(err_t))
		return err_t
	}
	cre.Token = token 

	is_exits := DB.Model(&cre).Where("owner_id = ?", user.Default.ID).Where("type = ?", TypePassword).Updates(&cre)
	if is_exits.Error != nil {
		log_service.Error("MODEL:[Credentials]" + fmt.Sprint(is_exits.Error))
		return is_exits.Error
	}

	if is_exits.RowsAffected == 0 {
		cre.Default.ID	= uuid.NewV4()
		if err := DB.Select("owner_id", "type", "token", "id").Create(&cre).Error; err != nil {
			log_service.Error("MODEL:[Credentials]" + fmt.Sprint(err))
			return err
		}
	}

	return nil
}