// Package loader provides functionality for loading and saving TIS-100 programs and puzzles.
// It supports serializing node code to `.tis` files and parsing puzzle definitions written in Lua.
// This includes extracting puzzle metadata, input/output stream definitions, and node layouts.
package loader

import (
	"errors"
	"fmt"

	"github.com/yuin/gopher-lua"

	"github.com/lekomish/tis-100/internal/model"
)

// LoadPuzzle loads and executes a Lua puzzle definition file and extracts
// the puzzle's title, description, streams, and layout by calling predefined Lua functions.
func LoadPuzzle(filePath string) (*model.Puzzle, error) {
	lState := lua.NewState()
	defer lState.Close()

	if err := lState.DoFile(filePath); err != nil {
		return nil, fmt.Errorf("unable to load lua script %s: %w", filePath, err)
	}

	// call individual fetchers for puzzle metadata and components
	title, err := fetchTitle(lState)
	if err != nil {
		return nil, err
	}
	description, err := fetchDescription(lState)
	if err != nil {
		return nil, err
	}
	streams, err := fetchStreams(lState)
	if err != nil {
		return nil, err
	}
	layout, err := fetchLayout(lState)
	if err != nil {
		return nil, err
	}

	return &model.Puzzle{
		Title:       title,
		Description: description,
		Streams:     streams,
		Layout:      layout,
	}, nil
}

// fetchTitle retrieves the puzzle title by calling the Lua function `GetTitle`.
func fetchTitle(lState *lua.LState) (string, error) {
	val, err := runLuaFunction(lState, "GetTitle")
	if err != nil {
		return "", err
	}
	return mustString(val, "GetTitle result")
}

// fetchDescription retrieves the puzzle description by calling `GetDescription`
// and converting the returned Lua table to a slice of strings.
func fetchDescription(lState *lua.LState) ([]string, error) {
	val, err := runLuaFunction(lState, "GetDescription")
	if err != nil {
		return nil, err
	}

	descTable, err := mustTable(val, "GetDescription result")
	if err != nil {
		return nil, err
	}

	var desc []string
	var iterErr error
	// iterate over each element in the Lua table
	descTable.ForEach(func(_, value lua.LValue) {
		if iterErr != nil {
			return
		}
		line, err := mustString(value, "description item")
		if err != nil {
			iterErr = err
			return
		}
		desc = append(desc, line)
	})

	if iterErr != nil {
		return nil, iterErr
	}
	return desc, nil
}

// fetchStreams retrieves a list of stream definitions by calling `GetStreams` in Lua.
// Each stream must be a 4-element table containing stream metadata and values.
func fetchStreams(lState *lua.LState) ([]*model.Stream, error) {
	val, err := runLuaFunction(lState, "GetStreams")
	if err != nil {
		return nil, err
	}

	streamsTable, err := mustTable(val, "GetStreams result")
	if err != nil {
		return nil, err
	}

	var streams []*model.Stream
	var iterErr error
	// iterate over each stream entry
	streamsTable.ForEach(func(_, streamVal lua.LValue) {
		sTab, err := mustTable(streamVal, "stream")
		if err != nil {
			iterErr = err
			return
		}
		if sTab.Len() != model.IOPositionsNumber {
			iterErr = fmt.Errorf(
				"stream must have %d elements, got %d",
				model.IOPositionsNumber,
				streamsTable.Len(),
			)
			return
		}

		// parse individaul stream fields
		typeNum, err := mustNumber(sTab.RawGetInt(1), "stream[1]")
		if err != nil || typeNum < 0 || typeNum >= model.StreamTypesNumber {
			iterErr = errors.New("stream[1] is not a valid StreamType")
			return
		}

		name, err := mustString(sTab.RawGetInt(2), "stream[2]")
		if err != nil {
			iterErr = err
			return
		}

		posNum, err := mustNumber(sTab.RawGetInt(3), "stream[3]")
		if err != nil || posNum < 0 || posNum > model.IOPositionsNumber {
			iterErr = fmt.Errorf("stream[3] out of range (0-%d)", model.IOPositionsNumber-1)
			return
		}

		valuesTable, err := mustTable(sTab.RawGetInt(4), "stream[4]")
		if err != nil {
			iterErr = err
			return
		}
		if valuesTable.Len() > model.MaxStreamValuesLength {
			iterErr = fmt.Errorf("stream[4]: too many values (max %d)", model.MaxStreamValuesLength)
			return
		}

		// collect individual stream values
		var values []int16
		valuesTable.ForEach(func(_, val lua.LValue) {
			if iterErr != nil {
				return
			}
			num, err := mustNumber(val, "stream[4]: stream value")
			if err != nil {
				iterErr = err
				return
			}
			if num < model.MinACC || num > model.MaxACC {
				iterErr = fmt.Errorf(
					"stream[4] value out of range (%d to %d)",
					model.MinACC,
					model.MaxACC,
				)
				return
			}
			values = append(values, int16(num))
		})

		if iterErr != nil {
			return
		}

		streams = append(streams, &model.Stream{
			Type:     model.StreamType(typeNum),
			Name:     name,
			Position: uint8(posNum),
			Values:   values,
		})
	})

	if iterErr != nil {
		return nil, iterErr
	}
	return streams, nil
}

// fetchLayout retrievese the puzzle's node layout by calling the Lua function `GetLayout`.
// It expects a table with a number of entries equal to `model.NodesNumber`.
func fetchLayout(lState *lua.LState) ([]model.NodeType, error) {
	val, err := runLuaFunction(lState, "GetLayout")
	if err != nil {
		return nil, err
	}

	layoutTable, err := mustTable(val, "GetLayout result")
	if err != nil {
		return nil, err
	}
	if layoutTable.Len() != model.NodesNumber {
		return nil, fmt.Errorf(
			"layout: expected %d items, got %d",
			model.NodesNumber,
			layoutTable.Len(),
		)
	}

	var layout []model.NodeType
	var iterErr error
	// parse each node type
	layoutTable.ForEach(func(_, val lua.LValue) {
		if iterErr != nil {
			return
		}
		num, err := mustNumber(val, "layout item")
		if err != nil {
			iterErr = err
			return
		}
		if num < 0 || num >= model.NodeTypesNumber {
			iterErr = fmt.Errorf("layout value out of range (0-%d)", model.NodeTypesNumber-1)
			return
		}
		layout = append(layout, model.NodeType(num))
	})

	if iterErr != nil {
		return nil, iterErr
	}
	return layout, nil
}
