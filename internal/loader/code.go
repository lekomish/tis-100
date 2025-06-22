package loader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lekomish/tis-100/internal/model"
)

const (
	fileExtension = ".tis"
	nodePrefix    = "@"
)

// SaveCode saves a Code object to a `.tis` file in the specified directory.
// Each node's instructions are written under a header like "@1", "@2", etc.
// Returns the full path to the created file or an error if the operation fails.
func SaveCode(dirPath string, code *model.Code) (string, error) {
	// ensure the target directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", dirPath)
	}

	// create file path using lowercased title
	fileName := fmt.Sprintf(
		"%s%s",
		strings.ReplaceAll(strings.ToLower(code.Title), "-", "_"),
		fileExtension,
	)
	filePath := filepath.Join(dirPath, fileName)

	// create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	// write each node's code
	for i, node := range code.Nodes {
		// check for too many nodes
		if i >= model.NodesNumber {
			return "", fmt.Errorf(
				"too many nodes (%d), expected max %d",
				len(code.Nodes),
				model.NodesNumber,
			)
		}
		// write node header
		if _, err := fmt.Fprintf(writer, "%s%d\n", nodePrefix, i+1); err != nil {
			return "", wrapWriterError("node header", filePath, err)
		}
		// write each instruction line
		for _, line := range node {
			if _, err := writer.WriteString(line + "\n"); err != nil {
				return "", wrapWriterError("node line", filePath, err)
			}
		}
		// write a newline to separate nodes
		if _, err := writer.WriteString("\n"); err != nil {
			return "", wrapWriterError("node separator", filePath, err)
		}
	}

	// flush buffered writer
	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush data to file %s: %w", filePath, err)
	}

	return filePath, nil
}

// LoadCode loads a `.tis` file into a Code object.
// The file is expected to have node sections prefixed with "@1", "@2", etc.
// Returns a parsed Code instance or an error if loading fails.
func LoadCode(filePath string) (*model.Code, error) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nodes := make([][]string, model.NodesNumber)

	curNode := -1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip blank lines
		if line == "" {
			continue
		}

		// detect node header
		if strings.HasPrefix(line, nodePrefix) {
			curNode++
			if curNode >= model.NodesNumber {
				return nil, fmt.Errorf("too many node headers in file %s", filePath)
			}
			continue
		}

		// chech that a node header has been seen before lines
		if curNode < 0 {
			return nil, fmt.Errorf("code line found before any node header in file %s", filePath)
		}

		// append instruction to current node
		nodes[curNode] = append(nodes[curNode], line)
	}

	// check for scanning error
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from file %s: %w", filePath, err)
	}

	// extract and format title from file name
	title := strings.ReplaceAll(
		strings.ToUpper(strings.TrimSuffix(filepath.Base(filePath), fileExtension)),
		"_",
		"-",
	)

	return &model.Code{
		Title: title,
		Nodes: nodes,
	}, nil
}
