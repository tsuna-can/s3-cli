package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tsuna-can/s3-cli/internal/aws"
	"github.com/tsuna-can/s3-cli/internal/model"
)

// UIModel represents the state for the terminal UI
type UIModel struct {
	s3Client    *aws.S3Client
	state       ViewState
	bucketModel model.BucketListModel
	objectModel model.ObjectListModel
	filterInput textinput.Model
	outputDir   string
	profile     string
	endpointURL string
	err         error
	msg         string
}

// StartUI initializes and starts the terminal UI
func StartUI(outputDir string, profile string, endpointURL string, debugMode bool) {
	// デバッグログを設定
	logFile, err := os.Create("/tmp/s3-cli-debug.log")
	if err == nil {
		defer logFile.Close()
		log.SetOutput(logFile)
		log.SetPrefix("DEBUG: ")
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	log.Println("アプリケーション起動")

	// outputDirが空の場合はカレントディレクトリを使う
	if outputDir == "" {
		outputDir = "."
	}

	filterInput := textinput.New()
	filterInput.Placeholder = "Filter buckets..."
	filterInput.Prompt = "🔍 "
	filterInput.Focus()

	initialModel := UIModel{
		state:       BucketsView,
		filterInput: filterInput,
		outputDir:   outputDir,
		profile:     profile,
		endpointURL: endpointURL,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running UI: %v\n", err)
	}
}

// Init initializes the UI model
func (m UIModel) Init() tea.Cmd {
	return m.initS3Client()
}

// initS3Client initializes the S3 client using AWS configuration
func (m *UIModel) initS3Client() tea.Cmd {
	return func() tea.Msg {
		log.Println("S3クライアント初期化開始")
		client, err := aws.NewS3Client(m.profile, m.endpointURL)
		if err != nil {
			log.Printf("S3クライアント初期化エラー: %v\n", err)
			return errorMsg{err}
		}
		log.Printf("S3クライアント初期化成功。プロファイル: %s, リージョン: %s, エンドポイント: %s\n",
			client.GetProfile(), client.GetRegion(), client.GetEndpointURL())
		return s3ClientInitMsg{client}
	}
}
