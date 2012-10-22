package key

import (
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	origKeyFile string
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	s.origKeyFile = authKey
}

func (s *S) SetUpTest(c *C) {
	changeAuthKey()
}

func (s *S) TearDownTest(c *C) {
	ok := clearAuthKeyFile()
	c.Assert(ok, Equals, true)
}

func changeAuthKey() {
	authKey = "testdata/authorized_keys"
}

func clearAuthKeyFile() bool {
	err := os.Truncate(authKey, 0)
	if err != nil {
		return false
	}
	return true
}

func (s *S) TestAuthKeysShouldBeAbsolutePathToUsersAuthorizedKeysByDefault(c *C) {
	home := os.Getenv("HOME")
	expected := path.Join(home, "authorized_keys")
	c.Assert(s.origKeyFile, Equals, expected)
}

func (s *S) TestShouldAddKeyWithoutError(c *C) {
	key := "somekey blaaaaaaa r2d2@host"
	err := Add(key)
	c.Assert(err, IsNil)
}

func (s *S) TestShouldWriteKeyInFile(c *C) {
	key := "somekey blaaaaaaa r2d2@host"
	err := Add(key)
	c.Assert(err, IsNil)
	b, err := ioutil.ReadFile(authKey)
	c.Assert(err, IsNil)
	got := string(b)
	c.Assert(got, Equals, formatKey(key))
}

func (s *S) TestShouldAppendKeyInFile(c *C) {
	key1 := "somekey blaaaaaaa r2d2@host"
	err := Add(key1)
	c.Assert(err, IsNil)
	key2 := "someotherkey fooo r2d2@host"
	err = Add(key2)
	c.Assert(err, IsNil)
	b, err := ioutil.ReadFile(authKey)
	c.Assert(err, IsNil)
	got := string(b)
	expected := fmt.Sprintf(".*%s\n.*%s", key1, key2)
	c.Assert(got, Matches, expected)
}

func (s *S) TestAddShouldWrapKeyWithRestrictions(c *C) {
	key := "somekey bleeeerh r2d2@host"
	expected := fmt.Sprintf("no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty %s", key)
	err := Add(key)
	c.Assert(err, IsNil)
	b, err := ioutil.ReadFile(authKey)
	c.Assert(err, IsNil)
	got := string(b)
	c.Assert(got, Equals, expected)
}

func (s *S) TestRemoveKey(c *C) {
	key1 := "somekey blaaaaaaa r2d2@host"
	err := Add(key1)
	c.Assert(err, IsNil)
	key2 := "someotherkey fooo r2d2@host"
	err = Add(key2)
	c.Assert(err, IsNil)
	err = Remove(key1)
	b, err := ioutil.ReadFile(authKey)
	c.Assert(err, IsNil)
	got := string(b)
	expected := fmt.Sprintf(".*%s", key2)
	c.Assert(got, Matches, expected)
}

func (s *S) TestRemoveWhenKeyDoesNotExists(c *C) {
	key1 := "somekey blaaaaaaa r2d2@host"
	err := Remove(key1)
	c.Assert(err, IsNil)
	b, err := ioutil.ReadFile(authKey)
	c.Assert(err, IsNil)
	got := string(b)
	c.Assert(got, Equals, "")
}

func (s *S) TestRemoveWhenExistsOnlyOneKey(c *C) {
	key1 := "somekey blaaaaaaa r2d2@host"
	err := Add(key1)
	c.Assert(err, IsNil)
	err = Remove(key1)
	c.Assert(err, IsNil)
	b, err := ioutil.ReadFile(authKey)
	c.Assert(err, IsNil)
	got := string(b)
	c.Assert(got, Equals, "")
}

func (s *S) TestFormatKeyShouldAddSshLoginRestrictionsAtBegining(c *C) {
	key := "somekeeey fooo bar@bar.com"
	got := formatKey(key)
	expected := fmt.Sprintf("no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty %s", key)
	c.Assert(got, Equals, expected)
}
