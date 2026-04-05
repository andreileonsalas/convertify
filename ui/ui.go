package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"convertify/ffmpeg"
	"convertify/formats"
	"convertify/notify"
)

// fileEntry represents a file queued for conversion
type fileEntry struct {
	path  string
	group formats.Group
}

// fileListModel is a walk-compatible model for the file list
type fileListModel struct {
	walk.ListModelBase
	mu    sync.Mutex
	items []*fileEntry
}

func (m *fileListModel) ItemCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.items)
}

func (m *fileListModel) Value(index int) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if index < 0 || index >= len(m.items) {
		return ""
	}
	e := m.items[index]
	return fmt.Sprintf("  %s  [%s]", filepath.Base(e.path), e.group.String())
}

func (m *fileListModel) add(paths []string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	added := false
	for _, p := range paths {
		g := formats.DetectGroup(p)
		if g == formats.Unknown {
			continue
		}
		// Avoid duplicates
		dup := false
		for _, existing := range m.items {
			if existing.path == p {
				dup = true
				break
			}
		}
		if !dup {
			m.items = append(m.items, &fileEntry{path: p, group: g})
			added = true
		}
	}
	return added
}

func (m *fileListModel) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = nil
}

func (m *fileListModel) snapshot() []*fileEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]*fileEntry, len(m.items))
	copy(cp, m.items)
	return cp
}

func (m *fileListModel) dominantGroup() formats.Group {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.items) == 0 {
		return formats.Unknown
	}
	// Use the group of the first file
	return m.items[0].group
}

// buildFilterString creates a Windows file dialog filter for all supported types
func buildFilterString() string {
	all := "*.jpg;*.jpeg;*.png;*.webp;*.gif;*.bmp;*.tiff;*.tif;*.avif;" +
		"*.mp4;*.avi;*.mkv;*.mov;*.webm;*.flv;*.wmv;*.m4v;*.mpg;*.mpeg;" +
		"*.mp3;*.wav;*.flac;*.aac;*.ogg;*.m4a;*.wma;*.opus"
	return "Archivos multimedia|" + all + "|Todos los archivos|*.*"
}

// Run launches the main application window
func Run() {
	model := &fileListModel{}

	var mw *walk.MainWindow
	var lb *walk.ListBox
	var formatCombo *walk.ComboBox
	var convertBtn *walk.PushButton
	var clearBtn *walk.PushButton
	var statusLabel *walk.Label

	var currentFormats []formats.OutputFormat

	refreshFormatCombo := func() {
		g := model.dominantGroup()
		currentFormats = formats.OutputsFor(g)
		if formatCombo == nil {
			return
		}
		labels := make([]string, len(currentFormats))
		for i, f := range currentFormats {
			labels[i] = f.Label
		}
		formatCombo.SetModel(labels)
		if len(labels) > 0 {
			formatCombo.SetCurrentIndex(0)
		}
	}

	addFiles := func(paths []string) {
		if model.add(paths) {
			model.PublishItemsReset()
			refreshFormatCombo()
		}
	}

	// Build the window
	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Convertify",
		MinSize:  Size{Width: 500, Height: 420},
		Size:     Size{Width: 540, Height: 460},
		Layout:   VBox{Margins: Margins{Left: 14, Top: 14, Right: 14, Bottom: 14}, Spacing: 10},
		OnDropFiles: func(files []string) {
			addFiles(files)
		},
		Children: []Widget{

			// Header
			Label{
				Text:      "Convertify  —  Conversor de archivos",
				TextColor: walk.RGB(60, 60, 60),
			},

			// File list
			ListBox{
				AssignTo: &lb,
				Model:    model,
				MinSize:  Size{Height: 200},
			},

			// Add / Clear buttons
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 6},
				Children: []Widget{
					PushButton{
						Text: "+ Agregar archivos",
						OnClicked: func() {
							dlg := new(walk.FileDialog)
							dlg.Title = "Seleccionar archivos multimedia"
							dlg.Filter = buildFilterString()
							if ok, err := dlg.ShowOpenMultiple(mw); err == nil && ok {
								addFiles(dlg.FilePaths)
							}
						},
					},
					PushButton{
						AssignTo: &clearBtn,
						Text:     "Limpiar",
						OnClicked: func() {
							model.clear()
							model.PublishItemsReset()
							statusLabel.SetText("")
							currentFormats = nil
							formatCombo.SetModel([]string{})
						},
					},
				},
			},

			// Separator hint
			Label{
				Text:      "Arrastra archivos directamente sobre esta ventana",
				TextColor: walk.RGB(140, 140, 140),
			},

			HSpacer{},

			// Format selector
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 8},
				Children: []Widget{
					Label{Text: "Convertir a:"},
					ComboBox{
						AssignTo: &formatCombo,
						Model:    []string{},
						MinSize:  Size{Width: 220},
					},
				},
			},

			// Convert button
			PushButton{
				AssignTo: &convertBtn,
				Text:     "Convertir",
				MinSize:  Size{Height: 32},
				OnClicked: func() {
					if model.ItemCount() == 0 {
						walk.MsgBox(mw, "Sin archivos", "Agrega al menos un archivo para convertir.", walk.MsgBoxIconWarning)
						return
					}
					idx := formatCombo.CurrentIndex()
					if idx < 0 || idx >= len(currentFormats) {
						walk.MsgBox(mw, "Sin formato", "Selecciona un formato de salida.", walk.MsgBoxIconWarning)
						return
					}
					targetFmt := currentFormats[idx]
					entries := model.snapshot()

					convertBtn.SetEnabled(false)
					clearBtn.SetEnabled(false)
					statusLabel.SetText(fmt.Sprintf("Convirtiendo %d archivo(s)...", len(entries)))

					go func() {
						var errList []string
						converted := 0
						for _, e := range entries {
							_, err := ffmpeg.Convert(e.path, targetFmt)
							if err != nil {
								errList = append(errList, filepath.Base(e.path)+": "+err.Error())
							} else {
								converted++
							}
						}

						mw.Synchronize(func() {
							convertBtn.SetEnabled(true)
							clearBtn.SetEnabled(true)
							if len(errList) > 0 {
								msg := strings.Join(errList, "\n")
								statusLabel.SetText(fmt.Sprintf("Errores en %d archivo(s)", len(errList)))
								walk.MsgBox(mw, "Errores durante la conversion", msg, walk.MsgBoxIconWarning)
							} else {
								statusLabel.SetText(fmt.Sprintf("✓  %d archivo(s) convertidos correctamente", converted))
								notify.Show("Convertify", fmt.Sprintf("%d archivo(s) → %s  listos", converted, strings.ToUpper(targetFmt.Ext)))
							}
						})
					}()
				},
			},

			// Status
			Label{
				AssignTo:  &statusLabel,
				Text:      "",
				TextColor: walk.RGB(0, 120, 60),
			},
		},
	}).Create(); err != nil {
		panic(err)
	}

	mw.Run()
}
