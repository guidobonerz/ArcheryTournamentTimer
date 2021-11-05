package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"time"

	localFonts "drazil/tournament/resources/fonts"
	localGraphics "drazil/tournament/resources/graphics"
	localSounds "drazil/tournament/resources/sounds"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth    = 1024
	screenHeight   = 768
	warnDuration   = 30
	actionDuration = 180
	zero           = "000"
)

type Stage int

const (
	Halt         Stage = 0
	InitPrepare        = 1
	StartPrepare       = 2
	InitAction         = 3
	StartAction        = 4
)

var (
	showTournamentView = false
	cancel             = false
	fullscreen         = true
	tournamentFont     font.Face
	infoFontLarge      font.Face
	infoFontSmall      font.Face
	roundFont          font.Face
	displayText        = ""
	colorWhite         = color.RGBA{255, 255, 255, 255}
	colorDarkGray      = color.RGBA{50, 50, 50, 255}
	colorBlack         = color.RGBA{0, 0, 0, 255}
	colorYellow        = color.RGBA{255, 255, 0, 255}
	colorRed           = color.RGBA{255, 0, 0, 255}
	colorGreen         = color.RGBA{0, 255, 0, 255}
	countDownColor     = colorYellow
	pairColor          = colorWhite
	startTime          time.Time
	endTime            time.Time
	pair               = [...]string{"A-B", "C-D", "C-D", "A-B"}
	prepareDuration    = [...]int{10, 20}
	duration           int
	half               int = 0
	round              int = 0
	stage              Stage
	displayFormat      string
	audioContext       *audio.Context
	testPlayer         *audio.Player
	signalPlayer1      *audio.Player
	signalPlayer2      *audio.Player
	signalPlayer3      *audio.Player
	buzzerPlayer       *audio.Player
	logo               *ebiten.Image
	red                *ebiten.Image
	green              *ebiten.Image
	yellow             *ebiten.Image
	signalLight        *ebiten.Image
)

func init() {

	var err error
	var img image.Image
	img, _, err = image.Decode(bytes.NewReader(localGraphics.LogoPNG))
	logo = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(localGraphics.Red2))
	red = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(localGraphics.Green2))
	green = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(localGraphics.Yellow2))
	yellow = ebiten.NewImageFromImage(img)

	signalLight = red

	audioContext = audio.NewContext(48000)
	signalSound1, err := wav.Decode(audioContext, bytes.NewReader(localSounds.CarHorn))
	signalPlayer1, err = audioContext.NewPlayer(signalSound1)
	signalPlayer1.SetVolume(20)

	signalSound2, err := wav.Decode(audioContext, bytes.NewReader(localSounds.CarHornDouble))
	signalPlayer2, err = audioContext.NewPlayer(signalSound2)
	signalPlayer2.SetVolume(20)

	signalSound3, err := wav.Decode(audioContext, bytes.NewReader(localSounds.CarHornTriple))
	signalPlayer3, err = audioContext.NewPlayer(signalSound3)
	signalPlayer3.SetVolume(20)

	testSound, err := wav.Decode(audioContext, bytes.NewReader(localSounds.Horn))
	testPlayer, err = audioContext.NewPlayer(testSound)
	testPlayer.SetVolume(20)

	buzzerSound, err := wav.Decode(audioContext, bytes.NewReader(localSounds.Buzzer))
	buzzerPlayer, err = audioContext.NewPlayer(buzzerSound)
	buzzerPlayer.SetVolume(20)

	digitalFont, err := opentype.Parse(localFonts.DigitalFont)
	tournamentFont, err = opentype.NewFace(digitalFont, &opentype.FaceOptions{
		Size:    450,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	roundFont, err = opentype.NewFace(digitalFont, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	textFont, err := opentype.Parse(localFonts.OspDin)
	infoFontLarge, err = opentype.NewFace(textFont, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	infoFontSmall, err = opentype.NewFace(textFont, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

}

type Tournament struct {
}

func (t *Tournament) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyH) {
		showTournamentView = false
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) && showTournamentView {
		countDownColor = colorYellow
		stage = Stage(InitPrepare)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyT) {
		showTournamentView = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEscape) && showTournamentView {
		cancel = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyN) && showTournamentView {
		round = 0
		half = 0
		stage = Stage(Halt)
		PlaySound(10)
		signalLight = red
		duration = 0
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		PlaySound(0)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyF11) {
		fullscreen = !fullscreen
		ebiten.SetFullscreen(fullscreen)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyX) {
		os.Exit(0)
	}

	return nil
}
func PlaySound(count int) {
	if count == 0 {
		if !testPlayer.IsPlaying() {
			testPlayer.Rewind()
			testPlayer.Play()
		}
	} else if count == 1 {
		if !signalPlayer1.IsPlaying() {
			signalPlayer1.Rewind()
			signalPlayer1.Play()
		}
	} else if count == 2 {
		if !signalPlayer2.IsPlaying() {
			signalPlayer2.Rewind()
			signalPlayer2.Play()
		}
	} else if count == 3 {
		if !signalPlayer3.IsPlaying() {
			signalPlayer3.Rewind()
			signalPlayer3.Play()
		}
	} else if count == 10 {
		if !buzzerPlayer.IsPlaying() {
			buzzerPlayer.Rewind()
			buzzerPlayer.Play()
		}
	} else {
	}

}

func (t *Tournament) Draw(screen *ebiten.Image) {
	screen.Fill(colorBlack)

	if showTournamentView {
		if stage == InitPrepare {
			stage = Stage(StartPrepare)
			signalLight = red
			duration = prepareDuration[half]
			PlaySound(2)
			startTime = time.Now()
			endTime = startTime
			endTime.Add(time.Second * time.Duration(duration))
		} else if stage == StartPrepare {
			duration = prepareDuration[half] - int(time.Now().Sub(endTime).Seconds())
			if duration == 0 {
				PlaySound(1)
				stage = Stage(InitAction)
			}
		} else if stage == InitAction {
			stage = Stage(StartAction)
			signalLight = green
			duration = actionDuration
			startTime = time.Now()
			endTime = startTime
			endTime.Add(time.Second * actionDuration)
		} else if stage == StartAction {
			duration = actionDuration - int(time.Now().Sub(endTime).Seconds())
			if duration <= warnDuration {
				countDownColor = colorRed
				signalLight = yellow
			}
			if duration == 0 || cancel {
				cancel = false
				countDownColor = colorYellow
				signalLight = red

				if round == 0 || round == 2 {
					half = 1
					PlaySound(2)
					stage = Stage(InitPrepare)
				} else {
					half = 0
					PlaySound(3)
					stage = Stage(Halt)
				}
				round++
				if round > 3 {
					round = 0
					half = 0
				}
			}
		}
		timeLeft := fmt.Sprintf("%3d", duration)
		roundText := fmt.Sprintf("ROUND:%1d", round+1)
		halfText := fmt.Sprintf("HALF :%1d", half+1)
		clockText := fmt.Sprintf("%02d:%02d:%02d", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		text.Draw(screen, zero, tournamentFont, 400, 350, colorDarkGray)
		text.Draw(screen, timeLeft, tournamentFont, 400, 350, countDownColor)
		text.Draw(screen, pair[round], tournamentFont, 400, 700, colorWhite)
		text.Draw(screen, roundText, roundFont, 440, 395, colorWhite)
		text.Draw(screen, halfText, roundFont, 640, 395, colorWhite)
		text.Draw(screen, clockText, roundFont, 840, 395, colorWhite)

		var op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(7), float64(13))
		op.GeoM.Translate(float64(0), float64(50))
		screen.DrawImage(signalLight, op)
	} else {
		text.Draw(screen, "Turnier Timer", infoFontLarge, 200, 50, colorWhite)
		text.Draw(screen, "BSV Eppinghoven 1743 e.V.", infoFontSmall, 200, 80, colorWhite)
		text.Draw(screen, "[T]urnier Ansicht (Start mit <RETURN>)\n[ESC] Passe vorzeitig beenden\n[N]eustart\n[S]oundcheck\n[F11] Vollbild\n[H]ilfe anzeigen\nE[x]it", infoFontLarge, 200, 150, colorWhite)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(0), float64(30))
		screen.DrawImage(logo, op)
	}
}

func getCenteredX(content string, screen ebiten.Image, df font.Face) int {
	rect := text.BoundString(df, content)
	sw, _ := screen.Size()
	x := (sw - rect.Dx()) / 2
	return x
}

func (t *Tournament) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Archery Tournament Timer")
	ebiten.SetFullscreen(fullscreen)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Tournament{}); err != nil {
		log.Fatal(err)
	}
}
