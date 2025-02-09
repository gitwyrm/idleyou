# IdleYou

A simple IDLE game written in Go, using Fyne for the GUI.

## How to play

- You need food to survive. (If food hits 0 it's GAME OVER)
- To buy food, you need money.
- To get money, you need to work.
- To work, you need energy.
- To get energy, you need to sleep.
- To have more energy, you need to be fit.
- To get more money, you need a better job.
- To get a better job, you need to gain work experience and have a good appearance.
- To get a better appearance, you need to be fit, charismatic, clean and in a good mood (everyone looks better with a smile).

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
