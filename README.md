# IRCd-Healthcheck

A simple IRCd health-check utility, written in go, which returns a non-zero status code if unable to connect to a given IRC server.  Perfect for use with mointoring systems.  [Download a binary of the latest release here](https://github.com/StormBit/ircd-healthcheck/releases/latest).

### Usage: `./healthcheck -server=irc.stormbit.net:6667`

Advanced usage:
```
$ ./healthcheck -server=ridley.stormbit.net:6697 -secure=true -skip-verification=true
$ echo $?
0 
```

Options:
```
  -secure
      Whether or not to use SSL/TLS when connecting.
  -server string
      Server and port to connect to. (format: irc.example.org:6667) (default "irc.stormbit.net:6667")
  -skip-verification
      Whether or not to skip verifying the certificate.
  -h
      Shows help.
```

---

Copyright (c) 2016, Alex Wilson <a@ax.gy>

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

