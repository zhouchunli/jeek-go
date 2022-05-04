package microerr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHttpErrMessage(t *testing.T) {
	got := GetHttpErrMessage(SamePassword)
	want := HttpErrMessages[SamePassword]
	assert.Equal(t, got, want)

	got2 := GetHttpErrMessage(ErrType(111111))
	want2 := ""
	assert.Equal(t, got2, want2)
}
