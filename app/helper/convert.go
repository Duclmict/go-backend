package helper

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	clone "github.com/huandu/go-clone"
)

// mysql
func ConvertQuery2Json(query interface{}) interface{} {
	arrayData, err := json.Marshal(query)

	if err != nil {
		return ""
	}

	jsonLength := len(arrayData)
	return  arrayData[:jsonLength]
}

func GenerateUUID(data map[string][]string) (map[string][]string) {
	temp := clone.Clone(data).(map[string][]string)
	temp["Default.ID"] = []string{uuid.NewV4().String()}
	return temp
}