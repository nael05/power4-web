package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var board [6][7]string
var turn = "R"

// Génère le HTML pour le plateau avec flèches
func renderBoard() string {
	html := "<table>"

	// Ligne des flèches
	html += "<tr>"
	for j := 0; j < 7; j++ {
		html += fmt.Sprintf(
			"<td><form action='/play' method='POST'><button class='arrow' type='submit' name='col' value='%d'>↓</button></form></td>",
			j,
		)
	}
	html += "</tr>"

	// Plateau avec pions
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
	html += fmt.Sprintf("<p>Tour de : %s</p>", turn)
	return html
}

// Page d'accueil
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Erreur chargement home.html", 500)
		fmt.Println("Erreur:", err)
		return
	}
	tmpl.Execute(w, nil)
}

// Page du jeu
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

// Jouer un coup
func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		col := r.FormValue("col")
		var c int
		fmt.Sscanf(col, "%d", &c)

		// Cherche la première case libre de la colonne
		for i := 5; i >= 0; i-- {
			if board[i][c] == "" {
				board[i][c] = turn
				if turn == "R" {
					turn = "J"
				} else {
					turn = "R"
				}
				break
			}
		}
	}
	http.Redirect(w, r, "/game", http.StatusSeeOther)
}

func main() {
	// Servir CSS statique
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", playHandler)

	fmt.Println("Serveur lancé sur http://localhost:4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Erreur serveur:", err)
	}
}
