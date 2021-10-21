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

	ypos map[int]int
	tab  [][]string
)

func start() {
	ypos = make(map[int]int)
	ypos[0] = 0
	ypos[1] = 0
	ypos[2] = 0
	ypos[3] = 0
	ypos[4] = 0
	ypos[5] = 0
	ypos[6] = 0

	tab = make([][]string, 7)
	for i := range tab {
		tab[i] = make([]string, 6)
	}
	printTab()
}

func main() {
	start()
	var color string
	i := 0
	for {
		if i%2 == 0 {
			color = "Red"
		} else {
			color = "Yellow"
		}
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("%s's turn  \n\n\n", color)
		fmt.Print("Press [E] To Exit\n\n")
		fmt.Print("Enter Position '[1-7]: ")
		strPos, _ := reader.ReadString('\n')

		if strings.EqualFold(strPos[0:1], "E") {
			syscall.Exit(syscall.F_OK)
		}

		pos, err := strconv.Atoi(strPos[0:1])
		if err != nil {
			fmt.Printf("Invalid value: %v", err)
			continue
		}

		if pos < 1 || pos > 7 {
			fmt.Print("valid values are 1, 2, 3, 4, 5, 6, 7 \n\n")
			continue
		}

		err = insert((pos - 1), color[0:1])

		if err != nil {
			continue
		}

		i++
		fmt.Printf("Total of peaces: %d  \n\n", i)
		if i >= 42 {
			fmt.Println("All positions are used. There is no winner")
			syscall.Exit(syscall.F_OK)
		}
	}

}

func insert(hpos int, color string) error {
	if hpos < 0 || hpos >= 7 {
		fmt.Println("outside of tab...")
		return fmt.Errorf("outside of tab")
	}
	y := ypos[hpos]

	if y == 6 {
		fmt.Printf("no more space on column %d \n", hpos+1)
		return fmt.Errorf("no more space to put in this column")
	}

	tab[hpos][y] = color
	y++
	ypos[hpos] = y

	printTab()

	winner, wcolor := checkWinner()

	if winner {
		fmt.Printf("%s is the winner\n\n", wcolor)
		printTab()
		syscall.Exit(syscall.F_OK)
	}

	return nil
}

func checkWinner() (winner bool, color string) {
	win, wColor := checkHorizontalLine()
	if win {
		fmt.Printf("Winner in horizontal line - ")
		return win, wColor
	}

	win, wColor = checkVerticalLine()
	if win {
		fmt.Printf("Winner in vertical line - ")
		return win, wColor
	}

	win, wColor = checkDiagonalLineUp()
	if win {
		fmt.Printf("Winner in diagonal up line - ")
		return win, wColor
	}

	win, wColor = checkDiagonalLineDown()
	if win {
		fmt.Printf("Winner in diagonal down line - ")
		return win, wColor
	}

	return win, wColor
}

func checkDiagonalLineUp() (isWinner bool, color string) {
	wins := false
	wColor := ""
	count := 0

	for x := 0; x <= 6; x++ {
		for y := 0; y <= 5; y++ {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			count++
			npy := y
			for npx := x + 1; npx <= 7; npx++ {
				npy = npy + 1

				if npx > 6 || npy < 0 {
					count = 0
					continue
				}

				r := checkNextPositionSameColor(npx, npy, value)
				if r {
					count++
					if count == 4 {
						return true, value
					}
				} else {
					count = 0
					break
				}

			}
		}
	}

	return wins, wColor
}

func checkDiagonalLineDown() (isWinner bool, color string) {
	wins := false
	wColor := ""
	count := 0

	for x := 0; x <= 6; x++ {
		for y := 5; y >= 0; y-- {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			count++
			npy := y
			for npx := x + 1; npx <= 7; npx++ {
				npy = npy - 1

				if debug {
					fmt.Printf("->next value [%d,%d]", npx, npy)
				}

				if npx > 6 || npy < 0 {
					count = 0
					break
				}

				r := checkNextPositionSameColor(npx, npy, value)
				if r {
					count++
					if count == 4 {
						return true, value
					}
				} else {
					count = 0
					break
				}

			}
		}
	}

	return wins, wColor
}

func checkVerticalLine() (isWinner bool, color string) {
	wins := false
	wColor := ""
	count := 0

	for y := 0; y <= 5; y++ {
		for x := 0; x <= 6; x++ {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			if debug {
				fmt.Printf("\n Found: [%d,%d]=%s > ", x, y, value)
			}

			count++
			for npy := y + 1; npy <= 6; npy++ {

				if npy > 6 || npy < 0 {
					count = 0
					break
				}

				r := checkNextPositionSameColor(x, npy, value)
				if r {
					count++
					if count == 4 {
						return true, value
					}
				} else {
					count = 0
					break
				}
			}
		}
	}

	return wins, wColor
}

func checkHorizontalLine() (isWinner bool, color string) {
	wins := false
	wColor := ""
	count := 0

	for y := 0; y <= 5; y++ {
		for x := 0; x <= 6; x++ {
			value := tab[x][y]

			if value == "" {
				count = 0
				continue
			}

			count++
			for npx := x + 1; npx <= 7; npx++ {

				if npx > 6 || npx < 0 {
					count = 0
					break
				}

				r := checkNextPositionSameColor(npx, y, value)
				if r {
					count++
					if count == 4 {
						return true, value
					}
				} else {
					count = 0
					break
				}
			}
		}
	}

	return wins, wColor
}

func checkNextPositionSameColor(x, y int, color string) bool {
	if color == "" {
		return false
	}

	if x < 0 || x > 6 {
		return false
	}

	if y < 0 || y > 5 {
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
	for y := 5; y >= 0; y-- {
		for x := 0; x <= 6; x++ {

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
