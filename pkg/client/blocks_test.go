package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

var addr, _ = proto.NewAddressFromString("3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh")

var blocksHeightBySignatureJson = `
{
  "height": 372306
}
`

func TestBlocks_HeightBySignature(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksHeightBySignatureJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.HeightBySignature(
			context.Background(),
			"2TXfMcQNHJVmkbNoznZrFRLaQHiBayFV9mzxt4VJkyXmxe9aGNn5A2unzUX4M2tqiHEfaWdfCBBo8zJQQpFrCKUY")
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.EqualValues(t, 372306, body.Height)
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/height/2TXfMcQNHJVmkbNoznZrFRLaQHiBayFV9mzxt4VJkyXmxe9aGNn5A2unzUX4M2tqiHEfaWdfCBBo8zJQQpFrCKUY", resp.Request.URL.String())
}

var blocksHeadersAtJson = `
{
  "version": 2,
  "timestamp": 1485532386874,
  "reference": "2DqubQMMBt4ot7y8F37JKWLV9J1Fvn35b4TBLGc3A9gzRvL4DweknWghxYJLYf8edDtDZujCbu1Cwqr19kC8jy12",
  "nxt-consensus": {
    "base-target": 279,
    "generation-signature": "GdXMcQzP99TJMsKX37v6BqVDcbC1xd26fgk5LRjhQUhR"
  },
  "generator": "3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh",
  "signature": "4rnYtWNE8WV4heso4q86Uwbcf1XZR5ShfszaKzyRg7aELP2Su3sFUhcCrQCyBA9SbE4T8pkd2AnLKnwBHwhUKaDq",
  "blocksize": 29882,
  "transactionCount": 100,
  "height": 370
}`

func TestBlocks_HeadersAt(t *testing.T) {
	reference, _ := crypto.NewSignatureFromBase58("2DqubQMMBt4ot7y8F37JKWLV9J1Fvn35b4TBLGc3A9gzRvL4DweknWghxYJLYf8edDtDZujCbu1Cwqr19kC8jy12")
	generator, _ := proto.NewAddressFromString("3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh")
	signature, _ := crypto.NewSignatureFromBase58("4rnYtWNE8WV4heso4q86Uwbcf1XZR5ShfszaKzyRg7aELP2Su3sFUhcCrQCyBA9SbE4T8pkd2AnLKnwBHwhUKaDq")

	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksHeadersAtJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.HeadersAt(context.Background(), 370)
	require.NoError(t, err)
	headers := &Headers{
		Version:   2,
		Timestamp: 1485532386874,
		Reference: reference,
		NxtConsensus: NxtConsensus{
			BaseTarget:          279,
			GenerationSignature: "GdXMcQzP99TJMsKX37v6BqVDcbC1xd26fgk5LRjhQUhR",
		},
		Features:         nil,
		Generator:        generator,
		Signature:        signature,
		Blocksize:        29882,
		TransactionCount: 100,
		Height:           370,
	}

	assert.NotNil(t, resp)
	assert.Equal(t, headers, body)
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/headers/at/370", resp.Request.URL.String())
}

var blocksHeadersLastJson = `
{
  "version": 3,
  "timestamp": 1542018438521,
  "reference": "5AP9TZaUXmK2M5dL2oLZVAF4wGdpB7fPoiQ6N7gPaZ6K6yc1uUpuqhJDNAGX5oYUjq8DUXh54h8vswBu5kye6Fb4",
  "nxt-consensus": {
    "base-target": 941,
    "generation-signature": "6fZ2coD4arq6jezTNHQXHLnUQqVvKwCTxQL8BzpqbYGZ"
  },
  "features": [],
  "generator": "3NBVqYXrapgJP9atQccdBPAgJPwHDKkh6A8",
  "signature": "3gTwUyQ995T1zmjsSZM6s6zCrWQh9JcFotfCetAaEMh6QYGFXPVgQsmudDtyLFsgrBfaS7GjPzAmpY2CkbLxEG5j",
  "blocksize": 225,
  "transactionCount": 0,
  "height": 372323
}
`

func TestBlocks_HeadersLast(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksHeadersLastJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.HeadersLast(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, &Headers{}, body)
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/headers/last", resp.Request.URL.String())
}

var blocksHeadersSetJson = `
[
  {
    "version": 2,
    "timestamp": 1485532386874,
    "reference": "2DqubQMMBt4ot7y8F37JKWLV9J1Fvn35b4TBLGc3A9gzRvL4DweknWghxYJLYf8edDtDZujCbu1Cwqr19kC8jy12",
    "nxt-consensus": {
      "base-target": 279,
      "generation-signature": "GdXMcQzP99TJMsKX37v6BqVDcbC1xd26fgk5LRjhQUhR"
    },
    "generator": "3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh",
    "signature": "4rnYtWNE8WV4heso4q86Uwbcf1XZR5ShfszaKzyRg7aELP2Su3sFUhcCrQCyBA9SbE4T8pkd2AnLKnwBHwhUKaDq",
    "blocksize": 29882,
    "transactionCount": 100,
    "height": 370
  }
]`

func TestBlocks_HeadersSeq(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksHeadersSetJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.HeadersSeq(context.Background(), 375500, 375500)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(body))
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/headers/seq/375500/375500", resp.Request.URL.String())
}

var blocksAtJson = `
{
  "version": 2,
  "timestamp": 1485530465594,
  "reference": "5Vwh1KEGqiBVG9ExuSKZwgwSEPbiU6CxvqL7TmtbpXd1eLQd3G4barxB161qLC3sDoVkTGwrhZEFtCBLqaRde5jt",
  "nxt-consensus": {
    "base-target": 450,
    "generation-signature": "AC94D2n1koQrY5NUtCHSfdeorxU213JNkLfJvRujmE1U"
  },
  "generator": "3My3KZgFQ3CrVHgz6vGRt8687sH4oAA1qp8",
  "signature": "2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt",
  "blocksize": 1402,
  "transactionCount": 1,
  "fee": 615366,
  "transactions": [
    {
      "type": 4,
      "id": "FYyDuMdFsJJinXcZhwdXvgnNgXKv7WnFiADxEAK2bE3j",
      "sender": "3Mv61qe6egMSjRDZiiuvJDnf3Q1qW9tTZDB",
      "senderPublicKey": "FkoFqtAeibv2E6Y86ZDRfAkZz61LwUMjLAP2gmS1j7xe",
      "fee": 189598,
      "timestamp": 1485530441535,
      "signature": "4AjgBor9GpaMd7sRg7XDMpLrTZam23XMuh7rWqTFKAzTaK3h7gPbLJQQWfWG5dM8yoZjyNDFFoLLPth4esRBz94w",
      "proofs": [
        "4AjgBor9GpaMd7sRg7XDMpLrTZam23XMuh7rWqTFKAzTaK3h7gPbLJQQWfWG5dM8yoZjyNDFFoLLPth4esRBz94w"
      ],
      "version": 1,
      "recipient": "3N5jhcA7R98AUN12ee9pB7unvnAKfzb3nen",
      "assetId": null,
      "feeAssetId": null,
      "feeAsset": null,
      "amount": 26,
      "attachment": "2escpYDq9RFWKxNYpyuAUdJ23N5wHBzybbE8zKJAREzppTvpZsDkCaSdyaJ6cmS7x2YmLTVRUwcyt43zWrWMrjS5MS3ZT8UMYHorETm8HUP5vuPkVzp5EQyukKCWSwKuw2GerfKm2qyHjBQnEXHt3Yx1ifydFLVN8xhg5qmJpe8hKBEFPnURto71hhMQCqU6"
    }],
  "height": 330
}`

func TestBlocks_At(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksAtJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.At(context.Background(), 330)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 1, len(body.Transactions))
	assert.EqualValues(t, proto.TransferTransaction, body.Transactions[0].(*proto.TransferV1).Type)
	assert.EqualValues(t, 330, body.Height)
	assert.Equal(t, "2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt", body.Signature.String())
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/at/330", resp.Request.URL.String())
}

var blocksDelayJson = `
{
  "delay": 33510
}`

func TestBlocks_Delay(t *testing.T) {
	sign, _ := crypto.NewSignatureFromBase58("2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt")
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksDelayJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Delay(context.Background(), sign, 1)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.EqualValues(t, 33510, body)
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/delay/2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt/1", resp.Request.URL.String())
}

func TestBlocks_Signature(t *testing.T) {
	sign, _ := crypto.NewSignatureFromBase58("2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt")
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksAtJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Signature(context.Background(), sign)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 1, len(body.Transactions))
	assert.EqualValues(t, proto.TransferTransaction, body.Transactions[0].(*proto.TransferV1).Type)
	assert.EqualValues(t, 330, body.Height)
	assert.Equal(t, "2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt", body.Signature.String())
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/signature/2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt", resp.Request.URL.String())
}

var blocksHeightJson = `
{
  "height": 375491
}`

func TestBlocks_Height(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksHeightJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Height(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.EqualValues(t, 375491, body.Height)
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/height", resp.Request.URL.String())
}

func TestBlocks_Child(t *testing.T) {
	sign, _ := crypto.NewSignatureFromBase58("2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt")
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksAtJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Child(context.Background(), sign)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 1, len(body.Transactions))
	assert.EqualValues(t, proto.TransferTransaction, body.Transactions[0].(*proto.TransferV1).Type)
	assert.EqualValues(t, 330, body.Height)
	assert.Equal(t, "2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt", body.Signature.String())
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/child/2WKKGrsL4kyqWPST9ZL4if198V9qYP5NMa92rv9mxGW56iqhseqaQYv15A74ThwtwZC2idj8C5px1b35oyQLzUKt", resp.Request.URL.String())
}

var blocksFirstJson = `
{
  "version": 1,
  "timestamp": 1460678400000,
  "reference": "67rpwLCuS5DGA8KGZXKsVQ7dnPb9goRLoKfgGbLfQg9WoLUgNY77E2jT11fem3coV9nAkguBACzrU1iyZM4B8roQ",
  "nxt-consensus": {
    "base-target": 153722867,
    "generation-signature": "11111111111111111111111111111111"
  },
  "generator": "3Mp6FarByk73bgv3CFnbrzMzWgLmMHAJnj2",
  "signature": "5uqnLK3Z9eiot6FyYBfwUnbyid3abicQbAZjz38GQ1Q8XigQMxTK4C1zNkqS1SVw7FqSidbZKxWAKLVoEsp4nNqa",
  "blocksize": 453,
  "transactionCount": 5,
  "fee": 0,
  "transactions": [
    {
      "type": 1,
      "id": "5G66c9GPn2egiM4bQBBF3gCkHS8sQZupRvWCpWKWGQTRRbtqdtZJ5Mt29exbHTDZW2RWygVKZ3oBNg4RwezN7wmA",
      "fee": 0,
      "timestamp": 1478000000000,
      "signature": "5G66c9GPn2egiM4bQBBF3gCkHS8sQZupRvWCpWKWGQTRRbtqdtZJ5Mt29exbHTDZW2RWygVKZ3oBNg4RwezN7wmA",
      "recipient": "3My3KZgFQ3CrVHgz6vGRt8687sH4oAA1qp8",
      "amount": 400000000000000
    },
    {
      "type": 1,
      "id": "3zpi4i5SeCoaiCBn1iuTUvCc5aahvtabqXBTrCXy1Y3ujUbJo56VVv6n4HQtcwiFapvg3BKV6stb5QkxsBrudTKZ",
      "fee": 0,
      "timestamp": 1478000000000,
      "signature": "3zpi4i5SeCoaiCBn1iuTUvCc5aahvtabqXBTrCXy1Y3ujUbJo56VVv6n4HQtcwiFapvg3BKV6stb5QkxsBrudTKZ",
      "recipient": "3NBVqYXrapgJP9atQccdBPAgJPwHDKkh6A8",
      "amount": 200000000000000
    },
    {
      "type": 1,
      "id": "3obfFPvsWXv2RyMYxjTT7owYGcpSGuSAm8fQVXeX5wErWYsgNSPPnQoFVV6nzuwm3RwGCbm8dfgvqwK9S8fVMpye",
      "fee": 0,
      "timestamp": 1478000000000,
      "signature": "3obfFPvsWXv2RyMYxjTT7owYGcpSGuSAm8fQVXeX5wErWYsgNSPPnQoFVV6nzuwm3RwGCbm8dfgvqwK9S8fVMpye",
      "recipient": "3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh",
      "amount": 200000000000000
    },
    {
      "type": 1,
      "id": "3TdE9G7V7fwED35981aGsWFM6aesxSS4W1XPfEx6p5xacwHLu7Kvf67Wzg73kgyU9gSFp1KsmPWqkFhaaR2S1fhp",
      "fee": 0,
      "timestamp": 1478000000000,
      "signature": "3TdE9G7V7fwED35981aGsWFM6aesxSS4W1XPfEx6p5xacwHLu7Kvf67Wzg73kgyU9gSFp1KsmPWqkFhaaR2S1fhp",
      "recipient": "3NCBMxgdghg4tUhEEffSXy11L6hUi6fcBpd",
      "amount": 200000000000000
    },
    {
      "type": 1,
      "id": "4hTrr7fqkujsGSH8AFN1qw7fJdfmKgwzoq3ByCCJwduHkgZPQZe1KgzG6oPBZXMuNr5ZQ6ErDSTiz2KGtxtkHpA5",
      "fee": 0,
      "timestamp": 1478000000000,
      "signature": "4hTrr7fqkujsGSH8AFN1qw7fJdfmKgwzoq3ByCCJwduHkgZPQZe1KgzG6oPBZXMuNr5ZQ6ErDSTiz2KGtxtkHpA5",
      "recipient": "3N18z4B8kyyQ96PhN5eyhCAbg4j49CgwZJx",
      "amount": 9000000000000000
    }
  ],
  "height": 1
}`

func TestBlocks_First(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksFirstJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.First(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 5, len(body.Transactions))
	assert.EqualValues(t, proto.GenesisTransaction, body.Transactions[0].(*proto.Genesis).Type)
	assert.EqualValues(t, 1, body.Height)
	assert.Equal(t, "5uqnLK3Z9eiot6FyYBfwUnbyid3abicQbAZjz38GQ1Q8XigQMxTK4C1zNkqS1SVw7FqSidbZKxWAKLVoEsp4nNqa", body.Signature.String())
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/first", resp.Request.URL.String())
}

var blocksLastJson = `
{
  "version": 3,
  "timestamp": 1542205356696,
  "reference": "z3TKjQhwhgntPm8zCwUjFzJK62k7K67rnZwgH9x8eGFajxSBrtpvFqEScUQA94vUWg6TNF4Hdt7fdAvHF1USW2X",
  "nxt-consensus": {
    "base-target": 750,
    "generation-signature": "4Lbbqe1D14ByNyq2Ej2D9BKoGMLrn7pD46HfvevqPZVY"
  },
  "features": [
    9
  ],
  "generator": "3MxTeL8dKLUGh9B1A2aaZxQ8BLL22bDdm6G",
  "signature": "3oNX2yLcKcPszzzA5CBMeNrt3p8i87AZ3eMZivkFzCut2ahGh95LZsoAQon6Qjs9XqfnTh9cTUC44o7WKWE47KzS",
  "blocksize": 227,
  "transactionCount": 0,
  "fee": 0,
  "transactions": [],
  "height": 375501
}`

func TestBlocks_Last(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com",
		Client:  NewMockHttpRequestFromString(blocksLastJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Last(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 0, len(body.Transactions))
	assert.EqualValues(t, 375501, body.Height)
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/last", resp.Request.URL.String())
}

var blocksSeqJson = `
[
  {
    "version": 1,
    "timestamp": 1460678400000,
    "reference": "67rpwLCuS5DGA8KGZXKsVQ7dnPb9goRLoKfgGbLfQg9WoLUgNY77E2jT11fem3coV9nAkguBACzrU1iyZM4B8roQ",
    "nxt-consensus": {
      "base-target": 153722867,
      "generation-signature": "11111111111111111111111111111111"
    },
    "generator": "3Mp6FarByk73bgv3CFnbrzMzWgLmMHAJnj2",
    "signature": "5uqnLK3Z9eiot6FyYBfwUnbyid3abicQbAZjz38GQ1Q8XigQMxTK4C1zNkqS1SVw7FqSidbZKxWAKLVoEsp4nNqa",
    "blocksize": 453,
    "transactionCount": 5,
    "fee": 0,
    "transactions": [
      {
        "type": 1,
        "id": "5G66c9GPn2egiM4bQBBF3gCkHS8sQZupRvWCpWKWGQTRRbtqdtZJ5Mt29exbHTDZW2RWygVKZ3oBNg4RwezN7wmA",
        "fee": 0,
        "timestamp": 1478000000000,
        "signature": "5G66c9GPn2egiM4bQBBF3gCkHS8sQZupRvWCpWKWGQTRRbtqdtZJ5Mt29exbHTDZW2RWygVKZ3oBNg4RwezN7wmA",
        "recipient": "3My3KZgFQ3CrVHgz6vGRt8687sH4oAA1qp8",
        "amount": 400000000000000
      },
      {
        "type": 1,
        "id": "3zpi4i5SeCoaiCBn1iuTUvCc5aahvtabqXBTrCXy1Y3ujUbJo56VVv6n4HQtcwiFapvg3BKV6stb5QkxsBrudTKZ",
        "fee": 0,
        "timestamp": 1478000000000,
        "signature": "3zpi4i5SeCoaiCBn1iuTUvCc5aahvtabqXBTrCXy1Y3ujUbJo56VVv6n4HQtcwiFapvg3BKV6stb5QkxsBrudTKZ",
        "recipient": "3NBVqYXrapgJP9atQccdBPAgJPwHDKkh6A8",
        "amount": 200000000000000
      },
      {
        "type": 1,
        "id": "3obfFPvsWXv2RyMYxjTT7owYGcpSGuSAm8fQVXeX5wErWYsgNSPPnQoFVV6nzuwm3RwGCbm8dfgvqwK9S8fVMpye",
        "fee": 0,
        "timestamp": 1478000000000,
        "signature": "3obfFPvsWXv2RyMYxjTT7owYGcpSGuSAm8fQVXeX5wErWYsgNSPPnQoFVV6nzuwm3RwGCbm8dfgvqwK9S8fVMpye",
        "recipient": "3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh",
        "amount": 200000000000000
      },
      {
        "type": 1,
        "id": "3TdE9G7V7fwED35981aGsWFM6aesxSS4W1XPfEx6p5xacwHLu7Kvf67Wzg73kgyU9gSFp1KsmPWqkFhaaR2S1fhp",
        "fee": 0,
        "timestamp": 1478000000000,
        "signature": "3TdE9G7V7fwED35981aGsWFM6aesxSS4W1XPfEx6p5xacwHLu7Kvf67Wzg73kgyU9gSFp1KsmPWqkFhaaR2S1fhp",
        "recipient": "3NCBMxgdghg4tUhEEffSXy11L6hUi6fcBpd",
        "amount": 200000000000000
      },
      {
        "type": 1,
        "id": "4hTrr7fqkujsGSH8AFN1qw7fJdfmKgwzoq3ByCCJwduHkgZPQZe1KgzG6oPBZXMuNr5ZQ6ErDSTiz2KGtxtkHpA5",
        "fee": 0,
        "timestamp": 1478000000000,
        "signature": "4hTrr7fqkujsGSH8AFN1qw7fJdfmKgwzoq3ByCCJwduHkgZPQZe1KgzG6oPBZXMuNr5ZQ6ErDSTiz2KGtxtkHpA5",
        "recipient": "3N18z4B8kyyQ96PhN5eyhCAbg4j49CgwZJx",
        "amount": 9000000000000000
      }
    ],
    "height": 1
  }
]`

func TestBlocks_Seq(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com/",
		Client:  NewMockHttpRequestFromString(blocksSeqJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Seq(context.Background(), 1, 1)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 5, len(body[0].Transactions))
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/seq/1/1", resp.Request.URL.String())
}

func TestNewBlocks(t *testing.T) {
	require.NotNil(t, NewBlocks(defaultOptions))
}

var blocksAddressJson = `
[
  {
    "version": 2,
    "timestamp": 1485530045905,
    "reference": "5B3tXxSPp8tmsP9QmD7TJjKTahgpcn7B4dTGpL1xh1A4rRpBMmpfAHVqYbMVCeMu8V3A1GvbGTFcpNVFjQdzjZxv",
    "nxt-consensus": {
      "base-target": 911,
      "generation-signature": "BXCMHMGpJzWPYxtt4m46DfVjoqHh3vnxmLmM66Zwb45x"
    },
    "generator": "3My3KZgFQ3CrVHgz6vGRt8687sH4oAA1qp8",
    "signature": "58yLiAeypuMr9og5WUfnWCygAo5ViL8RGjWfmht96oqxAyCkRxzmFKPa1QwvotF7t8Pkk2VHLYanKrwRiXTioVRc",
    "blocksize": 218,
    "transactionCount": 0,
    "fee": 0,
    "transactions": [],
    "height": 312
  }
]`

func TestBlocks_Address(t *testing.T) {
	client, err := NewClient(Options{
		BaseUrl: "https://testnode1.wavesnodes.com/",
		Client:  NewMockHttpRequestFromString(blocksAddressJson, 200),
	})
	require.Nil(t, err)
	body, resp, err :=
		client.Blocks.Address(context.Background(), addr, 1, 1)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	require.Equal(t, 1, len(body))
	assert.Equal(t, "https://testnode1.wavesnodes.com/blocks/address/3N5GRqzDBhjVXnCn44baHcz2GoZy5qLxtTh/1/1", resp.Request.URL.String())
}
