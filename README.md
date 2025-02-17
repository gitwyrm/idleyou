# IdleYou

A simple IDLE game written in Go, using Fyne for the GUI.

**This is currently in early development, so don't expect this to be more than a tech demo.**

It is mostly a learning project right now with next to no content, but might turn into an actual game. You can however already use it to create your own games or write mods for the official IdleYou game.

The code is licensed under the Mozilla Public License 2.0 (MPL-2.0), which in layman's terms (meaning, read the actual license; nothing I write here overrides what the license states) means that you only need to make modifications to the original source files open source, but you can keep your game closed source and sell it if you prefer.

## How to play

- You need food to survive. (If food hits 0 it's GAME OVER)
- To buy food, you need money.
- To get money, you need to work.
- To work, you need energy.
- To get energy, you need to sleep.
- To have more energy, you need to be fit.
- To get more money, you need a better job.
- To get a better job, you need to gain work experience and have a good appearance.
- To get a better appearance, you need to be fit, charismatic, clean and in a good mood.
- Events will happen at times, good or bad, sometimes they even offer you a choice.

## Running

Only tested on macOS right now, but should work on all platforms.

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

You can also easily create mods by creating a new folder inside of the scripts directory and adding one or more .txt files in there with new events. If you want to add images, create another directory (with the same name) in the images directory and put your images there.

The event names get automatically prefixed with the folder name when the script files are read by the game, so you don't need to worry about your mod's events interfering with other mods. Image paths are also automatically prefixed with the mod name, so you don't need to specify the mod's folder and can just use the file name in the script.

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

The print command prints a message to the message log.

The show command shows a picture (can also be an animated GIF). The pictures are stored in `~/Documents/IdleYou/images` and you only need to give the show command the name of the picture file, not the whole path.

Conditions check variables for their value, possible variables are:

```
Ticks
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

for an event that never fires, which comes in handy for multiple-choice events that have one additional command they can use, starting with a `*`.

```
=== My first multiple choice event
? true
! paused = true
! print "Choose wisely!"
* My first button -> Clicked first button
* mood > 5, appearance < 20: My second button -> Clicked second button
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

This creates three events, the first one is a multiple choice event, the other two are normal events which are never executed since we want to execute them manually through the first event.

The first event uses the condition `? true`, which means it always fires, but it returns true `> true`, so it only fires once and is then removed from the game.

The first event pauses the game, displays a message to the user and shows two buttons. Clicking on them executes one of the other two events, which print a message to the message log and unpause the game again. The second button has two conditions attached, the button is only shown if both conditions are true.

You can of course also link to another multiple-choice event and have a deep decision tree.

There is one other event type, a progress event. It's an event that shows a progress bar and only executes it's actions (`!`) once the progress has reached it's maximum value (`%`). Sleeping, Morning Routine and Watching TV are all progress events.

Here is how you can create one yourself:

```
=== My Progress Event
? ticks == 50
% 50
! print Done!
> true
```

This event fires when the game has run for 50 ticks, shows a progress bar with the label "My Progress Event" and finishes after 50 ticks (so at tick 100), then adds "Done!" to the message list.

## Creating a mod

Create a new folder for your mod in `~/Documents/IdleYou/scripts`, lets' call it `firefighter` since our example mod adds a firefighter job to the game.

In that folder `~/Documents/IdleYou/scripts/firefighter`, create one or more .txt files that contain your mod's script. The files get all concatenated together, so you can split everything up in as many files as you like.

If your mod has images, put them in `~/Documents/IdleYou/images/firefighter`. When you show an image with `! show image.png`, you don't need to add `firefighter` to the path, that is automatically added and inferred from the name of the folder. That's why both folders, the one inside images and the one inside scripts, should use the same name, `firefighter`.
