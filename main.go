package main

import (
	"letsgo/cli"
	//"letsgo/util"
	//"time"
	//"fmt"
	//"math/rand"
	//"math"
)


func main(){
	cli.Start()
}
/*
func init() {
    rand.Seed(time.Now().UnixNano())
}


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func multSha256(hash []byte, m int) []byte{
	for i :=0 ;i < m; i++ {
		hash = util.Sha256(hash)
	}
	return hash
} */

/*
func main(){
	h := fmt.Sprintf("%02x", 360000)
	fmt.Println([]byte(h))
}

func main(){
	start := time.Now()
	const MAX_HASH_COPY = 31250
	const initial_difficulty = 100
	const difficulty = 300
	const len = 83
	const perfect_nonce = 170000000

	mult := float64(perfect_nonce / MAX_HASH_COPY)
	mult *= (difficulty / initial_difficulty)
	mult = math.Ceil(mult)
	fmt.Println(mult)
	var ret [][]byte
	for i := 0; i < MAX_HASH_COPY; i++ {
		ret = append(ret, multSha256([]byte(RandStringRunes(len)), int(mult)))
	}
	fmt.Println(ret[MAX_HASH_COPY - 1])
	fmt.Println(time.Now().Sub(start))
}*/
/*
var t_a []time.Duration
const LEN = 100
func main(){
	const MAX_HASH_COPY = 31250
	const MULTY_HASH = 10000

	for j := 0; j < LEN; j++{
		start := time.Now()

		for i := 0; i < 100000; i++{
			util.Sha256([]byte(RandStringRunes(32)))
		}
		t_a = append(t_a, time.Now().Sub(start))
	}
	medium()
}

func medium(){
	var	total int64
	for _, t := range t_a{
		total += t.Nanoseconds()	
	}
	fmt.Println(total / LEN)
}*/