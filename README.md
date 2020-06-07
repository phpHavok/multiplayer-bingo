# Multiplayer Bingo

A simple, self-hosted, multiplayer bingo game for friends. Create your own
phrases and generate fun bingo boards to play in an online, real-time
environment.

Each player who joins the game will play on a randomized 4x4 bingo board,
generated from 16 phrases. In the future, we may support different-sized
bingo boards and better phrase support.

## Installation

Check out the [releases](https://github.com/phpHavok/multiplayer-bingo/releases) page to download a pre-compiled binary for Linux AMD64 systems.

If you want to build the application from scratch, just type `make`. You'll
need the Go programming language installed to compile.

## Running

To run, type `./bingo` to see the options:

```
Usage of ./bingo:
  -help
        print usage
  -html string
        path to the html directory for the game (default "./html")
  -phrases string
        the phrases file to use (required)
  -port string
        the port to listen on (default "8080")
  -room string
        the room code players will need to join this game (required)
  -topic string
        the topic for the game (default "Generic Bingo")
```

You'll need to create a text file containing 16 phrases (1 per line), and
specify the path to it using the `-phrases` option. These phrases will be used
to generate the bingo cards. You'll also need a room code, which acts as a sort
of password to prevent random/unwanted people from joining the game. It's also
possible to specify a topic for the room.
