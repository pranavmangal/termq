import json
import os

import requests

CEREBRAS_URL = "https://api.cerebras.ai/v1/models"
GROQ_URL = "https://api.groq.com/openai/v1/models"
GEMINI_URL = "https://generativelanguage.googleapis.com/v1beta/models"


def fetch_cerebras_models():
    api_key = os.environ["CEREBRAS_API_KEY"]
    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    }

    response = requests.get(CEREBRAS_URL, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Failed to fetch models from Cerebras: {response.text}")

    models = json.loads(response.text)
    return [model["id"] for model in models["data"]]


def fetch_groq_models():
    api_key = os.environ["GROQ_API_KEY"]
    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    }

    response = requests.get(GROQ_URL, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Failed to fetch models from Groq: {response.text}")

    models = json.loads(response.text)
    return [model["id"] for model in models["data"]]


def fetch_gemini_models():
    api_key = os.environ["GEMINI_API_KEY"]

    response = requests.get(GEMINI_URL, params={"key": api_key})
    if response.status_code != 200:
        raise Exception(f"Failed to fetch models from Gemini: {response.text}")

    models = json.loads(response.text)
    model_ids = [model["name"].rsplit("/", 1)[1] for model in models["models"]]
    return [model_id for model_id in model_ids if model_id.startswith("gemini")]


if __name__ == "__main__":
    cerebras_models = fetch_cerebras_models()
    groq_models = fetch_groq_models()
    gemini_models = fetch_gemini_models()

    models = {
        "cerebras": sorted(cerebras_models),
        "groq": sorted(groq_models),
        "gemini": sorted(gemini_models),
    }

    script_dir = os.path.dirname(os.path.abspath(__file__))
    output_path = os.path.join(script_dir, "available_models.json")

    with open(output_path, "w") as f:
        json.dump(models, f, indent=2)
