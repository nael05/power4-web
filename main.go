package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var board [6][7]string
var turn = "R"
var winner string

func renderBoard() string {
	html := "<table>"
	html += "<tr>"
	for j := 0; j < 7; j++ {
		disabled := winner != "" || board[0][j] != ""
		if !disabled {
			html += fmt.Sprintf(
				"<td><form action='/play' method='POST'><button class='arrow' type='submit' name='col' value='%d'>↓</button></form></td>",
				j,
			)
		} else {
			html += "<td></td>"
		}
	}
	html += "</tr>"

	for i := 0; i < 6; i++ {
		html += "<tr>"
		for j := 0; j < 7; j++ {
			cell := board[i][j]
			class := "empty"
			if cell == "R" {
				class = "red"
			} else if cell == "J" {
				class = "yellow"
			}
			html += fmt.Sprintf("<td><div class='cell %s'></div></td>", class)
		}
		html += "</tr>"
	}
	html += "</table>"

	if winner != "" {
		html += fmt.Sprintf("<h2>%s</h2>", winner)
		html += `<form action="/" method="GET"><button>Rejouer</button></form>`
	} else {
		html += fmt.Sprintf("<p>Tour de : %s</p>", turn)
	}
	return html
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if winner != "" || isBoardFull() {
		// réinitialiser seulement après victoire ou match nul
		board = [6][7]string{}
		turn = "R"
		winner = ""
	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Erreur chargement home.html", 500)
		fmt.Println("Erreur:", err)
		return
	}
	tmpl.Execute(w, nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		http.Error(w, "Erreur chargement game.html", 500)
		fmt.Println("Erreur:", err)
		return
	}
	data := struct {
		BoardHTML template.HTML
	}{BoardHTML: template.HTML(renderBoard())}
	tmpl.Execute(w, data)
}

func checkWin(player string) bool {
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == player && board[i][j+1] == player && board[i][j+2] == player && board[i][j+3] == player {
				return true
			}
		}
	}
	for j := 0; j < 7; j++ {
		for i := 0; i < 3; i++ {
			if board[i][j] == player && board[i+1][j] == player && board[i+2][j] == player && board[i+3][j] == player {
				return true
			}
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == player && board[i+1][j+1] == player && board[i+2][j+2] == player && board[i+3][j+3] == player {
				return true
			}
		}
	}
	for i := 3; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == player && board[i-1][j+1] == player && board[i-2][j+2] == player && board[i-3][j+3] == player {
				return true
			}
		}
	}
	return false
}

func isBoardFull() bool {
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || winner != "" {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	col := r.FormValue("col")
	var c int
	fmt.Sscanf(col, "%d", &c)

	placed := false
	for i := 5; i >= 0; i-- {
		if board[i][c] == "" {
			board[i][c] = turn
			placed = true
			break
		}
	}

	if !placed {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	if checkWin(turn) {
		if turn == "R" {
			winner = "Le joueur Rouge a gagné ! 🎉"
		} else {
			winner = "Le joueur Jaune a gagné ! 🎉"
		}
	} else if isBoardFull() {
		winner = "Match nul : la grille est remplie 🎯"
	} else {
		if turn == "R" {
			turn = "J"
		} else {
			turn = "R"
		}
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func main() {
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", playHandler)
	fmt.Println("Serveur lancé sur http://localhost:4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Erreur serveur:", err)
	}
}
