package tmux

import (
	"fmt"
	"os/exec"
	"strings"
)

type TmuxSession struct {
	Name    string
	Windows string
	Created string
	Size    string
}

func GetTmuxSessions() ([]TmuxSession, error) {
	cmd := exec.Command("tmux", "ls")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	sessions := []TmuxSession{}

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		rest := strings.TrimSpace(parts[1])

		session := TmuxSession{Name: name}

		if windowsIdx := strings.Index(rest, "windows"); windowsIdx != -1 {
			session.Windows = strings.TrimSpace(rest[:windowsIdx])
		}
		if createdIdx := strings.Index(rest, "(created"); createdIdx != -1 {
			endIdx := strings.Index(rest[createdIdx:], ")")
			if endIdx != -1 {
				session.Created = strings.TrimSpace(rest[createdIdx+9 : createdIdx+endIdx])
			}
		}
		if sizeIdx := strings.Index(rest, "["); sizeIdx != -1 {
			endIdx := strings.Index(rest[sizeIdx:], "]")
			if endIdx != -1 {
				session.Size = rest[sizeIdx+1 : sizeIdx+endIdx]
			}
		}

		if len(sessions) == 0 {
			return nil, fmt.Errorf("tmux not found or no sessions available")
		}
	}
	return sessions, nil
}
