# GoTelegram

Asks For a Chat ID and an Api Key.

Usage example:

```go
  chat := telegram.Chat("any name", "<chatID>", "apiKey")
```

To send a message:

```go
  chat.Send("Hi mom!")
```

Messages that are longer then the limit (4096 characters) will be sent in chunks.
