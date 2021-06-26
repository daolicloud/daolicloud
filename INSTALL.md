LotMint区块链基于[Cothority](https://github.com/dedis/cothority)项目框架，目前正在完善和原型实现阶段。

## Build

```
cd conode
go build -ldflags="-X main.gitTag=unknown" -gcflags="-N -l" .
cd app
go build -ldflags="-X main.gitTag=unknown" -gcflags="-N -l" .
```

## Run

```
./conode setup
./conode -d 3 -c /root/.config/daolicloud/private.toml server
```

## Configuration

1. 启动BlockChain

```
./app -c /root/.config/daolicloud/public.toml genesisblock
```

2. 添加Peer Node

```
./app peer add tls://e37e638fbe8b1c7e283d0e64e931ee37f43d78170c4687dca1f1321f1ebebbb8@localhost:7772
```

3. 查看Block

```
./app -c /root/.config/daolicloud/public.toml block --index 0
./app -c /root/.config/daolicloud/public.toml block --hash 0x8d157aca501e9161689a0561e0782e13e941c0ec0e25282c88d7becd5c94b656
```


**参考:**

https://github.com/dedis/cothority

https://github.com/bitcoin/bitcoin
