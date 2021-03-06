# certgen

A CLI tool to assist in generating and handling your own CA and certificates for servers and clients. The tool is exceptionally opinionated to my current personal use cases, but I hope others find it useful regardless.

## Assumptions

I may or not may not adjust the tool in the future to work around these assumptons, but they're listed here for now.

- All certificates are generated as ECDSA certificates.
- All private keys are `P384`
- P12 formatted client certs must have passwords

## Installation

No command that produces files will override any file that already exists, it's just not safe when it comes to certificates and chains.

You must set the `CERTGEN_SECRETS_ROOT_DIR` so the tool knows where to save the certificates it generates.

## Usage

You should always begin by generating a root CA. All other commands will fail until you have a CA saved.

When you run the CA command, it'll run an initiation, some folder creation and some templating. It is required you go into your certificate folder and update the template so the tool uses your common name and expiry dates over mine.

### `generate_root_ca`

`certgen gca`

### `generate_server_certificate`

`certgen gsc mydomain.com`

### `generate_client_certificate`

`certgen gcc computername --password {password}`

### `generate_key_pair`

`certgen gkp pairname`
