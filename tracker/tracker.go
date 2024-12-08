package tracker

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

var (
	instance  *Tracker
	timeRegex = regexp.MustCompile(`\[(.*?)\]`)
	zoneRegex = regexp.MustCompile(`You have entered (.*)`)
)

type Tracker struct {
	path          string
	onLineEvent   []func(time.Time, string)
	onZoneEvent   []func(time.Time, string) // zone name
	isLiveParse   bool
	trackerStart  time.Time
	isStarted     bool
	name          string
	pollCtx       context.Context
	pollCtxCancel context.CancelFunc
}

func New(path string) (*Tracker, error) {
	if instance != nil {
		return nil, fmt.Errorf("tracker already initialized")
	}

	t := &Tracker{
		path:         path,
		trackerStart: time.Now(),
	}
	instance = t

	if path != "" {
		if !strings.Contains(path, "eqlog_") {
			return nil, fmt.Errorf("invalid log file (expected eqlog_ prefix)")
		}

		pos := strings.Index(path, "eqlog_")
		name := path[pos+6:]
		pos = strings.Index(name, "_")
		if pos > 0 {
			name = name[:pos]
		}
		instance.name = name
	}
	return t, nil
}

func (t *Tracker) Start(isFromStart bool) error {
	if instance == nil {
		return fmt.Errorf("tracker not initialized")
	}

	if t.isStarted {
		return fmt.Errorf("tracker already started")
	}
	t.isStarted = true

	err := t.watchFile()
	if err != nil {
		return fmt.Errorf("watch file: %w", err)
	}

	return nil
}

func SetNewPath(path string) error {
	if instance == nil {
		return fmt.Errorf("tracker not initialized")
	}

	t := instance
	if !strings.Contains(path, "eqlog_") {
		return fmt.Errorf("invalid log file (expected eqlog_ prefix)")
	}

	t.path = path

	if !strings.Contains(path, "eqlog_") {
		return fmt.Errorf("invalid log file (expected eqlog_ prefix)")
	}

	pos := strings.Index(path, "eqlog_")
	name := path[pos+6:]
	pos = strings.Index(name, "_")
	if pos > 0 {
		name = name[:pos]
	}
	instance.name = name

	err := t.watchFile()
	if err != nil {
		return fmt.Errorf("watch file: %w", err)
	}

	return nil
}

func (t *Tracker) watchFile() error {
	if t.path == "" {
		return nil
	}

	if t.pollCtx != nil {
		t.pollCtxCancel()
	}

	config := tail.Config{
		Follow:    true,
		MustExist: true,
		Poll:      true,
	}
	config.Location = &tail.SeekInfo{Offset: 0, Whence: 2}
	t.isLiveParse = true

	config.Logger = tail.DiscardingLogger

	tailer, err := tail.TailFile(t.path, config)
	if err != nil {
		return fmt.Errorf("tail file %s: %w", t.path, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.pollCtx = ctx
	t.pollCtxCancel = cancel
	go t.poll(ctx, tailer)
	return nil
}

func (t *Tracker) poll(ctx context.Context, tailer *tail.Tail) {
	for line := range tailer.Lines {
		select {
		case <-ctx.Done():
			return
		default:
		}

		match := timeRegex.FindStringSubmatch(line.Text)
		if len(match) < 2 {
			continue
		}
		event, err := time.Parse("Mon Jan 02 15:04:05 2006", match[1])
		if err != nil {
			continue
		}

		if !t.isLiveParse && event.After(t.trackerStart) {
			t.isLiveParse = true
		}
		for _, fn := range t.onLineEvent {
			fn(event, line.Text)
		}
		t.onZone(event, line.Text)
	}
}

func (t *Tracker) onZone(event time.Time, line string) {
	match := zoneRegex.FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	if strings.Contains(match[1], "levitation effects") {
		return
	}

	for _, fn := range t.onZoneEvent {
		fn(event, match[1])
	}

	fmt.Println("You have entered", match[1])
}

func Subscribe(fn func(time.Time, string)) error {
	if instance == nil {
		return fmt.Errorf("tracker not initialized")
	}
	instance.onLineEvent = append(instance.onLineEvent, fn)
	return nil
}

func SubscribeToZoneEvent(fn func(time.Time, string)) error {
	if instance == nil {
		return fmt.Errorf("tracker not initialized")
	}
	instance.onZoneEvent = append(instance.onZoneEvent, fn)
	return nil
}

func IsLiveParse() bool {
	if instance == nil {
		return false
	}
	return instance.isLiveParse
}

func PlayerName() string {
	if instance == nil {
		return ""
	}
	return instance.name
}
