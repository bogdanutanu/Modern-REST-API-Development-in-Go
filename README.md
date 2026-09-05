<p align='center'><a href='https://www.packtpub.com/en-us/unlock?step=1'><img src='https://static.packt-cdn.com/assets/images/packt+events/finalGH_design_redeem.png'/></a></p>

# Modern-REST-API-Development-in-Go
Modern REST API Development in Go, Published by Packt

# TLS

In order to run the server locally with TLS certificates compatible with `localhost`, 
you'll need to generate them with [mkcert](https://github.com/filosottile/mkcert).

Follow the instalation instructions from the MKCert README, and make sure you **never** 
commit the certificates. 

To generate a certificate in the `certs` directory for `localhost`:
```
mkcert \
 -cert-file certs/local-cert.pem \
 -key-file certs/local-key.pem \
 localhost 127.0.0.1 ::1
```

Then install this local CA in the system trust store:
```
mkcert -install
```

# Example requests

Set the protocol once before running the requests. It defaults to `https`; use
`API_PROTOCOL=http` to connect to a server without TLS.

```sh
API_PROTOCOL="${API_PROTOCOL:-https}"
API_URL="${API_PROTOCOL}://localhost:8888"
```

## Login

Regular user:
```
curl -X POST \
  -d '{"username": "user", "password": "password"}' \
  -H "Content-Type: application/json" \
  "$API_URL/login"
```

Admin:

```
curl -X POST \
  -d '{"username": "admin", "password": "password"}' \
  -H "Content-Type: application/json" \
  "$API_URL/login"
```

Run `eval "$(./login.sh)"` to store the token into the TOKEN env var.

## List all lists
```
curl "$API_URL/lists" -H "Authorization: Bearer $TOKEN"
```

## Create List
```
curl -X POST "$API_URL/lists" \
  -H "Authorization: Bearer $TOKEN" \
  --data '{"id": 1, "name": "my first shopping list", "items": ["eggs", "milk"]}'
```

## Delete
```
curl -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  "$API_URL/v1/lists/1"
```

## Profiling
```
curl http://localhost:8888/debug/pprof/
```