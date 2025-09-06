# Pelican Art Gallery

A Go-based web application for generating AI-powered SVG artwork, inspired by [Simon Willison's exploration of SVG art generation with LLMs](https://simonwillison.net/2024/Oct/25/pelicans-on-a-bicycle/).

## What is Pelican Art Gallery?

Pelican Art Gallery is an interactive web application that serves as both an art gallery and workshop for creating custom SVG illustrations using various AI language models. The project was inspired by experiments with generating SVG art, particularly the famous "pelicans on a bicycle" example, and aims to make AI-powered SVG creation accessible and fun.

The application features:

- **Homepage**: A minimal landing page introducing the project
- **Workshop**: Full-featured interface for creating new SVG artwork
- **Gallery**: Display and browse your collection of generated artwork

## Features

- **Interactive Workshop**: Web interface for crafting SVG art prompts
- **Model Selection**: Choose from various OpenRouter AI models
- **Fine-tuning Controls**: Adjust temperature and max tokens for generation
- **Real-time Preview**: Instant SVG generation and display
- **YAML Configuration**: Customizable prompt templates
- **Go Backend**: Efficient server with embedded static files
- **Download & Copy**: Export generated SVGs easily
- **Brutalist Design**: Clean, minimal interface using Concrete.css
- **Comprehensive Logging**: Detailed error handling and debugging
- **Mobile Responsive**: Works on all device sizes

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
make install    # Install Go dependencies
make dev        # Run development server with hot reload
make build      # Build for production
make run        # Run production build
make clean      # Clean build artifacts
```

### Project Structure

```
├── main.go              # Application entry point and routing
├── internal/
│   ├── api/            # HTTP handlers
│   ├── config/         # Configuration management
│   ├── database/       # Database operations
│   └── models/         # Data structures
├── templates/          # HTML templates (homepage, workshop, gallery)
├── static/
│   ├── js/             # Modern ES6 modules
│   │   ├── main.js     # Application entry point
│   │   ├── modules/    # Feature modules
│   │   └── utils/      # Utility functions
│   ├── css/            # Component-based stylesheets
│   └── favicon.svg     # Site favicon
├── config/             # YAML configuration files
├── bin/               # Built binaries
└── tmp/               # Temporary build files
```

## Usage

1. **Homepage**: Start at the minimal landing page
2. **Workshop**: Create new artwork by:
   - Selecting an AI model
   - Adjusting generation parameters
   - Writing descriptive prompts
   - Generating and previewing SVGs
3. **Gallery**: Browse and manage your artwork collection

## Configuration

Customize the AI prompts and behavior in `config/prompt.yaml`. This file contains:

- System prompt templates
- User prompt formats
- Model-specific configurations

## API Endpoints

- `GET /` - Homepage
- `GET /workshop` - Art creation interface
- `GET /gallery` - Artwork gallery
- `POST /api/generate` - Generate SVG artwork
- `GET /api/models` - Available AI models
- `GET /health` - Health check

## Architecture

### Backend

- **Go 1.21+**: Standard library with custom routing
- **Templates**: Component-based HTML templates
- **Database**: SQLite for artwork storage
- **Build System**: Makefile-based development workflow

### Frontend

- **CSS**: Component-based architecture with CSS custom properties
- **JavaScript**: Modern ES6 modules with functional composition
- **Design**: Mobile-first, dark mode support, brutalist aesthetic

### Key Design Principles

- **Simplicity**: Vanilla solutions over complex frameworks
- **Performance**: Minimal bundle size, efficient rendering
- **Accessibility**: Semantic HTML, keyboard navigation
- **Maintainability**: Clear separation of concerns, modular architecture

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
