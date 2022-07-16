package websocket

import (
	"chess-backend/cancellableTimer"
	"chess-backend/data"
	"chess-backend/syncslice"
	"chess-backend/websocket/moveHelper"
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"chess-backend/websocketModels/modelErrors"
	"chess-backend/websocketModels/pieces"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func createNewGame(c websocketModels.CreateNewGameParams, p *websocketModels.Player) *websocketModels.Game {
	// Set a random color if its not set
	if c.Color == "" {
		rand := rand.Intn(2)
		switch rand {
		case 0:
			c.Color = "w"
		case 1:
			c.Color = "b"
		}
	}

	id := uuid.New().String()

	var white *websocketModels.PlayerState
	var black *websocketModels.PlayerState

	switch c.Color {
	case "w":
		white = &websocketModels.PlayerState{Player: p, TimeLeft: time.Duration(c.TimeControl) * time.Minute}
		black = &websocketModels.PlayerState{TimeLeft: time.Duration(c.TimeControl) * time.Minute}
	case "b":
		white = &websocketModels.PlayerState{TimeLeft: time.Duration(c.TimeControl) * time.Minute}
		black = &websocketModels.PlayerState{Player: p, TimeLeft: time.Duration(c.TimeControl) * time.Minute}
	}

	timeControl := time.Duration(c.TimeControl) * time.Minute
	addTimePerMove := time.Duration(c.Adder) * time.Second
	gameState := websocketModels.WaitsOpponent
	board := websocketModels.StringToBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")

	boardState := websocketModels.BoardState{Board: board, CastlingRights: []pieces.CastleSide{"K", "Q", "k", "q"}}
	commandChan := make(chan gameCommands.GameCommand)
	observers := syncslice.NewSyncSlice([]*websocketModels.Player{p})
	countDownTimer := cancellableTimer.NewCancellableTimer()

	g := &websocketModels.Game{
		Id:             id,
		TimeControl:    timeControl,
		AddTimePerMove: addTimePerMove,
		GameState:      gameState,
		BoardState:     boardState,
		CommandChan:    commandChan,
		Observers:      observers,
		Black:          black,
		White:          white,
		CountDownTimer: countDownTimer,
	}

	newGame(g)

	p.ObserverOf = append(p.ObserverOf, g)

	return g
}

func connectToGame(c websocketModels.ConnectToGameParams, p *websocketModels.Player) (*websocketModels.Game, *websocketModels.PlayerState, error) {

	var ps *websocketModels.PlayerState
	g, err := data.GetGame(c.GameId)

	if err != nil {
		return nil, ps, err
	}

	g.Observers.Add(p)

	whiteState := g.White
	blackState := g.Black

	// If player trying to reconnect
	switch {
	case whiteState.Player != nil && whiteState.Player.Nick == p.Nick:
		whiteState.Player = p
		ps = g.White
	case blackState.Player != nil && blackState.Player.Nick == p.Nick:
		blackState.Player = p
		ps = g.Black
	default:
	}

	// Whoever connects to the game first sets as player
	if g.GameState == websocketModels.WaitsOpponent {
		g.GameState = websocketModels.WaitsFirstMove
		switch {
		case whiteState.Player == nil:
			whiteState.Player = p
			ps = whiteState
		case blackState.Player == nil:
			blackState.Player = p
			ps = blackState
		default:
		}
		now := time.Now()
		g.CreatedAt = &now
		g.CountDownTimer.Start(30*time.Second, g)
	}

	return g, ps, nil
}

func playMove(g *websocketModels.Game, params *websocketModels.PlayMoveParams, now *time.Time) {

	// Active player changes in [BoardState.Play(...)]
	whoPlays := g.WhoPlays()

	play(&g.BoardState, params)

	latency := whoPlays.Player.Latency
	adder := g.AddTimePerMove

	// If game is in WaitingFirstMove state, set started time
	if g.GameState == websocketModels.WaitsFirstMove {
		// Both players had to play their first move
		if g.Black.Player.Nick == whoPlays.Player.Nick {
			g.GameState = websocketModels.Playing
		}

		// If its first move there is no need for latency balance, adder or delta
		g.LastPlayed = now
		latency = 0
		adder = 0
	}

	// How much time has the player lost
	delta := now.Sub(*g.LastPlayed)

	// Set new lastplayed value
	g.LastPlayed = now

	log.Printf("Old time: %d, latency: %d, adder: %d, delta: %d", whoPlays.TimeLeft.Milliseconds(), latency.Milliseconds(), adder.Milliseconds(), delta.Milliseconds())
	whoPlays.TimeLeft = whoPlays.TimeLeft + latency + adder - delta
	log.Printf("New time: %d", whoPlays.TimeLeft.Milliseconds())
}

func play(b *websocketModels.BoardState, params *websocketModels.PlayMoveParams) {
	// board := b.Board

	source := pieces.NewChessCoordFromString(params.Source)

	target := pieces.NewChessCoordFromString(params.Target)

	mH := moveHelper.MoveHelper{
		Board:          b.Board,
		Source:         source,
		Target:         target,
		CastlingRights: b.CastlingRights,
		PlayerColor:    b.ActiveColor,
		EnPassant:      b.EnPassant,
		Promote:        params.Promote,
	}

	mH.MakeMove(false)

	// Change active color as opposite
	b.ActiveColor = pieces.Color((int64(b.ActiveColor) + 1) % 2)
	b.EnPassant = mH.EnPassant
	b.CastlingRights = mH.CastlingRights
	return
}

func checkBoardStatePlayable(b *websocketModels.BoardState, params *websocketModels.PlayMoveParams) error {
	source := pieces.NewChessCoordFromString(params.Source)

	target := pieces.NewChessCoordFromString(params.Target)

	copyBoard := websocketModels.StringToBoard(websocketModels.BoardToString(b.Board))

	piece := copyBoard[source.Row][source.Column]

	if piece == nil {
		return errors.New("No piece at source")
	}

	ml := piece.CalculateMovableCoords(copyBoard)

	switch piece.(type) {
	case *pieces.PawnPiece:
		piece.(*pieces.PawnPiece).CalculateEnPassantLocations(ml, b.EnPassant)
	case *pieces.KingPiece:
		piece.(*pieces.KingPiece).CalculateCastlableLocations(copyBoard, ml, b.CastlingRights)
	}

	if !ml[target.Row][target.Column] {
		return errors.New("Move is not valid")
	}

	mH := moveHelper.MoveHelper{
		Board:          copyBoard,
		Source:         source,
		Target:         target,
		CastlingRights: b.CastlingRights,
		PlayerColor:    b.ActiveColor,
		EnPassant:      b.EnPassant,
		Promote:        params.Promote,
	}

	mH.MakeMove(true)

	if moveHelper.IsKingChecked(copyBoard, b.ActiveColor) {
		return errors.New("Move is not valid")
	}
	return nil
}

// We will gonna simulate everything in there
func checkGamePlayable(g *websocketModels.Game, params *websocketModels.PlayMoveParams, now *time.Time, player *websocketModels.Player) error {

	whoPlays := g.WhoPlays()

	if whoPlays.Player.Nick != player.Nick {
		return &modelErrors.IllegalTurn{}
	}

	// Locks the timer prevents it gets firing when processing the move
	// We call it there because nobody shouldn't lock it except the player
	if err := g.CountDownTimer.Lock(); err != nil {
		return &modelErrors.GameEnderLockTriggeredError{}
	}

	// Check game state is Playing or WaitsFirstMove
	if g.GameState != websocketModels.Playing && g.GameState != websocketModels.WaitsFirstMove {
		return &modelErrors.NotPlayableError{}
	}

	// Check move
	if err := checkBoardStatePlayable(&g.BoardState, params); err != nil {
		return err
	}

	latency := whoPlays.Player.Latency
	lastPlayed := g.LastPlayed

	if g.GameState == websocketModels.WaitsFirstMove {
		// If its first move there is no need for latency balance or delta
		lastPlayed = now
		latency = 0
	}

	delta := now.Sub(*lastPlayed)

	timeLeft := whoPlays.TimeLeft - delta + latency

	if (timeLeft / time.Millisecond) <= 0 {
		return &modelErrors.TimeOverError{}
	}

	return nil
}
