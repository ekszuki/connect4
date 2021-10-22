package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	debug bool = false

	colHeap          map[int]int
	xSize            int = 7
	ySize            int = 6
	maxPeacesOnBoard int = xSize * ySize
	tab              [][]string
	colorsMap        map[string]string
)

func start() {
	colorsMap = make(map[string]string, 2)
	colorsMap["R"] = "RED"
	colorsMap["Y"] = "YELLOW"

	colHeap = make(map[int]int)
	for x := 0; x <= ySize-1; x++ {
		colHeap[x] = 0
	}

	tab = make([][]string, xSize)
	for i := range tab {
		tab[i] = make([]string, ySize)
	}
	printTab()
}

func getTurn(i int) string {
	if i%2 == 0 {
		return colorsMap["R"]
	} else {
		return colorsMap["Y"]
	}
}

func checkInputPos(strPos string) (int, error) {
	p := strPos[0:1]

	if strings.EqualFold(p, "E") {
		syscall.Exit(syscall.F_OK)
	}

	if strings.EqualFold(p, "\n") {
		return -1, fmt.Errorf("please inform a column")
	}

	pos, err := strconv.Atoi(p)
	if err != nil {
		return -1, fmt.Errorf("invalid value: %v", err)
	}

	if pos < 1 || pos > xSize {
		fmt.Print(" \n\n")
		return -1, fmt.Errorf("valid values are 1, 2, 3, 4, 5, 6, 7")
	}

	// return (pos - 1) because slice in go start at 0 position
	return (pos - 1), nil
}

func isFullBoard(piecesOnBoard int) bool {
	return piecesOnBoard >= maxPeacesOnBoard
}

func main() {
	start()
	var color string
	piecesOnBoard := 0
	for {
		color = getTurn(piecesOnBoard)
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("%s's turn  \n\n\n", color)
		fmt.Print("Press [E] To Exit\n\n")
		fmt.Printf("Enter Position '[1-%d]: ", xSize)
		strPos, _ := reader.ReadString('\n')

		pos, err := checkInputPos(strPos)
		if err != nil {
			fmt.Printf("%v \n\n", err)
			printTab()
			continue
		}

		isWinner, err := insert(pos, color[0:1])
		if err != nil {
			continue
		}

		if isWinner {
			fmt.Printf(">>> %s <<< is the winner\n\n", color)
			printTab()
			syscall.Exit(syscall.F_OK)
		}

		piecesOnBoard++
		fmt.Printf("Total of peaces on board: %d  \n\n", piecesOnBoard)

		if isFullBoard(piecesOnBoard) {
			fmt.Println("All positions are used. There is no winner")
			syscall.Exit(syscall.F_OK)
		}
	}
}

func insert(hpos int, color string) (isWinner bool, err error) {
	if hpos < 0 || hpos >= xSize {
		fmt.Println("outside of tab...")
		return false, fmt.Errorf("outside of tab")
	}
	y := colHeap[hpos]

	if y == ySize {
		fmt.Printf("no more space on column %d \n", hpos+1)
		return false, fmt.Errorf("no more space to put in this column")
	}

	tab[hpos][y] = color
	y++
	colHeap[hpos] = y

	printTab()
	return checkWinner(), nil
}

func checkWinner() (winner bool) {
	win := checkHorizontalLine()
	if win {
		fmt.Printf("Winner in horizontal line: ")
		return win
	}

	win = checkVerticalLine()
	if win {
		fmt.Printf("Winner in vertical line: ")
		return win
	}

	win = checkDiagonalLineUp()
	if win {
		fmt.Printf("Winner in diagonal up line: ")
		return win
	}

	win = checkDiagonalLineDown()
	if win {
		fmt.Printf("Winner in diagonal down line: ")
		return win
	}

	return win
}

func checkDiagonalLineUp() (isWinner bool) {
	count := 0
	for x := 0; x <= xSize-1; x++ {
		for y := 0; y <= ySize-1; y++ {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			count++
			npy := y
			for npx := x + 1; npx <= xSize; npx++ {
				npy = npy + 1

				if npy < 0 || npx > xSize-1 {
					count = 0
					continue
				}

				r := checkNextPositionSameColor(npx, npy, value)
				if r {
					count++
					if count == 4 {
						return true
					}
				} else {
					count = 0
					break
				}

			}
		}
	}

	return false
}

func checkDiagonalLineDown() (isWinner bool) {
	count := 0
	for x := 0; x <= xSize-1; x++ {
		for y := 5; y >= ySize-1; y-- {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			count++
			npy := y
			for npx := x + 1; npx <= xSize; npx++ {
				npy = npy - 1

				if debug {
					fmt.Printf("->next value [%d,%d]", npx, npy)
				}

				if npy < 0 || npx > xSize-1 {
					count = 0
					break
				}

				r := checkNextPositionSameColor(npx, npy, value)
				if r {
					count++
					if count == 4 {
						return true
					}
				} else {
					count = 0
					break
				}

			}
		}
	}

	return false
}

func checkVerticalLine() (isWinner bool) {
	count := 0
	for y := 0; y <= ySize-1; y++ {
		for x := 0; x <= xSize-1; x++ {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			if debug {
				fmt.Printf("\n Found: [%d,%d]=%s > ", x, y, value)
			}

			count++
			for npy := y + 1; npy <= ySize; npy++ {
				if npy >= ySize || npy < 0 {
					count = 0
					break
				}

				r := checkNextPositionSameColor(x, npy, value)
				if r {
					count++
					if count == 4 {
						return true
					}
				} else {
					count = 0
					break
				}
			}
		}
	}

	return false
}

func checkHorizontalLine() (isWinner bool) {
	count := 0

	for y := 0; y <= ySize-1; y++ {
		for x := 0; x <= xSize-1; x++ {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			count++
			for npx := x + 1; npx <= xSize; npx++ {

				if npx > 6 || npx < 0 {
					count = 0
					break
				}

				r := checkNextPositionSameColor(npx, y, value)
				if r {
					count++
					if count == 4 {
						return true
					}
				} else {
					count = 0
					break
				}
			}
		}
	}

	return false
}

func checkNextPositionSameColor(x, y int, color string) bool {
	if color == "" {
		return false
	}

	if x < 0 || x > xSize-1 {
		return false
	}

	if y < 0 || y > ySize-1 {
		return false
	}

	value := tab[x][y]
	if debug {
		fmt.Printf("next value -> [%d,%d] = %s", x, y, value)
	}

	return value == color
}

func printTab() {
	fmt.Println("|  1  ||  2  ||  3  ||  4  ||  5  ||  6  ||  7  |")
	fmt.Println("|  |  ||  |  ||  |  ||  |  ||  |  ||  |  ||  |  |")
	fmt.Println("|     ||     ||     ||     ||     ||     ||     |")
	for y := ySize - 1; y >= 0; y-- {
		for x := 0; x <= xSize-1; x++ {

			value := tab[x][y]
			if value == "" {
				if debug {
					value = fmt.Sprintf("%d,%d [%s]", x, y, " ")
				} else {
					value = "   "
				}
			} else {
				if debug {
					value = fmt.Sprintf("%d,%d [%s]", x, y, value)
				} else {
					value = fmt.Sprintf(" %s ", value)
				}
			}

			fmt.Printf("| %s |", value)
		}
		fmt.Printf("\n")
	}
	fmt.Println("|_____||_____||_____||_____||_____||_____||_____|")
	fmt.Printf("\n\n\n")
}

// TODO: make tests using games below
// vertical winner
// insert(0, "r")
// insert(1, "b")
// insert(0, "r")
// insert(0, "b")
// insert(0, "b")
// insert(0, "b")
// insert(0, "b")

// horizontal winner
// insert(0, "b")
// insert(1, "r")
// insert(2, "b")
// insert(3, "r")
// insert(3, "b")
// insert(4, "r")
// insert(5, "r")
// insert(6, "r")

// diagonal up at 0
// insert(0, "r")
// insert(1, "b")
// insert(1, "r")
// insert(2, "r")
// insert(2, "r")
// insert(2, "r")
// insert(3, "b")
// insert(3, "b")
// insert(3, "b")
// insert(3, "r")
// insert(4, "b")
// insert(4, "r")
// insert(6, "r")

// diagonal up  at 1
// insert(1, "r")
// insert(2, "b")
// insert(2, "r")
// insert(3, "r")
// insert(3, "r")
// insert(3, "r")
// insert(4, "b")
// insert(4, "b")
// insert(4, "b")
// insert(4, "r")
// insert(5, "b")
// insert(5, "r")
// insert(6, "r")

// diagonal up  start at 2
// insert(2, "r")
// insert(3, "b")
// insert(3, "r")
// insert(4, "r")
// insert(4, "r")
// insert(4, "r")
// insert(5, "b")
// insert(5, "b")
// insert(5, "b")
// insert(5, "r")
// insert(6, "b")
// insert(6, "r")
// insert(6, "r")

// diagonal up start at 3
// insert(3, "r")
// insert(4, "b")
// insert(4, "r")
// insert(5, "r")
// insert(5, "r")
// insert(5, "r")
// insert(6, "b")
// insert(6, "b")
// insert(6, "b")
// insert(6, "r")
// insert(6, "b")
// insert(6, "r")
// insert(6, "r")

// diagonal down start at 0
// insert(0, "r")
// insert(0, "b")
// insert(0, "r")
// insert(0, "r")
// insert(0, "r")
// insert(1, "r")
// insert(1, "b")
// insert(1, "r")
// insert(1, "r")
// insert(2, "r")
// insert(2, "b")
// insert(2, "r")
// insert(3, "b")
// insert(3, "r")

// diagonal down start at 1
// insert(1, "r")
// insert(1, "b")
// insert(1, "r")
// insert(1, "r")
// insert(1, "r")
// insert(2, "r")
// insert(2, "b")
// insert(2, "r")
// insert(2, "r")
// insert(3, "r")
// insert(3, "b")
// insert(3, "r")
// insert(4, "b")
// insert(4, "r")

// diagonal down start at 2
// insert(2, "r")
// insert(2, "b")
// insert(2, "r")
// insert(2, "r")
// insert(2, "r")
// insert(3, "r")
// insert(3, "b")
// insert(3, "r")
// insert(3, "r")
// insert(4, "r")
// insert(4, "b")
// insert(4, "r")
// insert(5, "b")
// insert(5, "r")

// first line
// insert(1, "r")
// insert(2, "r")
// insert(3, "y")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")

// // second line
// insert(2, "r")
// insert(3, "y")
// insert(4, "y")
// insert(5, "y")
// insert(6, "r")

// // 3rd line
// insert(2, "y")
// insert(3, "r")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "y")
// insert(1, "r")
// insert(2, "y")
// insert(3, "r")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")

// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// winner using diagonal down
// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "y")
// insert(1, "r")
// insert(2, "y")
// insert(3, "r")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")

// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "y")
// insert(1, "r")
// insert(2, "y")
// insert(3, "r")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")

// no win with full tab
// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "y")
// insert(1, "r")
// insert(2, "y")
// insert(3, "r")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")

// insert(0, "y")
// insert(1, "r")
// insert(2, "y")
// insert(3, "r")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")

// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "r")
// insert(1, "y")
// insert(2, "r")
// insert(3, "y")
// insert(4, "r")
// insert(5, "y")
// insert(6, "r")

// insert(0, "y")
// insert(1, "r")
// insert(2, "y")
// insert(3, "r")
// insert(4, "y")
// insert(5, "r")
// insert(6, "y")
