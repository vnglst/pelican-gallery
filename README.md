# Pelican Art Gallery

What happens when you ask an AI to draw famous paintings as simple computer graphics? This project shows how language models—usually built for text—try to turn classic artworks into SVG images. Even though these AIs aren't trained to make vector graphics, it's fun to see how they mix what they know about art with their ability to write SVG code. This project is inspired by [Simon Willison's article](https://simonwillison.net/2025/Jun/6/six-months-in-llms/) about creative LLM benchmarks.

## Quick Start

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd genartifacts
   ```

2. **Set up environment**

   ```bash
   cp .env.example .env
   # Edit .env and add your OPENROUTER_API_KEY
   ```

3. **Install dependencies and run**

   ```bash
   make install
   make dev
   ```

4. **Open your browser** to `http://localhost:8080`

## Development

### Available Commands

```bash
make install      # Install dependencies and tools
make dev          # Development server with hot reload
make build        # Build for production
make run          # Run the built application
make clean        # Clean build artifacts
make test         # Run tests
make fmt          # Format Go code
make lint         # Lint Go code
make help         # Show help
```

### Development Workflow

The project uses Tailwind CSS v4 with a CSS-first configuration:

1. **Styling**: All styling uses Tailwind utility classes directly in templates and JavaScript
2. **CSS Build**: CSS is automatically built during development with `make dev`
3. **Theme**: Custom black/white theme defined in `static/css/input.css` with `@theme` directive
4. **Components**: Workshop interface built with Preact functional components using hooks
5. **No Build Step**: HTM provides JSX-like syntax without requiring compilation

### Project Structure

```
├── main.go              # Application entry point and routing
├── internal/
│   ├── api/            # HTTP handlers for groups, artworks, models
│   ├── config/         # Configuration management
│   ├── database/       # Database operations with normalized schema
│   ├── models/         # Data structures for groups and artworks
│   └── pages/          # Page rendering logic
├── templates/          # HTML templates (homepage, workshop, gallery)
├── static/
│   ├── js/             # Modern ES6 modules
│   │   └── workshop.js # Preact-based workshop interface
│   ├── css/
│   │   ├── input.css   # Tailwind CSS v4 configuration with @theme
│   │   └── output.css  # Generated Tailwind CSS
│   └── favicon.svg     # Site favicon
├── config/             # YAML configuration files
├── bin/               # Built binaries
├── artworks.db        # SQLite database with normalized schema
└── tmp/               # Temporary build files
```

## Usage

1. **Homepage**: Start at the minimal landing page
2. **Workshop**: Create new artwork by:
   - Creating an artwork group with a title and prompt
   - Selecting multiple AI models for the group
   - Adjusting generation parameters (temperature, max tokens)
   - Generating all artworks in the group at once
   - Previewing and managing individual artworks
3. **Gallery**: Browse and manage your artwork collection organized by groups

### Technology Stack

- **Backend**: Go 1.21+ with standard library routing
- **Frontend**: Preact with HTM (no build step required)
- **Styling**: Tailwind CSS v4 with CSS-first configuration
- **Database**: SQLite with normalized schema
- **Build**: Makefile-based development workflow

### Backend

- **Go 1.21+**: Standard library with custom routing
- **Templates**: Component-based HTML templates with custom functions
- **Database**: SQLite with normalized schema (artwork_groups ↔ artworks)
- **Build System**: Makefile-based development workflow

### Frontend

- **CSS**: Tailwind CSS v4 with CSS-first configuration and strict black/white design system
- **JavaScript**: Preact-based workshop interface with HTM for component rendering
- **Design**: Modern minimalist aesthetic, mobile-first responsive design
