# wekan-alerter
wekan-alerter permet de synthétiser une journée d'activité et de l'envoyer par mail selon les besoins de la MRE (en v1).

## Installation
```go
go get github.com/signaux-faibles/wekan-alerter
```

Créer le fichier `wekan-alerter.toml` en repartant du fichier exemple si besoin.

## Utilisation
```
$GOPATH/bin/wekan-alerter
```

Assurez vous que le fichier template et de configuration sont bien présents dans le répertoire où vous exécutez la commande.