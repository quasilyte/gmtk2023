package pathing_test

import (
	"testing"

	"github.com/quasilyte/gmtk2023/pathing"
)

func BenchmarkPathgridCheck(b *testing.B) {
	p := pathing.NewGrid(1856, 1856)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.CellIsFree(pathing.GridCoord{14, 5})
	}
}

func BenchmarkPathgridMark(b *testing.B) {
	p := pathing.NewGrid(1856, 1856)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.MarkCell(pathing.GridCoord{14, 5})
	}
}
