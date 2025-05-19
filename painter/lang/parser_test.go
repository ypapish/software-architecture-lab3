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
			name:  "move command updates figure position",
			input: "figure 0.5 0.5\nmove 0.1 0.1\nupdate",
			expected: []painter.Operation{
				painter.Figure{X: 0.6, Y: 0.6},
				painter.UpdateOp,
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
				painter.Figure{X: 0.6, Y: 0.6},
				painter.UpdateOp,
			},
		},
		{
			name:    "invalid command",
			input:   "invalidcommand",
			wantErr: "unknown command: invalidcommand",
		},
		{	
			name: "reset clears previous figures and background",
			input: `white
					figure 0.1 0.1
					reset
					figure 0.5 0.5
					update`,
			expected: []painter.Operation{
				painter.Reset{},
				painter.Figure{X: 0.5, Y: 0.5},
				painter.UpdateOp,
			},
		},
		{
			name: "background command overriding previous background",
			input: `white
					bgrect 0 0 1 1
					green
					update`,
			expected: []painter.Operation{
				painter.BgRect{X1: 0, Y1: 0, X2: 1, Y2: 1},
				painter.GreenFill{},
				painter.UpdateOp,
			},
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

			if !assert.Len(t, ops, len(tc.expected), "number of operations mismatch") {
				return
			}

			for i, expectedOp := range tc.expected {
				actualOp := ops[i]
				switch eop := expectedOp.(type) {
				case painter.BgRect:
					actual, ok := actualOp.(painter.BgRect)
					require.True(t, ok, "operation %d should be BgRect, got %T", i, actualOp)
					assert.Equal(t, eop.X1, actual.X1, "BgRect X1 mismatch at op %d", i)
					assert.Equal(t, eop.Y1, actual.Y1, "BgRect Y1 mismatch at op %d", i)
					assert.Equal(t, eop.X2, actual.X2, "BgRect X2 mismatch at op %d", i)
					assert.Equal(t, eop.Y2, actual.Y2, "BgRect Y2 mismatch at op %d", i)
				case painter.WhiteFill:
					_, ok := actualOp.(painter.WhiteFill)
					assert.True(t, ok, "operation %d should be WhiteFill, got %T", i, actualOp)
				case painter.GreenFill:
					_, ok := actualOp.(painter.GreenFill)
					assert.True(t, ok, "operation %d should be GreenFill, got %T", i, actualOp)
				case painter.Figure:
					actual, ok := actualOp.(painter.Figure)
					require.True(t, ok, "operation %d should be Figure, got %T", i, actualOp)
					assert.InDelta(t, eop.X, actual.X, 0.0001, "Figure X coordinate mismatch at op %d", i)
					assert.InDelta(t, eop.Y, actual.Y, 0.0001, "Figure Y coordinate mismatch at op %d", i)
				case painter.Reset:
					_, ok := actualOp.(painter.Reset)
					assert.True(t, ok, "operation %d should be Reset, got %T", i, actualOp)
				default:
					if expectedOp == painter.UpdateOp {
						assert.Equal(t, painter.UpdateOp, actualOp, "operation %d should be UpdateOp", i)
					} else {
						t.Fatalf("unexpected expected operation type %T at index %d, value: %+v", expectedOp, i, expectedOp)
					}
				}
			}
		})
	}
}
