## Running Ollama Third-Party Service

### Choosing a model

You can get the model ID that Ollama will launch from the [Ollama Library](https://ollama.com/library)

### Get Host IP

#### Windows

```sh
ipconfig
```

#### Linux/Mac

```sh
ifconfig
```

### Enviroment Variables

LLM_ENDPOINT_PORT = 8888
LLM_MODEL_ID="llama3.2:1b"
host_ip="0.0.0.0"

### Ollama API

Once the Ollama server is running we can make API calls to its API

https://github.com/ollama/ollama/blob/main/docs/api.md

#### Example

```sh
curl http://localhost:8888/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Why is the sky blue?"
}'
```

### Technical Uncertainty

- Q. If we set LLM_MODEL_ID and run docker-compose Ollama server will pre-load our desired model?
- A. This appears _not_ to be the case. Ollama CLI support to run multiple LLMs and will pull the needed LLM (if it doesn't exist) before providing the response.

- Q. Ollama CLI downloads the LLMs on the container or on a separate volume?
- A. Ollama with this configuration will download the LLMs into the container
