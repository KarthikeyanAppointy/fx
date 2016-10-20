// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package core

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/fx/core/config"
	"go.uber.org/fx/core/ulog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnCriticalError_NoObserver(t *testing.T) {
	err := errors.New("Blargh")
	sh := &serviceHost{
		serviceCore: serviceCore{
			log: ulog.Logger(),
		},
	}
	closeCh, ready, err := sh.Start(false)
	require.NoError(t, err, "Expected no error starting up")
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "Server failed to start up after 1 second")
	case <-ready:
		// do nothing
	}
	go func() {
		<-closeCh
	}()
	sh.OnCriticalError(err)
	assert.Equal(t, err, sh.shutdownReason.Error)
}

func TestSupportsRole_NoRoles(t *testing.T) {
	sh := &serviceHost{}
	assert.True(t, sh.supportsRole("anything"), "Empty host roles should pass any value")
}

func TestSuupportsRole_Matches(t *testing.T) {
	sh := &serviceHost{
		roles: map[string]bool{"chilling": true},
	}
	assert.True(t, sh.supportsRole("chilling"), "Should support matching role")
}

func TestSupportsRole_NoMatch(t *testing.T) {
	sh := &serviceHost{
		roles: map[string]bool{"business": true},
	}
	assert.False(t, sh.supportsRole("pleasure"), "Should not support non-matching role")
}

func TestServiceHost_Modules(t *testing.T) {
	mods := []Module{}
	sh := &serviceHost{modules: mods}

	copied := sh.Modules()
	assert.Equal(t, len(mods), len(copied), "Should have same amount of modules")
}

func TestTransitionState(t *testing.T) {
	sh := &serviceHost{}
	observer := ObserverStub().(*StubObserver)
	require.NoError(t, WithObserver(observer)(sh))

	cases := []struct {
		name     string
		from, to ServiceState
	}{
		{
			name: "Uninitialized to Initialized",
			from: Uninitialized,
			to:   Initialized,
		},
		{
			name: "Uninitialized to Starting",
			from: Uninitialized,
			to:   Starting,
		},
		{
			name: "Initialized to Stopping",
			from: Initialized,
			to:   Stopping,
		},
		{
			name: "Running to stopped",
			from: Running,
			to:   Stopped,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sh.state = c.from
			sh.transitionState(c.to)
			assert.Equal(t, observer.state, c.to)
		})
	}
}

func TestLoadInstanceConfig_NoField(t *testing.T) {
	cfg := config.StaticProvider(nil)
	instance := struct{}{}

	assert.False(t, loadInstanceConfig(cfg, "anything", &instance), "No field defined on struct")
}

func TestLoadInstanceConfig_WithServiceConfig(t *testing.T) {
	cfg := config.NewYAMLProviderFromBytes([]byte(`
foo:
  bar: 1
`))

	instance := struct {
		ServiceConfig struct {
			Bar int `yaml:"bar"`
		}
	}{}

	assert.True(t, loadInstanceConfig(cfg, "foo", &instance))
	assert.Equal(t, 1, instance.ServiceConfig.Bar)
}

func TestServiceHostStop_NoError(t *testing.T) {
	sh := &serviceHost{}
	assert.NoError(t, sh.Stop("testing", 1))
}

func TestOnCriticalError_ObserverShutdown(t *testing.T) {
	o := observerStub()
	sh := &serviceHost{
		observer: o,
		serviceCore: serviceCore{
			log: ulog.Logger(),
		},
	}

	sh.OnCriticalError(errors.New("simulated shutdown"))
	assert.True(t, o.criticalError)
}

func TestShutdownWithError_ReturnsError(t *testing.T) {
	sh := &serviceHost{
		closeChan: make(chan ServiceExit, 1),
	}
	exitCode := 1
	shutdown, err := sh.shutdown(errors.New("simulated"), "testing", &exitCode)
	assert.True(t, shutdown)
	assert.Error(t, err)
}

func TestServiceHostStart_InShutdown(t *testing.T) {
	sh := &serviceHost{
		inShutdown: true,
	}
	_, _, err := sh.Start(false)
	assert.Error(t, err)
}

func TestServiceHostStart_AlreadyRunning(t *testing.T) {
	sh := &serviceHost{
		closeChan: make(chan ServiceExit, 1),
	}
	_, _, err := sh.Start(false)
	assert.NoError(t, err)
}

func TestStartWithObserver_InitError(t *testing.T) {
	obs := observerStub()
	obs.initError = errors.New("can't touch this")
	sh := &serviceHost{
		observer: obs,
	}
	_, _, err := sh.Start(false)
	assert.Error(t, err)
	assert.True(t, obs.init)
}

func TestAddModule_Locked(t *testing.T) {
	sh := &serviceHost{
		locked: true,
	}
	assert.Error(t, sh.addModule(nil))
}

func TestAddModule_NotLocked(t *testing.T) {
	mod := NewStubModule()
	sh := &serviceHost{}
	assert.NoError(t, sh.addModule(mod))
	assert.Equal(t, sh, mod.Host)
}