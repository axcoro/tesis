package board

import "testing"

func TestCell_Alive(t *testing.T) {
	type fields struct {
		state bool
		next  bool
		x     int
		y     int
		board *Board
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cell{
				state: tt.fields.state,
				next:  tt.fields.next,
				x:     tt.fields.x,
				y:     tt.fields.y,
				board: tt.fields.board,
			}
			if got := c.Alive(); got != tt.want {
				t.Errorf("Cell.Alive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCell_neighborsAlive(t *testing.T) {
	type fields struct {
		state bool
		next  bool
		x     int
		y     int
		board *Board
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Cell{
				state: tt.fields.state,
				next:  tt.fields.next,
				x:     tt.fields.x,
				y:     tt.fields.y,
				board: tt.fields.board,
			}
			if got := c.neighborsAlive(); got != tt.want {
				t.Errorf("Cell.neighborsAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}
