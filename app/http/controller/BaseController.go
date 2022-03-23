package controller

import (
	"fmt"

	clone "github.com/huandu/go-clone"
	"github.com/Duclmict/go-backend/app/service/log_service"
)

func IgnoreRequestParams(data map[string][]string) (map[string][]string) {
	temp := clone.Clone(data).(map[string][]string)
	delete(temp , "Default.IsDeleted")
	delete(temp , "Default.CreatedAt")
	delete(temp , "Default.UpdatedAt")
	delete(temp , "Default.CreatedBy")
	delete(temp , "Default.UpdatedBy")
	delete(temp , "Default.CurrentVersion")
	log_service.Debug("IgnoreRequestParams data : " + fmt.Sprint(temp))
	return temp
}