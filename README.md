Puissance 4 projet Go

Ce repo contient un petit jeu "Puissance 4" écrit en Go. C'est un projet simple qui utilise des templates HTML et un CSS statique.

Structure du projet

- `main.go` : serveur HTTP, logique du jeu (plateau en mémoire).
- `templates/` : templates HTML (`home.html`, `game.html`).
- `style/` : fichiers CSS servis statiquement.
- `static/` : GIF de fond

Jeu et regle

- Le plateau est une grille 6x7
- Chaque joueur pose un pion chacun son tour
- Après chaque coup, le programme vérifi si un joueur a gagné (4 jetons alignés) ou si la grille est pleine.
- et renvoie les informations en conséquence (tour de..., gagant ou match nul)