# ynabler 

Go library to convert various bank/ccard CSV formats to fit YNAB's import function

## How to use it

To fetch, compile and install the ynabler binary, simply run

```
go install github.com/mschilli/go-ynabler/cmd/ynabler@latest
```

To install ynabler via Homebrew, use

```
$ brew tap mschilli/homebrew https://github.com/mschilli/homebrew
$ brew update
$ brew install ynabler
```

and then use

```
$ cat amex.csv
Date,Description,Card Member,Account #,Amount
03/02/2024,CHATGPT SUBSCRIPTIONSAN FRANCISCO       CA,BODO BALLERMANN,-12345,20.00
03/01/2024,SIRIUS COMMUNICATIONSAN FRANCISCO       CA,BODO BALLERMANN,-12345,105.00
02/29/2024,DIGITALOCEAN.COM    NEW YORK CITY       NY,BODO BALLERMANN,-12345,6.00
02/29/2024,MOBILE PAYMENT - THANK YOU,BODO BALLERMANN,-12345,-103.38

$ ynabler amex.csv
Date,Payee,Memo,Outflow,Inflow
03/04/2024,J.CREW FACTORY #30 NAPA CA,,$38.58,
03/04/2024,SAFEWAY #0667 SAN FRANCISCOCA,,$18.67,
...
```

to create a CSV format that YNAB's import function under the specific account will accept.

## Merge credit card data with Amazon order data

With `annotate`, there's an experimental utility to use a .csv file containing
Amazon orders (can be downloaded on Amazon under "Your Data") and a .csv file,
already converted to YNAB format.  Out comes a YNAB compatible file with the
"mem" fields set to item names from the order file. The match occurs by price,
so it might not be correct in all cases.

   ynabler-annotate --orders=amzn-orders.csv ynabler.csv 

## Installation

This Github repo contains the underlying code, if you want to develop with it:

`go get -u github.com/mschilli/go-ynabler`

## Author

Mike Schilli, m@perlmeister.com 2024

## License

Released under the [Apache 2.0](LICENSE)
