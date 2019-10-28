# Prunes

Prunes checks a list of subdomains against a list of resolvers. Successful responses are printed to stdout. By providing a non-existent domain name, this can be used to identify resolvers that do NXDOMAIN hijacking.

## Installation 

`$ go get -u github.com/leesoh/prunes`

## Usage

```sh
$ cat resolvers.txt 
1.1.1.1
8.8.8.8
9.9.9.9
50.49.243.135
50.120.215.2

$ cat subs.txt 
www
mail
remote
blog
webmail
...

$ cat resolvers.txt | ./prunes -s subs.txt -d xyxyxadfssa.com -c 10
50.49.243.135 ::: www.xyxyxadfssa.com => 23.217.138.109
50.49.243.135 ::: www.xyxyxadfssa.com => 23.195.69.108
50.120.215.2 ::: www.xyxyxadfssa.com => 23.217.138.109
50.120.215.2 ::: www.xyxyxadfssa.com => 23.202.231.168
...
```

## Thanks

- [BitQuark - dnspop](https://github.com/bitquark/dnspop) - Great resource for the list of popular subdomains I used while testing.
- [@tomnomnom](https://github.com/tomnomnom) for the help with custom resolvers.
