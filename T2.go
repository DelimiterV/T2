// T2.go
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

//------------------------------------
type ctt struct {
	str    []uint8
	strlen int
	id     int
}

var ustrings []ctt
var ctmu sync.Mutex

func FindStringPos(mstr []uint8, pie *[5]uint32) (r int, ind uint32) {

	pie[0] = uint32(len(ustrings))
	pie[1] = 0
	if pie[0] > 1 {
		pie[4] = pie[0] - 1

		for pie[3] = 0; pie[1] == 0; {

			if string(mstr) > string(ustrings[pie[3]].str) && string(mstr) < string(ustrings[pie[4]].str) {
				pie[2] = (pie[4] - pie[3]) >> 1
				if pie[2] == 0 {
					r = 100 //место для вставки
					ind = pie[4]
					pie[1] = 1
				} else {
					if string(mstr) < string(ustrings[pie[3]+pie[2]].str) {
						pie[4] = pie[4] - pie[2]
					} else {
						if string(mstr) == string(ustrings[pie[3]+pie[2]].str) {
							//нашли
							r = 0
							ind = pie[3] + pie[2]
							pie[1] = 1
						} else {
							if string(mstr) > string(ustrings[pie[3]+pie[2]].str) {
								pie[3] += pie[2]
							}
						}

					}
				}
			} else {
				//string(ustrings[l1].substr
				if string(mstr) < string(ustrings[pie[3]].str) {
					r = -1
					pie[1] = 1
					ind = pie[3]
				} else {
					if string(mstr) == string(ustrings[pie[3]].str) {
						r = 0
						pie[1] = 1
						ind = pie[3]
					} else {
						if string(mstr) > string(ustrings[pie[4]].str) {
							r = 1
							pie[1] = 1
							ind = pie[4]
						} else {
							r = 0
							pie[1] = 1
							ind = pie[4]
						}
					}
				}
			}
		}
	} else {
		if pie[0] == 1 {
			switch {
			case string(mstr) == string(ustrings[0].str):
				r = 0
				ind = 0
			case string(mstr) < string(ustrings[0].str):
				r = -1
				ind = 0
			case string(mstr) > string(ustrings[0].str):
				r = 1
				ind = 1
			}
		} else {
			r = 1
			ind = 0
		}

	}
	return
}
func paddString(substr []uint8, id int) {
	var pos uint32
	var us int
	var dString ctt
	var mbuf [5]uint32

	if len(substr) > 0 {

		dString.str = substr
		dString.strlen = len(substr)
		dString.id = id
		us, pos = FindStringPos(substr, &mbuf)

		switch us {
		case 0: //найден
			ustrings = append(ustrings, dString)
			copy(ustrings[pos+1:], ustrings[pos:])
			ustrings[pos] = dString
		case -1: //слева
			ustrings = append(ustrings, dString)
			copy(ustrings[1:], ustrings[0:])
			ustrings[0] = dString
		case 1: //справа
			ustrings = append(ustrings, dString)
		case 100: // между
			ustrings = append(ustrings, dString)
			copy(ustrings[pos+1:], ustrings[pos:])
			ustrings[pos] = dString
		}
	}
}

var pid [5]uint32
var wg sync.WaitGroup

type uff struct {
	l   uint8
	pos uint32
	id  int
}

func lookDec(pstr []uint8, a int, kk []uff) (rez int) {
	var r, f int
	var pid [5]uint32
	var pos uint32
	var un uff
	rez = 0
	for i := a; i > 0 && rez == 0; i-- {
		r, pos = FindStringPos(pstr[:i], &pid)
		if r == 0 {

			f = ustrings[pos].strlen
			un.id = ustrings[pos].id
			un.l = uint8(f)
			kk = append(kk, un)
			if a-f > 0 {

				rez = lookDec((pstr[f:]), a-f, kk)
			} else {
				for f = 0; f < len(kk); f++ {
					fmt.Print(kk[f].id)
					if f < len(kk)-1 {
						fmt.Print(",")
					}
				}
				fmt.Println("")
				rez = 1
			}
		}
	}
	return
}

func readStrings(filename string) bool {
	rez := true
	var l, k int
	r := ""
	dat, e := ioutil.ReadFile(filename)
	if e == nil {
		l = 0
		for i := 0; i < len(dat); i++ {
			switch dat[i] {
			case 0x0d:
				k = i
			case 0x0a:
				r = string(dat[l:k])
				kr := strings.Split(r, ",")
				mid, _ := strconv.Atoi(kr[0])
				paddString([]uint8(kr[1]), mid)
				l = i + 1
			}
		}
	} else {
		rez = false
	}
	return rez
}

var mainstring string

func readMainString(filename string) bool {
	rez := true
	mainstring = ""
	r := ""
	dat, e := ioutil.ReadFile(filename)
	if e == nil {
		zstr := string(dat[:])
		ustr := strings.Split(zstr, "\r")
		zstr = strings.Join(ustr, "")
		ustr = strings.Split(zstr, "\n")
		r = strings.Join(ustr, "")
		mainstring = r
	} else {
		rez = false
	}
	return rez
}
func decide() (string, bool) {
	rez := true
	t := make([]uff, 0)

	if lookDec([]uint8(mainstring), len(mainstring), t) == 0 {
		fmt.Println("No decission")
	} else {
		fmt.Println("---")
	}

	return "OK", rez
}

const MyStrings = "samples.txt"
const MainString = "mainstring.txt"

func main() {
	fmt.Println("Строчки!")

	fmt.Println("----prepare----")
	if readStrings(MyStrings) {
		if readMainString(MainString) {
			fmt.Println("----start-----")
			if t, e := decide(); e {
				fmt.Println("Success!", t)
			} else {
				fmt.Println("Нет решения")
			}
		} else {
			fmt.Println("Can't to open:", MainString)
		}
	} else {
		fmt.Println("Can't to open:", MyStrings)
	}
}
