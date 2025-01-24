# OAuth 2.0 Authorization Code Flow CLI

This is a utility program that you can use to get an access token using the authorization code flow without actually needing to build a frontend. You should provide a configuration file that looks like:

```json
{
  "server": {
    "auth_url": "http://localhost:8080/realms/your_client/protocol/openid-connect/auth",
    "token_url": "http://localhost:8080/realms/your_client/protocol/openid-connect/token",
    "client_id": "your_client",
    "client_secret": "your_client_secret",
    "redirect_uri": "http://localhost:8666/callback",
    "scope": "openid profile email"
  },
  "credentials_cache": "credentials.json",
  "skip_cache": false
}
```

Your auth provider of choice should expose all these parameters in an OpenID configuration document. Also, please ensure that you make `http://localhost:8666/callback` a valid `redirect_uri`. This could be changed to be a bit more flexible but as of right now it doesn't seem worth it.

To run the program, you need `Go`. Once you have it, you can run the command:

- `go build -ldflags "-w" -o ./bin/oauth-cli`

which will build the binary into `./bin`. To run it, you can `./bin/oauth-cli --config=/path/to/your/config`. You may also run it without a server file and instead pass your arguments as CLI flags, e.g. `--auth-url=your/auth/url`. Flags override config.

If you have `just`, which I highly recommend, there are two convenience recipes in the distributed `Justfile`. `build` will build your program and `buildrun` will build and then execute it.


#### License

BSD 3-Clause
