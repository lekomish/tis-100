package loader_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/loader"
	"github.com/lekomish/tis-100/internal/model"
)

const (
	errUnexpectedMsg   = "unexpected error"
	errCreatingFileMsg = "error while creating test file"
)

/* TESTS */

// --- LoadPuzzle ---
func TestLoadPuzzleWithCorrectScript(t *testing.T) {
	filePath, err := setupLua(t, newScript(), "test_load_puzzle_with_correct_script.lua")
	require.NoError(t, err, errCreatingFileMsg)

	puzzle, err := loader.LoadPuzzle(filePath)
	require.NoError(t, err, errUnexpectedMsg)

	expectedPuzzle := model.Puzzle{
		Title:       "TEST",
		Description: []string{"TEST LINE 1", "TEST LINE 2"},
		Streams: []model.Stream{
			{
				Type:     model.INPUT,
				Name:     "IN.TEST",
				Position: 0,
				Values:   []int16{1, 2, 3},
			},
			{
				Type:     model.OUTPUT,
				Name:     "OUT.TEST",
				Position: 0,
				Values:   []int16{1, 2, 3},
			},
		},
		Layout: []model.NodeType{
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
			model.COMPUTE,
		},
	}

	require.Equal(t, expectedPuzzle, *puzzle, "puzzle does not match expected result")
}

func TestLoadPuzzleWithWrongScript(t *testing.T) {
	s := newScript()
	s.Title = []string{"func GetTitle()", "return", ";"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_script.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "unable to load lua script")
}

func TestLoadPuzzleWithoutFunction(t *testing.T) {
	s := newScript()
	s.Title = []string{""}

	filePath, err := setupLua(t, s, "test_load_puzzle_without_function.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "lua function \"GetTitle\" not found")
}

// --- Title errors ---
func TestLoadPuzzleWithWrongTitle(t *testing.T) {
	s := newScript()
	s.Title = []string{"function GetTitle()", "return {}", "end"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_title.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "GetTitle result: expected string,")
}

// --- Description errors ---
func TestLoadPuzzleWithWrongDescriptionType(t *testing.T) {
	s := newScript()
	s.Description = []string{"function GetDescription()", "return 1", "end"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_description_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "GetDescription result: expected table,")
}

func TestLoadPuzzleWithWrongDescriptionLineType(t *testing.T) {
	s := newScript()
	s.Description = []string{"function GetDescription()", "return { 1, 2 }", "end"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_description_line_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "description item: expected string,")
}

// --- Streams errors ---
func TestLoadPuzzleWithWrongStreamsType(t *testing.T) {
	s := newScript()
	s.Streams = []string{"function GetStreams()", "return 1", "end"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_streams_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "GetStreams result: expected table,")
}

func TestLoadPuzzleWithWrongStreamType(t *testing.T) {
	s := newScript()
	s.Streams = []string{"function GetStreams()", "return { 1, 2 }", "end"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream: expected table,")
}

func TestLoadPuzzleWithWrongStreamArgumentsNumber(t *testing.T) {
	s := newScript()
	s.Streams = []string{"function GetStreams()", "return { { 0, \"IN.TEST\", 0 } }", "end"}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_arguments_number.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream must have")
}

func TestLoadPuzzleWithWrongStreamTypeValueType(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { {}, \"IN.TEST\", 0, { 1, 2, 3 } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_type_value_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[1] is not a valid StreamType")
}

func TestLoadPuzzleWithWrongStreamTypeValue(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 5, \"IN.TEST\", 0, { 1, 2, 3 } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_type_value.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[1] is not a valid StreamType")
}

func TestLoadPuzzleWithWrongStreamNameType(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, {}, 0, { 1, 2, 3 } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_name_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[2]: expected string,")
}

func TestLoadPuzzleWithWrongStreamPositionType(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", \"0\", { 1, 2, 3 } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_position_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[3] out of range")
}

func TestLoadPuzzleWithWrongStreamPosition(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 5, { 1, 2, 3 } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_position.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[3] out of range")
}

func TestLoadPuzzleWithWrongStreamValuesType(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, 0 } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_items_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[4]: expected table,")
}

func TestLoadPuzzleWithWrongStreamItemsLength(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, { 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31 } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_items_length.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[4]: too many values")
}

func TestLoadPuzzleWithWrongStreamItemType(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, { \"1\", \"2\" } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_item_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[4]: stream value: expected number")
}

func TestLoadPuzzleWithWrongStreamItems(t *testing.T) {
	s := newScript()
	s.Streams = []string{
		"function GetStreams()",
		"return { { 0, \"IN.TEST\", 0, { 500000, 2, } } }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_stream_items.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "stream[4] value out of range")
}

// --- Layout errors ---
func TestLoadPuzzleWithWrongLayoutType(t *testing.T) {
	s := newScript()
	s.Layout = []string{
		"function GetLayout()",
		"return 1",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_layout_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "GetLayout result: expected table")
}

func TestLoadPuzzleWithWrongLayoutLength(t *testing.T) {
	s := newScript()
	s.Layout = []string{
		"function GetLayout()",
		"return { 0, 0, 0, 0 }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_layout_length.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "layout: expected 12 items")
}

func TestLoadPuzzleWithWrongLayoutItemsType(t *testing.T) {
	s := newScript()
	s.Layout = []string{
		"function GetLayout()",
		"return { \"0\", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_layout_items_type.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "layout item: expected number")
}

func TestLoadPuzzleWithWrongLayoutItems(t *testing.T) {
	s := newScript()
	s.Layout = []string{
		"function GetLayout()",
		"return { 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }",
		"end",
	}

	filePath, err := setupLua(t, s, "test_load_puzzle_with_wrong_layout_items.lua")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadPuzzle(filePath)
	require.ErrorContains(t, err, "layout value out of range")
}

// fetchTitle -> covered in previous tests
// fetchDescription -> covered in previous tests
// fetchStreams -> covered in previous tests
// fetchLayout -> covered in previous tests
// mustString -> covered in previous tests
// mustNumber -> covered in previous tests
// mustTable -> covered in previous tests
// runLuaFunction -> covered in previous tests

/* BENCHMARKS */

// BenchmarkLoadPuzzle measures the performance of loading and parsing a Lua puzzle file.
// It uses a predefined script and runs `loader.LoadPuzzle` in a tight loop.
func BenchmarkLoadPuzzle(b *testing.B) {
	s := newScript()
	filePath, err := setupLua(b, s, "bench_load_puzzle.lua")
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loader.LoadPuzzle(filePath)
		require.NoError(b, err)
	}
}

/* UTILS */

// script represents a Lua script split into sections for modular construction.
// This is useful for generating test inputs with different puzzle definitions.
type script struct {
	Beginning   []string // constants and common Lua definitions
	Title       []string // lua function `GetTitle`
	Description []string // lua function `GetDescription`
	Streams     []string // lua function `GetStreams`
	Layout      []string // lua function `GetLayout`
}

// newScript creates a default test Lua script with valid content.
func newScript() *script {
	return &script{
		Beginning: []string{
			"local STREAM_INPUT = 0",
			"local STREAM_OUTPUT = 1",
			"local TILE_COMPUTE = 0",
			"local TILE_DAMAGED = 1",
		},
		Title: []string{"function GetTitle()", "return \"TEST\"", "end"},
		Description: []string{
			"function GetDescription()",
			"return { \"TEST LINE 1\", \"TEST LINE 2\" }",
			"end",
		},
		Streams: []string{
			"function GetStreams()",
			"return {",
			"{ STREAM_INPUT, \"IN.TEST\", 0, { 1, 2, 3 } },",
			"{ STREAM_OUTPUT, \"OUT.TEST\", 0, { 1, 2, 3 } },",
			"}",
			"end",
		},
		Layout: []string{
			"function GetLayout()",
			"return {",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"TILE_COMPUTE,",
			"}",
			"end",
		},
	}
}

// ToSlice returns the full Lua script as a slice of lines in proper order.
func (s *script) ToSlice() []string {
	fullScript := append(s.Beginning, s.Title...)
	fullScript = append(fullScript, s.Description...)
	fullScript = append(fullScript, s.Streams...)
	return append(fullScript, s.Layout...)
}

// writeScriptToFile writes the script contents to the given file.
// Each line is written with a newline character to ensure proper Lua formatting.
func writeScriptToFile(file *os.File, s *script) error {
	for _, line := range s.ToSlice() {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}
	return nil
}

// setupLua creates a temporary Lua file containing the provided script,
// registers a cleanup for test teardown, and returns the file path.
func setupLua(tb testing.TB, s *script, fileName string) (string, error) {
	tb.Helper()

	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := writeScriptToFile(file, s); err != nil {
		return "", err
	}

	tb.Cleanup(func() {
		os.Remove(file.Name())
	})

	return file.Name(), nil
}
