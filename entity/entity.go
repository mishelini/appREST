package entity

import (
	"fmt"
)

//Params application parameters
type Params struct {
	APPHost  string `json:"app_host" yaml:"app_host"`
	APPPort  string `json:"app_port" yaml:"app_port"`
	DBHost   string `json:"db_host" yaml:"db_host"`
	DBPort   string `json:"db_port" yaml:"db_port"`
	DBUser   string `json:"db_user" yaml:"db_user"`
	DBName   string `json:"db_name" yaml:"db_name"`
	DBPass   string `json:"db_pass" yaml:"db_pass"`
	SSLMode  string `json:"ssl_mode" yaml:"ssl_mode"`
	LogFile  string `json:"log_file" yaml:"log_file"`
	InitData bool   `json:"init_data" yaml:"init_data"`
}

func (p *Params) Validate() error {
	if p.APPPort == "" {
		return fmt.Errorf("invalid appport")
	}
	if p.DBHost == "" {
		return fmt.Errorf("invalid dbhost")
	}
	if p.DBPort == "" {
		return fmt.Errorf("invalid portdb")
	}
	if p.DBUser == "" {
		return fmt.Errorf("invalid dbuser")
	}
	if p.DBName == "" {
		return fmt.Errorf("invalid dbname")
	}
	if p.DBPass == "" {
		return fmt.Errorf("invalid dbpass")
	}
	if p.SSLMode == "" {
		return fmt.Errorf("invalid sslmode")
	}
	if p.LogFile == "" {
		return fmt.Errorf("invalid logfilename")
	}
	return nil
}

// TournamentIsFinished bool value to add test data
var TournamentIsFinished = 1

// Player system user
type Player struct {
	ID        int
	FirstName string
	Points    int64
}

// Tournament - competition events
type Tournament struct {
	ID      int
	Deposit int64
	Prize   int64
	Winner  int
	Status  int
}

// TournamentPlayer - player takes part in tournament
type TournamentPlayer struct {
	PlayerID     int
	TournamentID int
}

// Results JSON set
type Results struct {
	Winners []Winner `json:"winners"`
}

// Result JSON output
type Result struct {
	Winner Winner `json:"winner"`
}

// BalanceResults JSON output fro player balance
type BalanceResults struct {
	PlayerId int     `json:"playerId"`
	Balance  float64 `json:"balance"`
}

// Winner user JSON output
type Winner struct {
	PlayerID int     `json:"playerId"`
	Prize    float64 `json:"prize"`
	Balance  float64 `json:"balance"`
}
