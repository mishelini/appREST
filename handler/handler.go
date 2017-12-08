package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mishelini/controller"
)

var db *sql.DB

// Handler returns router mux
func Handler(db2 *sql.DB) *mux.Router {
	db = db2
	route := mux.NewRouter()
	route.HandleFunc("/fund", fundPlayerHandler).Queries("playerId", "{playerId:[0-9]+}", "points", "{points:[0-9]+}").Methods("GET")
	route.HandleFunc("/announceTournament", announceTournamentHandler).Queries("tournamentId", "{tournamentId:[0-9]+}", "deposit", "{deposit:[0-9]+}").Methods("GET")
	route.HandleFunc("/joinTournament", joinTournamentHandler).Queries("playerId", "{playerId:[0-9]+}", "tournamentId", "{tournamentId:[0-9]+}").Methods("GET")
	route.HandleFunc("/finishTournament", finishTournamentHandler).Queries("tournamentId", "{tournamentId:[0-9]+}").Methods("GET")
	route.HandleFunc("/resultTournament", resultTournamentHandler).Methods("GET")
	route.HandleFunc("/balance", playerBalanceHandler).Queries("playerId", "{playerId:[0-9]+}").Methods("GET")
	return route
}

func fundPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["playerId"])
	if err != nil {
		http.Error(w, "there was a missing or invalid playerId parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	point, err := strconv.ParseFloat(vars["points"], 64)
	if err != nil {
		http.Error(w, "there was a missing or invalid points parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = controller.FundPlayer(db, id, point)
	if err != nil {
		log.Println(err)
		http.Error(w, "this is Database Error", http.StatusInternalServerError)
	}
}

func announceTournamentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["tournamentId"])
	if err != nil {
		http.Error(w, "there was a missing or invalid tournamentId parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	deposit, err := strconv.ParseFloat(vars["deposit"], 64)
	if err != nil {
		http.Error(w, "there was a missing or invalid deposit parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = controller.AnnounceTournament(db, id, deposit)
	if err != nil {
		http.Error(w, "this is Database Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func joinTournamentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["playerId"])
	if err != nil {
		http.Error(w, "there was a missing or  invalid playerId  parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	tournamentID, err := strconv.Atoi(vars["tournamentId"])
	if err != nil {
		http.Error(w, "there was a missing or  invalid tournamentId  parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = controller.JoinTournament(db, userID, tournamentID)
	if err != nil {
		http.Error(w, "this is Database Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func resultTournamentHandler(w http.ResponseWriter, r *http.Request) {
	js, err := controller.GetFinishedTournamentSet(db)
	if err != nil {
		http.Error(w, "this is Database Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Fprintf(w, string(js))

}

func finishTournamentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tournamentID, err := strconv.Atoi(vars["tournamentId"])
	if err != nil {
		http.Error(w, "there was a missing or  invalid tournamentId  parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	js, err := controller.FinishTournament(db, tournamentID)
	if err != nil {
		http.Error(w, "this is Database Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Fprintf(w, string(js))
}

func playerBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["playerId"])
	if err != nil {
		http.Error(w, "there was a missing or invalid playerId parameter..", http.StatusBadRequest)
		log.Println(err)
		return
	}
	js, err := controller.GetUserBalance(db, id)
	if err != nil {
		http.Error(w, "there was a missing or  invalid  parameters from DB..", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Fprintf(w, string(js))
}
