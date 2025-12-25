package stealth

import (
	"math/rand"
	"strings"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/go-rod/rod"
)

type TypingSimulator struct {
	cfg    *config.TimingConfig
	random *rand.Rand
}

func NewTypingSimulator(cfg *config.TimingConfig) *TypingSimulator {
	return &TypingSimulator{
		cfg:    cfg,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ts *TypingSimulator) TypeIntoElement(element *rod.Element, text string, humanLike bool) error {
	if !humanLike {
		return element.Input(text)
	}

	if ts.random.Float64() < 0.15 {
		return ts.typeWithTypo(element, text)
	}

	return ts.typeRealistic(element, text)
}

func (ts *TypingSimulator) typeRealistic(element *rod.Element, text string) error {
	for _, char := range text {
		if err := element.Input(string(char)); err != nil {
			return err
		}

		delay := ts.getTypingDelay(char)
		time.Sleep(delay)

		if ts.random.Float64() < 0.05 {
			pauseDelay := time.Duration(200+ts.random.Intn(800)) * time.Millisecond
			time.Sleep(pauseDelay)
		}
	}

	return nil
}

func (ts *TypingSimulator) typeWithTypo(element *rod.Element, text string) error {
	if len(text) < 3 {
		return ts.typeRealistic(element, text)
	}

	typoPosition := 1 + ts.random.Intn(len(text)-2)

	typoChars := []rune{'a', 'e', 'i', 'o', 'u', 't', 's', 'r', 'n'}
	typoChar := typoChars[ts.random.Intn(len(typoChars))]

	for i, char := range text {
		if i == typoPosition {
			if err := element.Input(string(typoChar)); err != nil {
				return err
			}
			time.Sleep(ts.getTypingDelay(typoChar))

			recognitionDelay := time.Duration(200+ts.random.Intn(400)) * time.Millisecond
			time.Sleep(recognitionDelay)

			if err := element.Type(rod.KeyBackspace); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}

		if err := element.Input(string(char)); err != nil {
			return err
		}

		delay := ts.getTypingDelay(char)
		time.Sleep(delay)
	}

	return nil
}

func (ts *TypingSimulator) getTypingDelay(char rune) time.Duration {
	baseDelay := ts.cfg.TypingMin + ts.random.Intn(ts.cfg.TypingMax-ts.cfg.TypingMin)

	if char == ' ' {
		baseDelay = int(float64(baseDelay) * 1.5)
	}

	if strings.ContainsRune("asdfjkl;", char) {
		baseDelay = int(float64(baseDelay) * 0.8)
	}

	if strings.ContainsRune("qwertyzxcvb", char) {
		baseDelay = int(float64(baseDelay) * 1.2)
	}

	if char >= 'A' && char <= 'Z' {
		baseDelay = int(float64(baseDelay) * 1.3)
	}

	return time.Duration(baseDelay) * time.Millisecond
}

func (ts *TypingSimulator) TypeWithThinkDelay(element *rod.Element, text string) error {
	thinkDelay := time.Duration(ts.cfg.ThinkMin+ts.random.Intn(ts.cfg.ThinkMax-ts.cfg.ThinkMin)) * time.Millisecond
	time.Sleep(thinkDelay)

	return ts.TypeIntoElement(element, text, true)
}

func (ts *TypingSimulator) PasteText(element *rod.Element, text string, humanLike bool) error {
	if !humanLike {
		return element.Input(text)
	}

	thinkDelay := time.Duration(400+ts.random.Intn(800)) * time.Millisecond
	time.Sleep(thinkDelay)

	chunks := len(text) / 3
	if chunks < 1 {
		chunks = 1
	}

	chunkSize := len(text) / chunks
	for i := 0; i < chunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == chunks-1 {
			end = len(text)
		}

		chunk := text[start:end]
		if err := element.Input(chunk); err != nil {
			return err
		}

		if i < chunks-1 {
			chunkDelay := time.Duration(50+ts.random.Intn(150)) * time.Millisecond
			time.Sleep(chunkDelay)
		}
	}

	return nil
}
