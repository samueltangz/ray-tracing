package main_test

import (
	"testing"
	
	"github.com/stretchr/testify/assert"

	rt "mystiz.hk/ray-tracing"
)

func TestNorm(t *testing.T) {
	assert := assert.New(t)

	v := rt.NewVector(3, 4, 12)

	assert.Equal(v.Norm(), 13.0)
}

func TestCross(t *testing.T) {
	assert := assert.New(t)

	i := rt.NewVector(1, 0, 0)
	j := rt.NewVector(0, 1, 0)
	k := rt.NewVector(0, 0, 1)

	assert.True(rt.Equal(rt.Cross(i, j), k))
	assert.True(rt.Equal(rt.Cross(j, k), i))
	assert.True(rt.Equal(rt.Cross(k, i), j))
	assert.True(rt.Equal(rt.Cross(j, i), k.Neg()))
	assert.True(rt.Equal(rt.Cross(k, j), i.Neg()))
	assert.True(rt.Equal(rt.Cross(i, k), j.Neg()))
}