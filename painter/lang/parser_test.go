package lang_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ypapish/software-architecture-lab3/painter"
	"github.com/ypapish/software-architecture-lab3/painter/lang"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []painter.Operation
		wantErr  string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: []painter.Operation{},
		},
		{
			name:  "white command",
			input: "white",
			expected: []painter.Operation{
				painter.WhiteFill{},
			},
		},
		{
			name:  "green command",
			input: "green",
			expected: []painter.Operation{
				painter.GreenFill{},
			},
		},
		{
			name:  "bgrect with valid arguments",
			input: "bgrect 0.1 0.2 0.3 0.4",
			expected: []painter.Operation{
				painter.BgRect{X1: 0.1, Y1: 0.2, X2: 0.3, Y2: 0.4},
			},
		},
		{
			name:    "bgrect with insufficient arguments",
			input:   "bgrect 0.1 0.2",
			wantErr: "bgrect requires 4 arguments",
		},
		{
			name:  "figure with valid arguments",
			input: "figure 0.5 0.6",
			expected: []painter.Operation{
				painter.Figure{X: 0.5, Y: 0.6},
			},
		},
		{
			name:    "figure with invalid arguments",
			input:   "figure abc def",
			wantErr: "invalid number in figure",
		},
		{
			name:  "move with valid arguments",
			input: "move 0.2 0.3",
			expected: []painter.Operation{
				painter.Move{DX: 0.2, DY: 0.3},
			},
		},
		{
			name:  "reset command",
			input: "reset",
			expected: []painter.Operation{
				painter.Reset{},
			},
		},
		{
			name:  "update command",
			input: "update",
			expected: []painter.Operation{
				painter.UpdateOp,
			},
		},
		{
			name: "multiple commands in correct order",
			input: `bgrect 0 0 1 1
green
figure 0.5 0.5
move 0.1 0.1
update`,
			expected: []painter.Operation{
				painter.BgRect{X1: 0, Y1: 0, X2: 1, Y2: 1},
				painter.GreenFill{},
				painter.Figure{X: 0.5, Y: 0.5},
				painter.Move{DX: 0.1, DY: 0.1},
				painter.UpdateOp,
			},
		},
		{
			name:    "invalid command",
			input:   "invalidcommand",
			wantErr: "unknown command: invalidcommand",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := lang.Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.input))

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(tc.expected), len(ops), "number of operations mismatch")

			for i, expectedOp := range tc.expected {
				switch eop := expectedOp.(type) {
				case painter.BgRect:
					actual, ok := ops[i].(painter.BgRect)
					require.True(t, ok, "operation %d should be BgRect", i)
					assert.Equal(t, eop.X1, actual.X1, "X1 mismatch")
					assert.Equal(t, eop.Y1, actual.Y1, "Y1 mismatch")
					assert.Equal(t, eop.X2, actual.X2, "X2 mismatch")
					assert.Equal(t, eop.Y2, actual.Y2, "Y2 mismatch")
				case painter.WhiteFill:
					_, ok := ops[i].(painter.WhiteFill)
					assert.True(t, ok, "operation %d should be WhiteFill", i)
				case painter.GreenFill:
					_, ok := ops[i].(painter.GreenFill)
					assert.True(t, ok, "operation %d should be GreenFill", i)
				case painter.Figure:
					actual, ok := ops[i].(painter.Figure)
					require.True(t, ok, "operation %d should be Figure", i)
					assert.Equal(t, eop.X, actual.X, "X coordinate mismatch")
					assert.Equal(t, eop.Y, actual.Y, "Y coordinate mismatch")
				case painter.Move:
					actual, ok := ops[i].(painter.Move)
					require.True(t, ok, "operation %d should be Move", i)
					assert.Equal(t, eop.DX, actual.DX, "DX mismatch")
					assert.Equal(t, eop.DY, actual.DY, "DY mismatch")
				case painter.Reset:
					_, ok := ops[i].(painter.Reset)
					assert.True(t, ok, "operation %d should be Reset", i)
				default:
					if expectedOp == painter.UpdateOp {
						assert.Equal(t, painter.UpdateOp, ops[i], "operation %d should be UpdateOp", i)
					} else {
						t.Fatalf("unexpected operation type %T", expectedOp)
					}
				}
			}
		})
	}
}
