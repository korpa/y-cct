package retention

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

var (
	version = "development"
	goos    = "unknown"
	goarch  = "unknown"
	commit  = "unknown"
)

var (
	show       string
	verbose    bool
	dateFormat string
	dateRegex  string
)

type Entry struct {
	line   string
	date   time.Time
	keep   bool
	reason string
}

var Cmd = &cobra.Command{
	Use:   "retention [filename]",
	Short: "Processes entries with dates from a file or stdin",
	Example: `# Delete old backups. Keeps one backup per day of last 4 weeks. Delete all backups
# older than 4 weeks, but keep Monday backups and 1st of month.
  ls | retention -f 2006-01-02 | xargs rm -r

> Filename examples
> backup_2025-09-01.tar.gz
> backup_2025-09-02.tar.gz
> .
> .
> backup_2025-09-29.tar.gz
> backup_2025-09-30.tar.gz
> backup_2025-10-01.tar.gz
> backup_2025-10-02.tar.gz
> .
> .
> backup_2025-10-30.tar.gz
> backup_2025-10-31.tar.gz
> backup_2025-11-01.tar.gz
> backup_2025-11-02.tar.gz

`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var reader io.Reader

		if len(args) == 1 {
			file, err := os.Open(args[0])
			if err != nil {
				slog.Error("Failed to open file", "error", err)
				os.Exit(1)
			}
			defer file.Close()
			reader = file
		} else {
			stat, err := os.Stdin.Stat()
			if err != nil {
				slog.Error("Failed to read from stdin", "error", err)
				os.Exit(1)
			}
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				slog.Error("No file provided and no piped input")
				os.Exit(1)
			}
			reader = os.Stdin
		}

		if show == "both" {
			verbose = true
		}

		if dateRegex == "" {
			dateRegex = layoutToRegex(dateFormat)
		}

		entries, err := readEntries(reader, dateFormat, dateRegex)
		if err != nil {
			slog.Error("Failed to read entries", "error", err)
			os.Exit(1)
		}

		entries = classifyEntries(entries)
		entries = reduceToOnePerDay(entries, time.Now().AddDate(0, 0, -28))

		filtered := filterEntries(entries, show)
		printEntries(filtered, verbose)
	},
}

func init() {
	Cmd.Flags().StringVarP(&show, "show", "s", "delete", "What to show: keep, delete, both")
	Cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output (show reasons)")
	Cmd.Flags().StringVarP(&dateFormat, "date-format", "f", "2006-01-02_15-04-05", "Date format for parsing timestamps (golang date format)")
	Cmd.Flags().StringVar(&dateRegex, "date-regex", "", "Regex pattern to extract the timestamp (optional, auto-generated from format if empty)")
}

func appendReason(orig, add string) string {
	if orig == "" {
		return add
	}
	return orig + "; " + add
}

func readEntries(reader io.Reader, layout string, regexPattern string) ([]Entry, error) {
	re := regexp.MustCompile(regexPattern)

	var entries []Entry
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindString(line)
		if match == "" {
			slog.Warn("No date found in line", "line", line)
			continue
		}

		parsedTime, err := time.Parse(layout, match)
		if err != nil {
			slog.Warn("Failed to parse timestamp", "timestamp", match, "error", err)
			continue
		}

		entries = append(entries, Entry{line: line, date: parsedTime})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func classifyEntries(entries []Entry) []Entry {
	now := time.Now()
	fourWeeksAgo := now.AddDate(0, 0, -28)

	for i := range entries {
		e := &entries[i]
		if e.date.Weekday() == time.Monday {
			e.keep = true
			e.reason = appendReason(e.reason, "isMonday")
		}
		if e.date.Day() == 1 {
			e.keep = true
			e.reason = appendReason(e.reason, "isFirstOfMonth")
		}
		if e.date.After(fourWeeksAgo) {
			e.keep = true
			e.reason = appendReason(e.reason, "isWithinFourWeeks")
		}
		if !e.keep {
			e.reason = appendReason(e.reason, "olderThanFourWeeks")
		}
	}

	return entries
}

func reduceToOnePerDay(entries []Entry, cutoff time.Time) []Entry {
	latestPerDay := make(map[string]int)

	for i, e := range entries {
		dayKey := e.date.Format("2006-01-02")
		if e.date.After(cutoff) || !e.keep {
			continue
		}

		if idx, exists := latestPerDay[dayKey]; !exists {
			latestPerDay[dayKey] = i
		} else {
			if e.date.After(entries[idx].date) {
				entries[idx].keep = false
				entries[idx].reason = appendReason(entries[idx].reason, "deletedLaterEntrySameDay")
				latestPerDay[dayKey] = i
			} else {
				entries[i].keep = false
				entries[i].reason = appendReason(entries[i].reason, "deletedLaterEntrySameDay")
			}
		}
	}

	return entries
}

func filterEntries(entries []Entry, show string) []Entry {
	var filtered []Entry
	for _, e := range entries {
		switch show {
		case "keep":
			if e.keep {
				filtered = append(filtered, e)
			}
		case "delete":
			if !e.keep {
				filtered = append(filtered, e)
			}
		case "both":
			fallthrough
		default:
			filtered = append(filtered, e)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].date.Before(filtered[j].date)
	})

	return filtered
}

func printEntries(entries []Entry, verbose bool) {
	for _, e := range entries {
		if verbose {
			fmt.Printf("%v | keep=%t | %s | %s\n", e.date, e.keep, e.line, e.reason)
		} else {
			fmt.Printf("%s\n", e.line)
		}
	}
}

// layoutToRegex converts Go time layout to a regex pattern to extract timestamps
func layoutToRegex(layout string) string {
	replacements := []struct {
		token string
		regex string
	}{
		{"2006", `\d{4}`},
		{"01", `\d{2}`},
		{"02", `\d{2}`},
		{"15", `\d{2}`},
		{"04", `\d{2}`},
		{"05", `\d{2}`},
		{"06", `\d{2}`},
		{"3", `\d{1,2}`},
		{"PM", `(AM|PM)`},
		{"pm", `(am|pm)`},
		{"Mon", `[A-Z][a-z]{2}`},
		{"Monday", `[A-Z][a-z]+`},
		{"Jan", `[A-Z][a-z]{2}`},
		{"January", `[A-Z][a-z]+`},
		{".", `\.`},
		{"-", `-`},
		{":", `:`},
		{" ", `\s`},
		{"_", `_`},
	}

	regex := layout
	for _, r := range replacements {
		regex = replaceAllLiteral(regex, r.token, r.regex)
	}

	return regex
}

func replaceAllLiteral(s, token, replacement string) string {
	out := ""
	i := 0
	for i < len(s) {
		if len(s[i:]) >= len(token) && s[i:i+len(token)] == token {
			out += replacement
			i += len(token)
		} else {
			out += string(s[i])
			i++
		}
	}
	return out
}
