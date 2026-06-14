// Package yesno exposes the Yes/No API as a kit Domain.
//
// A multi-domain host (ant) enables it with a single blank import:
//
//	import _ "github.com/tamnd/yesno-cli/yesno"
//
// The same Domain also builds the standalone yesno binary.
package yesno

import (
	"context"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the yesno driver.
type Domain struct{}

// Info describes the scheme, hosts, and binary identity.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "yesno",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "yesno",
			Short:  "Random yes or no answers from yesno.wtf",
			Long: `yesno fetches random yes or no answers from yesno.wtf.
No API key or authentication required.`,
			Site: Host,
			Repo: "https://github.com/tamnd/yesno-cli",
		},
	}
}

// Register installs the client factory and all operations onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	kit.Handle(app, kit.OpMeta{
		Name:    "answer",
		Group:   "read",
		Summary: "Get a random yes or no answer",
	}, answerOp)
}

// newClient builds the client from host-resolved config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- inputs ---

type answerInput struct {
	Client *Client `kit:"inject"`
}

// --- handlers ---

func answerOp(ctx context.Context, in answerInput, emit func(Answer) error) error {
	a, err := in.Client.RandomAnswer(ctx)
	if err != nil {
		return mapErr(err)
	}
	return emit(a)
}

// --- Resolver ---

// Classify turns an input into the canonical (type, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if input == "" {
		return "", "", errs.Usage("empty yesno reference")
	}
	return "answer", input, nil
}

// Locate returns the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "answer":
		return BaseURL + "/api", nil
	default:
		return "", errs.Usage("yesno has no resource type %q", uriType)
	}
}

// mapErr converts library errors into kit error kinds.
func mapErr(err error) error {
	return err
}
