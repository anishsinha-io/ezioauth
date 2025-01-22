# OAuth 2.0 Authorization Code Flow CLI

This is a utility program that you can use to get an access token using the authorization code flow without actually needing to build a frontend. You should provide a server configuration file that looks like:

```json
{
  "auth_url": "http://localhost:8080/realms/your_realm/protocol/openid-connect/auth",
  "token_url": "http://localhost:8080/realms/your_realm/protocol/openid-connect/token",
  "client_id": "your_realm",
  "client_secret": "your_secret",
  "redirect_uri": "http://localhost:8666/callback",
  "scope": "openid profile email"
}
```

Your auth provider of choice should expose all these parameters in an OpenID configuration document. To run the program, you need `Go`. Once you have it, you can run the command:

- `go build -ldflags "-w" -o ./bin/oauth-cli`

which will build the binary into `./bin`. To run it, you can `./bin/oauth-cli --server-config=/path/to/your/config`. You may also run it without a server file and instead pass your arguments as CLI flags, e.g. `--auth-url=your/auth/url`. Flags override config.

If you have `just`, which I highly recommend, there are two convenience recipes in the distributed `Justfile`. `build` will build your program and `buildrun` will build and then execute it.


#### License

BSD 3-Clause
