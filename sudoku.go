package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 0: 未入力
// 1-9: 入っている
type Board [9][9]int

// 重複してるときはtrueを返す関数
func duplicated(c [10]int) bool {
	// cが{9, 0, 0, 0, 0, 0, 0, 0, 0, 0}であれば、その行は全て0が並んでいるということとなり、これはOK、つまりfalseにしたい。
	// cが{8, 1, 0, 0, 0, 0, 0, 0, 0, 0}であれば、その行は8個の0と、1個の1が並んでいる。これはOK、つまりfalseにしたい。
	// cが{7, 2, 0, 0, 0, 0, 0, 0, 0, 0}であれば、その行は7個の0と、2個の1が並んでいる。これはNG、つまりtrueにしたい。
	// 最終的な理想的な状態はcが{0, 1, 1, 1, 1, 1, 1, 1, 1, 1}となっていること。これはこれはOK、つまりfalseにしたい。

	for k, v := range c {
		if k == 0 {
			continue
		}
		if v >= 2 {
			return true
		}
	}
	return false
}

func verify(b Board) bool {
	// 行チェック
	for i := 0; i < 9; i++ {
		// 出現回数
		var c [10]int
		for j := 0; j < 9; j++ {
			c[b[i][j]]++
		}
		if duplicated(c) {
			return false
		}
	}

	// 列チェック
	for i := 0; i < 9; i++ {
		// 出現回数
		var c [10]int
		for j := 0; j < 9; j++ {
			c[b[j][i]]++
		}
		if duplicated(c) {
			return false
		}
	}

	// 3x3チェック
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			var c [10]int
			for row := i; row < i+3; row++ {
				for col := j; col < j+3; col++ {
					c[b[row][col]]++
				}
			}
			if duplicated(c) {
				return false
			}
		}
	}

	return true
}

func solved(b Board) bool {
	if !verify(b) {
		return false
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if b[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

// Boardを参照で渡すのはなぜか -> 値渡しだと呼び出し側のBoardとbacktrack側のBoardが別物となってしまうため。
// 正確には、呼び出し側のBoardのコピーがbacktrackに引数として渡されてしまう。なので、参照渡しを使う。
// そしてなんでコピーだと嫌かというと、今回この関数は再帰される(関数の中で自分自身を呼び出す)。
// そうすると、どんどんコピーができてしまうので、オリジナルのものでやろうという意図。
func backtrack(b *Board) bool {
	// time.Sleep(time.Second * 1)
	// fmt.Printf("%v\n", pretty(*b))
	if solved(*b) {
		return true
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			// ますが0だったら、ますに数値を入れる
			if b[i][j] == 0 {
				for c := 9; c >= 1; c-- {
					b[i][j] = c
					// 入れた上で、数値がルールに適合するならさらに深く探索
					if verify(*b) {
						if backtrack(b) {
							return true
						}
					}
					// 適合しないなら、0を入れる
					b[i][j] = 0
				}
				return false
			}
		}
	}
	return false
}

func pretty(b Board) string {
	var buf bytes.Buffer
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			buf.WriteString("+---+---+---+\n")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				buf.WriteString("|")
			}
			buf.WriteString(strconv.Itoa(b[i][j]))
		}
		buf.WriteString("|\n")
	}

	buf.WriteString("+---+---+---+\n")
	return buf.String()
}

// input :.5..83.17...1..4..3.4..56.8....3...9.9.8245....6....7...9....5...729..861.36.72.4
func short(input string) (*Board, error) {
	if len(input) != 81 {
		// fmt.Printf("input string length: %v\n", len(input))
		return nil, errors.New("input short string length must be 81")
	}
	s := bufio.NewScanner(strings.NewReader(input))
	s.Split(bufio.ScanRunes)

	var b Board
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if !s.Scan() {
				break
			}
			token := s.Text()
			if token == "." {
				b[i][j] = 0
				continue
			}
			n, err := strconv.Atoi(token)
			if err != nil {
				return nil, err
			}

			b[i][j] = n
		}
	}

	return &b, nil
}

func main() {
	flag.Parse()
	input := flag.Arg(0)
	// fmt.Printf("Arg(0) is %v, Arg(1) is %v\n", flag.Arg(0), flag.Arg(1))
	b, err := short(input)
	if err != nil {
		panic(err)
	}

	if backtrack(b) {
		fmt.Println(pretty(*b))
	} else {
		fmt.Fprintf(os.Stderr, "cannot solve")
	}
}
