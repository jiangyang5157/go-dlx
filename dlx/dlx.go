package dlx

/*
Dancing Links (Algorithm X) data struct.

                 |           |           |           |
- columns 0 - columns 1 - columns 2 - ......... - columns i -
                 |           |           |           |
            -  node(row)-  node     -  node     -  node     -
                 |           |           |           |
                                    -  node(row)-  node     -
                                         |           |
                        -  node(row)-  node     -  node     -
                             |           |           |
                        -  node(row)-  node     -  node     -
                             |           |           |
*/

type Column struct {
	Node
	Index int // index of the column
	Size  int // size of node of the column
}

type Node struct {
	Col                   *Column // column of the node
	Row                   *Node   // first node of the row of the node
	Up, Down, Left, Right *Node   // up, down, left, right of the node
}

// X holds the matrix-like columns and solution.
type X struct {
	Cols []Column
	Sol  Solution
}

type Solution []*Node

// NewX returns Dlx data struct with columns initialized.
func NewX(size int) *X {
	ret := &X{}
	ret.Cols = make([]Column, size+1)

	// use column 0 as head
	it := &ret.Cols[0]
	// head.column = head
	it.Col = it
	// head.left = last
	it.Left = &ret.Cols[size].Node
	// last.right = head
	ret.Cols[size].Right = &it.Node

	for i := 1; i < len(ret.Cols); i++ {
		curr := &ret.Cols[i]

		// curr index = i
		curr.Index = i
		// curr column = curr
		curr.Col = curr

		// curr up = curr
		curr.Up = &curr.Node
		// curr down = curr
		curr.Down = &curr.Node

		// curr left = prev
		curr.Left = &it.Node
		// prev.right = curr
		it.Right = &curr.Node

		// prev column = curr column
		it = curr
	}

	return ret
}

// AddRow adds a new row that contains multiple node associated with certain columns
// colIndexes: associated column indexes (excluded head index 0)
func (x *X) AddRow(colIndexes []int) {
	row := make([]Node, len(colIndexes))
	for i, ci := range colIndexes {
		col := &x.Cols[ci]
		nd := &row[i]

		col.Size++
		nd.Col = col
		nd.Row = &row[0]

		nd.Up = col.Up
		nd.Down = &col.Node
		nd.Left = &row[circularShiftLeft(i, len(colIndexes))]
		nd.Right = &row[circularShiftRight(i, len(colIndexes))]

		nd.Up.Down = nd
		nd.Down.Up = nd
		nd.Right.Left = nd
		nd.Left.Right = nd
	}
}

func circularShiftLeft(curr int, length int) int {
	left := curr - 1
	if left < 0 {
		return length - 1
	}
	return left
}

func circularShiftRight(curr int, length int) int {
	right := curr + 1
	if right >= length {
		return 0
	}
	return right
}

func cover(col *Column) {
	col.Right.Left = col.Left
	col.Left.Right = col.Right
	for i := col.Down; i != &col.Node; i = i.Down {
		for j := i.Right; j != i; j = j.Right {
			j.Down.Up = j.Up
			j.Up.Down = j.Down
			j.Col.Size--
		}
	}
}

func uncover(col *Column) {
	for i := col.Up; i != &col.Node; i = i.Up {
		for j := i.Left; j != i; j = j.Left {
			j.Down.Up = j
			j.Up.Down = j
			j.Col.Size++
		}
	}
	col.Right.Left = &col.Node
	col.Left.Right = &col.Node
}

// Search runs Algorithm X.
// Search Passes solution through f() when a completed solution found, Algorithm X will continue for next solution if f() returns false.
// Search ends normally, or the f() returns true.
func (x *X) Search(f func(Solution) bool) bool {
	head := &x.Cols[0]
	hrc := head.Right.Col
	if hrc == head {
		// circular search completed, the solution cache in x.sol
		return f(x.Sol)
	}

	// find the column has minimum size, it improves overall performance by compare with linear iterator
	it := hrc
	min := it.Size
	for {
		hrc = hrc.Right.Col
		if hrc == head {
			break
		}
		if hrc.Size < min {
			it = hrc
			min = hrc.Size
		}
	}

	ret := false
	cover(it)
	x.Sol = append(x.Sol, nil)
	oLen := len(x.Sol)
	for j := it.Down; j != &it.Node; j = j.Down {
		if ret {
			break
		}

		x.Sol[oLen-1] = j
		for i := j.Right; i != j; i = i.Right {
			cover(i.Col)
		}
		ret = x.Search(f)
		j = x.Sol[oLen-1]
		it = j.Col
		for i := j.Left; i != j; i = i.Left {
			uncover(i.Col)
		}
	}
	x.Sol = x.Sol[:oLen-1]
	uncover(it)
	return ret
}
