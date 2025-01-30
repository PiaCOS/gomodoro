# GOMODORO

`gomodoro` is a simple Pomodoro CLI tool built with *Go*. It helps you manage your work and rest intervals by playing sounds at the start and end of each session. You can customize the duration of the sessions using the `-w` (work) and `-r` (rest) flags.

You can stop the tool anytime by pressing `Ctrl+C`.

## Usage

```bash
Usage of gomodoro:
  -r int
        Time in minutes for the rest session (default 5)
  -w int
        Time in minutes for the work session (default 25)
```

## Installation

To install `gomodoro` you can clone this repo and run the following command inside it:

```bash
go install
```

Ensure that Go is installed on your machine. If not, follow the installation instructions [here](https://go.dev/learn/).

## Tips

You can customize the sound files used for the work and rest sessions by replacing them in the sounds folder. The tool currently only supports .mp3 files.

- `gomodoro_01.mp3`: Played when the work session starts.
- `gomodoro_02.mp3`: Played when the work session ends (break time).

Simply place your custom MP3 files in the sounds directory with the same filenames.
