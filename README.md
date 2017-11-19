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
message := client.Message(chatbase.MessageTypeAgent, "USER-ID", chatbase.PlatformTelegram)

// optional values are added after creation
message.SetMessage("I didn't understand that, sorry")

// this can also be chained
message.SetIntent("fallback").SetNotHandled(true)

// calling Submit will send the data to chatbase
response, err := message.Submit()
if err != nil || !response.Status.OK() {
	// handle errors
}
```

## Supported APIs

### Generic message API

The [generic message API](https://chatbase.com/documentation/generic) allows handling of `Message`, `Messages` and `Update` types.

#### `Message`

```go
message := client.Message(chatbase.UserType, "USER-ID", "messenger")
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
	message := client.AgentMessage("USER-ID", "messenger")
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

### Facebook Message API

The [Facebook Message API](https://chatbase.com/documentation/facebook) allows handling of `FacebookMessage`, `FacebookMessages`, `FacebookRequestResponse` and `FacebookRequestResponses` types.

#### `FacebookMessage`

```go
message := client.FacebookMessage(facebookPayload)
message.SetIntent("test-messenger")
response, err := message.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

#### `FacebookMessages`

```go
messages := chatbase.FacebookMessages{}
for _, msg := range listOfFacebookMessages {
	message := client.FacebookMessage(msg).SetVersion("0.0.1-beta")
	messages.Append(message)
}

response, err := messages.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

#### `FacebookRequestResponse`

```go
pair := client.FacebookRequestResponse(incomingMessage, respondingMessage)
pair.SetIntent("test-messenger")
response, err := pair.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

#### `FacebookRequestResponses`

```go
pairs := chatbase.FacebookRequestResponses{}
for _, msg := range listOfFacebookMessages {
	pair := client.FacebookRequestResponse(msg.request, msg.response).SetVersion("0.0.1-beta")
	pairs.Append(pair)
}

response, err := pairs.Submit()
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

### Link tracking API

The [link tracking](https://chatbase.com/documentation/taps) allows handling of `Link` types.

#### `Link`

```go
link := client.Link("https://golang.org/", chatbase.PlatformLine)
trackableHREF, err := link.Encode()
response, err := link.Submit()
if err != nil || !response.Status.OK() {
	// handle error
}
```

### License
MIT © [Frederik Ring](http://www.frederikring.com)
