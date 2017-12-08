package database

import (
	"database/sql"

	"github.com/mishelini/entity"
)

// InitData bool param for test data
var InitData = true

// CreateTablesIfNotExist   database initializing and adding test data.
func CreateTablesIfNotExist(db *sql.DB) error {

	createTablesQuery := `
	CREATE TABLE IF NOT EXISTS player
	(
	   id        SERIAL PRIMARY KEY,
	   first_name VARCHAR(30),
	   points     BIGINT DEFAULT 0
	);
  
    CREATE TABLE IF NOT EXISTS tournament
	(
	   id      SERIAL NOT NULL PRIMARY KEY,
	   deposit BIGINT NOT NULL DEFAULT 0,
	   prize   BIGINT NOT NULL DEFAULT 0,
	   winner   INT DEFAULT 0,
	   status   INT DEFAULT 0
	);
  
    CREATE TABLE IF NOT EXISTS tournament_player
	(
	   player_id     INT REFERENCES player (id) ON UPDATE CASCADE ON DELETE
	   CASCADE,
	   tournament_id INT REFERENCES tournament (id) ON UPDATE CASCADE,
	   CONSTRAINT tournament_player_pkey PRIMARY KEY (player_id, tournament_id)
	);
	`
	if InitData == true {
		addUserQuery := `
		INSERT INTO player (first_name, points) VALUES ('testuser', 0);
		INSERT INTO player (first_name, points) VALUES ('testuser2', 0);
		`
		createTablesQuery = createTablesQuery + addUserQuery
	}

	_, err := db.Query(createTablesQuery)
	return err
}

// FundPlayer  update user points.
func FundPlayer(db *sql.DB, playerID int, points int64) error {
	id := 0
	err := db.QueryRow("UPDATE player SET points = $1 WHERE id = $2 RETURNING id", points, playerID).Scan(&id)
	return err
}

// AnnounceTournaments  insert new  tournament.
func AnnounceTournaments(db *sql.DB, tournamentID int, deposit int64) error {
	id := 0
	err := db.QueryRow("INSERT INTO tournament (id, deposit)  VALUES($1, $2) RETURNING id", tournamentID, deposit).Scan(&id)
	return err
}

// SelectPlayer select player by id.
func SelectPlayer(db *sql.DB, playerID int) (entity.Player, error) {
	var player entity.Player
	row := db.QueryRow("SELECT * FROM player WHERE id = $1 ", playerID)
	err := row.Scan(&player.ID, &player.FirstName, &player.Points)
	return player, err
}

// SelectTournament select tournament by id.
func SelectTournament(db *sql.DB, tournamentID int) (entity.Tournament, error) {
	var tournament entity.Tournament
	row := db.QueryRow("SELECT * FROM tournament WHERE id = $1 ", tournamentID)
	err := row.Scan(&tournament.ID, &tournament.Deposit, &tournament.Prize, &tournament.Winner, &tournament.Status)
	return tournament, err
}

// SelectTournamentUsers select  tournament players by tournament id.
func SelectTournamentUsers(db *sql.DB, tournamentID int) ([]entity.TournamentPlayer, error) {
	players := make([]entity.TournamentPlayer, 0)
	rows, err := db.Query(`SELECT * FROM tournament_player WHERE tournament_id = $1 `, tournamentID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var player entity.TournamentPlayer
		if err := rows.Scan(&player.PlayerID, &player.TournamentID); err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}

// SelectFinishedTournaments select finished tournaments.
func SelectFinishedTournaments(db *sql.DB) ([]entity.Tournament, error) {
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

// ChangeTournamentsPrize update  tournament  prize.
func ChangeTournamentsPrize(db *sql.DB, tournamentID int, prize int64) error {
	id := 0
	err := db.QueryRow(`UPDATE tournament SET  prize = $1   WHERE id = $2 RETURNING id`, prize, tournamentID).Scan(&id)
	return err
}

// InsertUserIntoTournament  insert user and tournament into  tournament_player table.
func InsertUserIntoTournament(db *sql.DB, tournamentID int, playerID int) error {
	player_id := 0
	err := db.QueryRow("INSERT INTO tournament_player (player_id , tournament_id ) VALUES( $1 ,$2 )  RETURNING player_id", playerID, tournamentID).Scan(&player_id)
	return err
}

// FinishTournament update tournament status.
func FinishTournament(db *sql.DB, tournamentID int, playerID int) error {
	id := 0
	err := db.QueryRow("UPDATE tournament SET  winner = $1, status = $2 WHERE id = $3  RETURNING id", playerID, entity.TournamentIsFinished, tournamentID).Scan(&id)
	return err
}
