---
title: Remaining Gaps and Audit Notes
date: 2026-05-16
---

# nlm Remaining Gaps

This file was re-audited against the current tree on 2026-05-16. The local
release gate (`go test ./...`) passes. The items below are remaining
deferred, low-priority, or HAR-dependent gaps after checking the live CLI,
API, and MCP code paths; they are not current local test blockers.

### 1. Generated analytics proto remains scalar

`AUrzMb` returns time-series metrics, but the generated `ProjectAnalytics`
shape still expects scalar counts. The public API and CLI now bypass that
generated response model and parse the fixture-backed time-series shape into
typed metric rows.

Status: low-priority generated-proto cleanup, not HAR-blocked.

Next step: decide whether the generated proto should grow metric-series
messages or whether AUrzMb should remain a typed API-only path.

### 2. Video download is still manual-fallback only

`video download` still relies on the direct-RPC-only path in
`internal/notebooklm/api/client.go`. When the API response does not expose
direct media bytes or a CDN URL, the command now fails explicitly with a
manual-browser fallback instead of probing speculative RPC shapes.

Status: open. This is not an MCP gap. HAR would help automation, but the
current user-visible limitation is the lack of a verified automated download
path plus CDN browser-auth requirements. The CLI fallback prints the
NotebookLM notebook URL so the user can finish the download in a browser.

### 3. Weakly verified encoder paths still exist

Dead-path RPCs such as `xpWGLf`, `lCjAd`, and `BeTrYd` still exist in
codegen output / compatibility paths even though the current CLI does
not route through them. A few low-use argbuilder encoders are also still
weakly verified:

- `SubmitFeedback` (`uNyJKe`) works in practice but still uses the generic
  `[%project_id%, %feedback_type%, %feedback_text%]` shape.
- `DeleteNotes` (`AH0mwd`) works in practice but is not pinned by a
  HAR-backed encoder test.
- `GenerateNotebookGuide` (`VfAZjd`) has a hand-written encoder that emits
  `[%project_id%, %guide_type%]`, but both outline and mind-map variants
  still need HAR capture before the guard comment can call the shape verified.
- `GenerateFreeFormStreamed` exists as a gRPC-Web chat path in
  `internal/notebooklm/api`; the batchexecute method encoder is not a live
  CLI path.

Status: low priority cleanup / verification work.

### 4. `artifact get` direct path remains unverified

`nlm artifact get <artifact-id>` now tries the JS-bundle-canonical
`v9rmvd` direct RPC first, then falls back to scanning
`ListRecentlyViewedProjects` plus `ListArtifacts` when the direct path
fails or returns an unparsable response. That keeps the user-visible command
on the same list-scan behavior the CLI used before the direct probe.

Status: open for direct-path verification only. A fresh HAR would confirm
whether `v9rmvd` is callable on the live service or should remain a
best-effort fast path behind the list-scan fallback.

### 5. `chat config` server semantics unverified

`nlm chat config <id> <setting> [value]` rides on `MutateProject`
(`s0tc2d`) to apply chat goal/length settings. The CLI accepts
`goal default`, `goal custom "prompt"`, and `length default|longer|shorter`,
but none of these paths have been verified end-to-end against the live
service, and the `ChatGoal` enum values may not match server expectations.

Status: open. Low usage; verify when there is a real caller.

### 6. Auth-expiry mid-session auto-refresh

When the session cookies expire during a long-running command, the CLI now
detects auth-shaped API and transport failures and surfaces a targeted
"run `nlm auth`" message.

Status: low-priority polish. The remaining possible improvement is automatic
refresh after a 401/Unauthenticated response.

## Truly HAR-Blocked

### 1. `izAoDd` drag-drop bulk add shape

The generic bulk-add RPC still has no verified drag-drop capture. This no
longer blocks the main programmatic bulk-import use case, because
`nlm research --import` already uses the HAR-verified `LBwxtb`
bulk-import variant instead.

Status: HAR-blocked, low value until there is a real CLI caller.

### 2. Deep-research session state `6`

The active parser in `api.Client.pollResearch` safely treats unknown
states as still-running, but the semantics of observed state `6`
remain unknown.

Status: HAR-blocked for semantics only. The current fallback is safe.

## Next Work

1. Decide whether AUrzMb should stay typed API-only or get generated proto
   metric-series messages.
2. Decide whether `video download` should keep the current manual-fallback
   UX or get a real CDN capture and a browser-assisted path.
3. Remove dead generated RPC stubs so future audits do not
   mistake them for live command paths.
4. Re-capture `artifact get` against the live service and either verify
   `v9rmvd` or keep the list-scan fallback as the canonical path.
5. Verify `chat config` end-to-end (or hide it until there is a real caller).
6. Capture `izAoDd` only if a real bulk-add CLI caller is introduced.
