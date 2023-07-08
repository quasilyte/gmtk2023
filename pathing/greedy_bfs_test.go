package pathing_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/quasilyte/gmtk2023/pathing"
)

func BenchmarkGreedyBFS(b *testing.B) {
	for i := range bfsTests {
		test := bfsTests[i]
		if !test.bench {
			continue
		}
		numCols := len(test.path[0])
		numRows := len(test.path)
		b.Run(fmt.Sprintf("%s_%dx%d", test.name, numCols, numRows), func(b *testing.B) {
			parseResult := testParseGrid(b, test.path)
			bfs := pathing.NewGreedyBFS(parseResult.numCols, parseResult.numRows)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bfs.BuildPath(parseResult.grid, parseResult.start, parseResult.dest)
			}
		})
	}
}

func TestGreedyBFS(t *testing.T) {
	for i := range bfsTests {
		test := bfsTests[i]
		t.Run(test.name, func(t *testing.T) {
			m := test.path

			parseResult := testParseGrid(t, m)
			bfs := pathing.NewGreedyBFS(parseResult.numCols, parseResult.numRows)
			grid := parseResult.grid

			result := bfs.BuildPath(grid, parseResult.start, parseResult.dest)
			path := result.Steps

			pos := parseResult.start
			pathLen := 0
			for path.HasNext() {
				pathLen++
				d := path.Next()
				pos = pos.Move(d)
				marker := parseResult.haveRows[pos.Y][pos.X]
				switch marker {
				case 'A':
					parseResult.haveRows[pos.Y][pos.X] = 'A'
				case 'B':
					parseResult.haveRows[pos.Y][pos.X] = '$'
				case ' ':
					t.Fatal("visited one cell more than once")
				case '.':
					parseResult.haveRows[pos.Y][pos.X] = ' '
				default:
					panic(fmt.Sprintf("unexpected %c marker", marker))
				}
			}

			have := string(bytes.Join(parseResult.haveRows, []byte("\n")))
			want := strings.Join(test.path, "\n")

			if have != want {
				t.Fatalf("paths mismatch\nmap:\n%s\nhave (l=%d):\n%s\nwant (l=%d):\n%s",
					strings.Join(m, "\n"), pathLen, have, parseResult.pathLen, want)
			}

			wantPartial := test.partial
			havePartial := pos != parseResult.dest && result.Partial
			if havePartial != wantPartial {
				t.Fatalf("partial flag mismatch\nmap:\n%s\nhave: %v\nwant: %v", strings.Join(m, "\n"), havePartial, wantPartial)
			}
		})
	}
}

type testGrid struct {
	start    pathing.GridCoord
	dest     pathing.GridCoord
	grid     *pathing.Grid
	pathLen  int
	numCols  int
	numRows  int
	haveRows [][]byte
}

func testParseGrid(tb testing.TB, m []string) testGrid {
	tb.Helper()

	numCols := len(m[0])
	numRows := len(m)

	grid := pathing.NewGrid(pathing.CellSize*float64(numCols), pathing.CellSize*float64(numRows))

	pathLen := 0
	var startPos pathing.GridCoord
	var destPos pathing.GridCoord
	haveRows := make([][]byte, numRows)
	for row := 0; row < numRows; row++ {
		haveRows[row] = make([]byte, numCols)
		for col := 0; col < numCols; col++ {
			marker := m[row][col]
			cell := pathing.GridCoord{X: col, Y: row}
			haveRows[row][col] = marker
			switch marker {
			case 'x':
				grid.MarkCell(cell)
			case 'A':
				startPos = cell
			case 'B', '$':
				if marker == '$' {
					pathLen++
				}
				destPos = cell
				haveRows[row][col] = 'B'
			case ' ':
				pathLen++
				haveRows[row][col] = '.'
			}
		}
	}

	return testGrid{
		pathLen:  pathLen,
		start:    startPos,
		dest:     destPos,
		numRows:  numRows,
		numCols:  numCols,
		haveRows: haveRows,
		grid:     grid,
	}
}

type bfsTestCase struct {
	name    string
	path    []string
	partial bool
	bench   bool
}

var bfsTests = []bfsTestCase{
	{
		name: "trivial_short",
		path: []string{
			"..........",
			"...A   $..",
			"..........",
		},
		bench: true,
	},

	{
		name: "trivial_short2",
		path: []string{
			"..........",
			"...A......",
			"... ......",
			"... ......",
			"...  $....",
			"..........",
		},
		bench: true,
	},

	{
		name: "trivial",
		path: []string{
			".A..........",
			". ..........",
			". ..........",
			". ..........",
			". ..........",
			". ..........",
			".          $",
		},
		bench: true,
	},

	{
		name: "trivial_long",
		path: []string{
			".......................x........",
			"                               $",
			"A...............................",
			"..........................x.....",
		},
		bench: true,
	},

	{
		name: "simple_wall1",
		path: []string{
			"........",
			"...A....",
			"...   ..",
			"....x ..",
			"....x $.",
		},
		bench: true,
	},

	{
		name: "simple_wall2",
		path: []string{
			"...   ..",
			"...Ax ..",
			"....x ..",
			"....x ..",
			"....x $.",
		},
		bench: true,
	},

	{
		name: "simple_wall3",
		path: []string{
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			".............   ................",
			"..            x          $......",
			".. ...........x.................",
			"..A...........x.................",
			"....x...........................",
			"....x...........................",
			"....x...........................",
			"....x...........................",
		},
		bench: true,
	},

	{
		name: "simple_wall4",
		path: []string{
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			"..........x.....................",
			"................................",
			"..............x.................",
			"..............x.................",
			"..A...........x.................",
			".. .x...........................",
			".. .x...........................",
			".. .x...........................",
			".. .x...........................",
			".. .............................",
			".. .............................",
			".. ..................xxxxxxxx...",
			".. .............................",
			".. .............................",
			".. ...........x.................",
			".. ...........x.................",
			"..    ........x.................",
			"....x ..........................",
			"....x                      $....",
			"....x...........................",
			"....x...........................",
		},
		bench: true,
	},

	{
		name: "zigzag1",
		path: []string{
			"........",
			"   A....",
			" xxxxxx.",
			" .......",
			" .xxxxxx",
			" .......",
			" $......",
		},
		bench: true,
	},

	{
		name: "zigzag2",
		path: []string{
			"........",
			"...A    ",
			".xxxxxx ",
			".....   ",
			"..xxx xx",
			"..... ..",
			".....  $",
		},
		bench: true,
	},

	{
		name: "zigzag3",
		path: []string{
			"...   ....x.....",
			"..A x ....x.....",
			"....x ....x.....",
			"....x ....x.....",
			"....x        $..",
			"....x...........",
		},
		bench: true,
	},

	{
		name: "zigzag4",
		path: []string{
			"...   .x.   x...",
			"... x .x. x x...",
			"... x .x. x x...",
			"... x .x. x   ..",
			"..A x  x  x.x  $",
			"....x.   .x.x...",
		},
		bench: true,
	},

	{
		name: "zigzag5",
		path: []string{
			".A     ..",
			"xxxxxx ..",
			"..     ..",
			".. xxxxxx",
			"..   ....",
			"xxxx x...",
			"....    .",
			"...xxxx x",
			".......$.",
		},
		bench: true,
	},

	{
		name: "double_corner1",
		path: []string{
			".   .x  A.",
			". x .x ...",
			"x x .x ...",
			"  x .x ...",
			" xx    ...",
			" .xxxxxxxx",
			"   $......",
		},
		bench: true,
	},

	{
		name: "double_corner2",
		path: []string{
			".    x..A.",
			". x. x.. .",
			"x x. x.. .",
			"  x. x.. .",
			" xx.     .",
			" .xxxxxxxx",
			"        $.",
			"..........",
		},
	},

	{
		name: "double_corner3",
		path: []string{
			"    x..A.",
			" x. x.. .",
			" x. x.. .",
			" x.     .",
			" xxxxxxxx",
			"       $.",
		},
	},

	{
		name: "labyrinth1",
		path: []string{
			".........x.....",
			"xxxxxxxx.x.  $.",
			"x.     x.x. ...",
			"x. xxx x.x. ...",
			"x.   x x.x. ...",
			"x...Ax   xx .xx",
			"x....x.x x  ...",
			"xxxxxx.x x xxxx",
			"x......x x    .",
			"xxxxxxxx xxxx x",
			"........ x    .",
			"........   ....",
		},
		bench: true,
	},

	{
		name: "labyrinth2",
		path: []string{
			".x......x.......x............",
			".x......x.......x............",
			".x......x.......x............",
			".x......x.......xxxxxxxxxx...",
			".x....       ...x.....    ...",
			".x     xxx.x    x.....$.x  xx",
			"   .x..x...xxx. x.......x.  .",
			"A...x..x...x... xxxxxxxxxxx .",
			"..x.x..x.......     x       .",
			"..x.x..x....x...... x .......",
			"..x.x..x..xxxx...x.   .......",
			"..x.x.......x....x...........",
		},
		bench: true,
	},

	{
		name: "labyrinth3",
		path: []string{
			"...x......x........x............",
			"..Ax......x........x............",
			".. x......x........xxxxxxxxxx...",
			".. x...............x............",
			".. x.....xxx..x....x.......x..xx",
			".. ...x..x....xxx..x.......x....",
			".. ...x..x....x....xxxxxxxxxxx..",
			".. .x.x..x.....x...   .x........",
			".. .x.x..x...xxxx.  x         ..",
			"..        x....... xxxxxxxxxx ..",
			"xxxx.....       .. x........  ..",
			"...x.....xxx..x    x.......x .xx",
			"......x..x....xxx..x.......x   .",
			"......x..x....x....xxxxxx.xxxx .",
			"....x.x........x....x..x...$   .",
		},
		bench: true,
	},

	{
		// This is unfortunate.
		// TODO: can we adjust anything to make it better?
		name: "depth1",
		path: []string{
			"........................",
			".xxxxxxxxxxxxxxxxxxxx...",
			"                    x...",
			".xxxxxxxxxxxxxxxxxx x...",
			"..                  x...",
			".x xxxxxxxxxxxxxxxxxx...",
			"..                A.x.B.",
			".x.xxxxxxxxxxxxxxxxxx...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxx.x...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxxxx...",
			"........................",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "depth2",
		path: []string{
			"...................   ..",
			"..                  x ..",
			".x xxxxxxxxxxxxxxxxxx ..",
			"..                A.x $.",
			".x.xxxxxxxxxxxxxxxxxx...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxx.x...",
			"....................x...",
			".xxxxxxxxxxxxxxxxxxxx...",
			"........................",
		},
		bench: true,
	},

	{
		name: "nopath1",
		path: []string{
			"A    x.B",
			".....x..",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath2",
		path: []string{
			"....Ax.B",
			".....x..",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath3",
		path: []string{
			"........",
			".xxxxx..",
			".x...x..",
			".x.A.x..",
			".x.  x..",
			".xxxxx..",
			".......B",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath4",
		path: []string{
			".....x.....x..",
			".xxxxx.   .x..",
			".x...x. x .x..",
			".x.A.x. x  x..",
			".x.     xxxx..",
			".xxxxxxxx.....",
			".............B",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "nopath5",
		path: []string{
			".B...x.....x..",
			".xxxxx.....x..",
			".x  .x..x..x..",
			".x.A.x..x..x..",
			".x......xxxx..",
			".xxxxxxxx.....",
			"..............",
		},
		partial: true,
		bench:   true,
	},

	{
		name: "tricky1",
		path: []string{
			"               $",
			" .xxxxxxxxxxxx..",
			" ............x..",
			" ............x..",
			" ............x..",
			" ............x..",
			" ............x..",
			"A..xxxxxxxxxxx..",
			"................",
		},
		bench: true,
	},

	{
		name: "tricky2",
		path: []string{
			"...............",
			".             .",
			"  xxxxxxxxxxx $",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			"A.xxxxxxxxxxx..",
			"...............",
			"...............",
		},
		bench: true,
	},

	{
		name: "tricky3",
		path: []string{
			"...............",
			"...............",
			"..xxxxxxxxxxx A",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"............x .",
			"$ xxxxxxxxxxx .",
			".             .",
			"...............",
		},
		bench: true,
	},

	{
		name: "tricky4",
		path: []string{
			"...............",
			".             .",
			". xxxxxxxxxxx $",
			".     ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			"..... ......x..",
			".....A......x..",
			"..xxxxxxxxxxx..",
			"...............",
			"...............",
		},
		bench: true,
	},

	{
		name: "tricky5",
		path: []string{
			"...............",
			"...............",
			"A.xxxxxxxxxxx..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			" ...........x..",
			"  xxxxxxxxxxx $",
			".             .",
			"...............",
		},
	},

	{
		name: "tricky6",
		path: []string{
			"............$ .",
			"............. .",
			"..xxxxxxxxxxx .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..            .",
			"..A............",
		},
	},

	{
		name: "tricky7",
		path: []string{
			"..          A..",
			".  ............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". .............",
			". $............",
		},
	},

	{
		name: "tricky8",
		path: []string{
			". $............",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			".            ..",
			"............A..",
		},
	},

	{
		name: "tricky9",
		path: []string{
			". $............",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x         x..",
			". x .......Ax..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			".   ...........",
			"...............",
		},
	},

	{
		name: "tricky10",
		path: []string{
			". $............",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". xA........x..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			". x ........x..",
			".   ...........",
			"...............",
		},
	},

	{
		name: "tricky11",
		path: []string{
			".    $.........",
			". .............",
			". xxxxxxxxxxx..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.........x..",
			". x.        x..",
			". x  ...... x..",
			".   .......  ..",
			"............A..",
		},
	},

	{
		name: "tricky12",
		path: []string{
			"..........$   .",
			"............. .",
			"..xxxxxxxxxxx .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"..x.........x .",
			"............  .",
			"............A..",
		},
	},

	{
		name: "distlimit1",
		path: []string{
			"A                                                        ..........B",
		},
		bench:   true,
		partial: true,
	},

	{
		name: "distlimit2",
		path: []string{
			"A.............x......   ....            ......x.....x.....x....",
			" .............x...... x      xxxxxxxxxx ......x..x..x..x..x....",
			" ...xxxxxxxxxxx...... x...............x ......x..x..x..x..x....",
			"                      x...............x       ...x.....x......B",
		},
		bench:   true,
		partial: true,
	},
}
