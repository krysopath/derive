# derive 

> a simple kdf cli utility

## installation

```
DERIVE_SALT=$(openssl rand 32)
go install github.com/krysopath/derive/cmd/derive@v0.1.3
cat <<EOF >> ~/.bashrc
export DERIVE_SALT=$DERIVE_SALT
EOF

source ~/.bashrc
```

## usage

```
$ derive -b 32 [purpose]

Usage of derive:
  -b int
    	length of derived key in bytes (default 32)
  -c int
    	rounds for deriving key (default 4096)
  -f string
    	key output format: bytes|base64|hex|ascii|ascii@shell (default "bytes")
  -h string
    	hash for kdf function (default "sha512")
  -k string
    	kdf function for deriving key (default "pbkdf2")
  -v string
    	'versioned' key  (default "1000")

```

### ssh private keys with passphrases unlocking autmatically?

create `$HOME/bin/ssh_give_pass.sh`
```
cat <<EOF> ~/bin/ssh_give_pass.sh
#!/bin/bash
cat
EOC

```
> a script like an echo server: `cat /dev/stdin > /dev/stdout`


add the shell function below to `.bashrc` e.g.
```
add_keyfile_to_agent() {
    if [ -n "$1" -a -r "$1" ]; then
        derive -b 32 \
            -f base64 \
            -v $(basename $1) ssh \
        | DISPLAY=:0 SSH_ASKPASS=$HOME/bin/ssh_give_pass.sh ssh-add $1
    fi
}
```

> this function will try to open a private key file and it to the ssh-agent

> this invocation derives a key and passes it via the SSH_ASKPASS script into ssh-add

> this method leaves no passprhrases on disk and does not use secrets in command arguments just in pipes
