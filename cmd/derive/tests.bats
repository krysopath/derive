#!/usr/bin/env bats

# !!!derive failing these tests mean work, because downstream users generate different passwords
# go fix the code & make test green OR adapt the test and release a major version
#

@test run_derive_with_empty_salt_and_fail {
  unset DERIVE_SALT
  run ./derive
  [ "$status" -eq 2  ]
}

@test run_derive_with_short_salt_and_fail {
  export DERIVE_SALT="0123456789abcde"
  run ./derive
  [ "$status" -eq 2  ]
}

@test run_derive_hex {
  export DERIVE_SALT="0123456789abcdef"
  goal='66B43A1E74FEFFF062F06B2CFE8E65F7ABC51094BB13AD11A8BC1D1BBB52166B'
  value=$(echo -e "secret\n" | ./derive -b 32 -f hex)
  [ "$?" -eq 0  ] 
  [ "$value" = "$goal" ]  || { echo ERR: \"$value\" not matching \"$goal\"; exit 1; }
}

@test run_derive_base64 {
  export DERIVE_SALT="0123456789abcdef"
  goal='ZrQ6HnT+//Bi8Gss/o5l96vFEJS7E60RqLwdG7tSFms'
  value=$(echo -e "secret\n" | ./derive -b 32 -f base64)
  [ "$?" -eq 0  ] 
  [ "$value" = "$goal" ]  || { echo ERR: \"$value\" not matching \"$goal\"; exit 1; }
}

@test run_derive_ascii {
  export DERIVE_SALT="0123456789abcdef"
  goal='f5:tqbqk,ex,F<.)=<Rkh1wr}#\Ri1lf'
  value=$(echo -e "secret\n" | ./derive -b 32 -f ascii)
  [ "$?" -eq 0  ] 
  [ "$value" = "$goal" ]  || { echo ERR: \"$value\" not matching \"$goal\"; exit 1; }
}

@test run_derive_ascii_shell_escaped {
  export DERIVE_SALT="0123456789abcdef"
  goal="'f5:tqbqk,ex,F<.)=<Rkh1wr}#\Ri1lf'"
  value=$(echo -e "secret\n" | ./derive -b 32 -f ascii@escape)
  [ "$?" -eq 0  ] 
  [ "$value" = "$goal" ]  || { echo ERR: \"$value\" not matching \"$goal\"; exit 1; }
}

@test run_derive_ascii_length_16 {
  export DERIVE_SALT="0123456789abcdef"
  goal=17
  value=$(echo -e "secret\n" | ./derive -b 16 -f ascii)
  [ "$?" -eq 0  ] 
  [ "$(wc -c <<<"$value")" -eq "$goal" ]  || { echo ERR: \"$value\" not matching \"$goal\"; exit 1; }
}

@test run_derive_ascii_length_92 {
  export DERIVE_SALT="0123456789abcdef"
  goal=93
  value=$(echo -e "secret\n" | ./derive -b 92 -f ascii | wc -c)
  [ "$?" -eq 0  ] 
  [ "$value" -eq "$goal" ]  || { echo ERR: \"$value\" not matching \"$goal\"; exit 1; }
}
