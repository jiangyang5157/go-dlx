package dlx

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	x := NewX(10)
	for i := range x.Cols {
		fmt.Printf("%d ", x.Cols[i].Index)
	}
	fmt.Println()
	head := &x.Cols[0]
	for col := head.Right.Col; col.Node != head.Node; col = col.Right.Col {
		fmt.Printf("%d ", col.Index)
	}
	fmt.Println()
}
