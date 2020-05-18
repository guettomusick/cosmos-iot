# cosmos-iot
A Cosmos Blockchain implementation for small IoT devices

## build project

```
make install
```

It will generate qonicod and qonicli commands to run daemon and cli.

## run project

Remove any previous running node

```
rm -rf ~/.qonicod
rm -rf ~/.qonicocli
```

Init chain and a genesis account to mint coins
```
qonicod init node1 --chain-id qonico
qonicocli config chain-id qonico
qonicocli config output json
qonicocli config indent true
qonicocli config trust-node true
qonicocli config keyring-backend test
qonicocli keys add qonicogenesis
qonicod add-genesis-account $(qonicocli keys show qonicogenesis -a) 1000000000000000000qon
```

Use qon as currency denomination, you can change if you want

```
sed -i 's/stake/qon/gi' ~/.qonicod/config/genesis.json
```

Enable validator
```
qonicod gentx --name qonicogenesis --amount 1000000qon --keyring-backend test
qonicod collect-gentxs
```

Start daemon
```
qonicod start
```

Start rest server
```
qonicocli rest-server
```

## Related Projects
https://github.com/qonico/qonico-client
https://github.com/qonico/qonico-pi
https://github.com/qonico/cosmos-client
