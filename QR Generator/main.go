package main

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	a := app.New()
	w := a.NewWindow("QR Code Generator — AI Sikho")
	w.Resize(fyne.NewSize(400, 600))

	// Inputs
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL")

	nameEntry := widget.NewEntry()
	nameEntry.SetText("my_qr")

	sizeEntry := widget.NewEntry()
	sizeEntry.SetText("256")

	// Preview Image
	img := canvas.NewImageFromFile("")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))

	status := widget.NewLabel("")

	// Generate Button
	generateBtn := widget.NewButton("Generate QR", func() {
		url := urlEntry.Text
		name := nameEntry.Text
		size := 256

		if url == "" {
			status.SetText("❌ URL cannot be empty")
			return
		}

		if name == "" {
			name = "my_qr"
		}

		if url[:4] != "http" {
			url = "https://" + url
		}

		fileName := name + ".png"

		// Generate QR
		err := qrcode.WriteFile(url, qrcode.High, size, fileName)
		if err != nil {
			status.SetText("❌ Error generating QR")
			return
		}

		// Update Preview
		img.File = fileName
		img.Refresh()

		status.SetText("✅ Saved: " + fileName)
	})

	// Layout
	content := container.NewVBox(
		widget.NewLabelWithStyle("QR Code Generator", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),

		widget.NewLabel("URL"),
		urlEntry,

		widget.NewLabel("File Name"),
		nameEntry,

		widget.NewLabel("Size (pixels)"),
		sizeEntry,

		generateBtn,
		status,
		img,
	)

	w.SetContent(content)
	w.ShowAndRun()

	// Cleanup (optional)
	defer os.Remove("my_qr.png")
}
