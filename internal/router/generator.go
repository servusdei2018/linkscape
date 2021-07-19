// Copyright 2021 The Linkscape Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"sync"
)


// generator generates fixed length character sequences that fulfill
// the regex [a-zA-X]{x} where x is the length.
type generator struct {
	carry int
	indices []int
	sync.Mutex
}

// newGenerator returns a new generator.
func newGenerator(length int) *generator {
	return &generator{
		carry: 1,
		indices: make([]int, length),
	}
}

func (g *generator) Next() string {
	g.Lock()
	defer g.Unlock()

	uid := ""
	for _, val := range g.indices {
		uid += toChar(val)
	}

	g.carry = 1
	for i := len(g.indices)-1; i >= 0; i-- {
		if g.carry == 0 {
			break
		}
		
		g.indices[i] += g.carry
		g.carry = 0
		
		if g.indices[i] == arrlen {
			g.carry = 1
			g.indices[i] = 0
		}
	}

	return uid
}

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var arrlen = len(arr)
func toChar(i int) string { return arr[i] }
