# IdleYou

A simple IDLE game written in Go, using Fyne for the GUI.

## Running

```bash
go run .
```

## Building

```bash
go build .
```

## Packaging

You need to have the Fyne tools installed:

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

Then you can package for macOS, Windows and Linux:

```bash
fyne package -os darwin
fyne package -os windows
fyne package -os linux
```
