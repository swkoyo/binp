# binp

binp is an open-source pastebin service built with modern web technologies. It provides a simple and efficient way to share code snippets and text online.

## 🌐 Website

Visit [binp.io](https://binp.io) to use the service.

## 🚀 Features

- Fast and lightweight pastebin service
- Clean and responsive user interface
- Syntax highlighting for various programming languages using Chroma
- Easy sharing and collaboration
- Persistent storage using SQLite

## 🛠 Tech Stack

binp is built using the following technologies:

- [Go](https://golang.org/) (version 1.22 or later) - The core programming language
- [Echo](https://echo.labstack.com/) - High performance, minimalist Go web framework
- [templ](https://github.com/a-h/templ) - A HTML templating language for Go
- [HTMX](https://htmx.org/) - High power tools for HTML
- [Hyperscript](https://hyperscript.org/) - An approachable way to add interactivity to your web pages
- [Tailwind CSS](https://tailwindcss.com/) - A utility-first CSS framework
- [Chroma](https://github.com/alecthomas/chroma) - A syntax highlighter for Go
- [SQLite](https://www.sqlite.org/) - A C-language library that implements a small, fast, self-contained SQL database engine

## 🚀 Getting Started

### Prerequisites

- Go (version 1.22 or later)
- Templ
- Node.js and npm (for Tailwind CSS)
- SQLite

### Environment Variables
- `PORT` - The port number to run the server on (default: `8080`)
- `DB_PATH` - The path to the SQLite database file (default: `./db.sqlite`)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/swkoyo/binp.git
   cd binp
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Install Node.js dependencies:
   ```bash
   npm install
   ```

4. Build the CSS:
   ```bash
   npm run build
   ```

5. Set up your environment variables:
   ```bash
   cp .env.example .env
   ```
   Then edit the `.env` file with your specific configuration.

6. Build and run the application:
   ```bash
   go run cmd/api/main.go
   ```

7. Open your browser and navigate to `http://localhost:8080` (or the port you've configured).

## 🛠 Development

To watch for changes and automatically rebuild the CSS during development:

```bash
# Start the Go server (requires Air to be installed)
air

# Watch for changes in the CSS
npm run watch
```

## 🚦 Testing

To run the tests:

```bash
go test ./...
```

## CLI

binp also comes with a CLI tool that allows you to interact with the pastebin service from the command line.

### Building

To build the CLI tool:

```bash
go build -o tmp/binp cmd/cli/main.go
```

### Usage

To see the available commands and options:

```bash
./tmp/binp --help
```

To create a new paste:

```bash
# Options:
# -l, --language:  The language of the paste (default: "txt")
# -b, --burn:  Whether the paste should be deleted after viewing (default: false)
# -e, --expiry:  The expiry time of the paste (options: "1m", "1h", "1d". default: "1m")

./tmp/binp create <text>
```

To get a paste by its ID:

```bash
# Options:
# -j, --json:  Output the paste as JSON
# -p, --pretty:  Pretty print the JSON output (requires bat to be installed)

./tmp/binp get <id>
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is open source and available under the [GPLv3 license](LICENSE).

## 📞 Contact

If you have any questions or feedback, please open an issue on the GitHub repository.

---

Happy pasting with binp!
