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
	screenWidth     = 1024
	screenHeight    = 768
	warnDuration    = 30
	prepareDuration = 20
	actionDuration  = 180
	zero            = "000"
)

var (
	showTournamentMode = false
	startPrepareTimer  = false
	startActionTimer   = false
	initPrepareMode    = false
	initActionMode     = false
	tournamentFont     font.Face
	infoFontLarge      font.Face
	infoFontSmall      font.Face
	displayText        = ""
	colorWhite         = color.RGBA{255, 255, 255, 255}
	colorDarkGray      = color.RGBA{50, 50, 50, 255}
	colorBlack         = color.RGBA{0, 0, 0, 255}
	colorYellow        = color.RGBA{255, 255, 0, 255}
	colorRed           = color.RGBA{255, 0, 0, 255}
	colorGreen         = color.RGBA{0, 255, 0, 255}
	counterColor       = colorYellow
	pairColor          = colorWhite
	pair               = [...]string{"A-B", "C-D", "C-D", "A-B"}
	startTime          time.Time
	endTime            time.Time
	duration           int
	round              int
	displayFormat      string
	audioContext       *audio.Context
	testPlayer         *audio.Player
	signalPlayer1      *audio.Player
	signalPlayer2      *audio.Player
	signalPlayer3      *audio.Player
	logo               *ebiten.Image
	red                *ebiten.Image
	green              *ebiten.Image
	yellow             *ebiten.Image
	off                *ebiten.Image
	redLight           *ebiten.Image
	greenLight         *ebiten.Image
	yellowLight        *ebiten.Image
	signalLight        *ebiten.Image
)

func init() {

	var err error
	var img image.Image
	img, _, err = image.Decode(bytes.NewReader(localGraphics.LogoPNG))
	if err != nil {
		log.Fatal(err)
	}
	logo = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(localGraphics.Off))
	if err != nil {
		log.Fatal(err)
	}
	off = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(localGraphics.Red2))
	if err != nil {
		log.Fatal(err)
	}
	red = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(localGraphics.Green2))
	if err != nil {
		log.Fatal(err)
	}
	green = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(localGraphics.Yellow2))
	if err != nil {
		log.Fatal(err)
	}
	yellow = ebiten.NewImageFromImage(img)
	resetLights()

	audioContext = audio.NewContext(48000)
	signalSound1, err := wav.Decode(audioContext, bytes.NewReader(localSounds.CarHorn))
	if err != nil {
		log.Fatal(err)
	}
	signalPlayer1, err = audioContext.NewPlayer(signalSound1)
	signalPlayer1.SetVolume(20)
	if err != nil {
		log.Fatal(err)
	}

	signalSound2, err := wav.Decode(audioContext, bytes.NewReader(localSounds.CarHornDouble))
	if err != nil {
		log.Fatal(err)
	}
	signalPlayer2, err = audioContext.NewPlayer(signalSound2)
	signalPlayer2.SetVolume(20)
	if err != nil {
		log.Fatal(err)
	}

	signalSound3, err := wav.Decode(audioContext, bytes.NewReader(localSounds.CarHornTriple))
	if err != nil {
		log.Fatal(err)
	}
	signalPlayer3, err = audioContext.NewPlayer(signalSound3)
	signalPlayer3.SetVolume(20)
	if err != nil {
		log.Fatal(err)
	}

	testSound, err := wav.Decode(audioContext, bytes.NewReader(localSounds.Horn))
	if err != nil {
		log.Fatal(err)
	}
	testPlayer, err = audioContext.NewPlayer(testSound)
	testPlayer.SetVolume(20)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(localFonts.DigitalFont)
	if err != nil {
		log.Fatal(err)
	}

	tournamentFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    450,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	tt2, err := opentype.Parse(localFonts.OspDin)
	if err != nil {
		log.Fatal(err)
	}
	infoFontLarge, err = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	infoFontSmall, err = opentype.NewFace(tt2, &opentype.FaceOptions{
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

func resetLights() {
	//redLight = red
	//greenLight = off
	//yellowLight = off
	signalLight = red
}

func (t *Tournament) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		showTournamentMode = false
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		PlaySignal(2)
		counterColor = colorYellow
		initPrepareMode = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		showTournamentMode = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyP) {
	} else if inpututil.IsKeyJustReleased(ebiten.KeyT) {
		PlaySignal(0)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		showTournamentMode = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyX) {
		os.Exit(0)
	}

	return nil
}
func PlaySignal(count int) {
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
	} else {
		if !signalPlayer3.IsPlaying() {
			signalPlayer3.Rewind()
			signalPlayer3.Play()
		}
	}

}

func (t *Tournament) Draw(screen *ebiten.Image) {
	screen.Fill(colorBlack)

	if showTournamentMode {
		if initPrepareMode {
			startTime = time.Now()
			endTime = startTime
			endTime.Add(time.Second * prepareDuration)
			duration = prepareDuration
			initPrepareMode = false
			startPrepareTimer = true
			//redLight = red
			//greenLight = off
			//yellowLight = off
			signalLight = red

		} else if startPrepareTimer {
			duration = prepareDuration - int(time.Now().Sub(endTime).Seconds())
			if duration == 0 {
				PlaySignal(1)
				startPrepareTimer = false
				initActionMode = true
			}
		} else if initActionMode {

			startTime = time.Now()
			endTime = startTime
			endTime.Add(time.Second * actionDuration)
			duration = actionDuration
			initActionMode = false
			startActionTimer = true
			//redLight = off
			//greenLight = green
			//yellowLight = off
			signalLight = green
		} else if startActionTimer {
			duration = actionDuration - int(time.Now().Sub(endTime).Seconds())
			if duration <= warnDuration {
				counterColor = colorRed
				//redLight = off
				//greenLight = off
				//yellowLight = yellow
				signalLight = yellow
			}
			if duration == 0 {
				PlaySignal(3)
				startActionTimer = false
				counterColor = colorYellow
				resetLights()
				round++
				if round > 3 {
					round = 0
				}
			}
		}
		timeLeft := fmt.Sprintf("%3d", duration)
		couple := pair[round]
		text.Draw(screen, zero, tournamentFont, 400, 350, colorDarkGray)
		text.Draw(screen, timeLeft, tournamentFont, 400, 350, counterColor)
		text.Draw(screen, couple, tournamentFont, 400, 700, colorWhite)
		/*
			var op = &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(0.25), float64(0.25))
			op.GeoM.Translate(float64(100), float64(50))
			screen.DrawImage(redLight, op)

			op = &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(0.25), float64(0.25))
			op.GeoM.Translate(float64(100), float64(270))
			screen.DrawImage(greenLight, op)

			op = &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(0.25), float64(0.25))
			op.GeoM.Translate(float64(100), float64(490))
			screen.DrawImage(yellowLight, op)

		*/

		var op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(7), float64(13))
		op.GeoM.Translate(float64(0), float64(50))
		screen.DrawImage(signalLight, op)
	} else {
		text.Draw(screen, "Turnier Timer", infoFontLarge, 200, 50, colorWhite)
		text.Draw(screen, "BSV Eppinghoven 1743 e.V.", infoFontSmall, 200, 80, colorWhite)
		text.Draw(screen, "[S]tart\n[H]alt\n[T]est\n[M]ittagspause\n[R]eset\nE[x]it", infoFontLarge, 200, 150, colorWhite)
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
	ebiten.SetFullscreen(true)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(&Tournament{}); err != nil {
		log.Fatal(err)
	}
}
