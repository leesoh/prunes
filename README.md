# Prunes

Prunes checks a list of subdomains against a list of resolvers. Successful responses are printed to stdout. By providing a non-existent domain name, this can be used to identify resolvers that do NXDOMAIN hijacking.

## Installation 

`$ go get -u github.com/leesoh/prunes`

## Usage

By default, Prunes will only print the resolvers that **do not** respond to the queries provided. This allows you to use its output to build a list of good resolvers:

```sh
$ cat resolvers.txt
1.1.1.1
8.8.8.8
9.9.9.9
50.49.243.135
50.120.215.2

$ cat subdomains.txt
www
mail
remote
blog
webmail
server
ns1
ns2
smtp
secure

$ cat resolvers.txt | ./prunes
1.1.1.1
8.8.8.8
9.9.9.9
```

If you're more interested in finding out which name servers are hijacking invalid responses, use the `-details` flag:

```sh
$ cat resolvers.txt | ./prunes -details
50.49.243.135 ::: www.gCVHwgoDl55QLmHw.com => 23.217.138.109
50.49.243.135 ::: www.gCVHwgoDl55QLmHw.com => 23.195.69.108
50.120.215.2 ::: www.gCVHwgoDl55QLmHw.com => 23.217.138.109
50.120.215.2 ::: www.gCVHwgoDl55QLmHw.com => 23.202.231.168
50.49.243.135 ::: webmail.gCVHwgoDl55QLmHw.com => 23.217.138.109
50.49.243.135 ::: webmail.gCVHwgoDl55QLmHw.com => 23.195.69.108
50.120.215.2 ::: webmail.gCVHwgoDl55QLmHw.com => 23.217.138.109
50.120.215.2 ::: webmail.gCVHwgoDl55QLmHw.com => 23.202.231.168
```

## Thanks

- [BitQuark - dnspop](https://github.com/bitquark/dnspop) - Great resource for the list of popular subdomains I used while testing.
- [@tomnomnom](https://github.com/tomnomnom) for the help with custom resolvers.
