package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var board [6][7]string
var turn = "R"
var winner string

// renderBoard génère le HTML du plateau
func renderBoard() template.HTML {
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
	return template.HTML(html)
}

func resetBoard() { // Réinitialise le plateau
	board = [6][7]string{}
	turn = "R"
	winner = ""
}

func homeHandler(w http.ResponseWriter, r *http.Request) { // Handler de l'accueil
	log.Printf("%s requested %s", r.RemoteAddr, r.URL.Path)
	if winner != "" || isBoardFull() {
		resetBoard()
	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Erreur chargement home.html", 500)
		log.Println("Erreur template home:", err)
		return
	}
	tmpl.Execute(w, nil)
}

// Handler du plateau
func gameHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s requested %s", r.RemoteAddr, r.URL.Path)
	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		http.Error(w, "Erreur chargement game.html", 500)
		log.Println("Erreur template game:", err)
		return
	}
	data := struct {
		BoardHTML template.HTML
	}{BoardHTML: renderBoard()}
	tmpl.Execute(w, data)
}

// Vérifie victoire du joueur
func checkWin(player string) bool {
	// Horizontal
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == player && board[i][j+1] == player && board[i][j+2] == player && board[i][j+3] == player {
				return true
			}
		}
	}
	// Vertical
	for j := 0; j < 7; j++ {
		for i := 0; i < 3; i++ {
			if board[i][j] == player && board[i+1][j] == player && board[i+2][j] == player && board[i+3][j] == player {
				return true
			}
		}
	}
	// Diagonale \
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == player && board[i+1][j+1] == player && board[i+2][j+2] == player && board[i+3][j+3] == player {
				return true
			}
		}
	}
	// Diagonale /
	for i := 3; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == player && board[i-1][j+1] == player && board[i-2][j+2] == player && board[i-3][j+3] == player {
				return true
			}
		}
	}
	return false
}

// Vérifie si le plateau est plein
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

// Handler pour jouer un coup
func playHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s requested %s", r.RemoteAddr, r.URL.Path)

	if r.Method != http.MethodPost || winner != "" {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	c, err := strconv.Atoi(r.FormValue("col"))
	if err != nil || c < 0 || c > 6 {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	// Placer jeton
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

	// Vérification victoire
	if checkWin(turn) {
		if turn == "R" {
			winner = "Le joueur Rouge a gagné ! 🎉"
		} else {
			winner = "Le joueur Jaune a gagné ! 🎉"
		}
	} else if isBoardFull() {
		winner = "Match nul : la grille est remplie 🎯"
	} else {
		// Changer de tour
		if turn == "R" {
			turn = "J"
		} else {
			turn = "R"
		}
	}

	http.Redirect(w, r, "/game", http.StatusSeeOther)
}
func main() {
	// Fichiers statiques
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", playHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	fmt.Printf("Serveur lancé sur http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Erreur serveur:", err)
	}
}
