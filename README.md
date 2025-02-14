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

## Scripting

The game uses a simple scripting language for events that happen during the game. You can easily change the script by running the game once and then opening the `~/Documents/IdleYou/scripts/script.txt` file in a text editor.

You can also easily create mods by creating a new folder inside of the scripts directory and adding one or more .txt files in there with new events.

The event names get automatically prefixed with the folder name when the script files are read by the game, so you don't need to worry about your mod's events interfering with other mods.

### Syntax

An event always starts with the event declaration, which is three `=` signs followed by a space and then the name of the event.

Underneath, you can write commands. Every command starts with a single character that defines the type of command used, for example a `?` for a condition, a space and finally the command itself. At the end comes the return value of the event `> true`. If an event returns true, it is marked as done and removed from the game. If it returns false, it is executed again the next time the conditions are met.

```
# My events

=== It's a robbery!
? money > 100
? rand < 0.01
! print You were robbed and lost all your money!
! show GotRobbed.png
! money = 0
> true
```

The line starting with `#` is a comment and gets ignored. The other lines create a new event called "It's a robbery!" which fires if the player has more than $100, with a 1% chance (`rand < 0.01`), and only happens once (`> true`).

The event conditions (`?` lines) are checked on each game tick, which is about a tenth of a second at the default game speed. So there are 10 ticks happening per second where event conditions are checked.

The print command prints a message to the message log and the show command shows a picture (can also be an animated GIF). The pictures are stored in `~/Documents/IdleYou/images` and you only need to give the show command the path to the picture from that folder. So if it is directly in the images folder, the name of the image file is enough. If it is inside a subfolder, like `misc`, you would write `! show misc/GoodImage.png`

Conditions check variables for their value, possible variables are:

```
Work
WorkXP
Food
FoodMax
Energy
EnergyMax
Mood
Money
Charisma
Fitness
Job
Salary
Working
Paused
RoutineShower
RoutineShave
RoutineBrushTeeth
RoutineBonus
EventName
Appearance
```

And the operators you can use are:

```
< # lesser than
> # greater than
<= # lesser or equal to
>= # bigger or equal to
== # equal to
!= # not equal to
```

`!` commands can also modify variables, the supported operators are:

```
money -= 10 # subtracts 10 from money
money += 10 # adds 10 to money
money = 10 # sets money to 10
```

For conditions you can also just write:

```
? true
```

for an event that always fires, or:

```
? false
```

for an event that never fires, which comes in handy for multiple-choice events that have one additional command they can use, starting with a `-`.

```
=== My first multiple choice event
? true
! paused = true
! print "Choose wisely!"
- My first button -> Clicked first button
- My second button -> Clicked second button
> true

=== Clicked first button
? false
! paused = false
! print You clicked the first button
> true

=== Clicked second button
? false
! paused = false
! print You clicked the second button
> true
```

This creates three events, the first one is a multiple choice event, the other two are normal events which are never executed since we want to execute them manually through the first event. The first event uses the condition `? true`, which means it always fires, but it returns true `> true`, so it only fires once and is then removed from the game.

The first event pauses the game, displays a message to the user and shows two buttons. Clicking on them executes one of the other two events, which print a message to the message log and unpause the game again.

You can of course also link to another multiple-choice event and have a deep decision tree.
