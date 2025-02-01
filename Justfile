build:
  go build -ldflags "-w" -o ./bin/ezioauth

buildrun CONFIG:
  go build -ldflags "-w" -o ./tmp/ezioauth && ./tmp/ezioauth --config-file={{CONFIG}}


