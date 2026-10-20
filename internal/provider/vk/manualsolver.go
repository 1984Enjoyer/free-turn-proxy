package vk

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/samosvalishe/free-turn-proxy/internal/provider/vk/internal/browserprofile"
	"github.com/samosvalishe/free-turn-proxy/internal/provider/vk/internal/captcha"
	manualcaptcha "github.com/samosvalishe/free-turn-proxy/internal/provider/vk/internal/captcha/manual"
)

var proxyManualMu sync.Mutex

// ProxyManualSolver показывает UI после запуска локального captcha-прокси.
func ProxyManualSolver(show func(url string), hide func()) ManualSolverFunc {
	return func(ctx context.Context, e *captcha.Error, dialer net.Dialer, profile browserprofile.Profile) (string, error) {
		if e.RedirectURI == "" {
			return "", fmt.Errorf("manual captcha: no redirect_uri")
		}

		proxyManualMu.Lock()
		defer proxyManualMu.Unlock()
		if hide != nil {
			defer hide()
		}

		return manualcaptcha.SolveViaProxyWithPresenter(ctx, e.RedirectURI, dialer, profile, show)
	}
}

// ProxyPowCapture exposes the in-app browser to the automatic PoW capture flow.
func ProxyPowCapture(show func(url string), hide func()) PowCaptureFunc {
	return func(ctx context.Context, redirectURI string, dialer net.Dialer, profile browserprofile.Profile) (string, error) {
		if redirectURI == "" {
			return "", fmt.Errorf("browser pow capture: no redirect_uri")
		}
		proxyManualMu.Lock()
		defer proxyManualMu.Unlock()
		if hide != nil {
			defer hide()
		}
		return manualcaptcha.CapturePowViaProxyWithPresenter(ctx, redirectURI, dialer, profile, show)
	}
}
