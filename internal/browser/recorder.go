package browser

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
)

// FrameRecord holds screenshot image data and metadata for session playback.
type FrameRecord struct {
	StepIndex int    `json:"step_index"`
	DataURI   string `json:"data_uri"`
	FilePath  string `json:"file_path,omitempty"`
}

// Recorder handles frame-by-frame screenshot captures during persona navigation.
type Recorder struct {
	session    *Session
	storageDir string
}

// NewRecorder initializes a frame recorder. If storageDir is empty, frames are kept in-memory as data URIs.
func NewRecorder(session *Session, storageDir string) *Recorder {
	if storageDir != "" {
		_ = os.MkdirAll(storageDir, 0755)
	}
	return &Recorder{
		session:    session,
		storageDir: storageDir,
	}
}

// CaptureFrame takes a screenshot of the visible viewport, returning a FrameRecord.
func (r *Recorder) CaptureFrame(stepIndex int) (*FrameRecord, error) {
	ctx := r.session.Context()
	buf, err := chromedp.Run(ctx, chromedp.CaptureScreenshot())
	if err != nil {
		return nil, fmt.Errorf("capture screenshot: %w", err)
	}

	dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf)
	filePath := ""

	if r.storageDir != "" {
		fileName := fmt.Sprintf("frame_%03d.png", stepIndex)
		filePath = filepath.Join(r.storageDir, fileName)
		if writeErr := os.WriteFile(filePath, buf, 0644); writeErr != nil {
			return nil, fmt.Errorf("save frame file: %w", writeErr)
		}
	}

	return &FrameRecord{
		StepIndex: stepIndex,
		DataURI:   dataURI,
		FilePath:  filePath,
	}, nil
}
