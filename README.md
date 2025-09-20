# Pelican Art Gallery

What happens when you ask an AI to draw famous paintings as simple computer graphics? This project shows how language models—usually built for text—try to turn classic artworks into SVG images. Even though these AIs aren't trained to make vector graphics, it's fun to see how they mix what they know about art with their ability to write SVG code. This project is inspired by [Simon Willison's article](https://simonwillison.net/2025/Jun/6/six-months-in-llms/) about creative LLM benchmarks.

## What is Pelican Art Gallery?

Pelican Art Gallery is an interactive web application that serves as both an art gallery and workshop for creating custom SVG illustrations using various AI language models. The project was inspired by experiments with generating SVG art, particularly the famous "pelicans on a bicycle" example, and aims to make AI-powered SVG creation accessible and fun.

The application features:

- **Homepage**: A minimal landing page introducing the project
- **Workshop**: Full-featured interface for creating new SVG artwork
- **Gallery**: Display and browse your collection of generated artwork

## Features

- **Interactive Workshop**: Web interface for crafting SVG art prompts with group-based workflow
- **Model Selection**: Choose from various OpenRouter AI models with batch generation
- **Group Management**: Create artwork groups with multiple models and prompts
- **Fine-tuning Controls**: Adjust temperature and max tokens for generation
- **Real-time Preview**: Instant SVG generation and display
- **YAML Configuration**: Customizable prompt templates
- **Go Backend**: Efficient server with embedded static files
- **Download & Copy**: Export generated SVGs easily
- **Modern Design**: Clean, minimalist black-and-white aesthetic with Tailwind CSS
- **Comprehensive Logging**: Detailed error handling and debugging
- **Mobile Responsive**: Works on all device sizes
- **Database Normalization**: Efficient storage with artwork groups and individual artworks

## Prerequisites

- Go 1.21 or higher
- OpenRouter API key (for AI model access)

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

## Configuration

Customize the AI prompts and behavior in `config/prompt.yaml`. This file contains:

- System prompt templates
- User prompt formats
- Model-specific configurations

## API Endpoints

### Pages

- `GET /` - Homepage
- `GET /workshop` - Art creation interface
- `GET /gallery` - Artwork gallery
- `GET /gallery/` - Artwork gallery (with trailing slash)

### Artwork Groups

- `GET /api/groups` - List all artwork groups
- `POST /api/groups` - Create new artwork group
- `GET /api/groups/{id}` - Get specific artwork group
- `PUT /api/groups/{id}` - Update artwork group
- `DELETE /api/groups/{id}` - Delete artwork group

### Artworks

- `GET /api/artworks` - List all artworks
- `POST /api/artworks` - Create new artwork
- `GET /api/artworks/{id}` - Get specific artwork
- `PUT /api/artworks/{id}` - Update artwork
- `DELETE /api/artworks/{id}` - Delete artwork
- `POST /api/generate` - Generate SVG artwork for a group

### Models & Configuration

- `GET /api/models` - Available AI models
- `GET /api/config` - Get current configuration
- `POST /api/config` - Update configuration

### System

- `GET /health` - Health check

## Architecture

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

### Database Schema

The application uses a normalized SQLite database with two main tables:

- **artwork_groups**: Stores group metadata (title, prompt, category, created_at)
- **artworks**: Stores individual artworks (group_id FK, model, params, SVG, created_at)

This structure enables efficient storage and retrieval of related artworks while maintaining data integrity.

### Key Design Principles

- **Simplicity**: Vanilla solutions over complex frameworks, utility-first CSS
- **Performance**: Minimal bundle size, efficient rendering with Preact
- **Accessibility**: Semantic HTML, keyboard navigation, ARIA attributes
- **Maintainability**: Component-based architecture, clear separation of concerns
- **Data Integrity**: Normalized database schema with foreign key constraints
- **Design Aesthetic**: Strict black-and-white color palette with flat interactions

## Inspiration

This project draws direct inspiration from [Simon Willison's blog post](https://simonwillison.net/2024/Oct/25/pelicans-on-a-bicycle/) where he benchmarked 16 different AI models on their ability to generate SVG code for "a pelican riding a bicycle." He chose this creative prompt specifically because:

- He likes pelicans
- It's unlikely to exist in training data, making it a true test of generative capabilities

The experiment tested models from OpenAI (GPT-4o, o1), Anthropic (Claude), Google (Gemini), and Meta (Llama on Cerebras), demonstrating the wide range of SVG generation quality possible with different LLMs. Pelican Art Gallery builds upon this concept by providing an accessible, user-friendly interface for exploring AI-generated SVG art beyond just the pelican example.

## Contributing

Contributions are welcome! Please feel free to:

- Report bugs or suggest features
- Submit pull requests for improvements
- Share your generated artwork examples

## License

[Add your license information here]

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

## Configuration

Customize the AI prompts and behavior in `config/prompt.yaml`. This file contains:

- System prompt templates
- User prompt formats
- Model-specific configurations

## API Endpoints

### Pages

- `GET /` - Homepage
- `GET /workshop` - Art creation interface
- `GET /gallery` - Artwork gallery
- `GET /gallery/` - Artwork gallery (with trailing slash)

### Artwork Groups

- `GET /api/groups` - List all artwork groups
- `POST /api/groups` - Create new artwork group
- `GET /api/groups/{id}` - Get specific artwork group
- `PUT /api/groups/{id}` - Update artwork group
- `DELETE /api/groups/{id}` - Delete artwork group

### Artworks

- `GET /api/artworks` - List all artworks
- `POST /api/artworks` - Create new artwork
- `GET /api/artworks/{id}` - Get specific artwork
- `PUT /api/artworks/{id}` - Update artwork
- `DELETE /api/artworks/{id}` - Delete artwork
- `POST /api/generate` - Generate SVG artwork for a group

### Models & Configuration

- `GET /api/models` - Available AI models
- `GET /api/config` - Get current configuration
- `POST /api/config` - Update configuration

### System

- `GET /health` - Health check

## Architecture

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

### Database Schema

The application uses a normalized SQLite database with two main tables:

- **artwork_groups**: Stores group metadata (title, prompt, category, created_at)
- **artworks**: Stores individual artworks (group_id FK, model, params, SVG, created_at)

This structure enables efficient storage and retrieval of related artworks while maintaining data integrity.

### Key Design Principles

- **Simplicity**: Vanilla solutions over complex frameworks, utility-first CSS
- **Performance**: Minimal bundle size, efficient rendering with Preact
- **Accessibility**: Semantic HTML, keyboard navigation, ARIA attributes
- **Maintainability**: Component-based architecture, clear separation of concerns
- **Data Integrity**: Normalized database schema with foreign key constraints
- **Design Aesthetic**: Strict black-and-white color palette with flat interactions

## Inspiration

This project draws direct inspiration from [Simon Willison's blog post](https://simonwillison.net/2024/Oct/25/pelicans-on-a-bicycle/) where he benchmarked 16 different AI models on their ability to generate SVG code for "a pelican riding a bicycle." He chose this creative prompt specifically because:

- He likes pelicans
- It's unlikely to exist in training data, making it a true test of generative capabilities

The experiment tested models from OpenAI (GPT-4o, o1), Anthropic (Claude), Google (Gemini), and Meta (Llama on Cerebras), demonstrating the wide range of SVG generation quality possible with different LLMs. Pelican Art Gallery builds upon this concept by providing an accessible, user-friendly interface for exploring AI-generated SVG art beyond just the pelican example.

## Contributing

Contributions are welcome! Please feel free to:

- Report bugs or suggest features
- Submit pull requests for improvements
- Share your generated artwork examples

## License

[Add your license information here]
