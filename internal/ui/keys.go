package ui

import (
	"os"

	"golang.org/x/sys/unix"
)

type keyKind int

const (
	keyUnknown keyKind = iota
	keyCancel
	keyEnter
	keyBackspace
	keyUp
	keyDown
	keyRune
)

type keyEvent struct {
	kind keyKind
	r    rune
}

const (
	byteCtrlC     = 3
	byteEscape    = 27
	byteEnter     = 13
	byteBackspace = 127

	byteCtrlJ = 10
	byteCtrlN = 14
	byteCtrlK = 11
	byteCtrlP = 16

	byteArrowPrefix = 91
	byteArrowUp     = 65
	byteArrowDown   = 66
)

func readKeyEvent() (keyEvent, error) {
	b := make([]byte, 3)
	n, err := os.Stdin.Read(b)
	if err != nil {
		return keyEvent{}, err
	}

	// Arrow keys are typically delivered as 3 bytes.
	if n == 3 && b[0] == byteEscape && b[1] == byteArrowPrefix {
		switch b[2] {
		case byteArrowUp:
			return keyEvent{kind: keyUp}, nil
		case byteArrowDown:
			return keyEvent{kind: keyDown}, nil
		}
	}

	if n == 1 {
		switch b[0] {
		case byteCtrlC:
			return keyEvent{kind: keyCancel}, nil
		case byteEscape:
			handled, ev, err := tryReadEscapeSequence()
			if err != nil {
				return keyEvent{}, err
			}
			if handled {
				return ev, nil
			}
			return keyEvent{kind: keyCancel}, nil
		case byteEnter:
			return keyEvent{kind: keyEnter}, nil
		case byteBackspace:
			return keyEvent{kind: keyBackspace}, nil
		case byteCtrlJ, byteCtrlN:
			return keyEvent{kind: keyDown}, nil
		case byteCtrlK, byteCtrlP:
			return keyEvent{kind: keyUp}, nil
		default:
			if b[0] >= 32 && b[0] < 127 {
				return keyEvent{kind: keyRune, r: rune(b[0])}, nil
			}
		}
	}

	return keyEvent{kind: keyUnknown}, nil
}

func tryReadEscapeSequence() (bool, keyEvent, error) {
	fd := int(os.Stdin.Fd())
	ready, err := hasPendingInput(fd)
	if err != nil {
		return false, keyEvent{}, err
	}
	if !ready {
		return false, keyEvent{}, nil
	}
	extra := make([]byte, 2)
	total := 0
	for total < len(extra) {
		n, err := os.Stdin.Read(extra[total:])
		if err != nil {
			return false, keyEvent{}, err
		}
		total += n
		if n == 0 {
			break
		}
	}
	if total >= 2 && extra[0] == '[' {
		switch extra[1] {
		case 'A':
			return true, keyEvent{kind: keyUp}, nil
		case 'B':
			return true, keyEvent{kind: keyDown}, nil
		}
	}
	return false, keyEvent{}, nil
}

func hasPendingInput(fd int) (bool, error) {
	var set unix.FdSet
	fdSet(fd, &set)
	var tv unix.Timeval
	n, err := unix.Select(fd+1, &set, nil, nil, &tv)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func fdSet(fd int, set *unix.FdSet) {
	set.Bits[fd/64] |= 1 << (uint(fd) % 64)
}
