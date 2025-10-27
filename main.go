package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// État du jeu (simple, en mémoire)
var (
	board  [6][7]string
	turn   = "R"
	winner string
)

// On charge les templates au démarrage (simplifie les handlers)
var templates = template.Must(template.ParseFiles("templates/home.html", "templates/game.html"))

// renderBoard construit le HTML du plateau (identique à l'original)
func renderBoard() template.HTML {
	var sb strings.Builder
	sb.WriteString("<table>")

	// ligne des flèches (jouer)
	sb.WriteString("<tr>")
	for c := 0; c < 7; c++ {
		disabled := winner != "" || board[0][c] != ""
		if disabled {
			sb.WriteString("<td></td>")
		} else {
			sb.WriteString(fmt.Sprintf("<td><form action='/play' method='POST'><button class='arrow' name='col' value='%d'>↓</button></form></td>", c))
		}
	}
	sb.WriteString("</tr>")

	// cases
	for r := 0; r < 6; r++ {
		sb.WriteString("<tr>")
		for c := 0; c < 7; c++ {
			cls := ""
			switch board[r][c] {
			case "R":
				cls = "red"
			case "J":
				cls = "yellow"
			}
			sb.WriteString(fmt.Sprintf("<td><div class='cell %s'></div></td>", cls))
		}
		sb.WriteString("</tr>")
	}
	sb.WriteString("</table>")

	if winner != "" {
		sb.WriteString(fmt.Sprintf("<h2>%s</h2>", winner))
		sb.WriteString(`<form action="/" method="GET"><button>Rejouer</button></form>`)
	} else {
		sb.WriteString(fmt.Sprintf("<p>Tour de : %s</p>", turn))
	}
	return template.HTML(sb.String())
}

func resetBoard() {
	board = [6][7]string{}
	turn = "R"
	winner = ""
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Si une partie est terminée, on réinitialise avant d'afficher l'accueil
	if winner != "" || isBoardFull() {
		resetBoard()
	}
	templates.ExecuteTemplate(w, "home.html", nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	data := struct{ BoardHTML template.HTML }{BoardHTML: renderBoard()}
	templates.ExecuteTemplate(w, "game.html", data)
}

// checkWin vérifie les 4 en ligne horizontaux, verticaux et diagonaux
func checkWin(p string) bool {
	// horizontal
	for r := 0; r < 6; r++ {
		for c := 0; c <= 3; c++ {
			if board[r][c] == p && board[r][c+1] == p && board[r][c+2] == p && board[r][c+3] == p {
				return true
			}
		}
	}
	// vertical
	for c := 0; c < 7; c++ {
		for r := 0; r <= 2; r++ {
			if board[r][c] == p && board[r+1][c] == p && board[r+2][c] == p && board[r+3][c] == p {
				return true
			}
		}
	}
	// diagonal \ and /
	for r := 0; r <= 2; r++ {
		for c := 0; c <= 3; c++ {
			if board[r][c] == p && board[r+1][c+1] == p && board[r+2][c+2] == p && board[r+3][c+3] == p {
				return true
			}
		}
	}
	for r := 3; r < 6; r++ {
		for c := 0; c <= 3; c++ {
			if board[r][c] == p && board[r-1][c+1] == p && board[r-2][c+2] == p && board[r-3][c+3] == p {
				return true
			}
		}
	}
	return false
}

func isBoardFull() bool {
	for r := 0; r < 6; r++ {
		for c := 0; c < 7; c++ {
			if board[r][c] == "" {
				return false
			}
		}
	}
	return true
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || winner != "" {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
	colStr := r.FormValue("col")
	col, err := strconv.Atoi(colStr)
	if err != nil || col < 0 || col > 6 {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
	// pose du jeton
	placed := false
	for r := 5; r >= 0; r-- {
		if board[r][col] == "" {
			board[r][col] = turn
			placed = true
			break
		}
	}
	if !placed {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
	// vérifications
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
	// fichiers statiques
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", playHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	fmt.Printf("Serveur lancé sur http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
