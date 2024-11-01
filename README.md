# termq (tq)

A simple tool to query LLMs from the terminal.

[![GitHub Release](https://img.shields.io/github/v/release/pranavmangal/termq?sort=semver&style=for-the-badge&color=%2346A758)](https://github.com/pranavmangal/termq/releases/latest)

![Demo](docs/demo.gif)

## Usage

```bash
tq '<query>'
```

## Supported Providers & Models

- [Cerebras](https://cerebras.ai/inference) (very fast inference)
  - Llama 3.1
- [Groq](https://groq.com/) (fast inference)
  - Llama 3, 3.1 & 3.2
  - Gemma 1 & 2
  - Mixtral
- [Google Gemini](https://ai.google.dev/gemini-api)
  - Gemini 1.0, 1.5 Flash & Pro

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
