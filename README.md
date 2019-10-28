# Prunes

Prunes checks a list of subdomains against a list of resolvers. Successful responses are printed to stdout. By providing a non-existent domain name, this can be used to identify resolvers that do NXDOMAIN hijacking.

## Installation 

`$ go get -u github.com/leesoh/prunes`

## Usage

```sh
$ cat subdomains.txt
www
search
mail
...

$ cat resolvers.txt
1.1.1.1
4.4.4.4
8.8.8.8
...

$ cat resolvers.txt | prunes -s subdomains.txt -d nonexistent.com -c 100
1.1.1.1: www.nonexistent.com => 192.168.1.1
4.4.4.4: search.nonexistent.com => 192.168.2.2
...
```

## Thanks

- [BitQuark - dnspop](https://github.com/bitquark/dnspop) - Great resource for the list of popular subdomains I used while testing.
- [@tomnomnom](https://github.com/tomnomnom] for the help with custom resolvers.
