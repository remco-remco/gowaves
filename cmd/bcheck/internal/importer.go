package internal

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/mr-tron/base58/base58"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"io"
	"os"
	"time"
)

type Importer struct {
	storage         *Storage
	transactionType proto.TransactionType
}

func NewImporter(storage *Storage, transactionType int) *Importer {
	return &Importer{storage: storage, transactionType: proto.TransactionType(transactionType)}
}

func (im *Importer) Import(n string) {
	start := time.Now()

	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("Import took %s\n", elapsed)
	}()

	f, err := os.Open(n)
	if err != nil {
		fmt.Printf("Failed to open blockchain file '%s': %s\n", n, err.Error())
		return
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("Failed to close blockchain file: %s\n\n", err.Error())
		}
	}()

	st, err := f.Stat()
	if err != nil {
		fmt.Printf("Failed to get file info: %s\n\n", err.Error())
		return
	}
	fmt.Printf("Importing blockchain file '%s' of size %d bytes\n", n, st.Size())

	sb := make([]byte, 4)
	buf := make([]byte, 2*1024*1024)
	r := bufio.NewReader(f)
	h := 2
	sh, err := im.storage.GetHeight()
	if err != nil {
		fmt.Printf("Failed to read state height: %s\n", err.Error())
	}
	for {
		_, err := io.ReadFull(r, sb)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Unable to read data size: %s\n\n", err.Error())
				return
			}
			fmt.Printf("EOF received while reading size\n")
			return
		}
		s := binary.BigEndian.Uint32(sb)
		bb := buf[:s]
		_, err = io.ReadFull(r, bb)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Unable to read block: %s\n\n", err.Error())
				return
			}
			fmt.Printf("EOF received while reading block\n")
			return
		}
		if h > sh {
			var block proto.Block
			err = block.UnmarshalBinary(bb)
			if err != nil {
				fmt.Printf("Failed to unmarshal block: %s\n\n", err.Error())
				return
			}
			if !crypto.Verify(block.GenPublicKey, block.BlockSignature, bb[:len(bb)-crypto.SignatureSize]) {
				fmt.Printf("Block %s has invalid signature\n", block.BlockSignature.String())
				return
			}
			d := block.Transactions
			for i := 0; i < int(block.TransactionCount); i++ {
				s := int(binary.BigEndian.Uint32(d[0:4]))
				txb := d[4 : s+4]
				switch txb[0] {
				case 0:
					switch txb[1] {
					case byte(proto.IssueTransaction):
						var tx proto.IssueV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract IssueV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.TransferTransaction):
						var tx proto.TransferV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract TransferV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.ReissueTransaction):
						var tx proto.ReissueV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract ReissueV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.BurnTransaction):
						var tx proto.BurnV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract BurnV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.ExchangeTransaction):
						var tx proto.ExchangeV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract ExchangeV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.LeaseTransaction):
						var tx proto.LeaseV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract LeaseV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.LeaseCancelTransaction):
						var tx proto.LeaseCancelV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract LeaseCancelV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.CreateAliasTransaction):
						var tx proto.CreateAliasV2
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract CreateAliasV2 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.DataTransaction):
						var tx proto.DataV1
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract DataV1 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.SetScriptTransaction):
						var tx proto.SetScriptV1
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract SetScriptV1 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					case byte(proto.SponsorshipTransaction):
						var tx proto.SponsorshipV1
						err := tx.UnmarshalBinary(txb)
						if err != nil {
							fmt.Printf("%d: Failed to extract SponsorshipV1 transaction: %s\n", h, err.Error())
						}
						im.checkAndPut(h, tx.ID[:], tx.Type)
					default:
						fmt.Printf("ALARM 2 AT %d\n", h)
					}
				case byte(proto.GenesisTransaction):
					var tx proto.Genesis
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract Genesis transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.PaymentTransaction):
					var tx proto.Payment
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract Payment transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.IssueTransaction):
					var tx proto.IssueV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract IssueV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.TransferTransaction):
					var tx proto.TransferV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract TransferV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.ReissueTransaction):
					var tx proto.ReissueV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract ReissueV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.BurnTransaction):
					var tx proto.BurnV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract BurnV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.ExchangeTransaction):
					var tx proto.ExchangeV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract ExchangeV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.LeaseTransaction):
					var tx proto.LeaseV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract LeaseV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.LeaseCancelTransaction):
					var tx proto.LeaseCancelV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract LeaseCancelV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.CreateAliasTransaction):
					var tx proto.CreateAliasV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract CreateAliasV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				case byte(proto.MassTransferTransaction):
					var tx proto.MassTransferV1
					err := tx.UnmarshalBinary(txb)
					if err != nil {
						fmt.Printf("%d: Failed to extract MassTransferV1 transaction: %s\n", h, err.Error())
					}
					im.checkAndPut(h, tx.ID[:], tx.Type)
				default:
					fmt.Printf("ALARM 1 AT %d\n", h)
				}

				err := im.storage.PutHeight(h)
				if err != nil {
					fmt.Printf("Unable to update state heigh: %s\n", err.Error())
					return
				}

				d = d[4+s:]
			}
		}
		if h%100000 == 0 {
			elapsed := time.Since(start)
			fmt.Printf("%d: %s\n", h, elapsed)
		}
		h++
	}
}

func (im *Importer) checkAndPut(h int, id []byte, tt proto.TransactionType) {
	has, err := im.storage.HasID(id)
	if err != nil {
		fmt.Printf("Failed to check ID: %s\n", err.Error())
	}
	if has {
		fh, err := im.storage.GetID(id)
		if err != nil {
			fmt.Printf("Failed to check ID: %s\n", err.Error())
		}
		s := base58.Encode(id)
		fmt.Printf("%d: ALARM %s ALREADY EXISTS AT %d\n", h, s, fh)
	} else {
		err := im.storage.PutID(id, h)
		if err != nil {
			fmt.Printf("Failed to put ID: %s\n", err.Error())
		}
	}
	if im.transactionType != 0 && im.transactionType == tt {
		fmt.Printf("Transaction of type %d with ID %s at height %d\n", tt, base58.Encode(id), h)
	}
}
