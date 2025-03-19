# Process Monitor

**Process Monitor** is a cross-platform desktop application built with [Go](https://golang.org/) and [Wails](https://wails.io/) that lists system processes and uses a large language model (LLM) to provide detailed descriptions of each process. The generated descriptions help determine if a process is malicious or safe.

## Features

- **Process Listing:** Enumerates all running system processes.
- **AI-Powered Analysis:** Uses LLMs to generate detailed descriptions for each process.
- **Security Insights:** Highlights whether a process is considered safe or potentially malicious.
- **Flexible AI Services:** Supports both Groq API and Ollama API for LLM integration.

## AI Service Configuration

Process Monitor currently supports two AI services:
- **Groq API**
- **Ollama API**

Configure your desired AI service by setting the environment variables in the `.env` file. Below is an example configuration file (`env.example`):

```env
GROQ_API_URL=https://api.groq.com/openai/v1/chat/completions
GROQ_API_KEY=
GROQ_MODEL=llama-3.3-70b-versatile

OLLAMA_API_URL=http://localhost:11434/api/chat
OLLAMA_API_KEY=
OLLAMA_MODEL=llama3.2:1b

AI_SERVICE= # Set this to "groq" or "ollama" to select the desired service
```

> **Note:** Rename the file to `.env` and fill in your API keys before building the application.

## Folder Structure

The project is organized as follows:

```
.
├── ai
│   ├── ai.go         # Main AI interface logic
│   ├── groq.go       # Groq API integration
│   └── ollama.go     # Ollama API integration
├── app.go            # Application entry point
├── backend           # Additional backend logic
├── build             # Build assets and platform-specific files
│  
├── config
│   └── env.go        # Environment variable configuration
├── frontend          # Wails frontend source files
│
├── go.mod            # Go module file
├── go.sum
├── main.go           # Application main entry point
├── process
│   └── process.go    # System process enumeration logic
├── scripts           # Build and installation scripts
│   ├── build-macos-arm.sh
│   ├── build-macos-intel.sh
│   ├── build-macos.sh
│   ├── build.sh
│   ├── build-windows.sh
│   └── install-wails-cli.sh
└── wails.json        # Wails configuration file
```

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or later)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) installed globally
- Node.js (for frontend development)

### Setup Instructions

1. **Clone the repository:**

   ```bash
   git clone https://your-repo-url.git
   cd process-monitor
   ```

2. **Configure the Environment:**

   - Copy the `env.example` file to `.env`:
     ```bash
     cp env.example .env
     ```
   - Edit the `.env` file to provide your API keys and choose your AI service by setting the `AI_SERVICE` variable (`groq` or `ollama`).

3. **Install Dependencies:**

   - For Go dependencies:
     ```bash
     go mod tidy
     ```
   - For frontend dependencies:
     ```bash
     cd frontend
     npm install
     cd ..
     ```

4. **Build the Application:**

   Use the Wails CLI to build your project:
   ```bash
   wails build
   ```

5. **Run the Application:**

   After building, run the executable generated in the build folder.

## Usage

- **Listing Processes:** The application will automatically list all active system processes.
- **AI Analysis:** Each process is analyzed using the configured LLM service, and a detailed description is provided along with an indication of whether the process is safe or malicious.
- **Settings:** You can update your API credentials and switch between AI services using the settings panel in the application.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgments

- [Wails](https://wails.io/)
- [Go](https://golang.org/)
- [Groq API](https://groq.com/)
- [Ollama API](https://ollama.com/)