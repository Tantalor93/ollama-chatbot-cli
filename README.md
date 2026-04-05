# Ollama chatbot CLI

just a test project to get on hands experience on LLM concepts, it is a simple 
chatbot able to communicate with Ollama server

## Local development

### Run Ollama server

install Ollama

```
brew install Ollama
```

run Ollama
```
ollama serve
```

pull model
```
ollama pull llama3.2
```


### Build & Run chatbot

```
go build
```

```
./ollama-chatbot-cli
```
