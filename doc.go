/*
Package chatbase provides a client and types for interacting with the
Chatbase bot analytics platform - https://chatbase.com

The primary entrypoint for interacting with Chatbase is a `Client` instance:

	// a client is a wrapper around an chatbase API key
	client := chatbase.NewClient("MY_CHATBASE_API_KEY")
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

For further usage instructions refer to the README.

*/
package chatbase
