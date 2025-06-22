package loader_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"

	"github.com/lekomish/tis-100/internal/loader"
	"github.com/lekomish/tis-100/internal/model"
)

/* TESTS */

// --- SaveCode ---
func TestSaveCodeWithCorrectInput(t *testing.T) {
	code := newCode("TEST-SAVE-CODE-WITH-CORRECT-INPUT")

	dirPath, err := setupDir(t, "test_save_code")
	require.NoError(t, err, errCreatingFileMsg)

	filePath, err := loader.SaveCode(dirPath, code)
	require.NoError(t, err, errUnexpectedMsg)
	require.Equal(
		t,
		filePath,
		filepath.Join(dirPath, "test_save_code_with_correct_input.tis"),
		"unexpected path for file",
	)

	fileContent, err := os.ReadFile(filePath)
	require.NoError(t, err)
	content := string(fileContent)
	require.Equal(t, content, codeToString(code.Nodes), "data in file is not equal expected result")
}

func TestSaveCodeWithNotExistingDirectory(t *testing.T) {
	code := newCode("TEST-SAVE-CODE-WITH-NOT-EXISTING-DIRECTORY")

	_, err := loader.SaveCode("/notexistingdirectory", code)
	require.ErrorContains(t, err, "directory does not exist:")
}

// -- LoadCode ---

func TestLoadCodeWithCorrectInput(t *testing.T) {
	expectedCode := newCode("TEST-LOAD-CODE-WITH-CORRECT-INPUT")
	filePath, err := setupCode(t, expectedCode, "test_load_code_with_correct_input")
	require.NoError(t, err, errCreatingFileMsg)

	code, err := loader.LoadCode(filePath)

	// remove digits from tmp file name
	var builder strings.Builder
	for _, r := range code.Title {
		if !unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	code.Title = builder.String()

	require.NoError(t, err, errUnexpectedMsg)
	require.Equal(t, expectedCode, code)
}

func TestLoadCodeWithNotExistingFile(t *testing.T) {
	_, err := loader.LoadCode("/notexistingfile.tis")
	require.ErrorContains(t, err, "failed to open file")
}

func TestLoadCodeWithWrongHeadersNumber(t *testing.T) {
	code := newCode("TEST-LOAD-CODE-WITH-WRONG-HEADERS-NUMBER")
	code.Nodes = append(code.Nodes, []string{"MOV UP DOWN"})

	filePath, err := setupCode(t, code, "test_load_code_with_wrong_headers_number")
	require.NoError(t, err, errCreatingFileMsg)

	_, err = loader.LoadCode(filePath)
	require.ErrorContains(t, err, "too many node headers")
}

// wrapWriterError -> covered in previous tests

/* BENCHMARKS */

// BenchmarkSaveCode measures the performance of saving a Code struct to disk
// using the SaveCode function. It writes each iteration to a unique filename
// within a temporary directory to avoid overwritting or OS-level caching.
func BenchmarkSaveCode(b *testing.B) {
	code := newCode("bench_save_code")
	dir, err := setupDir(b, "bench_save_code_dir")
	require.NoError(b, err)
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code.Title = fmt.Sprintf("bench_save_code_%d", i)
		_, err := loader.SaveCode(dir, code)
		require.NoError(b, err)
	}
}

// BenchmarkLoadCode measures the performance of loading a `.tis` code file
// using the LoadCode function. The file is written once before benchmarking,
// and then read repeatedly during the timed portion of the benchmark.
func BenchmarkLoadCode(b *testing.B) {
	code := newCode("bench_load_code")
	file, err := setupCode(b, code, "bench_load_code")
	require.NoError(b, err)
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loader.LoadCode(file)
		require.NoError(b, err)
	}
}

/* UTILS */

// setupDir creates a temporary directory with the given name,
// registers it for automatic cleanup after the test, and returns its path.
func setupDir(tb testing.TB, dirName string) (string, error) {
	tb.Helper()

	dir, err := os.MkdirTemp("", dirName)
	if err != nil {
		return "", err
	}

	tb.Cleanup(func() { os.RemoveAll(dir) })

	return dir, err
}

// newCode returns a mock `*model.Code` object with the given title,
// where each node contains two identical test instructions.
func newCode(title string) *model.Code {
	nodes := make([][]string, model.NodesNumber)
	for i := range model.NodesNumber {
		nodes[i] = []string{"MOV UP DOWN", "MOV UP DOWN"}
	}
	return &model.Code{
		Title: title,
		Nodes: nodes,
	}
}

// codeToString converts a Code's node instructions into a `.tis` file format string.
// Each node is written an "@<index>" header followed by its instructions.
func codeToString(nodes [][]string) string {
	var builder strings.Builder
	for i, node := range nodes {
		fmt.Fprintf(&builder, "@%d\n", i+1)
		for _, line := range node {
			builder.WriteString(line + "\n")
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// setupCode writes a `model.Code` to a temporary `.tis` file using TIS-100 formatting,
// and returns the file path. It registers the file for automatic cleanup.
func setupCode(tb testing.TB, code *model.Code, fileName string) (string, error) {
	tb.Helper()

	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := fmt.Fprint(file, codeToString(code.Nodes)); err != nil {
		return "", err
	}

	tb.Cleanup(func() { os.Remove(file.Name()) })

	return file.Name(), nil
}
