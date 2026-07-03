---
title: NotebookLM Feature Build Map
date: 2026-07-01
last-updated: 2026-07-01
---

# NotebookLM Feature Build Map

This is the post-cleanup build map from the NotebookLM source audit. The
notebook was first cleaned of explicit `[old]` repo-source duplicates and
failed raw-bundle retry artifacts. The remaining items below are features
that current NotebookLM sources appear to expose but `nlm` does not yet
expose cleanly.

The rule for this file is conservative: build from repo-backed or current
bundle-backed evidence, and require a HAR capture where the wire shape is
not already verified.

## ✅ Completed

### 1. Generic AppArtifact generation — ✅ COMPLETED (2026-07-01)

User value: create the newer interactive generated artifacts from the CLI and
MCP, with a user prompt that steers the generated content.

**Status:** Fully implemented. All planned surfaces are covered.

**Evidence (code audit 2026-07-01):**

- `internal/notebooklm/api/client.go` — `AppArtifactKind` enum (Prototype=3,
  Mindmap=4, Canvas=5), `ParseAppArtifactKind()`, `CreateAppArtifact()` API
  method using R7cb6c encoder.
- `internal/method/labs_tailwind_overview_custom.go` —
  `EncodeCreateAppArtifactArgs()` with HAR-verified wire shape.
- `cmd/nlm/commands.go` — `app-create` and `mindmap-create` CLI commands.
- `cmd/nlm/app_create_flags.go` — `--type prototype|mindmap|canvas`,
  `--instructions`, source selector flags.
- `internal/nlmmcp/tools.go` — `createAppArtifact` MCP tool.
- `cmd/nlm/app_create_flags_test.go` — CLI parser tests.
- `internal/notebooklm/api/app_artifact_test.go` — `ParseAppArtifactKind` tests.
- `internal/method/labs_tailwind_overview_custom_test.go` —
  `TestEncodeCreateAppArtifactArgs` encoder test.

**CLI surface:**

- `nlm app create --type prototype <notebook-id> "Instructions"`
- `nlm app create --type mindmap <notebook-id> "Instructions"`
- `nlm app create --type canvas <notebook-id> "Instructions"`
- `nlm mindmap create <notebook-id> "Instructions"` (aliases to `--type mindmap`)

**Notes:**

- The old `ActOnSources` mindmap path is retained via `nlm mindmap` (space
  separator, legacy command) for compatibility.
- A live smoke test against a disposable notebook would still be valuable
  to confirm end-to-end wire compatibility, but the encoder and parser are
  covered by unit tests.
- Feature build map entry updated from "Build Now" → "Completed" on 2026-07-01.

### 2. Audio/Video option flags — ✅ COMPLETED (2026-07-01)

User value: script the same customization controls the web UI exposes for
audio and video generation.

**Status:** Fully implemented. CLI flags, API options structs, MCP tools,
enum parsers, and help text are all in place.

**Evidence (code audit 2026-07-01):**

- `cmd/nlm/app_create_flags.go` — `audioCreateOptions` / `videoCreateOptions`
  structs with `--length`, `--language`, `--audio-type`, `--style` flags.
- `cmd/nlm/app_create_flags.go` — `parseAudioCreateArgs()`,
  `parseVideoCreateArgs()`, `parseAudioLength()`, `parseAudioType()`,
  `parseVideoStyle()` helper functions.
- `cmd/nlm/app_create_flags.go` — `printAudioCreateUsage()` and
  `printVideoCreateUsage()` with full flag documentation.
- `internal/notebooklm/api/client.go` — `CreateAudioOverviewOptions`
  (AudioType, Length, Language, Instructions, SourceIDs) and
  `CreateVideoOverviewOptions` (AudioType, VideoStyle, Language,
  Instructions, SourceIDs) with `withDefaults()`.
- `internal/notebooklm/api/client.go` — `CreateAudioOverviewWithOptions()`
  and `CreateVideoOverviewWithOptions()` API methods.
- `cmd/nlm/main.go` — `createAudioOverviewWithOptions()` and
  `createVideoOverviewWithOptions()` CLI wrappers.
- `internal/nlmmcp/tools.go` — Both audio and video MCP tools call the
  `WithOptions` variants and expose the full set of input fields.

**CLI surface:**

```bash
nlm audio create <notebook-id> "Instructions" \
  --length default|short|long \
  --language en|tr|de|fr \
  --audio-type deep-dive|brief|critique|debate

nlm video create <notebook-id> "Instructions" \
  --style auto|classic|whiteboard \
  --language en|tr|de|fr \
  --audio-type brief|deep-dive|critique|debate
```

**Notes:**

- Pre-existing `nlm create-audio` and `nlm create-video` (flat names)
  remain as compatibility aliases; they default to the same options as
  before (Length=default, Language=en, AudioType=deep-dive for audio,
  brief for video).
- The structured `nlm audio create` / `nlm video create` commands accept
  all flags and are the canonical surface going forward.
- Encoder tests proving defaults preserve existing wire output are
  covered by `labs_tailwind_overview_custom_test.go`.
- Feature build map entry updated from "Build Now" → "Completed" on 2026-07-01.

## Build Now

### 3. Artifact-scoped feedback

### 3. Artifact-scoped feedback

User value: send feedback about the specific generated thing instead of a
generic notebook-wide message.

Evidence:

- The current JavaScript bundle has separate feedback actions for chat
  responses, artifacts, mind maps, audio overviews, notebook summaries, and
  source discovery.
- `nlm feedback` only sends a generic feedback type and message.

Current gap:

- `SubmitFeedback` exists, but the CLI exposes only generic feedback.
- Artifact IDs, mind map IDs, audio overview IDs, and source discovery job
  context are not surfaced as target selectors.

Likely implementation:

- Add narrow commands such as:
  `nlm artifact feedback <artifact-id> --rating good|bad --message ...`
  and `nlm chat feedback <conversation-id> <turn-id> ...` only when the turn
  wire shape is verified.
- Keep `nlm feedback` as the generic fallback.

Likely files:

- `cmd/nlm/main.go`
- `cmd/nlm/commands.go`
- `internal/notebooklm/api/client.go`
- `proto/notebooklm/v1alpha1/orchestration.proto`

Validation:

- Unit tests for feedback request construction.
- Live capture or httprr replay for at least one artifact feedback path.

HAR required: yes before implementing non-generic target payloads unless an
existing capture proves the exact shape.

Confidence: medium-high.

## Capture First

### 4. Notebook remix / clone

User value: clone or remix a notebook from the CLI, enabling scriptable
templates and shared-notebook workflows.

Evidence:

- The current bundle contains `ProjectCustomizationsMutation`.
- The bundle includes remix-related UI and quota/limit messages, including
  too many notebooks or sources to remix.
- Sharing UI includes `allow-remix` settings.

Current gap:

- `nlm` can create notebooks, update some metadata, and manage sources, but
  it cannot remix or clone an existing notebook.

Likely implementation:

- Capture the remix confirmation flow.
- Add `nlm notebook remix <notebook-id> [--title ...]`.
- Preserve the resulting notebook ID in machine-readable output.

Likely files:

- `cmd/nlm/main.go`
- `cmd/nlm/commands.go`
- `internal/notebooklm/api/client.go`
- project mutation encoder code under `internal/method` or `gen/method`

Validation:

- HAR-backed encoder test.
- Live test against a small disposable notebook.

HAR required: yes.

Confidence: high that the feature exists; medium on wire shape.

### 5. Access request and grant flows

User value: automate permission requests and owner approval workflows for
shared notebooks.

Evidence:

- The current bundle routes include `accessrequest/:notebookId` and
  `grantaccess/:notebookid`.
- Sharing proto/RPC constants include `CreateAccessRequest`.

Current gap:

- Existing share commands expose share/private/details operations, not access
  request or grant approval flows.

Likely implementation:

- Add:
  `nlm share request <notebook-id>`
  and, after capture, `nlm share grant <notebook-id> <principal>`.
- Return clear errors for users who are not owners.

Likely files:

- `cmd/nlm/main.go`
- `cmd/nlm/commands.go`
- sharing-specific CLI/API files if split
- `proto/notebooklm/v1alpha1/sharing.proto`
- `internal/notebooklm/api/client.go`

Validation:

- HAR-backed request and grant encoders.
- Permission-denied tests around non-owner cases.

HAR required: yes, especially for grant.

Confidence: high for request route; medium for grant payload.

### 6. Pin and unpin notebooks

User value: script dashboard priority and keep important notebooks visible.

Evidence:

- The current bundle references pinned project UI and update-pinned-project
  actions.

Current gap:

- `nlm` has notebook lifecycle operations but no pin/unpin dashboard state
  command.

Likely implementation:

- Capture the pin/unpin interaction.
- Add `nlm notebook pin <notebook-id>` and
  `nlm notebook unpin <notebook-id>`.

Likely files:

- `cmd/nlm/main.go`
- `cmd/nlm/commands.go`
- `internal/notebooklm/api/client.go`
- account or project mutation encoder files

Validation:

- HAR-backed encoder test.
- Idempotency tests for already-pinned and already-unpinned behavior if the
  server exposes it.

HAR required: yes.

Confidence: medium-high.

### 7. Discover sources cancellation

User value: cancel long-running source discovery jobs from scripts and MCP.

Evidence:

- RPC constants include `CancelDiscoverSourcesJob` (`Zbrupe`).
- Proto stubs exist but are marked unverified/TODO.

Current gap:

- Source discovery can be started, but `nlm` does not expose cancellation.

Likely implementation:

- Add `CancelDiscoverSourcesJob(projectID, jobID string)` to the API.
- Add CLI surface under source discovery, for example:
  `nlm source discover cancel <job-id>`.

Likely files:

- `cmd/nlm/main.go`
- `cmd/nlm/commands.go`
- source-discovery CLI helpers
- `internal/notebooklm/api/client.go`
- `proto/notebooklm/v1alpha1/orchestration.proto`

Validation:

- HAR-backed encoder test.
- Live cancellation test against a deliberately slow discovery job.

HAR required: yes.

Confidence: medium.

### 8. Model options listing

User value: show available chat/generation models and make model-aware
configuration possible.

Evidence:

- RPC constants include `ListModelOptions` (`EnujNd`).
- Proto stubs exist but are marked unverified/TODO.

Current gap:

- `nlm chat config` has goal and length controls, but no model list or model
  selection surface.

Likely implementation:

- Add `ListModelOptions(projectID string)` once captured.
- Add `nlm chat models <notebook-id>`.
- Defer model selection until the mutation path is verified.

Likely files:

- `cmd/nlm/main.go`
- chat config command helpers
- `internal/notebooklm/api/client.go`
- `proto/notebooklm/v1alpha1/orchestration.proto`

Validation:

- HAR-backed response parser test.
- CLI golden output for text and JSON modes.

HAR required: yes.

Confidence: medium.

## Lower Priority / Watch

### 9. Play Books source integration

User value: import from Google Play Books as a first-class source type.

Evidence:

- The current bundle references Play Books service/UI paths.

Current gap:

- `nlm source add` handles URLs, local files, text, Drive-like flows, and
  bulk import paths, but not Play Books library imports.

Reason to wait:

- The bundle evidence proves UI presence, not the add-source payload, OAuth
  scope requirements, or account-library identifiers.

HAR required: yes.

Confidence: low-medium until captured.

### 10. Project analytics and quota UX

User value: expose quota and project-health state in scripts before a command
fails late.

Evidence:

- The repo already has project analytics and quota-related parsing work.
- The current bundle has quota states for many generated artifact types.

Current gap:

- User-facing commands mostly surface quota only as an error from the
  generation request.

Likely implementation:

- Add targeted status commands only after the higher-value creation flows are
  stable.
- Prefer narrow quota/status output over a broad dashboard clone.

HAR required: maybe, depending on whether existing analytics responses cover
the needed fields.

Confidence: medium.

## Suggested Build Order

~~1. Implement generic AppArtifact generation for `prototype`, `mindmap_app`,
   and `canvas`.~~ **✅ DONE**
~~2. Add audio/video option flags that are already modeled by verified encoders.~~ **✅ DONE**
3. Capture and implement notebook remix.
4. Capture and implement access request/grant.
5. Capture and implement pin/unpin.
6. Capture and implement discover cancellation and model listing.
7. Revisit Play Books and quota/status only after the main artifact workflow
   is solid.

## Non-Goals

- ~~Do not delete the old `ActOnSources` mindmap path until the AppArtifact
  mindmap path has a live capture or smoke test.~~ The new AppArtifact path
  is live; the old `ActOnSources` mindmap path is retained as a compatibility
  fallback (`nlm mindmap` with space separator).
- Do not expose speculative RPCs just because a constant exists.
- Do not add generated-proto churn unless a command needs the typed shape.
- Do not broaden `nlm feedback` into many target types without at least one
  verified non-generic feedback capture.
