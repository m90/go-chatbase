/*
Package chatbase provides a client and types for interacting with the
Chatbase bot analytics platform - https://chatbase.com

The primary entrypoint for interacting with Chatbase is a `Client` instance:

	// a client is a wrapper around an chatbase API key
	client := chatbase.New("MY_CHATBASE_API_KEY")

	// calling Message requires passing of all required values
	message := client.Message(chatbase.MessageTypeAgent, "USER-ID", chatbase.PlatformTelegram)

	// optional values are added after creation
	message.SetMessage("I didn't understand that, sorry")

	// this can also be chained
	message.SetIntent("fallback").SetNotHandled(true)

	// calling Submit will send the data to chatbase
	response, err := message.Submit()
	if err != nil {
		// an error submitting the data occurred
		fmt.Println(err)
	} else if !response.Status.OK() {
		// the data was submitted to ChatBase, but
		// the response contained an error code
		fmt.Println(response.Reason)
	}

For further usage instructions refer to the README.

*/
package chatbase
