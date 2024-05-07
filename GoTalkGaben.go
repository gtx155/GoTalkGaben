package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os/exec"
	"strings"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type game struct {
	ui *ebitenui.UI
}

type ListEntry struct {
	Word string
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GoTalkGaben")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	validListEntries := []interface{}{"gaben", "two", "welcome", "i'm", "dota", "blood", "double", "first", "fun", "have", "kill", "thanks", "achived", "have", "playing", "this", "you"}
	var invalidInputs []string

	// valid words taken from validListEntries to be used by playAudio
	validValues := make(map[string]bool)

	// Iterate over the validListEntries slice
	for _, entry := range validListEntries {
		// Convert the entry to lowercase and add it to the map
		validValues[strings.ToLower(entry.(string))] = true
	}

	for key, value := range validValues {
		fmt.Printf("%s: %t\n", key, value)
	}

	var bgColor = color.NRGBA{
		R: 56, G: 53, B: 48, A: 255,
	}

	var textColor = color.NRGBA{
		R: 251, G: 241, B: 215, A: 255,
	}

	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(bgColor)),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	// Needed for errors???
	var err error

	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(22)
	faceHeader, _ := loadFont(32)

	labelTitle := widget.NewText(
		widget.TextOpts.Text("GoTalkGaben", faceHeader, textColor),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	// Add the first Text as a child of the container
	rootContainer.AddChild(labelTitle)

	// String used to pass the inputbox text to playAudio() button when the submit button is pressed
	var userInput string

	//var invalidInputs []string

	// construct a standard textinput widget
	TextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),

		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     e_image.NewNineSliceColor(color.NRGBA{R: 41, G: 37, B: 38, A: 255}),
			Disabled: e_image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(face),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          textColor,
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{R: 163, G: 152, B: 132, A: 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(20)),

		//Set the font and width of the caret
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Enter text here"),

		//This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
			userInput = args.InputText
		}),
	)

	rootContainer.AddChild(TextInput)

	audioButton := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Talk", face, &widget.ButtonTextColor{
			Idle: textColor,
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   45,
			Right:  45,
			Top:    15,
			Bottom: 15,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("audio button clicked")
			wordsLower := strings.ToLower(userInput)
			words := strings.Split(wordsLower, " ")

			for _, word := range words {
				if validValues[word] {
					fmt.Printf("User input %s is valid\n", word)
					// playAudio from audio.go
					playAudio("assets/" + word + ".mp3")
				} else {
					fmt.Printf("User input %s is invalid\n", word)
					invalidInputs = append(invalidInputs, word)
					fmt.Println("Invalid", invalidInputs)
					windowPopup(face, invalidInputs, eui)

					// func to add the invalids
				}
			}
			// Clears the slice after use
			invalidInputs = nil
		}),
	)

	rootContainer.AddChild(audioButton)

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(10, 0))))
	rootContainer.AddChild(listsContainer)

	listValid := newList(validListEntries, face, buttonImage, widget.WidgetOpts.LayoutData(widget.GridLayoutData{
		MaxHeight: 150,
	}))
	listsContainer.AddChild(listValid)

	// add the button as a child of the container
	// To display the text widget, we have to add it to the root container.

	// Loads the only context before we can play audio from audio.go
	initOtoContext()

	game := game{
		ui: eui,
	}

	err = ebiten.RunGame(&game)
	if err != nil {
		log.Print(err)
	}
}

func (g *game) Update() error {
	// ui.Update() must be called in ebiten Update function, to handle user input and other things
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// ui.Draw() should be called in the ebiten Draw function, to draw the UI onto the screen.
	// It should also be called after all other rendering for your game so that it shows up on top of your game world.
	g.ui.Draw(screen)
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := e_image.NewNineSliceColor(color.NRGBA{R: 119, G: 107, B: 95, A: 255})

	hover := e_image.NewNineSliceColor(color.NRGBA{R: 145, G: 73, B: 59, A: 255})

	pressed := e_image.NewNineSliceColor(color.NRGBA{R: 129, G: 57, B: 43, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func openLink(url string) {
	var cmd *exec.Cmd
	// Open the link in the default web browser
	cmd = exec.Command("cmd", "/c", "start", url)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
	}
}

func newList(entries []interface{}, face font.Face, buttonImage *widget.ButtonImage, widgetOpts ...widget.WidgetOpt) *widget.List {
	return widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widgetOpts...)),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle:     e_image.NewNineSliceColor(color.NRGBA{R: 41, G: 37, B: 38, A: 255}),
			Disabled: e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			Mask:     e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
		})),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  e_image.NewNineSliceColor(color.NRGBA{R: 56, G: 53, B: 48, A: 255}),
				Hover: e_image.NewNineSliceColor(color.NRGBA{R: 56, G: 53, B: 48, A: 255}),
			}, buttonImage),
			widget.SliderOpts.MinHandleSize(10),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(3)),
		),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntryFontFace(face),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{R: 251, G: 241, B: 215, A: 255}, // Foreground color for the unfocused selected entry
			Unselected:                 color.NRGBA{R: 251, G: 241, B: 215, A: 255}, // Foreground color for the unfocused unselected entry
			SelectedBackground:         color.NRGBA{R: 145, G: 73, B: 59, A: 255},   // Background color for the unfocused selected entry
			SelectingBackground:        color.NRGBA{R: 129, G: 57, B: 43, A: 255},   // Background color for the unfocused being selected entry
			SelectingFocusedBackground: color.NRGBA{R: 145, G: 73, B: 59, A: 255},   // Background color for the focused being selected entry
			SelectedFocusedBackground:  color.NRGBA{R: 135, G: 73, B: 59, A: 255},   // Background color for the focused selected entry
			FocusedBackground:          color.NRGBA{R: 145, G: 73, B: 59, A: 255},   // Background color for the focused unselected entry
		}),
		// REMAKE AS A STRUCT ADD BOOL TYPE FOR INVALID/VALID AND ON SELECT ADD TO INPUT BOX
		/*widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(ListEntry)
			fmt.Println("Entry Selected: ", entry)
		}),*/
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
	)
}

// Create window func for invalid inputs
func windowPopup(face font.Face, invalidInputs []string, eui *ebitenui.UI) {
	var textColorWindow = color.NRGBA{
		R: 251, G: 241, B: 215, A: 255,
	}

	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{0x2a, 0x25, 0x26, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	windowContainer.AddChild(widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("%s", invalidInputs), face, textColorWindow),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	// Create the titlebar for the window
	titleContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{0x77, 0x6b, 0x5f, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	if len(strings.TrimSpace(invalidInputs[0])) == 0 {
		titleContainer.AddChild(widget.NewText(
			widget.TextOpts.Text("Empty Input", face, textColorWindow),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
		))
	} else {
		titleContainer.AddChild(widget.NewText(
			widget.TextOpts.Text("Contains Invalid Input", face, textColorWindow),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
		))
	}

	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(windowContainer),
		//Set the titlebar for the window (Optional)
		widget.WindowOpts.TitleBar(titleContainer, 25),
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Set the minimum size the window can be
		widget.WindowOpts.MinSize(250, 150),
		//Set the maximum size a window can be
		widget.WindowOpts.MaxSize(400, 400),
	)
	//Get the preferred size of the content
	x, y := window.Contents.PreferredSize()
	//Create a rect with the preferred size of the content
	r := image.Rect(0, 0, x, y)
	//Use the Add method to move the window to the specified point
	r = r.Add(image.Point{200, 150})
	fmt.Println(r)
	//Set the windows location to the rect.
	window.SetLocation(r)
	//Add the window to the UI.
	//Note: If the window is already added, this will just move the window and not add a duplicate.
	eui.AddWindow(window)
}
