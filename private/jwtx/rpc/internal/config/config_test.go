package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestConfig(t *testing.T) {
	var c Config
	conf.MustLoad("../../../../../public/jwtx/example.yaml", &c)

	if _, ok := c.JWTX["nothing"]; assert.Equal(t, ok, false) {
		if _, ok := c.JWTX["admin"]; assert.Equal(t, ok, true) {
			assert.Equal(t, c.JWTX["admin"].CheckIP, true)
		}
	}

}
