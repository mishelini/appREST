package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mishelini/database"

	"github.com/go-yaml/yaml"
	// Pure Go Postgres driver for database/sql
	_ "github.com/lib/pq"
	"github.com/mishelini/entity"
	"github.com/mishelini/handler"
)

var (
	configFile string
	appParams  entity.Params
)

func main() {
	initDataFromFile()

	flag.StringVar(&appParams.LogFile, "logfile", appParams.LogFile, "Log File")
	flag.Parse()
	f, err := os.OpenFile(appParams.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	err = processFlags(&appParams)
	if err != nil {
		log.Printf("parcing flags: %s", err)
		return
	}
	db, err := initializeDB(appParams)
	if err != nil {
		log.Printf("initialize DB: %s", err)
		return
	}
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", appParams.APPHost, appParams.APPPort), handler.Handler(db))
	if err != nil {
		log.Printf("initialize DB: %s", err)
	}
}

func initDataFromFile() {
	flag.StringVar(&configFile, "config_file", "./conf.yaml", "Set current file for parsing  Data from YAML file")
	flag.Parse()

	if configFile != "" {
		yamlFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal([]byte(yamlFile), &appParams)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initializeDB(appParams entity.Params) (*sql.DB, error) {
	postgresConfig := fmt.Sprintf("host=%s port=%s  user=%s dbname=%s sslmode=%s  password=%s",
		appParams.DBHost, appParams.DBPort, appParams.DBUser, appParams.DBName, appParams.SSLMode, appParams.DBPass)
	dbConn, err := sql.Open("postgres", postgresConfig)
	if err != nil {
		return nil, err
	}
	database.InitData = appParams.InitData
	err = database.CreateTablesIfNotExist(dbConn)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func processFlags(appParams *entity.Params) error {
	flag.BoolVar(&appParams.InitData, "initdata", appParams.InitData, "Set default data")
	flag.StringVar(&appParams.DBHost, "dbhost", appParams.DBHost, "Data Base Host")
	flag.StringVar(&appParams.DBName, "dbname", appParams.DBName, "Data Base Name")
	flag.StringVar(&appParams.DBPass, "dbpass", appParams.DBPass, "Data Base Password")
	flag.StringVar(&appParams.DBPort, "dbport", appParams.DBPort, "Data Base Port")
	flag.StringVar(&appParams.DBUser, "dbuser", appParams.DBUser, "Data Base User")
	flag.StringVar(&appParams.APPPort, "appport", appParams.APPPort, "APP Port")
	flag.StringVar(&appParams.SSLMode, "sslmode", appParams.SSLMode, "Data Base SSL Mode")
	flag.Parse()

	return appParams.Validate()
}

func Add(value1 int, value2 int) int {
	return value1 + value2
}
