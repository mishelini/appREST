package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/mishelini/database"
	"github.com/mishelini/entity"
)

// FundPlayer convert player points from float64  to int64 ,
// and set parameters to database layer.
func FundPlayer(db *sql.DB, id int, points float64) error {
	return database.FundPlayer(db, id, int64(points*100))
}

// AnnounceTournament  convert tournament deposit from float64  to int64 ,
// and set parameters to database layer.
func AnnounceTournament(db *sql.DB, id int, deposit float64) error {
	return database.AnnounceTournaments(db, id, int64(deposit*100))
}

// JoinTournament checks enough points for the user to participate in the tournament adds user to the tournament
// and set parameters to database layer.
func JoinTournament(db *sql.DB, userID int, tournamentID int) error {
	tournamentData, err := database.SelectTournament(db, tournamentID)
	if err != nil {
		return err
	}
	userData, err := database.SelectPlayer(db, userID)
	if err != nil {
		return err
	}

	if userData.Points < tournamentData.Deposit {
		return fmt.Errorf("user " + string(userID) + " does not have enough points")
	}
	if tournamentData.Status == entity.TournamentIsFinished {
		return fmt.Errorf("tournment is closed")
	}
	newUserPoints := userData.Points - tournamentData.Deposit
	newTormentPrize := tournamentData.Deposit + tournamentData.Prize
	err = database.ChangeTournamentsPrize(db, tournamentID, newTormentPrize)
	if err != nil {
		return err
	}
	err = database.FundPlayer(db, userID, newUserPoints)
	if err != nil {
		return err
	}
	err = database.InsertUserIntoTournament(db, tournamentID, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetFinishedTournamentSet  get list of finished tournaments from database layer
// convert tournament prize and player points  from int64  to float64.
func GetFinishedTournamentSet(db *sql.DB) ([]byte, error) {
	tournaments, err := database.SelectFinishedTournaments(db)
	if err != nil {
		return nil, err
	}
	if len(tournaments) == 0 {
		return nil, fmt.Errorf("tournaments have not yet been created")
	}
	winnersSet := make([]entity.Winner, 0)
	for i := range tournaments {
		tournament := tournaments[i]
		winnerUserID := tournament.Winner
		player, err := database.SelectPlayer(db, winnerUserID)
		if err != nil {
			return nil, err
		}
		win := entity.Winner{winnerUserID, float64(tournament.Prize / 100), float64(player.Points / 100)}
		winnersSet = append(winnersSet, win)
	}
	res := entity.Results{
		Winners: winnersSet,
	}
	js, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return js, err
}

// FinishTournament checks tournament status and if it is not finished
// randomly chooses the winner and set parameters to database layer.
func FinishTournament(db *sql.DB, tournamentID int) ([]byte, error) {
	var playerBalance int64
	tournament, err := database.SelectTournament(db, tournamentID)
	if err != nil {
		return nil, err
	}
	if tournament.Status == entity.TournamentIsFinished {
		return nil, fmt.Errorf("tournment is finished")
	}
	winnerID := tournament.Winner
	if winnerID != 0 {
		winnerPlayer, err := database.SelectPlayer(db, winnerID)
		if err != nil {
			return nil, err
		}
		playerBalance = winnerPlayer.Points
		tournamentPlayerSet, err := database.SelectTournamentUsers(db, tournamentID)
		if err != nil {
			return nil, err
		}
		if len(tournamentPlayerSet) > 0 {
			var totalUsers []int
			for i := 0; i < len(tournamentPlayerSet); i++ {
				playerID := tournamentPlayerSet[i].PlayerID
				totalUsers = append(totalUsers, playerID)
			}
			rand.Seed(time.Now().Unix())
			winnerID = totalUsers[rand.Intn(len(totalUsers))]

			winnerPlayer, err := database.SelectPlayer(db, winnerID)
			if err != nil {
				return nil, err
			}
			playerBalance = tournament.Prize + winnerPlayer.Points
			err = database.FundPlayer(db, winnerID, playerBalance)
		}
	}

	err = database.FinishTournament(db, tournamentID, winnerID)
	if err != nil {
		return nil, err
	}
	win := entity.Winner{winnerID, float64(tournament.Prize / 100), float64(playerBalance / 100)}

	res := entity.Result{
		Winner: win,
	}
	js, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// GetUserBalance  get user from database layer and convert user balanse from int64  to float64.
func GetUserBalance(db *sql.DB, id int) ([]byte, error) {
	player, err := database.SelectPlayer(db, id)
	if err != nil {
		return nil, err
	}
	res2 := entity.BalanceResults{
		PlayerId: id,
		Balance:  float64(player.Points / 100),
	}
	js, err := json.Marshal(res2)
	return js, err
}
