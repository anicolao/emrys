# Emrys

> *"Your personal AI assistant, running entirely on your Mac."*

Emrys is a personal AI assistant inspired by fiction like Jane (Ender's Game) and Jarvis (Iron Man). It runs on a dedicated Mac Mini, using open-source AI models to provide intelligent assistance while maintaining complete privacy and user control.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## Overview

Unlike cloud-based assistants, Emrys runs entirely on your own hardware using local AI models. It has complete control of your Mac through terminal access, AppleScript, browser automation, and other system tools, enabling it to truly act as a capable personal assistant.

**Key Features:**
- 🏠 **100% Local**: All AI inference runs on your Mac Mini—no cloud dependencies
- 🔒 **Privacy First**: Your data never leaves your machine
- 🤖 **Agentic Intelligence**: Uses advanced task planning and tool orchestration
- 🎯 **Complete Control**: Full access to terminal, AppleScript, browser, and system APIs
- 🌐 **Web Presence**: Automated browser for online tasks and research
- 🔓 **Open Source**: Built on open-source models and tools (GPLv3)

## Vision

Emrys aims to be the kind of AI assistant portrayed in science fiction—capable, trustworthy, and truly helpful. See [VISION.md](VISION.md) for our long-term goals and philosophy.

## Prerequisites

### Hardware
- **Mac Mini** (M1/M2/M3 recommended for optimal AI performance)
  - Minimum 16GB RAM (32GB+ recommended)
  - 100GB+ free storage for AI models
- Dedicated exclusively to Emrys, like you'd set up a computer for a human assistant with its own credentials and accounts

### Software
- **macOS** 12.0 (Monterey) or later
- **nix-darwin** for system configuration and package management

## Installation

> **Note:** Emrys is currently in early development. Full installation instructions will be provided as the project matures.

## Architecture

Emrys consists of several key components:

- **AI Core**: Local LLM inference using Ollama/llama.cpp
- **Agent Framework**: Task planning and tool orchestration
- **System Interface**: Terminal, AppleScript, and system API access
- **Browser Automation**: Playwright-based web interaction
- **Knowledge Base**: Local document indexing and retrieval
- **Interface**: Internet-based communication (FaceTime, Messages, Email) - Emrys operates as a separate entity

## Usage Examples

Once installed, Emrys can help with tasks like:

```
"Check my email and summarize anything urgent" (if you've given Emrys access)
"Schedule a meeting with John next week"
"Research the best price for a Mac Studio"
"Organize my downloads folder" (if you've shared it with Emrys)
"Monitor system resources and alert me if anything unusual happens"
"Write a Python script to parse these log files"
```

## Development Status

🚧 **Early Development** - Emrys is currently in the early stages of development. Core architecture and foundational components are being built.

**Current Priorities:**
1. Core agentic framework
2. LLM integration with local models
3. Basic tool implementations (terminal, AppleScript, browser)
4. Task planning and execution engine

See [VISION.md](VISION.md) for our development roadmap and long-term goals.

## Contributing

Emrys is open source and welcomes contributions! Whether you're interested in:
- Core framework development
- Tool integrations
- Documentation
- Testing and bug reports
- Use case exploration

Please see `CONTRIBUTING.md` (coming soon) for guidelines.

## Privacy & Security

Emrys takes privacy and security seriously:

- ✅ All AI processing happens locally—no data sent to external servers
- ✅ Open source and auditable
- ✅ You maintain complete control over your data
- ⚠️ **Be cautious**: If you share access with Emrys, it has control and could potentially destroy it
- ✅ Comprehensive logging of all assistant actions

**Note:** Emrys operates as its own entity with its own internet accounts and credentials (not your personal accounts). Only share access to resources you're comfortable having Emrys manage.

## License

Emrys is licensed under the GNU General Public License v3.0. See [LICENSE](LICENSE) for details.

This means:
- ✅ Free to use, modify, and distribute
- ✅ Must remain open source if you distribute it
- ✅ No warranty or liability

## Acknowledgments

**Inspired by:**
- **Jane** from Orson Scott Card's *Ender's Game* series
- **Jarvis** from Marvel's *Iron Man*
- The open-source AI community

**Built with:**
- Open-source language models (Llama, Mistral, Qwen, etc.)
- Ollama and llama.cpp for inference
- Playwright for browser automation
- The amazing macOS developer community

## Contact & Community

- **GitHub**: [anicolao/emrys](https://github.com/anicolao/emrys)
- **Issues**: [GitHub Issues](https://github.com/anicolao/emrys/issues)
- **Discussions**: [GitHub Discussions](https://github.com/anicolao/emrys/discussions)

## Disclaimer

Emrys is experimental software. It has the ability to execute commands and control its Mac. Use it at your own risk, always on a dedicated machine. Always review what your assistant is doing and maintain appropriate backups.

---

*"Any sufficiently advanced technology is indistinguishable from magic."* — Arthur C. Clarke [sic]
