package main

import (
	"image/color"
	"os"

	"example.com/myproject/mypack"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Pear Service Manager")
	win.Resize(fyne.NewSize(900, 600))

	
	data, err := os.ReadFile("pear.png")
	if err == nil {
		icon := fyne.NewStaticResource("pear.png", data)
		win.SetIcon(icon)
	}

	title := widget.NewLabelWithStyle("Pear Service Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	title.TextStyle.Bold = true

	subtitle := widget.NewLabel("Servis işlemleri için aşağıdaki butonları kullanın:")
	subtitle.Alignment = fyne.TextAlignCenter

	logVBox := container.NewVBox()
	
	darkBackground := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	
	logWithBackground := container.NewStack(darkBackground, container.NewPadded(logVBox))
	scrollableLog := container.NewVScroll(logWithBackground)
	scrollableLog.SetMinSize(fyne.NewSize(0, 400))

	appendLog := func(msg string) {
		logLabel := widget.NewLabel(msg)
		logLabel.Wrapping = fyne.TextWrapWord
		logVBox.Add(logLabel)
		scrollableLog.ScrollToBottom()
	}

	btnLoadEnv := widget.NewButtonWithIcon(".env Load", theme.FileTextIcon(), func() {
		mypack.LoadEnv()
		appendLog("✅ ENV. Loaded.")
	})
	btnLoadEnv.Importance = widget.MediumImportance

	btnServices := widget.NewButtonWithIcon("Servisleri Başlat", theme.MailSendIcon(), func() {
		services := mypack.GETpearOperation()
		if len(services) == 0 {
			appendLog("⚠️ Service Empty.")
			return
		}
		mypack.GETJobQuotes(services, appendLog)
		appendLog("✅ Service Done Successfully.")
	})
	btnServices.Importance = widget.HighImportance

	btnGrid := container.NewGridWithColumns(2,
		btnLoadEnv,
		btnServices,
	)


	header := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(subtitle),
		container.NewPadded(btnGrid),
		widget.NewSeparator(),
	)


	logLabel := widget.NewLabelWithStyle("📋 LOGS:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	
	logSection := container.NewBorder(
		container.NewPadded(logLabel),
		nil,
		nil,
		nil,
		container.NewPadded(scrollableLog),
	)


	content := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		logSection,
	)

	win.SetContent(content)
	win.ShowAndRun()
}