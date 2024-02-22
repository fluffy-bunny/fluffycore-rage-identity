# Traefik Setup

## Thanks

[traefik-v2-https-ssl-localhost](https://github.com/Heziode/traefik-v2-https-ssl-localhost)

Next, go to the root of the repo and generate certificates using [mkcert](https://github.com/FiloSottile/mkcert) :

```powershell
# If it's the firt install of mkcert, run
.\mkcert.exe -install

# Generate certificate for domain "localhost.dev"
.\mkcert.exe -cert-file certs/local-cert.pem -key-file certs/local-key.pem "localhost.dev" "*.localhost.dev"

# incase you want a CRT as well.  Not needed though
cd certs
openssl x509 -outform DER -in .\local-cert.pem -out .\local-cert.crt

```

Create networks that will be used by Traefik:

```bash
docker network create smtp4dev
```

_Note: you can access to Træfik dashboard at: [traefik.localhost.dev](https://traefik.localhost.dev)_  
_Note: you can access to SMTP portal at: [smtp.localhost.dev](https://smtp.localhost.dev)_

for just http:

_Note: you can access to SMTP portal at: [localhost:4000](http://localhost:4000)_

## hosts file Windows

```txt
127.0.0.1 localhost.dev traefik.localhost.dev smtp.localhost.dev
```

## License

MIT
