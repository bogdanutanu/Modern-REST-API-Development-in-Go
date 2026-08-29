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