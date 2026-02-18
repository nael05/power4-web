# Power 4 - Web Game

A simple, lightweight web implementation of the classic strategy game **Connect 4** (Puissance 4), built with **Go** (Golang) for the backend and **HTML/CSS** for the frontend.

##  Project Overview

This project serves a playable Connect 4 grid where two players can compete locally. The server handles the game logic, including move validation, win detection (horizontal, vertical, diagonal), and draw conditions.

##  Features

- **Classic Gameplay**: 6x7 grid implementation respecting standard rules.
- **Two Player Mode**: Local multiplayer with alternating turns (🔴 Red vs 🟡 Yellow).
- **Game Logic**:
  - Automatic win detection (4 tokens aligned).
  - Draw detection (full board).
  - Invalid move prevention (full columns).
- **Interactive UI**:
  - Hover effects on drop buttons.
  - Visual indicators for the current turn.
  - **Fireworks animation** upon victory.
  - Animated background (GIF/Video support).
- **Reset Functionality**: Ability to restart the game immediately after a match ends.

##  Tech Stack

- **Backend**: Go (Golang) 1.25+
- **Frontend**: HTML5, CSS3
- **Templating**: Go `html/template` package

##  Project Structure

```text
.
├── main.go            # Main entry point: HTTP server and game logic
├── go.mod             # Go module definition
├── templates/         # HTML templates
│   ├── home.html      # Welcome screen
│   └── game.html      # Main game board interface
├── style/             # Static CSS files
│   └── style.css      # Styling for the board and tokens
└── static/            # Static assets (background images/gifs)
