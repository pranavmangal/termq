# termq (tq)

[![GitHub Release](https://img.shields.io/github/v/release/pranavmangal/termq?sort=semver&style=for-the-badge&color=%2346A758)](https://github.com/pranavmangal/termq/releases/latest)

A simple tool to query LLMs from the terminal.

![Demo](docs/demo.gif)

## Installation

Using **Homebrew**:

```bash
brew install pranavmangal/tap/termq
```

## Usage

```bash
tq '<query>'
```

## Supported Providers

- [Cerebras](https://cerebras.ai/inference) (very fast inference)
- [Groq](https://groq.com/) (fast inference)
- [Google Gemini](https://ai.google.dev/gemini-api)

## Configuration

The config file is located at `~/.config/termq/config.toml` (macOS and Linux) or `~\AppData\Roaming\termq\config.toml` (Windows).

`termq` should automatically create a skeletal config for you at first run and ask you to fill in your API keys and preferred models.

Example:

```toml
system_prompt = "You are a helpful assistant."

[groq]
model = "llama-3.1-70b-versatile"
api_key = "<your-api-key>"
```
