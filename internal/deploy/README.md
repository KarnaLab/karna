# Karna Deploy

## How it works

Karna Deploy will package the target folder and place it in the .karna / <functionName> / <alias> / <karna.json>.file folder, then upload the code,
either on S3 on a bucket is specified in karna.json, or directly on lambda, then publish the function and finally tag the version
specified in karna.json.
Karna Deploy can also remove aliases not mentioned in the configuration file, range of versions, and create new aliases
via the prune option.

## Options

`{ ..., "prune": { "keep": <int>, "alias": <bool> } }`
If alias is specified, it will destroy ALL aliases which dit not match with aliases in karna.json.
If keep is specified, it will destroy all versions which dit not match to the pattern: <each-alias-version> + range to <each-alias-version> - range

## Commands

Karna deploy requires two flags:

- --alias or -a => Alias name
- --target or -t => Function name

You can map the alias with the following options:

- a string that corresponds to an existing version
- "latest" will tag the alias on the \$LATEST version
- "fixed@update": tag the alias on the LAST version that the user is deploying and not on the \$LATEST version. It
  prevents other deployments from overwriting the current version on which the alias is placed.
