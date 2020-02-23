---
title: "Coinbase Feed"
date: 2019-11-09T15:01:25+01:00
draft: false
---

# Things to do while bitcoind is syncing: looking at a feed of coinbase messages

This is the first one (hopefully) in a series of blog posts where I will be exploring Bitcoin Core via 
command line scripts, as an effort to deepen my understanding of it.

So I'm using a VPS for hosting my Bitcoin Core full node.  I had some issues with it lately
and partly because of this, partly to do some experimenting I decided to resync the chain. Now, when you have set 
up everything and started our beloved piece of software with `bitcoind`, what is there to do? Depending on your 
system it may take hours/days/weeks before it has synced the chain.

One of the more intriguing things in Bitcoin (in my opinion) are the coinbase messages from the miners. A miner has made an enormous effort to produce a block and this
is there one chance to share something with the world, forever to be engraved in the chain. First we will look at some coinbase messages using the command
line, and then we will write a script in order to get a feed from them as our node is syncing the chain. Finally we will see a script you can use to see the 
coinbase messages from the last X blocks (useful for when your node has actually synced and you want to investigate eg. miner decentralization).

## Prerequisites
- Bitcoin Core installation running and preferably in the process of syncing the chain
- `txindex=1` option enabled in `.bitcoin/bitcoin.conf`
- [jq](https://github.com/stedolan/jq)
- the program `xxd`, to decode hex code to something readable

## The coinbase transaction from the Genesis block

Of course, the most famous coinbase transaction is the one from the genesis block. It also comes as no surprise 
that this is actually an exception and it is not retrievable using the standard methods. To get the coinbase message, we will first look at the
Genesis block

```
bitcoin-cli getblockhash 0 | xargs -I {}  bitcoin-cli getblock {} 2
{
  "hash": "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
  "confirmations": 195541,
  "strippedsize": 285,
  "size": 285,
  "weight": 1140,
  "height": 0,
  "version": 1,
  "versionHex": "00000001",
  "merkleroot": "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
  "tx": [
    {
      "txid": "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
      "hash": "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
      "version": 1,
      "size": 204,
      "vsize": 204,
      "weight": 816,
      "locktime": 0,
      "vin": [
        {
          "coinbase": "04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73",
          "sequence": 4294967295
        }
      ],
      "vout": [
        {
          "value": 50.00000000,
          "n": 0,
          "scriptPubKey": {
            "asm": "04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f OP_CHECKSIG",
            "hex": "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac",
            "reqSigs": 1,
            "type": "pubkey",
            "addresses": [
              "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
            ]
          }
        }
      ],
      "hex": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
    }
  ],
  "time": 1231006505,
  "mediantime": 1231006505,
  "nonce": 2083236893,
  "bits": "1d00ffff",
  "difficulty": 1,
  "chainwork": "0000000000000000000000000000000000000000000000000000000100010001",
  "nTx": 1,
  "nextblockhash": "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048"
}

```

In my attempt to flex my knowledge of the command line I am writing this command in one line 
using `xargs` to use the output of a command (the blockhash of the genesis block) as an argument of the following command
(the `2` at the end is for getting detailed output and is the hack required to read the coinbase message of the Genesis block).

If we want to look at the message encoded in the genesis block we can now do:

```
bitcoin-cli getblockhash 0 | xargs -I {}  bitcoin-cli getblock {} 2 | jq -r  .tx[0].hex | xxd -r -p
����M��EThe Times 03/Jan/2009 Chancellor on brink of second bailout for banks�����*CAg����UH'g��q0�\֨(�9	�yb��a޶I��?L�8��U���\8M��
                                                                                                                                                     �W�Lp+k�_�
```

And there you have it, the most famous headline of the century.

## Getting a stream of all coinbase messages
As you all know, the coinbase transaction is the first transaction in every block and has 1 input, with the tag `coinbase`.
We can use this to get all other coinbase messages via the `getrawtransaction` and `decoderawtransaction` CLI calls: 

```
bitcoin-cli getblockhash <BLOCKHEIGHT> | xargs bitcoin-cli getblock | jq -r .tx[0] | xargs bitcoin-cli getrawtransaction | xargs bitcoin-cli decoderawtransaction | jq -r .vin[0].coinbase | xxd -r -p
```

Let's say we want to have this in a bash script and also have the date of the block and some newlines in there (and using `getbestblockhash` as a shortcut for the first 2 commands):

```
#!/bin/bash
blockhash=$(bitcoin-cli getbestblockhash)
txid=$(bitcoin-cli getblock "$blockhash" | jq -r '.tx[0]')
timestamp=$(bitcoin-cli getblock "$blockhash" | jq -r '.time')
hexcb=$(bitcoin-cli getrawtransaction "$txid" | xargs bitcoin-cli decoderawtransaction | jq -r '.vin[0].coinbase')
echo "Block timestamp:"
date -d @"$timestamp"
echo "coinbase text"
echo $hexcb | xxd -r -p
echo
echo
```

Then we can look at the messages as our node is syncing:

```
while true; do ./getcurrentcoinbasetext; done
Block timestamp:
Sam Okt 13 18:10:55 EEST 2012
coinbase text
~s4/P2SH/BIP16/slush/R,��mm���\�q����� ��Q��se�}r0C\Z�C8

Block timestamp:
Sam Okt 13 18:49:03 EEST 2012
coinbase text
Py�o

Block timestamp:
Sam Okt 13 19:42:22 EEST 2012
coinbase text
xPy��+��mmrT��`�ًP�����g㴨����9��Ax�9]jEclipseMC: Aluminum Falcons 1bfgminerK
```
Behold these messages from ancient times.


If you are synced to the chain and you'd want 
to analyze the mining distribution, you can view the messages from the last `$1` blocks:

```
#!/bin/bash
depth=$1
block_count=$(bitcoin-cli getblockcount)
let counter=$block_count-$depth
while [ "$counter" -lt "$block_count" ]
do
	blockhash=$(bitcoin-cli getblockhash "$counter")
	txid=$(bitcoin-cli getblock "$blockhash"| jq -r '.tx[0]')
	timestamp=$(bitcoin-cli getblock "$blockhash" | jq -r '.time')
	hexcb=$(bitcoin-cli getrawtransaction "$txid" | xargs bitcoin-cli decoderawtransaction | jq -r '.vin[0].coinbase')
	echo "Block timestamp:"
	date -d @"$timestamp"
	echo "coinbase text"
	echo $hexcb | xxd -r -p
	echo
	echo
	let counter++
done
```

That's it for now, I hope someone found this useful. Please try it out for yourself and if you have any questions, get in touch using some of the links on the left.