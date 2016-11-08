# Kugelblitz -- Controlling Lightning

Kugelblitz is a simple UI for the c-lightning daemon `lightningd` and `bitcoind`.

If you have Docker installed then Kugelblitz is just one `docker run` away:

```bash
sudo docker run -p 19735:19735 cdecker/kugelblitz:latest
```

This will download the docker image, start `bitcoind`, `lightningd` and the kugelblitz UI bound to port 19735. Opening http://localhost:19735 should show you the interface.

Notice that the `bitcoind` instance is synchronizing with testnet, which may take a few hours. In order not to have to do that all the time (and not losing your hard earned testcoins) we recommend that you use this instead:

```bash
sudo docker run -p 19735:19735
	-v `pwd`/bitcoin:/bitcoin
	-v `pwd`/lightning:/lightning
	cdecker/kugelblitz:latest
```

This will actually persist both the state of `lightningd` and `bitcoind` in the local directory.

Once your `bitcoind` is fully synched with the testnet, you can go ahead and add funds to it.
Kugelblitz will show you an address, just copy&paste that into a [faucet](http://tpfaucet.appspot.com/).
Once you have the funds you can go ahead and create channels, and start making transfers.

Don't know what to spend your hard(ly) earned testcoins on? Here are a few ideas:

 - http://128.199.80.48/ the original cat picture server 
 - http://159.203.218.14:8000/ a cookie clicker game

Just connect to the respective node, get an invoice from the webserver and send the payment.

If you have more funny ways to spend testcoins, let me know.
