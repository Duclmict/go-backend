package users

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Duclmict/go-backend/app/model"
	"github.com/Duclmict/go-backend/app/http/controller/api"
)

var (
	ModelName string = model.UsersModelName
	ModelOrder string = model.UsersOrder
	ModelSearch []model.SearchParams = model.UsersSearch
	ModelStore api.StoreStruct = api.StoreStruct { 
		true, model.UsersBeforeStore, 			// before store users action define
		true, model.UsersAfterStore, 			// after store users action define
		false, nil,								// check store users action define
	}
	ModelUpdate api.UpdateStruct = api.UpdateStruct { 
		false, nil, 							// before update users action define
		true, model.UsersAfterUpdate, 			// after update users action define
		true, model.UsersCheckUpdate,			// check update users action define
	}
)

type StoreValidate struct {
	Name         			string   	`form:"name" json:"name" validate:"required,max=20,min=6"`
	Email        			string   	`form:"email" json:"email" validate:"required,email,max=80"`
	Age          			int	  		`form:"age" json:"age" validate:"max=200,min=0"`
	Birthday     			time.Time 	
	Password	 			string		`form:"password" json:"password" validate:"required,max=20,min=6"`
	PasswordConfirmation	string		`form:"password_confirmation" json:"password_confirmation" validate:"required"`
}

type UpdateValidate struct {
	Name         			string   	`form:"name" json:"name" validate:"required,max=20,min=6"`
	Email        			string   	`form:"email" json:"email" validate:"required,email,max=80"`
	Age          			int	  		`form:"age" json:"age" validate:"max=200,min=0"`
	Birthday     			time.Time 	
}

func Index(c *gin.Context) {
	api.Index(c, ModelName, ModelOrder)
}

func Get(c *gin.Context) {
	api.Get(c, ModelName)
}

func Search(c *gin.Context) {
	api.Search(c, ModelName, ModelOrder, ModelSearch)
}

func Store(c *gin.Context) {
	api.Store(c, ModelName, ModelStore, new(StoreValidate))
}

func Update(c *gin.Context) {
	api.Update(c, ModelName, ModelUpdate, new(UpdateValidate))
}

func Delete(c *gin.Context) {
	api.Delete(c, ModelName)
}