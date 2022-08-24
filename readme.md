# derive 

A

> hollow anti token container wrapping noise

based on true random wisdom from `xkcd-pass`

The idea of this is to implement a scriptable way for adding secret keys to
keyring agents. Solutions recommended on the web often forget, that on
multiuser systems the process tree discloses command arguments to every logged
user. They may end up in remote logs and endanger your private keys. Also it is
not a solution to write key pass phrases into config files. However a fair
amount of this sin has been committed. Even by accident, it may happen that an
encrypted file is not encrypted properly on disk and then checked into source
control. In infrastructure teams often passwordless and shared private keys are
used. These are all desaster scenarios or severe anti pattern.

It is recommended to setup an ssh-agent securely, such that it integrates into
the os keyring system. People can choose to even integrate a smartcard to host
the private keys. But how to deal with private keys, that can not be hosted as
such? You absolutely should store them encrypted (that means they have a
passphrase guarding their usage)

However supplying passphrases to private keys when they are added to the agent
is a manual effort. The default tooling is build for interactive input. The
goal of this project is to find a way to programmatically derive key phrases in
such a way that they can be fed to ssh-agent. This method should be more secure
than keys without encryption passphrase.


- This tool can derive new keys based on salt and secret key factors via pbkdf2.
- Use if you cant place a key on smartcard, but also dont want to use passwordless 
  keys
- Programmatically create and use encrypted key material in pipelines
- Avoid communicating secret passphrases between departments. Procedurally
  derive one twice.
- better even: use a dedicated crypto host and smartcards for when it matters!


## installation
### go

```
go install github.com/krysopath/derive/cmd/derive@v1
```

### compile

```
git clone git@github.com:krysopath/derive.git
cd derive

# checking deps for build
make deps

# testing and building and installing code in  workdir
make install
```


## roadmap

- a better method to receive a kdf result
    - to not leak the secret to the consumer OS
    - but run a KDF inside the smartcard, with a secret from the smartcard
- statefile for count of operations per key topic
- blinking yubikey lights
- more KDF juice
- when outputting as ascii, it might happen that several hundred bytes of input do not contain printable characters, this MUST be mitigated
    - when such an unprintable char is received it will be rejected
    - currently there is a mitigation but no prevention:
        - we generate 2x as much bytes and discard what does not fit
        - this is assuming we can fit the requested bytes, though
        - unlikely as it sounds:
        - secret with only empty bytes is possible and that would break the security
        - for such a case a prevention is planned, but not properly implemented
- audit

## contributions 
> are always welcome

- leave an issue if you are missing something
- raise a problem if you see one
- suggest expected usecases
- `<3`

## setup

```
DERIVE_SALT=$(openssl rand -base64 48)
go install github.com/krysopath/derive/cmd/derive@v1
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
> `derive` reads a salt from the environment and waits for a newline character
> on `/dev/stdin`. It uses those values and the passed arguments to derive a
> key and emit it to `dev/stdout

> the author used a usb HID rubberduck to output a static secret string

> a yubikey may be programmed to emit a static key of 38 bytes via HID

> a master password can be remembered and concatened with the rubberduck static key, too!

> consider using `| xclip -i -selection clipboard` to capture the results

### derive with yubikey?

Yes, you need to configure slot 2 for emitting a static secret:

`ykman otp static --generate --length 38  --keyboard-layout US 2`

> maxlength is 38 bytes

> slot 1 hosts u2f, lets not overwrite it.

> static OTP bytes are static :.(


- Such a static key has no smooth rotation
- if it leaked once (remember it is emitted plaintext stdout), then your secrets are void
- though together with `derive` that static key can be used to derive many more keys
- and at the same time it prevents disclosure of that key by accidental stdout shell
  blooper (because it emits after a long press and is hashed)
- You only trust your host, not a remote system.
- However you trust your host enough with your static key. This is always a risk.
- Be sure to disable the static code feature for the NFC channel tho. Would be embarassing.
- Be sure to keep a backup rubberduck, yubikey or else with the same static
  key, else it would embarassing too.


### secure ssh pkeys with phrases unlocking automatically?

> quickly, hold your terminals!


create `$HOME/bin/ssh_give_pass.sh`
```
cat <<EOF> ~/bin/ssh_give_pass.sh
#!/bin/bash
cat
EOF
```
> magic script is implementing the API of /bin/cat, like an echo server: `cat </dev/stdin >/dev/stdout`

> key derivation could happen in here, but would be less flexible then.

After you add this shell function below to `.bashrc` e.g.
```
cat <<EOF>> ~/.bashrc
add_keyfile_to_agent() {
    if [ -n "\$1" -a -r "\$1" ]; then
        derive -b 32 \\
            -f base64 \\
            -v \$(basename $1) ssh \\
        | DISPLAY=:0 SSH_ASKPASS=\$HOME/bin/ssh_give_pass.sh ssh-add \$1
    fi
}
EOF

# source and run
source ~/.bashrc
add_keyfile_to_agent ~/.ssh/id_rsa
```
> this function will try to open a private key file and add it to the ssh-agent

> this invocation derives a key and passes it via the SSH_ASKPASS script into ssh-add

> this method leaves no passphrases on disk and does not disclose exec arguments in `ps`


