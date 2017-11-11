# go-chatbase

[![Build Status](https://travis-ci.org/m90/go-chatbase.svg?branch=master)](https://travis-ci.org/m90/go-chatbase)
[![godoc](https://godoc.org/github.com/m90/go-chatbase?status.svg)](http://godoc.org/github.com/m90/go-chatbase)

> Golang client for interacting with the [chatbase.com](https://chatbase.com) bot analytics API

## Installation

Use `go get`:

```sh
$ go get github.com/m90/go-chatbase
```

## Example

Send a single message to Chatbase:

```go
// a client is a wrapper around an chatbase API key
client := chatbase.NewClient("MY-API-KEY")

// calling Message requires passing of all required values
message := client.Message(chatbase.MessageTypeAgent, "USER-ID", "telegram")

// optional values are attached lateron
message.SetMessage("I didn't understand that, sorry")

// this can also be chained
message.SetIntent("fallback").SetNotHandled(true)

// calling Submit will send the data to chatbase
response, err := message.Submit()
```

## Supported APIs

### Generic message API

The [generic message API](https://chatbase.com/documentation/generic) allows handling of `Message`, `Messages` and `Update` types.

#### `Message`

```go
message := client.UserMessage("USER-ID", "messenger")
message.SetMessage("How are you today?")
response, err := message.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

#### `Messages`

```go
messages := chatbase.Messages{}
for i := 1; i < 4; i++ {
	message := client.UserMessage("USER-ID", "messenger")
	message.SetMessage(fmt.Sprintf("Counting: %d", i))
	messages.Append(message)
}
response, err := messages.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

#### `Update`

```go
update := client.Update("ID-OF-MESSAGE-TO-UPDATE")
update.SetIntent("this-changed")
response, err := update.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

### Events API

The [Events API](https://chatbase.com/documentation/events) allows handling of `Event` and `Events` types.

#### `Event`

```go
event := client.Event("USER-ID", "intent-name")
event.SetPlatform("line").SetVersion("2.2.1")
event.AddProperty("success", true)
if err := event.Submit(); err != nil {
	// handle error
}
```

#### `Events`

```go
events := chatbase.Events{}
for i := 0; i < 99; i ++ {
	event := client.Event("USER-ID", "counting-intent")
	event.AddProperty("is-even", i%2 == 0)
	events.Append(event)
}
if err := events.Submit(); err != nil {
	// handle error
}
```

### License
MIT © [Frederik Ring](http://www.frederikring.com)
