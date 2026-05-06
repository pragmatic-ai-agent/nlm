package api

import "errors"

// Typed error sentinels for states the batchexecute error classification
// cannot disambiguate on its own. The batchexecute layer only sees RPC-level
// codes ("Failed precondition", "Invalid argument") that are polysemic —
// source-cap, artifact-still-generating, and long-poll-not-ready can all
// surface as the same batchexecute code. The callsite knows which state
// it's in, and wraps the underlying error with one of these sentinels so
// cmd/nlm's exit-code classifier can map them to distinct exit codes per
// the CLI can map them to distinct exit codes.
//
// Callers wrap via fmt.Errorf("...: %w: %w", ErrX, underlying) and consumers
// check via errors.Is(err, ErrX).
var (
	// ErrSourceCapReached indicates an AddSource* call was rejected because
	// the notebook is at the per-notebook source-count cap (NotebookLM
	// enforces ~300). The wire code 9 ("Failed precondition") carries no
	// machine-readable discriminator and is *not* by itself sufficient to
	// classify a failure as cap-reached — code-9 also appears for oversize
	// payloads, malformed envelopes, and server policy. Wrap with this
	// sentinel only when out-of-band evidence (e.g. a fresh ListSources
	// count at or near the cap) confirms the state. Maps to exit code 5
	// (permanent precondition).
	ErrSourceCapReached = errors.New("notebook source cap reached")

	// ErrSourceTooLarge indicates a single source payload exceeded the per-
	// request limit the server accepts. The observed failure band is 13MB+
	// for the text path; client-side we trip at MaxTextSourceBytes to keep
	// a deterministic error ahead of the server's ambiguous response (which
	// it otherwise mislabels as code 9 "failed precondition"). Split the
	// content or use `nlm sync`/`nlm sync-pack` which chunks automatically.
	// Maps to exit code 5 (permanent precondition).
	ErrSourceTooLarge = errors.New("source exceeds per-request size limit")

	// ErrNotebookCapReached indicates CreateProject was rejected because the
	// account is at the NotebookLM notebook cap. ZwVcOc currently reports a
	// limit of 500 notebooks, but create failures still arrive as generic RPC
	// errors, so callers should wrap this only after checking account status.
	// To attach the observed count/limit so the user-facing message can
	// surface "492/500", use NotebookCapError instead — it satisfies
	// errors.Is(err, ErrNotebookCapReached) so existing classifiers still
	// match.
	ErrNotebookCapReached = errors.New("notebook cap reached")

	// ErrArtifactGenerating indicates an artifact is still in the
	// ARTIFACT_STATUS_GENERATING transient state and a caller that wanted a
	// finished artifact should retry. Maps to exit code 7 (resource busy).
	ErrArtifactGenerating = errors.New("artifact is still generating")

	// ErrResearchPolling indicates a deep-research request is still being
	// polled via e3bVqc and the final report has not arrived. Maps to exit
	// code 7 (resource busy).
	ErrResearchPolling = errors.New("research is still in progress")
)

// NotebookCapError carries the observed account state alongside an
// ErrNotebookCapReached classification so the user-facing rewriter can
// surface the actual numbers ("492/500") instead of a generic "at the
// notebook limit" message. Count and Limit are -1 when the value was not
// available at classification time.
//
// The type is exposed (vs. an unexported struct) so cmd/nlm can extract
// the numbers via errors.As. errors.Is(err, ErrNotebookCapReached) still
// matches, so the exit-code classifier is unaffected.
//
// Count comes from ListRecentlyViewedProjects taken just after the
// CreateProject failure, so it can lag the server's true notebook count —
// most visibly right after a batch of deletes, where the message may still
// read "492/500" for a few seconds while the server's index catches up.
// The classification itself is unaffected (the wrapping happened because
// CreateProject already failed); only the displayed numbers are advisory.
type NotebookCapError struct {
	Count int
	Limit int
	Err   error
}

func (e *NotebookCapError) Error() string {
	if e.Err != nil {
		return ErrNotebookCapReached.Error() + ": " + e.Err.Error()
	}
	return ErrNotebookCapReached.Error()
}

func (e *NotebookCapError) Unwrap() error { return e.Err }

// Is matches both the sentinel (so existing errors.Is checks pass) and
// other *NotebookCapError values (so callers can introspect the type
// without losing the count/limit when re-wrapping).
func (e *NotebookCapError) Is(target error) bool {
	return target == ErrNotebookCapReached
}
