package game

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestState_Copy2(t *testing.T) {
	var state1 = NewState()
	var state2 = *state1
	state1.PlayMove(Coordinate{
		Col: 0,
		Row: 0,
	})
	fmt.Println(state2)
}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func TestState_GetMoves(t *testing.T) {
	defer timeTrack(time.Now(), "GetFreeRows")
	var s = NewState()
	s.GetFreeRows()
}
func TestState_HasConnectedFour(t *testing.T) {
	var s = NewState()
	s.PlayMove(Coordinate{
		Col: 3,
		Row: 5,
	})
	fmt.Println(s)

	s.PlayMove(Coordinate{
		Col: 0,
		Row: 5,
	})
	fmt.Println(s)

	s.PlayMove(Coordinate{
		Col: 4,
		Row: 5,
	})
	fmt.Println(s)

	s.PlayMove(Coordinate{
		Col: 0,
		Row: 4,
	})
	fmt.Println(s)

	s.PlayMove(Coordinate{
		Col: 2,
		Row: 5,
	})
	fmt.Println(s)

	s.PlayMove(Coordinate{
		Col: 0,
		Row: 3,
	})
	fmt.Println(s)

	s.PlayMove(Coordinate{
		Col: 1,
		Row: 5,
	})
	fmt.Println(s)
	fmt.Println(s.Outcome)
	if s.Outcome != Player1Won {
		t.Errorf("not good")
	}
}
func TestState_Playout(t *testing.T) {
	var s = NewState()
	int := s.Playout()
	//fmt.Printf("%+v\n", s)
	if int == None {
		t.Errorf("not good")
	}
}
func TestState_GetFreeColumns(t *testing.T) {
	var s = NewState()
	for i := 0; i < Cols; i++ {
		for row := 0; row < Rows; row++ {
			s.PlayMove(Coordinate{
				Col: i,
				Row: row,
			})
		}
	}

	freeColumns := s.GetFreeRows()
	_, err := s.GetRandomFreeRow(freeColumns)
	if err == nil {
		t.Errorf("a column is empty after a playout")
	}
	t.Log("succeeded")
}

func BenchmarkState_Playout(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	var dic = make(map[int]int)
	for n := 0; n < b.N; n++ {
		var s = NewState()
		int := s.Playout()
		dic[int] += 1
	}
	fmt.Println(dic)
}
func BenchmarkState_GetMoves(b *testing.B) {
	var s = NewState()
	for n := 0; n < b.N; n++ {
		s.GetFreeRows()
	}
}
func BenchmarkState_GetFreeRowCell(b *testing.B) {
	var s = NewState()
	for n := 0; n < b.N; n++ {
		s.GetFreeRow(0)
	}
}
func BenchmarkState_HasConnectedFour(b *testing.B) {
	var s = NewState()
	//source := rand.NewSource(time.Now().Unix())
	//r := rand.New(source)
	//s.Playout(r)
	for n := 0; n < b.N; n++ {
		move := Coordinate{Col: 0, Row: 0}
		s.HasConnectedFour(move)
	}
}
func BenchmarkState_GetFreeRows(b *testing.B) {
	var s = NewState()
	s.PlayMove(Coordinate{
		Col: 0,
		Row: 0,
	})
	s.PlayMove(Coordinate{
		Col: 0,
		Row: 1,
	})
	s.PlayMove(Coordinate{
		Col: 0,
		Row: 2,
	})
	s.PlayMove(Coordinate{
		Col: 0,
		Row: 3,
	})
	s.PlayMove(Coordinate{
		Col: 0,
		Row: 4,
	})
	s.PlayMove(Coordinate{
		Col: 0,
		Row: 5,
	})

	for n := 0; n < b.N; n++ {
		freeRows := s.GetFreeRows()
		s.GetRandomFreeRow(freeRows)
	}
}
