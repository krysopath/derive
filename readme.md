# derive 

A
> hollow anti token container wrapping noise
^based on true randomness

## installation

```
DERIVE_SALT=$(openssl rand 32)
go install github.com/krysopath/derive/cmd/derive@v0
cat <<EOF >> ~/.bashrc
export DERIVE_SALT=$DERIVE_SALT
EOF

source ~/.bashrc
```

> we assume your ssh-agent is properly setup

## usage

```
derive [FLAGS] [purpose]

FLAGS:
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

simple run:
```
$ derive -b 12 -f hex
! Enter Secret Token (hold Yubikey 5secs) ...OK
16F61AD0160EE71CAC668FC3
```

> consider using `| xclip -i -selection clipboard` to capture the results


### secure ssh pkeys with phrases unlocking automatically?

> quickly, hold your terminals!


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

> this method leaves no passphrases on disk and does not disclose exec arguments in `ps`


