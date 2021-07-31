package main

import (
	"log"
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
)

type Loc struct {
	X, Y uint8
}

func (l Loc)String() string {
	return fmt.Sprintf("(%d,%d)", l.X, l.Y)
}

func (l Loc) AdjacentLocations() (ret []Loc) {
	ret = make([]Loc, 0, 4)
	if l.X > 0 {
		ret = append(ret, Loc{l.X - 1, l.Y})
	}
	if l.Y > 0 {
		ret = append(ret, Loc{l.X, l.Y - 1})
	}
	ret = append(ret, Loc{l.X + 1, l.Y})
	ret = append(ret, Loc{l.X, l.Y + 1})
	//log.Printf("Locations adjacent to %+v: %+v", l, ret)
	return
}

type LetterMatrix map[Loc]rune

type LetterLoc struct {
	Location Loc
	Rune rune
}

func (ll LetterLoc)String() string {
	return fmt.Sprintf("%s %s", ll.Location, string(ll.Rune))
}

type MatrixWord []LetterLoc

func (ms *MatrixWord)Word() (s string) {
	for _, lloc := range *ms {
		s = s + string(lloc.Rune)
	}
	return
}

func (ms *MatrixWord)HasLoc(loc Loc) bool {
	for _, lloc := range *ms {
		if lloc.Location == loc {
			return true
		}
	}
	return false
}

type Words []string

func (w *Words) IsWord(s string) bool {
	for _, i := range *w {
		if i == s {
			return true
		}
	}
	return false
}

type WordSearch struct {
	matrix LetterMatrix
	words  Words
}

func (ws *WordSearch) FindWords() (words []MatrixWord) {
	words = make([]MatrixWord, 0)
	for loc, _ := range ws.matrix {
		words = append(words, ws.findWordsAt(loc, MatrixWord{})...)
	}
	return
}

func WordList(words []MatrixWord) (out []interface{}) {
	out = make([]interface{},len(words))
	for idx, word := range words {
		out[idx] = word.Word()
	}
	return
}

// recursive next-letter word searcher with previous path included
func (ws *WordSearch) findWordsAt(loc Loc, sofar MatrixWord) (words []MatrixWord) {
	orig_prefix := log.Prefix()
	new_prefix := fmt.Sprintf("findWordsAt(%s,%+v) ", loc, sofar)
	//log.Print(new_prefix)
	log.SetPrefix(orig_prefix + " " + new_prefix)
	defer log.SetPrefix(orig_prefix)

	words = make([]MatrixWord, 0)
	sofar = append(sofar, LetterLoc{loc, ws.matrix[loc]})
	if ws.words.IsWord(sofar.Word()) {
		// if this finishes off a word, add it to our list of words so far
		log.Printf("Found word %s!  %s", sofar.Word(), sofar)
		word := make(MatrixWord,len(sofar))
		copy(word, sofar)
		words = append(words, word)
	}
	for _, nextloc := range loc.AdjacentLocations() { // potentially valid next locations base soley on Loc
		if _, ok := ws.matrix[nextloc]; !ok {				// location is not defined in matrix
			//log.Print("%s is not in the matrix", nextloc)
			continue
		}
		// if we already went to that location , can't go there again
		if sofar.HasLoc(nextloc) {
			//log.Printf("%s is already in the word candidate %+v", nextloc, sofar)
			continue
		}
		// find words continuing from our `nextloc` onwards
		words = append(words, ws.findWordsAt(nextloc, sofar)...)
	}
	return
}

func main() {
	ws := WordSearch{
		matrix: LetterMatrix{
			Loc{0, 0}: 'a',
			Loc{1, 0}: 'h',
			Loc{2, 0}: 't',
			Loc{3, 0}: 'l',
			Loc{0, 1}: 's',
			Loc{1, 1}: 'n',
			Loc{2, 1}: 'o',
			Loc{3, 1}: 'a',
			Loc{0, 2}: 'e',
			Loc{1, 2}: 'f',
			Loc{2, 2}: 's',
			Loc{3, 2}: 't',
			Loc{0, 3}: 'w',
			Loc{1, 3}: 'a',
			Loc{2, 3}: 's',
			Loc{3, 3}: 'e',
		},
		words: Words{ },
	}
	if f, err := os.Open("words.txt"); err != nil{
		panic(err)
	} else {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
				ws.words = append(ws.words, strings.TrimSpace(scanner.Text()))
		}
	}
	fmt.Println(WordList(ws.FindWords())...)
	ws.matrix.Draw(os.Stdout)
}

func (lm LetterMatrix)Draw(out io.Writer) {
	var maxx, maxy uint8
	// find the highest loc
	for loc, _ := range lm {
		if loc.X > maxx {
				maxx = loc.X
		}
		if loc.Y > maxy {
			maxy = loc.Y
		}
	}
	for y := uint8(0); y <= maxy; y++ {
		for x := uint8(0); x <= maxx; x++ {
			str := ""
			if r, ok := lm[Loc{x,y}]; ok {
				str = string(r)
			}
			fmt.Fprintf(out, " %s ", str)
		}
		fmt.Fprintf(out, "\n\n")
	}
}
