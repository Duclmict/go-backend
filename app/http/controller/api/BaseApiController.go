package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Duclmict/go-backend/app/http/controller"
	"github.com/Duclmict/go-backend/app/model"
	"github.com/Duclmict/go-backend/app/service/log_service"
	"github.com/go-playground/validator/v10"
	"github.com/Duclmict/go-backend/app/helper"
)

// type
type CheckStore func(map[string][]string) (error)
type BeforeStore func(map[string][]string) (map[string][]string, error)
type AfterStore func(interface{}, map[string][]string) (error)

type CheckUpdate func(string, map[string][]string) (error)
type BeforeUpdate func(string, map[string][]string) (map[string][]string, error)
type AfterUpdate func(string, interface{}, map[string][]string) (error)

type StoreStruct struct {
	BeforeStoreFlag bool
    BeforeStoreFunc BeforeStore
	AfterStoreFlag  bool
	AfterStoreFunc  AfterStore
	CheckStoreFlag  bool
	CheckStoreFunc	CheckStore
}

type UpdateStruct struct {
	BeforeUpdateFlag bool
    BeforeUpdateFunc BeforeUpdate
	AfterUpdateFlag  bool
	AfterUpdateFunc  AfterUpdate
	CheckUpdateFlag  bool
	CheckUpdateFunc	 CheckUpdate
}

var (
	pagination string = "10"
)

// public
func Index(c *gin.Context, ModelName string, ModelOrder string) {

	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[INDEX]-> START")

	index_model,index_model_data,_,err_m := model.StringToStruct(ModelName)
	if err_m != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_m))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_m)})
		return
	}

	use_pagi := pagination

	limit,ok := c.GetQuery("limit")
	if ok {
		use_pagi = limit
	}

	list, _, err := model.List(index_model, "0", use_pagi, ModelOrder, index_model_data)

	if err != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err)})
		return
	}
	
	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[INDEX]-> END")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func Get(c *gin.Context, ModelName string) {
	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[GET]-> START")

	get_model,get_model_data,_,err_m := model.StringToStruct(ModelName)
	if err_m != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_m))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_m)})
		return
	}

	id := c.Param("id")

	_, err_s := model.Get(get_model, id, get_model_data)
	if err_s != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_s))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_s)})
		return
	}

	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[GET]-> END")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": get_model_data})
}

func Store(c *gin.Context, ModelName string, ModelStore StoreStruct, StoreValidate interface{}) {
	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[STORE]-> START")

	store_model,_,_,err_m := model.StringToStruct(ModelName)
	if err_m != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_m))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_m)})
		return
	}

	// get and assign data
	var request_data = make(map[string][]string)
	
	c.MultipartForm()
    for key, value := range c.Request.PostForm {
		if key == "IsDeleted" ||
		key == "CreatedBy" || key == "UpdatedBy" || key == "CurrentVersion" {
            request_data[string("Default." + key)] = value
        } else {
            request_data[key] = value
        }
    }
	log_service.Debug("Request Data: " + fmt.Sprint(request_data))

	// assign UUID
	request_data = helper.GenerateUUID(request_data)

	// remove field not allow update

	// validate
	model.Decoder.Decode(StoreValidate, request_data)
	log_service.Debug("Validate model: " + fmt.Sprint(StoreValidate))

	validate := validator.New()
	if err_v := validate.Struct(StoreValidate);err_v != nil {

		// for _, err_i := range err_v.(validator.ValidationErrors) {
		// 	log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_i))
		// 	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_i)})
		// 	// fmt.Println(err)//Key: 'Users.Passwd' Error:Field validation for 'Passwd' failed on the 'min' tag
		// 	return
		// }

		// if _, ok := err_v.(*validator.InvalidValidationError); ok {
		// 	fmt.Println(err_v)
		// 	return
		// }

		// for _, err_i := range err_v.(validator.ValidationErrors) {

		// 	fmt.Println(err_i.Namespace())
		// 	fmt.Println(err_i.Field())
		// 	fmt.Println(err_i.StructNamespace())
		// 	fmt.Println(err_i.StructField())
		// 	fmt.Println(err_i.Tag())
		// 	fmt.Println(err_i.ActualTag())
		// 	fmt.Println(err_i.Kind())
		// 	fmt.Println(err_i.Type())
		// 	fmt.Println(err_i.Value())
		// 	fmt.Println(err_i.Param())
		// 	fmt.Println()
		// }

		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_v))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_v)})
		return
	}

	// Check releation or error data after save
	if(ModelStore.CheckStoreFlag) {
		if err_c := ModelStore.AfterStoreFunc(store_model, request_data); err_c != nil {
			log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_c))
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_c)})
			return
		}
	}

	// begin transaction

	// before store and remove field not allow update
	save_data := controller.IgnoreRequestParams(request_data)

	if(ModelStore.BeforeStoreFlag) {
		store_data,err_b := ModelStore.BeforeStoreFunc(request_data)
		if  err_b != nil {
			log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_b))
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_b)})
			return
		}
		save_data = controller.IgnoreRequestParams(store_data)
	}

	model.Decoder.Decode(store_model, save_data)

	// map data and save
	log_service.Debug("Data model: " + fmt.Sprint(store_model))

	result, err_s := model.Store(store_model)
	if err_s != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_s))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_s)})
		return
	}

	fmt.Printf("%+v\n", store_model)

	// afterstore
	if(ModelStore.AfterStoreFlag) {
		if err_a := ModelStore.AfterStoreFunc(store_model, request_data); err_a != nil {
			log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_a))
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_a)})
			return
		}
	}

	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[STORE]-> END")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func Update(c *gin.Context, ModelName string, ModelUpdate UpdateStruct, UpdateValidate interface{}) {
	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[UPDATE]-> START")

	update_model,_,_,err_m := model.StringToStruct(ModelName)
	if err_m != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_m))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_m)})
		return
	}

	id := c.Param("id")

	// get and assign data
	var request_data = make(map[string][]string)
	
	c.MultipartForm()
    for key, value := range c.Request.PostForm {
		if key == "IsDeleted" || key == "ID" ||
		key == "CreatedBy" || key == "UpdatedBy" || key == "CurrentVersion" {
            request_data[string("Default." + key)] = value
        } else {
            request_data[key] = value
        }
    }
	log_service.Debug("Request Data: " + fmt.Sprint(request_data))

	// validate
	model.Decoder.Decode(UpdateValidate, request_data)
	log_service.Debug("Validate model: " + fmt.Sprint(UpdateValidate))

	validate := validator.New()
	if err_v := validate.Struct(UpdateValidate);err_v != nil {

		// for _, err_i := range err_v.(validator.ValidationErrors) {
		// 	log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_i))
		// 	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_i)})
		// 	// fmt.Println(err)//Key: 'Users.Passwd' Error:Field validation for 'Passwd' failed on the 'min' tag
		// 	return
		// }

		// if _, ok := err_v.(*validator.InvalidValidationError); ok {
		// 	fmt.Println(err_v)
		// 	return
		// }

		// for _, err_i := range err_v.(validator.ValidationErrors) {

		// 	fmt.Println(err_i.Namespace())
		// 	fmt.Println(err_i.Field())
		// 	fmt.Println(err_i.StructNamespace())
		// 	fmt.Println(err_i.StructField())
		// 	fmt.Println(err_i.Tag())
		// 	fmt.Println(err_i.ActualTag())
		// 	fmt.Println(err_i.Kind())
		// 	fmt.Println(err_i.Type())
		// 	fmt.Println(err_i.Value())
		// 	fmt.Println(err_i.Param())
		// 	fmt.Println()
		// }

		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_v))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_v)})
		return
	}

	// Check releation or error data after save
	if(ModelUpdate.CheckUpdateFlag) {
		err_c := ModelUpdate.CheckUpdateFunc(id, request_data)
		if  err_c != nil {
			log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_c))
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_c)})
			return
		}
	} 
	// begin transaction
	
	// before update and remove field not allow update
	save_data := controller.IgnoreRequestParams(request_data)

	if(ModelUpdate.BeforeUpdateFlag) {
		update_data,err_b := ModelUpdate.BeforeUpdateFunc(id, request_data)
		if  err_b != nil {
			log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_b))
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_b)})
			return
		}
		save_data = controller.IgnoreRequestParams(update_data)
	}

	// map data and save
	model.Decoder.Decode(update_model, save_data)
	log_service.Debug("Update data load: " + fmt.Sprint(update_model))

	model.Update(update_model, id)
	fmt.Printf("%+v\n", update_model)

	// afterstore
	if(ModelUpdate.AfterUpdateFlag) {
		if err_a := ModelUpdate.AfterUpdateFunc(id, update_model, request_data); err_a != nil {
			log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_a))
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_a)})
			return
		}
	}
	
	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[UPDATE]-> END")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": update_model})
}

func Delete(c *gin.Context, ModelName string) {
	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[DELETE]-> START")

	delete_model,_,_,err_m := model.StringToStruct(ModelName)
	if err_m != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_m))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_m)})
		return
	}

	id := c.Param("id")

	err_s := model.Delete(delete_model, id)
	if err_s != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_s))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_s)})
		return
	}

	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[DELETE]-> END")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

func Search(c *gin.Context, ModelName string, ModelOrder string, ModelSearch []model.SearchParams) {

	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[SEARCH]-> START")

	use_pagi := pagination

	limit,ok := c.GetQuery("limit")
	if ok {
		use_pagi = limit
	}

	query := c.Request.URL.Query() // query is empty

	search_model,search_model_data,_,err_m := model.StringToStruct(ModelName)
	if err_m != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_m))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_m)})
		return
	}

	result, _, err_s := model.Search(search_model, "0", use_pagi, ModelOrder, ModelSearch, query, search_model_data)
	if err_s != nil {
		log_service.Error("MODEL:[" + ModelName + "]" + fmt.Sprint(err_s))
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err_s)})
		return
	}

	log_service.Info("MODEL:[" + ModelName + "]  " + "FUNC:[SEARCH]-> END")

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"data": result})
}