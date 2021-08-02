package oraivisor_test

import (
	"bufio"
	"io"
	"testing"

	"github.com/oraichain/orai/oraivisor"

	"github.com/stretchr/testify/require"
)

func TestWaitForInfo(t *testing.T) {
	cases := map[string]struct {
		write         []string
		expectUpgrade *oraivisor.UpgradeInfo
		expectErr     bool
	}{
		"no match": {
			write: []string{"some", "random\ninfo\n"},
		},
		"match name with no info": {
			write: []string{"first line\n", `UPGRADE "myname" NEEDED at height: 123: `, "\nnext line\n"},
			expectUpgrade: &oraivisor.UpgradeInfo{
				Name: "myname",
				Info: "",
			},
		},
		"match name with info": {
			write: []string{"first line\n",
				`UPGRADE "take2" NEEDED at height: 123:   https://ipfs.io/ipfs/Qmahj5DWvXanBji73YywYuDs9dXA2Cm3dfsLMtvxL7GsJC`, "\nnext line\n"},
			expectUpgrade: &oraivisor.UpgradeInfo{
				Name: "take2",
				Info: "https://ipfs.io/ipfs/QmSymqJhmDAqa5CkfDpri8Pjt3rfL1qbGT3CC7XqXEnaBo",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r, w := io.Pipe()
			scan := bufio.NewScanner(r)

			// write all info in separate routine
			go func() {
				for _, line := range tc.write {
					n, err := w.Write([]byte(line))
					require.NoError(t, err)
					require.Equal(t, len(line), n)
				}
				w.Close()
			}()

			// now scan the info
			info, err := oraivisor.WaitForUpdate(scan)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectUpgrade, info)
		})
	}
}
