# wekan-alerter
wekan-alerter permet de synthétiser une journée d'activité et de l'envoyer par mail selon les besoins de la MRE (en v1).

## Installation
```go
go get github.com/signaux-faibles/wekan-alerter
```

## Utilisation
```
$GOPATH/bin/wekan-alerter --smtp smtp.host --port 25 --template template.html
```
