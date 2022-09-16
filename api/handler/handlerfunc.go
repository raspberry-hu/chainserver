package handler

func TokenAddrToName(addr string) (tokenName string) {
	if addr == "0x0000000000000000000000000000000000000000" {
		tokenName = "BNB"
	} else if addr == "0x5BF475eB0c37508D40dFb0B99d3dF828B906B6cC" {
		tokenName = "CACA"
	} else if addr == "0x337610d27c682e347c9cd60bd4b3b107c9d34ddd" {
		tokenName = "BUSD-T"
	}
	return tokenName
}
