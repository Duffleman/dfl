rm -fr build/

echo '💿 Building Windows amd64 binaries'
GOOS=windows GOARCH=amd64 go build -o build/windows/win64-auth.exe ./cmd/auth/...
GOOS=windows GOARCH=amd64 go build -o build/windows/win64-certgen.exe ./cmd/certgen/...
GOOS=windows GOARCH=amd64 go build -o build/windows/win64-short.exe ./cmd/short/...
GOOS=windows GOARCH=amd64 go build -o build/windows/win64-update.exe ./cmd/update/...
echo '✅ Done'

echo '💻 Building macOS amd64 binaries'
GOOS=darwin GOARCH=amd64 go build -o build/macos/mac64-auth ./cmd/auth/...
GOOS=darwin GOARCH=amd64 go build -o build/macos/mac64-certgen ./cmd/certgen/...
GOOS=darwin GOARCH=amd64 go build -o build/macos/mac64-short ./cmd/short/...
GOOS=darwin GOARCH=amd64 go build -o build/macos/mac64-update ./cmd/update/...
echo '✅ Done'
