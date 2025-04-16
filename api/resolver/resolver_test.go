package resolver

import (
	"strings"
	"testing"
)

func TestGetRootServers(t *testing.T) {
	if len(strings.Split(ROOT_SERVERS, ",")) != len(GetRootServers()) {
		t.Error("invalid length of root servers parsed")
	}
}
