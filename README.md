# BOTLORD 

**BOTLORD** is a modular Discord bot built in Go, featuring an optional API and a graphical client using [raylib](https://github.com/gen2brain/raylib-go). The bot is designed to run standalone or integrate seamlessly with other applications via API.

## Features

### Server (Main Discord Bot)
*Handles events and executes bot commands*
- standalone go server with clean modularization
- SQLite database integration
- fully functional out of the box
- optional API wrapper support for extending functionality with external applications or services
- Dockerfile for easy containerization
- Makefile with various build and run targets

#### Environment Variables

| Variable             | Description                                               | Required |
|:---------------------|:----------------------------------------------------------|:---------|
| `DISCORD_BOT_TOKEN`  | Token used for bot authentication                         | Yes      |
| `TEXT_CHANNEL_ID`    | Id of the text channel for voice join messages            | No       |
| `ENABLE_API`         | Enables the API server for external interaction           | No       |

#### Available Commands

| Command                | Description                              |
|:-----------------------|:-----------------------------------------|
| `!hi`                  | Replies with a greeting.                 |
| `!addQuote [text]`     | Adds a new quote.                        |
| `!deleteQuote [id]`    | Deletes a quote by Id.                   |
| `!quotes`              | Lists all stored quotes.                 |
| `!quote`               | Displays a random quote.                 |
| `!addGif [url]`        | Adds a GIF by URL.                       |
| `!deleteGif [id]`      | Deletes a GIF by Id.                     |
| `!gifs`                | Lists all stored GIFs.                   |
| `!gif`                 | Displays a random GIF.                   |
| `!commands`            | Shows all available commands.            |

#### Other functionality
- voice channel watcher: sends a message to a specified text channel when a user joins a voice channel (requires `TEXT_CHANNEL_ID` to be set)

### Client (Raylib Application)
*A GUI for managing and interacting with the bot.*

- written in Go using Raylib
- provides status information and bot control (start/stop)
- logs events and actions
- requires `ENABLE_API=true` on the server to function
