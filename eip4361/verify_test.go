package eip4361

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const verifyData = `
{
  "message": {
    "domain": "login.xyz",
    "address": "0x5d9de0318BeF0c3B81C46aeC31450Ffa54aa6906",
    "statement": "Sign-In With Ethereum Example Statement",
    "uri": "https://login.xyz",
    "version": "1",
    "nonce": "risxcddc",
    "issued_at": "2023-02-16T09:48:07.667Z",
    "expiration_time": "2023-02-18T09:48:07.665Z",
    "chain_id": 1
  },
  "signature": "0x2180d5ee60082d09ff990ff3c3a8f6e857f61812b5b4ef05730cf081decf4bc43ac45716c627675e0c9635410b3f76960e143fd789a8f45ffda52cad13b7ee571c"
}
`

const rawMsg = `login.xyz wants you to sign in with your Ethereum account:
0x5d9de0318BeF0c3B81C46aeC31450Ffa54aa6906

Sign-In With Ethereum Example Statement

URI: https://login.xyz
Version: 1
Chain ID: 1
Nonce: risxcddc
Issued At: 2023-02-16T09:48:07.667Z
Expiration Time: 2023-02-18T09:48:07.665Z`

func TestRecover(t *testing.T) {
	var body struct {
		Message   Message `json:"message"`
		Signature string  `json:"signature"`
	}

	if err := json.Unmarshal([]byte(verifyData), &body); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, rawMsg, body.Message.String())
	require.NoError(t, Verify(&body.Message, body.Signature))
}

func TestVerify(t *testing.T) {
	const content = "localhost:6006 wants you to sign in with your Ethereum account:\n0xf98d50cC5A1c711eA7C9C1d1fC28338B6cEa2eA9\n\nThis is statement\n\nURI: http://localhost:6006\nVersion: 1\nChain ID: 1\nNonce: FDEZrXJ1TDYQHWz5d\nIssued At: 2023-02-17T13:23:51.280Z"
	const signature = `0x6ad450bb4529948641ca1c4c42482fd37f077f4e6feb7ee5fa687f10976a32025c1e54339d0954d162ad546fcb4d0aca2ca51c1a96089147fbaaa147382261bd1b`

	msg, err := Parse(content)
	if err != nil {
		t.Fatal("parse failed", err)
	}

	assert.Equal(t, msg.String(), content)

	if err := msg.Validate(time.Now()); err != nil {
		t.Fatal("validate failed", err)
	}

	if err := Verify(msg, signature); err != nil {
		t.Fatal("verify failed", err)
	}
}
