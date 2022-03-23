package model

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"github.com/gorilla/schema"
	"github.com/Duclmict/go-backend/app/helper"
	"github.com/Duclmict/go-backend/app/service/log_service"
)

type SearchParams struct {
	Name 	string					`json:"name"`
	ID	 	uint					`json:"id"`
	Field 	string 					`json:"field"`
	Sign 	string 					`json:"sign"`
	Option []interface{}			`json:"option"`
}

// type ComplexOption struct

const (
	SearchMatch 				= 1
	SearchLike     				= 2
	SearchTime     				= 3
	SearchDate					= 4 
	SearchJson     				= 5
	SearchComplex     			= 6
)

var (
    DB         *gorm.DB
	Decoder    *schema.Decoder
)

func updateCurrentVersion(id string, req interface{}, version int) {
	DB.Model(req).Where("id = ?", id).Update("CurrentVersion", version)
}

// public
func List(req interface{}, s_page string, s_pagi string, order string, 
		datatype interface{}) (res interface{}, totalRows int64, err error) {

	i_page, err := strconv.Atoi(s_page)
	if err != nil {
		return nil, -1, err
	}

	i_pagi, err := strconv.Atoi(s_pagi)
	if err != nil {
		return nil, -1, err
	}

	list := DB.Model(req).Where("is_deleted = ?", "0")
	list.Count(&totalRows)

	if i_pagi > 0 {
		offset := (i_page - 1) * i_pagi
		list = list.Offset(offset).Limit(i_pagi)
	} else {
		list = list.Limit(i_pagi)
	}

	if order != "" {
		list = list.Order(order)
	}

	if err = list.Find(datatype).Error; err != nil {
		return nil, -1, err
	}

	return datatype, totalRows, nil
}

func Get(req interface{}, s_id string, datatype interface{}) (interface{}, error) {

	result := DB.Model(req).Where("is_deleted = ?", "0").Where("id = ?", s_id)

	if err := result.Find(datatype).Error; err != nil {
		return nil, err
	}

	return datatype, nil
}

func Store(req interface{}) (interface{}, error) {
	
	result := DB.Create(req);

	if  result.Error != nil {
		return nil, result.Error
	}

	return req, nil
}

func Update(req interface{}, s_id string) {
	DB.Model(req).Where("is_deleted = ?", "0").Where("id = ?", s_id).Updates(req)
}

func Delete(req interface{}, s_id string) (error) {

	DB.Model(req).Where("is_deleted = ?", "0").Where("id = ?", s_id).Update("is_deleted", "1")
	return nil
}

func Search(req interface{}, s_page string, s_pagi string, order string, 
		params []SearchParams, query map[string][]string,
		datatype interface{}) (res interface{}, totalRows int64, err error) {

	i_page, err := strconv.Atoi(s_page)
	if err != nil {
		return nil, -1, err
	}

	i_pagi, err := strconv.Atoi(s_pagi)
	if err != nil {
		return nil, -1, err
	}

	search := DB.Model(&req).Where("is_deleted = ?", "0")
	if(len(params) > 0 && len(query) > 0) {
		for _, element := range params {
			// check query[search_field]  exits or no
			if search_value,ok := query[element.Name]; ok {
				search_field := element.Name
				if element.Field != "" {
					search_field = element.Field
				}

				log_service.Debug("Query Field: " + fmt.Sprint(search_field + " = ?"))
				log_service.Debug("Query value: " + fmt.Sprint(string(search_value[0])))

				switch element.ID {
					case SearchMatch:
						search = search.Where(search_field + " = ?", string(search_value[0]))
					case SearchLike:
						search = search.Where(search_field + " LIKE ?", "%" + search_value[0] + "%")
					case SearchTime:
						layout := "2006-01-02 15:04:05"
						if s_time, err := time.Parse(layout, search_value[0]);err != nil {
						} else {
							log_service.Debug("S time: " + fmt.Sprint(s_time.String()))
							search = search.Where(search_field + " "+  element.Sign + " ?", s_time.String())
						}
					case SearchDate:
						layout := "2006-01-02 15:04:05"
						if s_date, err := time.Parse(layout, search_value[0] + " 00:00:00");err != nil {
							return nil, -1, err
						} else {
							log_service.Debug("S date: " + fmt.Sprint(s_date.String()))
							search = search.Where(search_field + " "+  element.Sign + " ?", s_date.String())
						}
				}
			}
			
		}
	}

	search.Count(&totalRows)
	
	if i_pagi > 0 {
		offset := (i_page - 1) * i_pagi
		search = search.Offset(offset).Limit(i_pagi)
	} else {
		search = search.Limit(i_pagi)
	}

	if order != "" {
		search = search.Order(order)
	}

	if err := search.Find(datatype).Error; err != nil {
		return nil, -1, err
	}

	return datatype, totalRows, nil
}

func StringToStruct(name string) (interface{},interface{} ,interface{}, error) {
    switch name {
		case "Users":
			return new(Users),new([]Users), new(*Users), nil
		case "Credentials":
			return new(Credentials),new([]Credentials), new(*Credentials), nil
		case "Roles":
			return new(Roles),new([]Roles), new(*Roles), nil
		default:
			return nil, nil, nil, helper.ErrCanNotConvertStruct
    }
}