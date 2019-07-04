package handlepacket

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHashKey(info *PacketInfo) string {
	h := sha256.New()
	str := string(info.SrcIP) + string(info.DstIP) + string(info.SrcPort) + string(info.DstPort)
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
