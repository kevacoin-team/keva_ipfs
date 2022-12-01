# keva_ipfs
This is the backend for the Keva mobile app to upload and pin IPFS files.

# Prerequsites:

### Download and install ElectrumX
It assumes that you already have an EletrumX server running on the same server.


### Download and install IPFS
You must run an IPFS node on the server. To install follow the instructions here: https://docs.ipfs.tech/install/command-line/


### Download and install Golang
The backend is written in Golang and you need Golang to build the server: To install follow the instructions here: https://go.dev/doc/install


# Build keva_ipfs
Clone this repo and build the server:

```
git clone https://github.com/kevacoin-project/keva_ipfs
cd keva_ipfs
go build .
```

# Configuration variables:
### config.json
#### *example Payment_address value set, ensure to update.*
```
Electrum_host
Electrum_port
Min_payment
Payment_address
Tls_enabled
Tls_key
Tls_cert
```


# To start the server:
This server will only pin an IPFS file if a payment is made to `Payment_address`, with the minimal Kevacoin amount defined by `Min_payment`.

The backend listens to port `Electrum_port + 10`. E.g. if the ElectrumX server is listening on port `50001`, this backend will listen to port `50011`.

```
./keva_ipfs
```
