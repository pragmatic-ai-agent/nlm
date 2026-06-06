package method

import notebooklmv1alpha1 "github.com/tmc/nlm/gen/notebooklm/v1alpha1"

// artifactTypeDescriptor is the constant arg[0] for R7cb6c calls.
// Verified against HAR captures for audio, video, and slides — identical across all types.
var artifactTypeDescriptor = []interface{}{
	2, nil, nil,
	[]interface{}{1, nil, nil, nil, nil, nil, nil, nil, nil, nil, []interface{}{1}},
	[]interface{}{[]interface{}{1, 4, 2, 3, 6, 5}},
}

// encodeOverviewSourceRefs returns 3-level nesting: [[[id1]], [[id2]], ...]
// Used for the outer source refs at arg[2][3].
func encodeOverviewSourceRefs(sourceIDs []string) []interface{} {
	refs := make([]interface{}, 0, len(sourceIDs))
	for _, id := range sourceIDs {
		refs = append(refs, []interface{}{[]interface{}{id}})
	}
	return refs
}

// encodeInnerSourceRefs returns 2-level nesting: [[id1], [id2], ...]
// Used for source refs inside audioConfig/videoConfig inner arrays.
func encodeInnerSourceRefs(sourceIDs []string) []interface{} {
	refs := make([]interface{}, 0, len(sourceIDs))
	for _, id := range sourceIDs {
		refs = append(refs, []interface{}{id})
	}
	return refs
}

// EncodeCreateAudioOverviewArgs encodes the observed R7cb6c audio-overview payload.
func EncodeCreateAudioOverviewArgs(req *notebooklmv1alpha1.CreateAudioOverviewRequest) []interface{} {
	// Wire format verified against HAR capture (2026-04-14) — do not regenerate.
	sourceRefs := encodeOverviewSourceRefs(req.GetSourceIds())
	innerSourceRefs := encodeInnerSourceRefs(req.GetSourceIds())
	var instructions interface{}
	if req.GetCustomInstructions() != "" {
		instructions = req.GetCustomInstructions()
	}
	return []interface{}{
		artifactTypeDescriptor,
		req.GetProjectId(),
		[]interface{}{
			nil,
			nil,
			1, // artifact type = audio
			sourceRefs,
			nil,
			nil,
			[]interface{}{
				nil,
				[]interface{}{
					instructions,              // [0] custom instructions or nil
					2,                         // [1] constant
					nil,                       // [2]
					innerSourceRefs,           // [3] 2-level nesting
					req.GetLanguage(),         // [4] language
					nil,                       // [5] nil (not true)
					int32(req.GetAudioType()), // [6] audio style enum
				},
			},
		},
	}
}

// SlideDeckFormat selects the slide-deck layout NotebookLM generates.
//
// The R7cb6c slide config is [instructions, language, format, length]. Only the
// default capture (format=1, length=3) is HAR-verified; the format/length
// integers below are inferred from that single capture and the web UI's two
// deck choices, so the non-default values are EXPERIMENTAL until a presenter
// HAR is captured. Detailed is the verified default.
type SlideDeckFormat int32

const (
	// SlideDeckFormatDetailed is the dense, standalone-handout deck. This is
	// the HAR-verified default (format=1, length=3).
	SlideDeckFormatDetailed SlideDeckFormat = iota
	// SlideDeckFormatPresenter is the sparse, talk-along deck (fewer slides,
	// one idea per slide). EXPERIMENTAL: wire values not yet HAR-verified.
	SlideDeckFormatPresenter
)

// slideConfig returns the trailing [format, length] integers for the deck
// format. Detailed reproduces the verified default exactly; presenter is a
// best-effort guess pending a captured presenter request.
func (f SlideDeckFormat) slideConfig() (format, length int) {
	switch f {
	case SlideDeckFormatPresenter:
		return 2, 1
	default:
		return 1, 3
	}
}

// EncodeCreateSlideDeckArgs encodes the observed R7cb6c slide-deck payload.
func EncodeCreateSlideDeckArgs(projectID string, sourceIDs []string, instructions, language string, format SlideDeckFormat) []interface{} {
	// Wire format verified against HAR capture (2026-04-14) — do not regenerate.
	// The default (SlideDeckFormatDetailed) reproduces the capture exactly;
	// see SlideDeckFormat for the experimental presenter encoding.
	sourceRefs := encodeOverviewSourceRefs(sourceIDs)
	fmtCode, lenCode := format.slideConfig()
	return []interface{}{
		artifactTypeDescriptor,
		projectID,
		[]interface{}{
			nil,
			nil,
			8, // artifact type 8 = slide deck
			sourceRefs,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			[]interface{}{[]interface{}{instructions, language, fmtCode, lenCode}},
		},
	}
}

// EncodeCreateReportArgs encodes the R7cb6c report artifact payload.
// Mode 4 = ARTIFACT_TYPE_REPORT. Report config goes at field index 7
// (not 16 like slides). The generation options contain:
//   - reportType: topic name from suggestions (e.g. "Technical Specification")
//   - reportDescription: detailed description from suggestions
//   - instructions: custom user prompt (or nil)
//   - innerSourceRefs: 2-level nested source IDs
func EncodeCreateReportArgs(projectID string, sourceIDs []string, reportType, reportDescription, instructions string) []interface{} {
	sourceRefs := encodeOverviewSourceRefs(sourceIDs)
	innerSourceRefs := encodeInnerSourceRefs(sourceIDs)
	var inst interface{}
	if instructions != "" {
		inst = instructions
	}
	var rtype interface{}
	if reportType != "" {
		rtype = reportType
	}
	var rdesc interface{}
	if reportDescription != "" {
		rdesc = reportDescription
	}
	return []interface{}{
		artifactTypeDescriptor,
		projectID,
		[]interface{}{
			nil,
			nil,
			4, // artifact type = report
			sourceRefs,
			nil, nil, nil,
			[]interface{}{ // field 8: TailoredReport message
				nil,
				[]interface{}{ // GenerationOptions
					rtype,           // report type / topic
					rdesc,           // report description
					inst,            // custom instructions
					innerSourceRefs, // 2-level nested source IDs
				},
			},
		},
	}
}

// EncodeCreateVideoOverviewArgs encodes the observed R7cb6c video-overview payload.
func EncodeCreateVideoOverviewArgs(req *notebooklmv1alpha1.CreateVideoOverviewRequest) []interface{} {
	// Wire format verified against HAR capture (2026-04-14) — do not regenerate.
	sourceRefs := encodeOverviewSourceRefs(req.GetSourceIds())
	innerSourceRefs := encodeInnerSourceRefs(req.GetSourceIds())
	return []interface{}{
		artifactTypeDescriptor,
		req.GetProjectId(),
		[]interface{}{
			nil,
			nil,
			3, // artifact type = video
			sourceRefs,
			nil, nil, nil, nil,
			[]interface{}{
				nil,
				nil,
				[]interface{}{
					innerSourceRefs,            // [0] 2-level nesting
					nil,                        // [1]
					nil,                        // [2]
					nil,                        // [3]
					int32(req.GetVideoStyle()), // [4] video style enum
				},
			},
		},
	}
}

// EncodeCreateAppArtifactArgs encodes the R7cb6c AppArtifact payload.
//
// Bundle evidence (2026-05-31 LabsTailwindUi) maps appType values
// 3=prototype, 4=mindmap_app, and 5=canvas. The same AppArtifact generation
// options object carries the user prompt at field 3 ("Tp" in the compiled JS).
// This follows the HAR-verified R7cb6c envelope used by audio, video, slides,
// and reports; capture a targeted HAR before changing the nested app shape.
func EncodeCreateAppArtifactArgs(projectID string, sourceIDs []string, appType int32, instructions string) []interface{} {
	sourceRefs := encodeOverviewSourceRefs(sourceIDs)
	return []interface{}{
		artifactTypeDescriptor,
		projectID,
		[]interface{}{
			nil,
			nil,
			5, // artifact type = app
			sourceRefs,
			nil, nil, nil, nil, nil,
			[]interface{}{
				nil,
				[]interface{}{
					appType,      // [0] appType: 3 prototype, 4 mindmap_app, 5 canvas
					nil,          // [1] app-specific options
					instructions, // [2] prompt / Tp
				},
			},
		},
	}
}
