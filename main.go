package main

import (
	// "fmt"
	"math"
	// "math/rand"
	"chessAI/chess"
	"log"
	"os"
	"strconv"
	"runtime"
	"sort"
	
)

const(
	posInf = 10000000.0
	negInf = -10000000.0
)


func logCurrentCaller() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func str(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func logger(filePath string, description string, message string) error {
	// Open or create a log file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)
	log.SetPrefix("|" + logCurrentCaller() + "| ")
	// Write the message to the log file
	log.Println(" " + description + " {" + message + "}")

	return nil
}









































func evaluation(p *chess.Position) float64{

	eval := p.Board().Material()
	// logger("logs/control.log", "evaluations", str(eval))
	//temp eval function
	return eval
}




type positionMoves struct{
	pos *chess.Position
	moves []*chess.Move
}

func (a *positionMoves) Len() int           { return len(a.moves) }
func (a *positionMoves) Swap(i, j int)      { a.moves[i], a.moves[j] = a.moves[j], a.moves[i] }
func (a *positionMoves) Less(i, j int) bool { return exchangeValue(a.pos, a.moves[i]) > exchangeValue(a.pos, a.moves[j])}



func exchangeValue(pos *chess.Position, move *chess.Move) float64{
	value := math.Abs(chess.SpecificPieceValue(pos.Board().Piece(move.S2()))) - math.Abs(chess.SpecificPieceValue(pos.Board().Piece(move.S1())))
	return value
}


func generateChildren(pos *chess.Position) []*chess.Move{
	//maybe generate the move and the heuristic at the same time, maybe use less sorting
	moves := pos.ValidMoves()
	
	sortableBind := positionMoves{pos, moves}
	sort.Sort(&sortableBind)

	return sortableBind.moves
	
}



func qSearch(node *chess.Position, maximizingPlayer bool)float64{

	return 0.0
}























type Bound struct {
    
    upperBound float64
    lowerBound float64
    
}

type TranspositionTable struct{
	tt map[[16]byte]Bound
}

func createTable() *TranspositionTable {

	tt := TranspositionTable{}
	tt.tt = make(map[[16]byte]Bound)

	return &tt
}

func (table *TranspositionTable) cachedZeroWindow(root *chess.Position, alpha float64, beta float64, depth int, maximizingPlayer bool) float64{
	
	//search for the cached move
	cachedMove, exists := table.tt[root.Hash()]
	if  exists {

		logger("logs/control.log", "cached bounds (upper , lower) ", str(cachedMove.upperBound) + " , " + str(cachedMove.lowerBound))
		if cachedMove.lowerBound >= beta {return cachedMove.lowerBound}
		if cachedMove.upperBound <= alpha {return cachedMove.upperBound}
		alpha = math.Max(alpha, cachedMove.lowerBound)
		beta = math.Min(beta, cachedMove.upperBound)
	}

	//evaluations
	if depth == 0 {
		//cache this too
		return evaluation(root)
	}

	//initialize gamma
	gamma := 0.0
	if maximizingPlayer {
		gamma = negInf
		a := alpha


		for _, child := range generateChildren(root) {

			gamma = math.Max(gamma, table.cachedZeroWindow(root.Update(child), a, beta, depth - 1, false))
			if gamma >= beta {
				break
			}
			a = math.Max(a, gamma)

		}
	}else{
		gamma = posInf
		b := beta


		for _, child := range generateChildren(root) {

			gamma = math.Min(gamma, table.cachedZeroWindow(root.Update(child), alpha, b, depth - 1, true))
			if gamma <= alpha {
				break
			}
			b = math.Min(b, gamma)

		}

	}


	// failing low
	bound, exists := table.tt[root.Hash()]

	if !exists {
		bound = Bound {
			upperBound: posInf,
			lowerBound: negInf,
		}
	}

	if gamma <= alpha {
		logger("logs/control.log", "fail low ", str(gamma))
		bound.lowerBound = gamma
		table.tt[root.Hash()] = bound
	}
	

	// failing high
	if gamma >= beta {
		logger("logs/control.log", "fail high", str(gamma))
		bound.upperBound = gamma
		table.tt[root.Hash()] = bound
	} 

	return gamma
	
	
}



func MDTF(root *chess.Position, f float64, depth int, maximizingPlayer bool) float64{

	
	gamma := f
	upperBound := posInf
	lowerBound := negInf



	for lowerBound < upperBound {
		beta := math.Max(float64(gamma), lowerBound + 1)
		
		logger("logs/control.log", "zero window beta", str(beta))

		transpositionTable := createTable()
		gamma = transpositionTable.cachedZeroWindow(root, beta - 1, beta, depth, maximizingPlayer)
		logger("logs/control.log", "zero window gamma", str(gamma))
		if gamma < beta {
			upperBound = gamma

		}else{
			lowerBound = gamma

		}

		logger("logs/control.log", "Upperbound-Lowerbound", str(upperBound) + " , " + str(lowerBound))
	}

	logger("logs/control.log", "Upperbound-Lowerbound", str(upperBound) + " , " + str(lowerBound))

	return gamma
	

}







func test(stockfishEval float64, fenStr string){


	logger("logs/test.log", "START TEST", "---------------")
	fen, _ := chess.FEN(fenStr)
	game := chess.NewGame(fen)

	
	maximizingPlayer := true
	if game.Position().Turn() == 2 {
		maximizingPlayer = false
	}

	score := MDTF(game.Position(), 0, 5, false)
	logger("logs/test.log", "Game: ", game.Position().Board().Draw())
	logger("logs/test.log", "MAX player ", strconv.FormatBool(maximizingPlayer))
	logger("logs/test.log", "Score ", str(score))
	logger("logs/test.log", "Stockfish Score ", str(stockfishEval))
	logger("logs/test.log", "END TEST", "---------------")
}



func main() {


	logger("logs/control.log", "Main", "START")
	// game := chess.NewGame()
	
	// i := 0
	// for i < 15 {
	// 	// logger("logs/control.log", "Main", "random move: " + strconv.Itoa(i))
	// 	// select a random move
	// 	moves := game.ValidMoves()
	// 	move := moves[rand.Intn(len(moves))]
	// 	game.Move(move)

	// 	i = i + 1
	// }
		
	// score := MDTF(game.Position(), 0, 5, false)

	// fmt.Println(score)
	// fmt.Println(game.Position().Board().Draw())
	// // fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
	// fmt.Println(game.String())



	// test(1.4, "r2q1rk1/3nbppp/p2pbn2/1p2p1P1/4P3/1NN1BP2/PPPQ3P/2KR1B1R b - - 0 12")

	test(9.1, "5rk1/1pq2rp1/p1pR4/P3p1p1/4p1Q1/4B3/3R1PP1/6K1 w - - 4 41")


	test(-1.0, "r4r2/3nppkp/3p2p1/qppP4/1P2P3/2N5/P2Q1PPP/1R3RK1 b - - 0 19")

	test(4.5, "r4r1k/1p2Q1pp/p1p4n/P3p3/4p1B1/2q1B2P/5PP1/RR4K1 b - - 2 29")

	logger("logs/control.log", "Main", "HALT")
	
}