package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wilmoore/p/internal/clierr"
	"github.com/wilmoore/p/internal/history"
	"github.com/wilmoore/p/internal/i18n"
	"github.com/wilmoore/p/internal/tmux"
	"github.com/wilmoore/p/internal/ui"
)

// Version is set at build time via -ldflags "-X main.Version=vX.Y.Z"
var Version = "dev"

const usage = `p - minimal tmux session switcher

Usage:
  p                          Show interactive session selector
  p <path>                   Create new session in directory (use . for current directory)
  p <path> --name <custom>   Create session with a custom name
  p --log                    Browse session history ledger
  p --version                Show version information
  p --help                   Show this help message

Navigation:
  Type           Filter sessions by name
  Arrow keys     Navigate up/down
  Enter          Attach to selected session
  Esc/Ctrl+C     Cancel

Examples:
  p              Select from existing sessions
  p .            Create session in current directory
  p ~/projects   Create session in ~/projects
  p ./revenue --name savvy-revenue
  p --log        Inspect or relaunch recent sessions
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, clierr.Format(err))
		os.Exit(1)
	}
}

func run() error {
	cmd, err := parseArgs(os.Args[1:])
	if err != nil {
		return err
	}

	switch cmd.kind {
	case commandVersion:
		fmt.Println(Version)
		return nil
	case commandHelp:
		fmt.Print(usage)
		return nil
	case commandHistory:
		return showHistory()
	case commandCreate:
		return createSessionFromPath(cmd.path, cmd.sessionName)
	case commandSelector:
		return showSessionSelector()
	default:
		return fmt.Errorf("unknown command")
	}
}

func showSessionSelector() error {
	sessions, err := tmux.ListSessions()
	if err != nil && !tmux.IsNoServerError(err) {
		return fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	if len(sessions) == 0 {
		return fmt.Errorf(i18n.ErrNoTmuxSessionsAvailable)
	}

	choice, err := ui.ShowSelector(sessions)
	if err != nil {
		return err
	}
	if choice == nil {
		return nil
	}
	return attachAndLog(choice.Name, history.ActionAttach)
}

func showHistory() error {
	entries, err := history.List(200)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println(i18n.MsgNoSessionHistoryYet)
		return nil
	}
	choice, err := ui.ShowHistory(entries)
	if err != nil {
		return err
	}
	if choice == nil {
		return nil
	}
	if choice.TargetDir == "" {
		return fmt.Errorf(i18n.ErrHistoryMissingTargetDir)
	}
	return createSessionFromPath(choice.TargetDir, choice.SessionName)
}

// createSessionFromPath creates (or attaches to) a tmux session in the specified directory.
func createSessionFromPath(path, overrideName string) error {
	spec, err := buildSessionSpec(path, overrideName)
	if err != nil {
		return err
	}
	action, err := tmux.CreateSession(spec.sessionName, spec.workingDir)
	if err != nil {
		return err
	}
	logAction := history.ActionCreate
	if action == tmux.LaunchActionAttachExisting {
		logAction = history.ActionAttachExisting
	}
	logLaunch(logAction, spec.sessionName, spec.workingDir)
	return tmux.AttachToSession(spec.sessionName)
}

type sessionSpec struct {
	sessionName string
	workingDir  string
}

func buildSessionSpec(path, overrideName string) (*sessionSpec, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}
	resolved, err := resolvePath(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory does not exist: %s", resolved)
		}
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", resolved)
	}
	sessionName := overrideName
	if sessionName == "" {
		sessionName = filepath.Base(resolved)
	}
	return &sessionSpec{sessionName: sessionName, workingDir: resolved}, nil
}

func resolvePath(path string) (string, error) {
	switch {
	case path == ".":
		return os.Getwd()
	case strings.HasPrefix(path, "~"):
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Abs(filepath.Join(home, path[1:]))
	default:
		return filepath.Abs(path)
	}
}

func attachAndLog(sessionName string, action history.Action) error {
	targetDir, err := tmux.GetSessionPath(sessionName)
	if err != nil {
		targetDir = ""
	}
	logLaunch(action, sessionName, targetDir)
	return tmux.AttachToSession(sessionName)
}

func logLaunch(action history.Action, sessionName, targetDir string) {
	if sessionName == "" {
		return
	}
	invokeDir, _ := os.Getwd()
	if targetDir != "" {
		abs, err := filepath.Abs(targetDir)
		if err == nil {
			targetDir = abs
		}
	}
	entry := history.Entry{
		Timestamp:   time.Now(),
		Action:      action,
		SessionName: sessionName,
		InvokeDir:   invokeDir,
		TargetDir:   targetDir,
	}
	if err := history.Append(entry); err != nil {
		fmt.Fprintf(os.Stderr, i18n.WarnWriteHistoryFailedFmt, err)
	}
}

type commandKind int

const (
	commandSelector commandKind = iota
	commandCreate
	commandHistory
	commandVersion
	commandHelp
)

type command struct {
	kind        commandKind
	path        string
	sessionName string
}

func parseArgs(args []string) (*command, error) {
	if len(args) == 0 {
		return &command{kind: commandSelector}, nil
	}
	switch args[0] {
	case "--version", "-v":
		return &command{kind: commandVersion}, nil
	case "--help", "-h":
		return &command{kind: commandHelp}, nil
	case "--log":
		if len(args) > 1 {
			return nil, fmt.Errorf("--log cannot be combined with other arguments")
		}
		return &command{kind: commandHistory}, nil
	}

	cmd := &command{kind: commandCreate, path: args[0]}
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--name=") {
			cmd.sessionName = strings.TrimPrefix(arg, "--name=")
			continue
		}
		if arg == "--name" {
			if i+1 >= len(args) {
				return nil, errors.New("--name requires a value")
			}
			cmd.sessionName = args[i+1]
			i++
			continue
		}
		return nil, fmt.Errorf("unknown option: %s", arg)
	}
	return cmd, nil
}
