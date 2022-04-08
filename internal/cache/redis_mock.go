package cache

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m Mock) VerificationKeyAndSetGet(key string, defaultValue string) (string, error) {
	args := m.Called(key, defaultValue)
	return args.String(0), args.Error(1)
}

func (m Mock) Drop(s string) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m Mock) Increment(s string) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m Mock) SetMul(key ...string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m Mock) Get(s string) (string, error) {
	args := m.Called(s)
	return args.String(0), args.Error(1)
}

func (m Mock) Set(a, b string, t time.Duration) error {
	args := m.Called(a, b, t)
	return args.Error(0)
}

func (m Mock) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m Mock) Ping() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
