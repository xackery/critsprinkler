package sound

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/library"
	"github.com/xackery/quail/pfs"
	"golang.org/x/exp/rand"
)

type SoundEffect int

const (
	SoundEffectBuyItem SoundEffect = iota
	SoundEffectMoneyBounce1
	SoundEffectMoneyBounce2
	SoundEffectMoneyBounce3
	SoundEffectMoneyBounce4
	SoundEffectMoneyBounce5
	SoundEffectMoneyBounce6
	SoundEffectMoneyBounce7
	SoundEffectMoneyBounce8
)

var (
	cfg     *config.CritSprinklerConfiguration
	sounds  = make(map[SoundEffect][]byte)
	context = audio.NewContext(44100)
)

func (e SoundEffect) String() string {
	switch e {
	case SoundEffectBuyItem:
		return "SoundEffectBuyItem"
	case SoundEffectMoneyBounce1:
		return "SoundEffectMoneyBounce1"
	case SoundEffectMoneyBounce2:
		return "SoundEffectMoneyBounce2"
	case SoundEffectMoneyBounce3:
		return "SoundEffectMoneyBounce3"
	case SoundEffectMoneyBounce4:
		return "SoundEffectMoneyBounce4"
	case SoundEffectMoneyBounce5:
		return "SoundEffectMoneyBounce5"
	case SoundEffectMoneyBounce6:
		return "SoundEffectMoneyBounce6"
	case SoundEffectMoneyBounce7:
		return "SoundEffectMoneyBounce7"
	case SoundEffectMoneyBounce8:
		return "SoundEffectMoneyBounce8"

	default:
		return "Unknown"
	}
}

// Play plays a sound effect
func Play(key SoundEffect) {

	data, ok := sounds[key]
	if !ok {
		fmt.Println("sound", key, "not found")
		return
	}

	fmt.Println("Playing", key)

	p := context.NewPlayerFromBytes(data)
	p.Play()
}

func PlayBounceRandom() {
	i := rand.Intn(8) + 1
	switch i {
	case 1:
		Play(SoundEffectMoneyBounce1)
	case 2:
		Play(SoundEffectMoneyBounce2)
	case 3:
		Play(SoundEffectMoneyBounce3)
	case 4:
		Play(SoundEffectMoneyBounce4)
	case 5:
		Play(SoundEffectMoneyBounce5)
	case 6:
		Play(SoundEffectMoneyBounce6)
	case 7:
		Play(SoundEffectMoneyBounce7)
	case 8:
		Play(SoundEffectMoneyBounce8)
	}
}

func New(config *config.CritSprinklerConfiguration) error {
	cfg = config
	rand.Seed(uint64(time.Now().UnixNano()))
	OnEQPathLoad()

	return nil
}

func OnEQPathLoad() {
	if cfg == nil {
		return
	}
	var err error
	if cfg.EQPath == "" {
		return
	}

	elements := []struct {
		isArchiveFile   bool
		path            string
		archiveFileName string
		key             SoundEffect
	}{
		{true, "/snd2.pfs", "buyitem.wav", SoundEffectBuyItem},
		{false, "assets/sounds/clink1.mp3", "", SoundEffectMoneyBounce1},
		{false, "assets/sounds/clink2.mp3", "", SoundEffectMoneyBounce2},
		{false, "assets/sounds/clink3.mp3", "", SoundEffectMoneyBounce3},
		{false, "assets/sounds/clink4.mp3", "", SoundEffectMoneyBounce4},
		{false, "assets/sounds/clink5.mp3", "", SoundEffectMoneyBounce5},
		{false, "assets/sounds/clink6.mp3", "", SoundEffectMoneyBounce6},
		{false, "assets/sounds/clink7.mp3", "", SoundEffectMoneyBounce7},
		{false, "assets/sounds/clink8.mp3", "", SoundEffectMoneyBounce8},
	}

	for _, element := range elements {
		if element.isArchiveFile {
			err = loadEQSound(element.path, element.archiveFileName, element.key)
			if err != nil {
				fmt.Println("load eq sound", element.key, "error:", err)
			}
			continue
		}

		err = loadSound(element.path, element.key)
		if err != nil {
			fmt.Println("load sound", element.key, "error:", err)
		}

	}

	fmt.Println("Loaded", len(sounds), "sounds")

}

func loadSound(path string, key SoundEffect) error {
	var err error
	sounds[key], err = library.SoundByPath(path)
	if err != nil {
		return fmt.Errorf("sound by path: %w", err)
	}

	return nil
}

func loadEQSound(archivePath, fileName string, key SoundEffect) error {
	r, err := os.Open(cfg.EQPath + archivePath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer r.Close()

	archive, err := pfs.New("foo")
	if err != nil {
		return fmt.Errorf("new archive: %w", err)
	}
	err = archive.Read(r)
	if err != nil {
		return fmt.Errorf("read archive: %w", err)
	}

	data, err := archive.File(fileName)
	if err != nil {
		return fmt.Errorf("file: %w", err)
	}

	stream, err := wav.DecodeWithSampleRate(context.SampleRate(), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return fmt.Errorf("read from stream: %w", err)
	}

	sounds[key] = buf.Bytes()
	return nil
}
