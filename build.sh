rm -fr build/

echo '💿 Building Windows amd64 binaries'
GOOS=windows GOARCH=amd64 go build -o build/windows/auth.exe ./cmd/auth/...
GOOS=windows GOARCH=amd64 go build -o build/windows/certgen.exe ./cmd/certgen/...
GOOS=windows GOARCH=amd64 go build -o build/windows/short.exe ./cmd/short/...
echo '✅ Done'

echo '💻 Building macOS amd64 binaries'
GOOS=darwin GOARCH=amd64 go build -o build/macos/amd64/auth ./cmd/auth/...
GOOS=darwin GOARCH=amd64 go build -o build/macos/amd64/certgen ./cmd/certgen/...
GOOS=darwin GOARCH=amd64 go build -o build/macos/amd64/short ./cmd/short/...
echo '✅ Done'
