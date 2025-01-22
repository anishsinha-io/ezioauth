build:
  go build -ldflags "-w" -o ./bin/oauth-cli

buildrun CONFIG:
  go build -ldflags "-w" -o ./tmp/oauth-cli && ./tmp/oauth-cli --server-config={{CONFIG}}


