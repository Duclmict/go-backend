package database

import (
    "fmt"
    "reflect"
    "strconv"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/Duclmict/go-backend/app/model"
    "github.com/Duclmict/go-backend/config"
    "github.com/Duclmict/go-backend/app/service/log_service"
    "github.com/gorilla/schema"
)

func Init()  {
    db := Connection()
    model.DB = db

    timeConverter := func(value string) reflect.Value {
		tstamp, err := strconv.ParseInt(value, 10, 64)
        log_service.Debug("value:" + string(value))
        if err == nil {
            return reflect.ValueOf(time.Unix(tstamp, 0))
        }

        // date 
        layout := "2006-01-02"
        t_date, err := time.Parse(layout, value)
        log_service.Debug("t_date:" + t_date.String())
        if err == nil {
            return reflect.ValueOf(t_date)
        }
        
        // time 
        layout = "2006-01-02 13:03:23"
        t_time, err := time.Parse(layout, value)
        log_service.Debug("t_time:" + t_time.String())
        if err == nil {
            return reflect.ValueOf(t_time)
        }
        
        return reflect.ValueOf(time.Now)
	}
	decoder := schema.NewDecoder()
	decoder.RegisterConverter(time.Time{}, timeConverter)
    model.Decoder = decoder

    db.AutoMigrate(&model.Users{}, &model.Roles{}, &model.Credentials{})
    return
}

func Connection() *gorm.DB{

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s", 
        config.DB_MYSQL_USER, config.DB_MYSQL_PASSWORD, config.DB_MYSQL_HOST, config.DB_MYSQL_PORT, 
        config.DB_MYSQL_DATABASE, config.DB_MYSQL_CHARSET, config.DB_MYSQL_PARSETIME, config.DB_MYSQL_LOC)
    
    log_service.Debug("MYSQL config DSN:" + dsn)

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    
    if err != nil {
        panic("Failed to create a connection to database")
    }

    return db
}

func Close(db *gorm.DB) {
    dbSQL, err := db.DB()
    if err != nil {
        panic("Failed to close connection from database")
    }
    dbSQL.Close()
}