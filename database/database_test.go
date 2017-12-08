package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/mishelini/entity"
	"github.com/stretchr/testify/assert"
	// Pure Go Postgres driver for database/sql
	_ "github.com/lib/pq"
)

//Prepare statement
var testUser = entity.Player{
	ID:        1,
	FirstName: "testuser",
	Points:    200,
}
var testUser2 = &entity.Player{
	ID:        2,
	FirstName: "testuser2",
	Points:    300,
}
var testTournament = entity.Tournament{
	ID:      1,
	Deposit: 200,
	Status:  0,
	Prize:   0,
	Winner:  0,
}

func getDBConnection() (*sql.DB, error) {
	postgresConfig := fmt.Sprintf("host=%s port=%s   user=%s dbname=%s sslmode=%s  password=%s",
		"localhost", "5432", "postgres", "postgres", "disable", "postgres")
	dbConn, err := sql.Open("postgres", postgresConfig)
	return dbConn, err
}
func initTestDb(db *sql.DB) error {
	new_schema_strings := []string{
		`CREATE TABLE IF NOT EXISTS player
		(
		   id        SERIAL PRIMARY KEY,
		   first_name VARCHAR(30),
		   points     BIGINT DEFAULT 0
		);`,
		` CREATE TABLE IF NOT EXISTS tournament
		(
		   id      SERIAL NOT NULL PRIMARY KEY,
		   deposit BIGINT NOT NULL DEFAULT 0,
		   prize   BIGINT NOT NULL DEFAULT 0,
		   status   INT DEFAULT 0,
		   winner   INT DEFAULT 0
		   
		);`,
		` CREATE TABLE IF NOT EXISTS tournament_player
		(
		   player_id     INT REFERENCES player (id) ON UPDATE CASCADE ON DELETE
		   CASCADE,
		   tournament_id INT REFERENCES tournament (id) ON UPDATE CASCADE,
		   CONSTRAINT tournament_player_pkey PRIMARY KEY (player_id, tournament_id)
		);`,
		` INSERT INTO player (first_name, points) VALUES ('testuser', 0);`,
		` INSERT INTO player (first_name, points) VALUES ('testuser2', 0);`,
	}
	for _, qstr := range new_schema_strings {
		_, err := db.Exec(qstr)
		if err != nil {
			return err
		}
	}
	return nil

}

func prepareTestEnv() (*sql.DB, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}
	err = dropTestSchema(db)
	if err != nil {
		return nil, err
	}
	err = setTestSchema(db)
	if err != nil {
		return nil, err
	}
	err = initTestDb(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func announceTestTournament(db *sql.DB, tournamentID int, deposit int64) error {
	id := 0
	err := db.QueryRow("INSERT INTO tournament (id, deposit)  VALUES($1, $2) RETURNING id", tournamentID, deposit).Scan(&id)
	return err

}

func insertUserIntoTournament(db *sql.DB, userID int, tournamentID int) error {
	player_id := 0
	err := db.QueryRow("INSERT INTO tournament_player (player_id , tournament_id ) VALUES( $1 ,$2 ) RETURNING player_id", userID, tournamentID).Scan(&player_id)
	return err
}
func selectFinishedTournaments(db *sql.DB) ([]entity.Tournament, error) {
	tournaments := make([]entity.Tournament, 0)
	rows, err := db.Query(`SELECT * FROM tournament WHERE status = $1 `, entity.TournamentIsFinished)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.Tournament
		if err := rows.Scan(&t.ID, &t.Deposit, &t.Prize, &t.Winner, &t.Status); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil

}

func setTestSchema(db *sql.DB) error {
	new_schema_strings := []string{
		"CREATE SCHEMA IF NOT EXISTS test_schema ;",
		"SET search_path TO test_schema;",
	}
	for _, qstr := range new_schema_strings {
		_, err := db.Exec(qstr)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropTestSchema(db *sql.DB) error {
	new_schema_strings := "DROP  SCHEMA IF EXISTS test_schema CASCADE;"
	_, err := db.Exec(new_schema_strings)
	return err
}
func TestDBConnection(t *testing.T) {
	db, err := getDBConnection()
	defer db.Close()
	actual := err
	if actual != nil {
		t.Errorf("test DB connection failed", err)
	}
}
func TestCreateTables(t *testing.T) {
	db, err := getDBConnection()
	assert.NoError(t, err, "test DB connection failed")
	err = setTestSchema(db)
	assert.NoError(t, err, "func setTestSchema failed")
	defer db.Close()

	err = CreateTablesIfNotExist(db)
	assert.Equal(t, err, nil, "test func CreateTablesIfNotExist  failed")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}
func TestFundPlayer(t *testing.T) {
	var player entity.Player
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	assert.NoError(t, err, "func initTestDb failed")
	err = FundPlayer(db, testUser.ID, testUser.Points)
	assert.NoError(t, err, "fuc FundPlayer return error")
	row := db.QueryRow("SELECT * FROM player WHERE id = $1 ", testUser.ID)
	err = row.Scan(&player.ID, &player.FirstName, &player.Points)
	assert.NoError(t, err, "select player return error")
	assert.Equal(t, testUser.Points, player.Points, "player points after funding should be equal")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}
func TestAnnounceTournaments(t *testing.T) {
	var tournament entity.Tournament
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	err = AnnounceTournaments(db, testTournament.ID, testTournament.Deposit)
	assert.NoError(t, err, "func AnnounceTournaments failed")
	row := db.QueryRow("SELECT * FROM tournament WHERE id = $1 ", 1)
	err = row.Scan(&tournament.ID, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status)
	assert.NoError(t, err, "select tournament return error")
	assert.Equal(t, testTournament.ID, tournament.ID, "no test tournament in db")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}
func TestSelectPlayer(t *testing.T) {
	var player entity.Player
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	testPlayer, err := SelectPlayer(db, testUser.ID)
	assert.NoError(t, err, "func SelectPlayer failed")
	row := db.QueryRow("SELECT * FROM player WHERE id = $1 ", testUser.ID)
	err = row.Scan(&player.ID, &player.FirstName, &player.Points)
	assert.NoError(t, err, "selecting player return error")
	assert.Equal(t, testPlayer.FirstName, player.FirstName, "test user not selected")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}

func TestSelectTournament(t *testing.T) {
	var tournament entity.Tournament
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	err = announceTestTournament(db, 1, 200)
	assert.NoError(t, err, "func announceTestTournament failed")
	firstTournament, err := SelectTournament(db, testTournament.ID)
	assert.NoError(t, err, "func SelectTournament failed")
	row := db.QueryRow("SELECT * FROM tournament WHERE id = $1 ", testTournament.ID)
	err = row.Scan(&tournament.ID, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status)
	assert.NoError(t, err, "selecting tournament return error")
	assert.Equal(t, firstTournament.ID, tournament.ID, "test tournament not selected")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}

func TestInsertUserIntoTournament(t *testing.T) {
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	err = announceTestTournament(db, testTournament.ID, testTournament.Deposit)
	assert.NoError(t, err, "func announceTestTournament1 failed")
	err = InsertUserIntoTournament(db, testTournament.ID, testUser.ID)
	assert.NoError(t, err, "func InsertUserIntoTournament failed")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}

func TestSelectTournamentUsers(t *testing.T) {
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()
	err = announceTestTournament(db, testTournament.ID, testTournament.Deposit)
	assert.NoError(t, err, "func announceTestTournament failed")
	err = insertUserIntoTournament(db, testUser.ID, testTournament.ID)
	assert.NoError(t, err, "func insertUserIntoTournament failed")

	tournamentUserSet, err := SelectTournamentUsers(db, testTournament.ID)
	assert.NoError(t, err, "func SelectTournament failed")
	players := make([]entity.TournamentPlayer, 0)
	rows, err := db.Query(`SELECT * FROM tournament_player WHERE tournament_id = $1 `, testTournament.ID)
	assert.NoError(t, err, "select tournament_player set failed")
	for rows.Next() {
		var player entity.TournamentPlayer
		err = rows.Scan(&player.PlayerID, &player.TournamentID)
		assert.NoError(t, err, "select tournament_player set failed")
		players = append(players, player)
	}
	assert.Equal(t, players, tournamentUserSet, "test tournament not selected")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}
func TestChangeTournamentsPrize(t *testing.T) {
	var newPrize int64 = 100 + testTournament.Prize
	var tournament entity.Tournament
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	err = announceTestTournament(db, testTournament.ID, testTournament.Deposit)
	assert.NoError(t, err, "func announceTestTournament failed")
	err = ChangeTournamentsPrize(db, testTournament.ID, newPrize)
	assert.NoError(t, err, "func ChangeTournamentsPrize failed")
	row := db.QueryRow("SELECT * FROM tournament WHERE id = $1 ", testTournament.ID)
	err = row.Scan(&tournament.ID, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status)
	assert.NoError(t, err, "select tournament return error")
	assert.Equal(t, newPrize, tournament.Prize, "no test tournament in db")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}

func TestFinishTournament(t *testing.T) {
	var tournament entity.Tournament
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	err = announceTestTournament(db, testTournament.ID, testTournament.Deposit)
	assert.NoError(t, err, "func announceTestTournament failed")
	err = FinishTournament(db, testTournament.ID, testUser.ID)
	assert.NoError(t, err, "func FinishTournament failed")
	row := db.QueryRow("SELECT * FROM tournament WHERE id = $1 ", testTournament.ID)
	err = row.Scan(&tournament.ID, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status)
	assert.NoError(t, err, "select tournament return error")
	assert.Equal(t, entity.TournamentIsFinished, tournament.Status, "no test tournament in db")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}
func TestSelectFinishedTournaments(t *testing.T) {
	db, err := prepareTestEnv()
	assert.NoError(t, err, "func prepareTestEnv failed")
	defer db.Close()

	err = announceTestTournament(db, testTournament.ID, testTournament.Deposit)
	assert.NoError(t, err, "func announceTestTournament failed")
	err = FinishTournament(db, testTournament.ID, testUser.ID)
	assert.NoError(t, err, "func FinishTournament failed")
	testTornaments, err := selectFinishedTournaments(db)
	tornaments, err := SelectFinishedTournaments(db)
	assert.NoError(t, err, "select tournament return error")
	assert.Equal(t, tornaments, testTornaments, "no test tournament in db")

	err = dropTestSchema(db)
	assert.NoError(t, err, "func dropTestSchema faild")
}
