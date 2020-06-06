# Multiplayer Bingo

A simple, self-hosted, multiplayer bingo game for friends. Create your own
phrases and generate fun bingo boards to play in an online, real-time
environment.

## Compilation and Running

To build the application, type `make`. You'll need the Go language installed to
compile.

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
