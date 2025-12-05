package replacer

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type Replacement struct {
	RegExp string
	Line   string
}

type ReplacementCompiled struct {
	RegExp *regexp.Regexp
	Line   string
	Found  bool
}

var (
	file             string
	replacementsArgs []string
	appendMissing    bool
)

var Cmd = &cobra.Command{
	Use:   "replacer",
	Short: "Replace lines in a file using regexp patterns.",
	Long: `replacer is a CLI tool for replacing or appending lines in a file based on regular expression patterns.

For each --replace argument, replacer searches for a matching line. If a match is found, the line is replaced.
If no match is found, the replacement line can optionally be appended at the end of the file (this is the default).

- Each --replace must be in the format REGEXP:LINE.
- You can specify --replace multiple times.
- You are responsible for adding ^ (line start) or $ (line end) anchors to your regular expressions if desired.
- For replacements where the pattern is not found, the line is appended only if --append-missing is true (default: true).

`,
	Example: `# Replace all lines starting with "Hello" with "Greetings!"
  replacer --file text.txt --replace "^Hello:Greetings!"

  # Replace lines starting with "But" or "Hello"
  replacer --file text.txt --replace "^But:HOWEVER" --replace "^Hello:GOOD DAY"

  # Append the replacement line at the end if no match is found (default behavior)
  replacer --file text.txt --replace "^NotFound:New line"

  # Do NOT append if no match is found (replace only)
  replacer --file text.txt --replace "^NotFound:New line" --append-missing=false

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if file == "" {
			return fmt.Errorf("--file is required")
		}
		if len(replacementsArgs) == 0 {
			return fmt.Errorf("At least one --replace is required")
		}

		// Parse replacements from arguments
		replacements := []Replacement{}
		for _, repArg := range replacementsArgs {
			parts := strings.SplitN(repArg, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("Invalid --replace format: %q (expected REGEXP:LINE)", repArg)
			}
			replacements = append(replacements, Replacement{RegExp: parts[0], Line: parts[1]})
		}

		// Open the original file
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("Failed to open file %q: %w", file, err)
		}
		defer f.Close()

		// Create temp file
		tmp, err := os.CreateTemp(".", ".tmp-replace-*")
		if err != nil {
			return fmt.Errorf("Failed to create temp file: %w", err)
		}
		defer tmp.Close()

		// Replace while copying from f to tmp
		fileChanged, missingLines, err := replace(f, tmp, replacements)
		if err != nil {
			return fmt.Errorf("Replace failed: %w", err)
		}

		// If enabled: append missing lines at end
		if appendMissing && len(missingLines) > 0 {
			for _, line := range missingLines {
				if _, err := io.WriteString(tmp, line+"\n"); err != nil {
					return fmt.Errorf("Failed to append missing line: %w", err)
				}
				fileChanged = true
			}
		}

		// Make sure the temp file was successfully written to
		if err := tmp.Close(); err != nil {
			return fmt.Errorf("Failed to close temp file: %w", err)
		}

		// Close the input file
		if err := f.Close(); err != nil {
			return fmt.Errorf("Failed to close input file: %w", err)
		}

		if fileChanged {
			// Overwrite the original file with the temp file if changed
			if err := os.Rename(tmp.Name(), file); err != nil {
				return fmt.Errorf("Failed to overwrite original file: %w", err)
			}
		} else {
			slog.Info("File not changed - nothing to do")
		}

		// Delete temp file
		if err := os.Remove(tmp.Name()); err != nil {
			return fmt.Errorf("Failed to remove temp file: %w", err)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&file, "file", "f", "", "Path to the file to process (required)")
	Cmd.Flags().StringArrayVarP(&replacementsArgs, "replace", "r", []string{}, "Replacement in the format 'REGEXP:LINE'. Can be specified multiple times.")
	Cmd.Flags().BoolVar(&appendMissing, "append-missing", true, "Append replacement line at end of file if regexp did not match anywhere")
}

func replace(r io.Reader, w io.Writer, replacements []Replacement) (bool, []string, error) {
	fileChanged := false

	repls := []ReplacementCompiled{}
	for _, replacement := range replacements {
		re, err := regexp.Compile(replacement.RegExp)
		if err != nil {
			return fileChanged, nil, err
		}
		repls = append(repls, ReplacementCompiled{RegExp: re, Line: replacement.Line, Found: false})
	}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		lineChanged := false

		for i, replacement := range repls {
			if replacement.RegExp.MatchString(line) {
				fileChanged = true
				line = replacement.Line
				repls[i].Found = true
				lineChanged = true
				break
			}
		}

		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return fileChanged, nil, err
		}
		if lineChanged {
			// no-op, info kept in .Found
		}
	}
	// Collect all lines whose RegExp was not found:
	missingLines := []string{}
	for _, replacement := range repls {
		if !replacement.Found {
			missingLines = append(missingLines, replacement.Line)
		}
	}

	return fileChanged, missingLines, sc.Err()
}
