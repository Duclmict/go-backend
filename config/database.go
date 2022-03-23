package config

// prefix DB

var (
	DB_CONNECTMODE	string   											// [sql, nosql, other]
	DB_CONNECT	string 											 		// default mysql

	// mysql
	DB_MYSQL_DRIVER string   = "mysql"									// default mysql
	DB_MYSQL_HOST string	
	DB_MYSQL_PORT string  										 		// default 3306
	DB_MYSQL_DATABASE string  
	DB_MYSQL_USER string  
	DB_MYSQL_PASSWORD string  
	DB_MYSQL_UNIX_SOCKET string	
	DB_MYSQL_CHARSET string	  = "utf8"
	DB_MYSQL_PARSETIME string =	"True"									// default true
	DB_MYSQL_LOC	string 	  = ""									//  default local
)