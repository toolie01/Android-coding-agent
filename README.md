# Android Coding Agent

A high-performance, customizable autonomous coding agent for Android development with user-defined safeguards and integrated LLM capabilities.

## Features

- **Autonomous Code Generation**: Multi-file code generation and modification
- **Project Understanding**: Full project context awareness
- **Git Operations**: Commit, branch, PR management
- **Testing & Validation**: Automated testing and validation
- **Custom Safeguards**: JSON-based rule engine for complete control
- **Integrated LLM**: Built-in language model with custom interaction patterns
- **Android-Optimized**: Gradle, manifest, resource handling

## Installation

### Prerequisites

- Go 1.22 or higher
- Git
- Make (optional, for build automation)

### Clone the Repository

```bash
git clone https://github.com/toolie01/android-coding-agent.git
cd android-coding-agent
```

### Install Dependencies

```bash
go mod download
go mod tidy
```

### Build the Agent

```bash
make build
```

Or manually:

```bash
go build -o agent .
```

This creates an executable `agent` binary in the current directory.

## Quick Start

### Basic Usage

```bash
./agent --config safeguards.json --prompt "your task here"
```

### Example Commands

Generate a simple Android activity:
```bash
./agent --config safeguards.json --prompt "Create a simple Android activity with a button"
```

Create a REST API client:
```bash
./agent --config safeguards.json --prompt "Generate a Retrofit HTTP client for API calls"
```

### Custom Configuration

Create your own `safeguards.json` with custom rules:

```bash
./agent --config my-safeguards.json --prompt "your task here"
```

## Configuration

The agent's behavior is controlled entirely by `safeguards.json`. Key sections:

- **file_system**: Allowed/forbidden paths and file types
- **code_generation**: Code size limits, languages, restricted imports
- **git_operations**: Allowed git operations and approval requirements
- **testing**: Test framework and coverage requirements
- **api_calls**: Whitelisted services and rate limits
- **security**: Malicious pattern scanning and review triggers
- **android_specific**: Android SDK versions and permission restrictions
- **llm_interaction**: LLM model settings and custom instructions

See [safeguards.json](safeguards.json) for full configuration details.

## Build Commands

```bash
# Build the agent
make build

# Run with default safeguards
make run

# Show custom run example
make run-custom

# Clean build artifacts
make clean

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Download dependencies
make deps

# Show all available commands
make help
```

## Architecture

The agent consists of four main components:

1. **Safeguard Engine** - JSON-based rule validation
2. **Agent Core** - Main logic and orchestration
3. **Custom LLM** - Local language model integration
4. **Validation Layer** - Security and risk assessment

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture.

## Usage Examples

### Example 1: Generate Android Fragment

```bash
./agent --config safeguards.json --prompt "Create a Fragment with RecyclerView for displaying a list of items"
```

### Example 2: Create Gradle Dependencies

```bash
./agent --config safeguards.json --prompt "Add Jetpack Compose dependencies to build.gradle"
```

### Example 3: Git Operations

```bash
./agent --config safeguards.json --prompt "Create a feature branch and commit changes"
```

## Safeguard Rules

The agent enforces strict safeguards:

- ✅ Only modifies allowed file types (.java, .kt, .xml, .gradle, etc.)
- ✅ Restricted to allowed directories (app/, src/, res/, etc.)
- ✅ Prevents dangerous operations (force push, history rewrites)
- ✅ Scans for malicious code patterns
- ✅ Enforces test coverage requirements
- ✅ Whitelists API services
- ✅ Resource limits (memory, CPU, timeout)

**You control everything** via `safeguards.json`.

## Troubleshooting

### Agent Won't Start

```bash
# Check Go version
go version

# Verify safeguards.json exists and is valid JSON
cat safeguards.json | python -m json.tool
```

### Build Fails

```bash
# Clean and rebuild
make clean
make deps
make build
```

### Permission Denied

```bash
# Make binary executable
chmod +x agent
./agent --config safeguards.json --prompt "test"
```

## Development

### Project Structure

```
.
├── main.go                          # Agent core
├── safeguards.json                  # Configuration
├── go.mod                           # Dependencies
├── llm/
│   └── llm.go                       # Custom LLM
├── safeguard_validator/
│   └── validator.go                 # Validation engine
├── docs/
│   └── ARCHITECTURE.md              # Architecture docs
└── Makefile                         # Build automation
```

### Adding Custom Safeguards

1. Add new rules to `safeguards.json`
2. Implement validation in `safeguard_validator/validator.go`
3. Update `Agent.ValidateOperation()` in `main.go`

## Security & Safeguards

### Whitelist-Based Security

The agent operates on a whitelist model:
- Only explicitly allowed operations are permitted
- All file paths must be in allowed directories
- API calls restricted to approved services
- Code patterns scanned for malicious content

### Risk Assessment

Each operation is scored (0.0 to 1.0):
- **0.0-0.3**: Safe - Automatic execution
- **0.3-0.7**: Medium - Log and monitor
- **0.7-1.0**: High - Require human review

### Audit Trail

All operations are logged:
- Decision logs with timestamps
- Risk scores and rationale
- API calls and modifications
- Failed operations and reasons

## Performance

- **Validation overhead**: <100ms per operation
- **LLM response**: Depends on model (50-5000ms typical)
- **Memory usage**: ~50MB base + model size
- **Max parallel operations**: 4 (configurable)

## Customization

### Modify Safeguards

Edit `safeguards.json` to:
- Add/remove allowed file paths
- Adjust code size limits
- Change allowed frameworks
- Configure API whitelists
- Modify permission restrictions

### Extend the Agent

Add new validation rules:
1. Define rule in safeguards.json
2. Implement checker in validator.go
3. Hook into ValidateOperation()

## Advanced Usage

### Custom LLM Integration

The agent supports custom LLM backends. Modify `llm/llm.go` to integrate:
- OpenAI GPT models
- Local language models
- Custom ML models
- Proprietary LLMs

### Risk Assessment Customization

Adjust risk scoring in `safeguard_validator/validator.go`:
- File operation weights
- Code pattern severity
- Git operation risks
- API call thresholds

## Support & Documentation

- [Architecture Documentation](docs/ARCHITECTURE.md)
- [Safeguards Configuration](safeguards.json)
- [Build Instructions](Makefile)

## License

Proprietary - toolie01

## Contributing

For internal use. Contact owner for modifications.

---

**Last Updated**: 2026-06-20
**Version**: 1.0
**Status**: Production Ready
