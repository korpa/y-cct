package rousego

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/korpa/y-cct/global/signalhandler"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vipercfg *viper.Viper

// rootCmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "rousego",
	Short: "rousego startes multiple processes in parallel",
	Long: `rousego startes multiple processes in parallel


Processes have to be defined in rousego.toml

Example rousego.toml
----------------------------------------

[[cmds]]
label = "Backend2"
cmd = "sleep 3 ; echo 'hello from backend via stdout'; sleep 4"


[[cmds]]
label = "Frontend"
cmd = "sleep 2 ; echo 'hello from frontend via stdout' ; sleep 1"

[[cmds]]
label = "Frontend2"
cmd = "sleep 1 ; echo 'hello from frontend2 via stderr' 1>&2 ; sleep 1"

----------------------------------------

`,
	// BashCompletionFunction: bashCompletionFunc,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run(); err != nil {
			return err
		}
		return nil
	},
}

var colors []string

func init() {
	colors = append(colors, "#BB00BB")
	colors = append(colors, "#00BBBB")
	colors = append(colors, "#BBBB00")
	colors = append(colors, "#00BB00")
}

var processes []*process

type process struct {
	Name    string
	Cmd     *exec.Cmd
	Running bool
	Command string
	Style   lipgloss.Style
}

type cfgMain struct {
	Cmds []cfgCommands `toml:"cmds"`
}

type cfgCommands struct {
	Label string `toml:"label"`
	Cmd   string `toml:"cmd"`
}

func run() error {
	file, err := os.Open("rousego.toml")
	if err != nil {
		return errors.New("no rousego.toml found")
	}
	defer file.Close()

	var cfg cfgMain

	b, err := io.ReadAll(file)
	if err != nil {
		if err != nil {
			return errors.New("could not read rousego.toml")
		}
	}

	err = toml.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}

	if len(cfg.Cmds) < 1 {
		return errors.New("no commands defined in rousego.toml")
	}

	ctx := context.Background()
	c, cancel := signalhandler.Init(ctx)

	var wgProcesses sync.WaitGroup
	var wgMain sync.WaitGroup

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	for i, cmd := range cfg.Cmds {
		p := runCmd(&wgProcesses, int64(i), cmd.Label, cmd.Cmd)
		processes = append(processes, p)
	}

	wgMain.Add(1)
	go func() {
		d := 2 * time.Second
		for {
			select {

			case <-time.After(d):
				aliveMessage()
			case <-ctx.Done():
				for _, p := range processes {
					go p.shutdown()
				}
				return
				// default:
			}
		}
	}()

	// Wait untill all started process have stopped, then release main waitgroup
	go func() {
		wgProcesses.Wait()
		time.Sleep(300 * time.Millisecond)
		slog.Info("All processes stopped")
		wgMain.Done()
	}()
	wgMain.Wait()

	slog.Info("Stopping main process")
	return nil
}

func aliveMessage() {
	// slog.Info("Still alive")
}

func runCmd(wg *sync.WaitGroup, i int64, name string, command string) *process {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colors[i]))
		// Background(lipgloss.Color("#333333"))
	wg.Add(1)
	run := process{Name: name, Command: command, Running: true, Style: style}
	slog.Info("Starting: " + style.Render("["+name+"]") + " " + command)

	go func() {
		var args []string
		args = append(args, "-c")
		args = append(args, command)
		// fmt.Printf("%+v\n", args)
		// cmd := exec.CommandContext(ctx, c[0], strings.Split(c[1], " ")...)
		cmd := exec.Command("sh", args...)
		run.Cmd = cmd
		// cmd.Dir = ""
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()

		err := cmd.Start()
		if err != nil {
			// Run could also return this error and push the program
			// termination decision to the `main` method.
			log.Fatal(err)
		}

		go func() {
			scanner := bufio.NewScanner(stderr)
			// scanner.Split(bufio.ScanWords)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(style.Render("["+name+"]") + " " + m)
			}
		}()
		go func() {
			scanner := bufio.NewScanner(stdout)
			// scanner.Split(bufio.ScanWords)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(style.Render("["+name+"]") + " " + m)
			}
		}()

		err = cmd.Wait()
		if err != nil {
			slog.Warn("Stopped: " + style.Render("["+name+"]") + " via: " + fmt.Sprint(err))
		}

		run.Running = false
		time.Sleep(1 * time.Second)
		wg.Done()
		slog.Info("Finished: " + style.Render("["+name+"]"))
	}()
	return &run
}

func (p process) shutdown() {
	if !p.Running {
		slog.Warn("Shutting down " + p.Style.Render("["+p.Name+"]") + ": nothing todo. Process already finished")
		return
	}
	slog.Info("Shutting down " + p.Style.Render("["+p.Name+"]"))
	err := p.Cmd.Process.Kill()
	if err != nil {
		slog.Error(fmt.Sprint(err))
	}
}
