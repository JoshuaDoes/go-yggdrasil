# go-yggdrasil

[![GoDoc](https://godoc.org/github.com/JoshuaDoes/go-yggdrasil?status.svg)](https://godoc.org/github.com/JoshuaDoes/go-yggdrasil)
[![Go Report Card](https://goreportcard.com/badge/github.com/JoshuaDoes/go-yggdrasil)](https://goreportcard.com/report/github.com/JoshuaDoes/go-yggdrasil)

Single file library for Mojang's Yggdrasil API written in Golang without extra dependencies

# Installing
`go get github.com/JoshuaDoes/go-yggdrasil`

# Recommended authentication procedure
**Note: This procedure is based on the official Minecraft launcher's authentication logic.**

When authenticating with Yggdrasil, you should use ``*yggdrasil.Client.Authenticate()`` only if you do not already have an access/client token pair stored somewhere, as Yggdrasil's ``/authenticate`` endpoint is severely rate-limited. If you have an access/client token pair, you should use ``*yggdrasil.Client.Validate()`` on a pre-initialized ``*yggdrasil.Client`` variable to make sure it is still valid. If it is not, you should then use ``*yggdrasil.Client.Refresh()`` to get a new access token and have it automatically be stored in the ``*yggdrasil.Client`` variable. Should this fail, however, then it is safe to use ``*yggdrasil.Client.Authenticate()`` again and store the returned access/client token pair for future use.

# Additional notes
When an authentication is successful, your ``*yggdrasil.Client`` variable will contain the client's access/client token pair, the current selected profile (changing this is not yet implemented in Yggdrasil), and the current user.

When using any available go-yggdrasil functions, any internal errors will be returned as ``*yggdrasil.Error.FuncError`` rather than be classically available as an ``error`` type. Any errors returned from Yggdrasil itself will be returned as ``*yggdrasil.Error``.

# Authentication example
```go
package main

import "fmt"
import "github.com/JoshuaDoes/yggdrasil"

var yggdrasilClient *yggdrasil.Client

func main() {
	yggdrasilClient = &yggdrasil.Client{ClientToken: "your client token here"}

	//Auth with Minecraft version 1
	authResponse, err := yggdrasilClient.Authenticate("your email/username here", "your password here", "Minecraft", 1)
	if err != nil {
		fmt.Println("Error: " + fmt.Sprintf("%v", err))
	} else {
		//Print access/client token pair
		fmt.Println("Access Token: " + authResponse.AccessToken)
		fmt.Println("Client Token: " + authResponse.ClientToken)
	}
}
```
### Output
```
> go run main.go
Access Token: 32 char hexadecimal access token
Client Token: specified client toen
```

## License
The source code for go-yggdrasil is released under the MIT License. See LICENSE for more details.

## Donations
All donations are appreciated and helps me stay awake at night to work on this more. Even if it's not much, it helps a lot in the long run!
You can find the donation link here: [Donation Link](https://paypal.me/JoshuaDoes)