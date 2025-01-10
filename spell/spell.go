package spell

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mux      sync.Mutex
	instance *spell
)

type spell struct {
	spells map[int]*spellEntry
}

type spellEntry struct {
	ID   int
	Name string
	Icon int
}

// Load initializes the spell package
func Load(path string) error {
	mux.Lock()
	defer mux.Unlock()

	instance = &spell{
		spells: make(map[int]*spellEntry),
	}

	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	start := time.Now()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "^")
		val, err := strconv.Atoi(records[0])
		if err != nil {
			return fmt.Errorf("parse spell id: %w", err)
		}

		icon, err := strconv.Atoi(records[144])
		if err != nil {
			return fmt.Errorf("parse spell icon: %w", err)
		}

		instance.spells[val] = &spellEntry{
			ID:   val,
			Name: records[1],
			Icon: icon,
		}
	}
	fmt.Println("Loaded", len(instance.spells), "spells in", time.Since(start).Seconds(), "seconds")
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

// NameByID returns the name of a spell
func NameByID(id int) string {
	mux.Lock()
	defer mux.Unlock()
	if instance == nil {
		return fmt.Sprintf("Spell %d", id)
	}
	spell, ok := instance.spells[id]
	if !ok {
		return fmt.Sprintf("Spell %d", id)
	}

	return spell.Name
}

func IsDoT(id int) bool {
	return false
}

func SpellIDByName(name string) int {
	mux.Lock()
	defer mux.Unlock()
	if instance == nil {
		return -1
	}
	for _, spell := range instance.spells {
		if spell.Name == name {
			return spell.ID
		}
	}

	return 0
}

func IconBySpellID(id int) int {
	mux.Lock()
	defer mux.Unlock()
	if instance == nil {
		return 0
	}
	spell, ok := instance.spells[id]
	if !ok {
		return 0
	}

	return spell.Icon
}
