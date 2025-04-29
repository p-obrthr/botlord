# BOTLORD 

This project consists of a modular Discord Bot (server) with optional API connection and a graphical client based on [raylib](https://github.com/gen2brain/raylib-go). Both are written in Go.

## 🚀 Features

### Server (Main Discord Bot)
*Handles events and executes bot commands*
- Standalone go server with clean modularization
- SQLite database integration
- Can be run and used immediately without the need for an API
- Optional API wrapper for extending functionality with external applications or services
- Dockerfile for easy containerization
- Makefile with various build and run targets

### Client (Raylib application)
*Visualization and interaction with the bot, e.g., check the status or start/stop it*
- Standalone Go client with Raylib for a graphical interface
