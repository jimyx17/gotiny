package bst

// You only need to import one library!
import "testing"

// This is called a test table. It's a way to easily
// specify tests while avoiding boilerplate.
// See https://github.com/golang/go/wiki/TableDrivenTests
var tests = []struct {
	input  uint64
	output bool
}{
	{6, true},
	{16, false},
	{3, true},
}

func TestSearch(t *testing.T) {
	//     6
	//    /
	//   3
	tree := &Node{Key: 6, Index: 12, Left: &Node{Key: 3, Index: 23}}

	for i, test := range tests {
		if res := tree.Search(test.input); (res != nil) != test.output {
			t.Errorf("%d: got %v, expected %v", i, res, test.output)
		}
	}

}

func BenchmarkSearch(b *testing.B) {
	tree := &Node{Key: 6, Index: 234}

	for i := 0; i < b.N; i++ {
		tree.Search(6)
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := &Node{}

	// data := []uint64{34, 5, 34523542345, 1, 2, 6, 82424}

	for i := 0; i < b.N; i++ {
		tree.Insert(42, 23)
	}
}
