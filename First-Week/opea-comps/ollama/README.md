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

#### Preload Models

Using the query suggested [here](https://stackoverflow.com/questions/78500319/how-to-pull-model-automatically-with-container-creation) we have pre-loaded Ollama with llama3.2:1b.

#### Example

```sh
curl http://localhost:8888/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Why is the sky blue?"
}'
```

### Open Web-UI

We're using Open Web-UI here to interact with the Ollama container. This will allow us to have human readeable output form ollama as well as a better medium to comunicate with it.

### Technical Uncertainty

- Q. If we set LLM_MODEL_ID and run docker-compose Ollama server will pre-load our desired model?
- A. This appears _not_ to be the case. Ollama CLI support to run multiple LLMs and will pull the needed LLM (if it doesn't exist) before providing the response.

- Q. Ollama CLI downloads the LLMs on the container or on a separate volume?
- A. Ollama with this configuration will download the LLMs into the container

- Q. Will WebUI be able to download a model for us?
- A. Ollama does have an endpoint (api/pull) that's suppose to allow external requests for models, by default this api don't work. So we're preloading a model by telling the container to download llama3.2:1b.

- Q. Will Web-UI be available while ollama is loading?
- A. Web-UI components run independently from Ollama, so the UI itself will be available even though the interaction with ollama will be limited until the container finish downloading the model.
