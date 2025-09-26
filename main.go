package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type model struct {
	files         []os.DirEntry
	currentDir    string
	cursor        int
	input         string
	err           error
	height        int
	offset        int
	message       string
	navMode       bool
	searchInput   string
	lastDir       string
	sortBy        string
	filteredFiles []os.DirEntry
}

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8")).
			Background(lipgloss.Color("#626262")).
			PaddingLeft(2)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Background(lipgloss.Color("#3C3C3C")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Background(lipgloss.Color("#3C3C3C")).
			Padding(0, 1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Background(lipgloss.Color("#3C3C3C")).
			Padding(0, 1)

	navModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1)

	searchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1)
)

// Initial model
func initialModel() model {
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "/"
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		files = []os.DirEntry{}
	}

	return model{
		files:         files,
		currentDir:    currentDir,
		cursor:        0,
		input:         "",
		err:           err,
		height:        20, // Default height, will be updated on resize
		offset:        0,
		message:       "",
		navMode:       false,
		searchInput:   "",
		lastDir:       "",
		sortBy:        "name",
		filteredFiles: files,
	}
}

// Commands
type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type dirChangedMsg struct {
	dir   string
	files []os.DirEntry
	err   error
}

func changeDir(dir string) tea.Cmd {
	return func() tea.Msg {
		files, err := os.ReadDir(dir)
		return dirChangedMsg{dir: dir, files: files, err: err}
	}
}

// Execute system command
func executeCommand(command string, args []string, workingDir string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir
	
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Sort files by different criteria
func sortFiles(files []os.DirEntry, sortBy string) []os.DirEntry {
	switch sortBy {
	case "name":
		// Sort by name (directories first, then files)
		for i := 0; i < len(files)-1; i++ {
			for j := i + 1; j < len(files); j++ {
				if files[i].IsDir() && !files[j].IsDir() {
					continue
				} else if !files[i].IsDir() && files[j].IsDir() {
					files[i], files[j] = files[j], files[i]
				} else if files[i].Name() > files[j].Name() {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
	case "size":
		// Sort by size (largest first)
		for i := 0; i < len(files)-1; i++ {
			for j := i + 1; j < len(files); j++ {
				infoI, _ := files[i].Info()
				infoJ, _ := files[j].Info()
				if infoI.Size() < infoJ.Size() {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
	case "time":
		// Sort by modification time (newest first)
		for i := 0; i < len(files)-1; i++ {
			for j := i + 1; j < len(files); j++ {
				infoI, _ := files[i].Info()
				infoJ, _ := files[j].Info()
				if infoI.ModTime().Before(infoJ.ModTime()) {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
	case "type":
		// Sort by type (directories first, then files alphabetically)
		for i := 0; i < len(files)-1; i++ {
			for j := i + 1; j < len(files); j++ {
				if files[i].IsDir() && !files[j].IsDir() {
					continue
				} else if !files[i].IsDir() && files[j].IsDir() {
					files[i], files[j] = files[j], files[i]
				} else if files[i].Name() > files[j].Name() {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
	}
	return files
}

// Filter files based on search input
func filterFiles(files []os.DirEntry, searchInput string) []os.DirEntry {
	if searchInput == "" {
		return files
	}
	
	var filtered []os.DirEntry
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(searchInput)) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// Find cursor position for a specific directory name
func findCursorPosition(files []os.DirEntry, targetDir string) int {
	for i, file := range files {
		if file.Name() == targetDir {
			return i
		}
	}
	return 0
}

// Helper function to update scroll position
func (m *model) updateScroll() {
	if len(m.files) == 0 {
		m.offset = 0
		return
	}

	// Calculate how many items can fit in the visible area
	// Reserve space for title, directory, and input (about 6 lines)
	visibleItems := m.height - 6
	if visibleItems < 1 {
		visibleItems = 1
	}

	// If cursor is above visible area, scroll up
	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	// If cursor is below visible area, scroll down
	if m.cursor >= m.offset+visibleItems {
		m.offset = m.cursor - visibleItems + 1
	}

	// Ensure offset doesn't go negative
	if m.offset < 0 {
		m.offset = 0
	}

	// Ensure offset doesn't exceed file count
	if m.offset >= len(m.files) {
		m.offset = len(m.files) - 1
		if m.offset < 0 {
			m.offset = 0
		}
	}
}

// Init function
func (m model) Init() tea.Cmd {
	return nil
}

// Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.updateScroll()
		return m, nil

	case tea.KeyMsg:
		if m.navMode {
			// Navigation mode key handling
			switch msg.String() {
			case "esc":
				m.navMode = false
				m.searchInput = ""
				return m, nil
			case "backspace":
				if len(m.searchInput) > 0 {
					m.searchInput = m.searchInput[:len(m.searchInput)-1]
					m.filteredFiles = filterFiles(m.files, m.searchInput)
					m.cursor = 0
					m.updateScroll()
				}
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m.updateScroll()
				}
			case "down", "j":
				if m.cursor < len(m.filteredFiles)-1 {
					m.cursor++
					m.updateScroll()
				}
			case "enter":
				if len(m.filteredFiles) > 0 && m.cursor < len(m.filteredFiles) {
					selected := m.filteredFiles[m.cursor]
					if selected.IsDir() {
						newDir := filepath.Join(m.currentDir, selected.Name())
						m.lastDir = filepath.Base(m.currentDir)
						m.searchInput = "" // Clear search box
						m.filteredFiles = m.files // Reset filtered files
						return m, changeDir(newDir)
					}
				}
			case "1":
				m.sortBy = "name"
				m.files = sortFiles(m.files, m.sortBy)
				m.filteredFiles = filterFiles(m.files, m.searchInput)
				m.cursor = 0
				m.updateScroll()
			case "2":
				m.sortBy = "size"
				m.files = sortFiles(m.files, m.sortBy)
				m.filteredFiles = filterFiles(m.files, m.searchInput)
				m.cursor = 0
				m.updateScroll()
			case "3":
				m.sortBy = "time"
				m.files = sortFiles(m.files, m.sortBy)
				m.filteredFiles = filterFiles(m.files, m.searchInput)
				m.cursor = 0
				m.updateScroll()
			case "4":
				m.sortBy = "type"
				m.files = sortFiles(m.files, m.sortBy)
				m.filteredFiles = filterFiles(m.files, m.searchInput)
				m.cursor = 0
				m.updateScroll()
			default:
				if len(msg.String()) == 1 {
					m.searchInput += msg.String()
					m.filteredFiles = filterFiles(m.files, m.searchInput)
					m.cursor = 0
					m.updateScroll()
				}
			}
		} else {
			// Normal mode key handling
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.navMode = true
				m.searchInput = ""
				m.filteredFiles = m.files
				return m, nil
			case "up":
				if m.cursor > 0 {
					m.cursor--
					m.updateScroll()
				}
			case "down":
				if m.cursor < len(m.files)-1 {
					m.cursor++
					m.updateScroll()
				}
			case "left":
				// Go to parent directory
				parentDir := filepath.Dir(m.currentDir)
				m.lastDir = filepath.Base(m.currentDir)
				return m, changeDir(parentDir)
			case "right":
				// Enter selected directory
				if len(m.files) > 0 && m.cursor < len(m.files) {
					selected := m.files[m.cursor]
					if selected.IsDir() {
						newDir := filepath.Join(m.currentDir, selected.Name())
						m.lastDir = filepath.Base(m.currentDir)
						return m, changeDir(newDir)
					}
				}
			case "enter":
				// Execute command
				cmd := strings.TrimSpace(m.input)
				m.input = "" // Clear input immediately
				
				if cmd == "" {
					return m, nil
				}
				
				// Handle commands
				parts := strings.Fields(cmd)
				if len(parts) == 0 {
					return m, nil
				}

				command := parts[0]
				args := parts[1:]

				// Special handling for built-in commands
				switch command {
				case "cd":
					newDir := ""
					if len(args) > 0 {
						newDir = args[0]
					} else {
						newDir = os.Getenv("HOME")
					}
					if !filepath.IsAbs(newDir) {
						newDir = filepath.Join(m.currentDir, newDir)
					}
					return m, changeDir(newDir)

				case "ls", "dir":
					return m, changeDir(m.currentDir)

				case "pwd":
					m.message = "Current directory: " + m.currentDir
					return m, nil

				case "clear":
					m.message = ""
					return m, nil

				case "quit", "exit":
					return m, tea.Quit

				default:
					// Execute system command
					output, err := executeCommand(command, args, m.currentDir)
					if err != nil {
						m.message = fmt.Sprintf("Error executing %s: %v\n%s", command, err, output)
					} else {
						if strings.TrimSpace(output) == "" {
							m.message = fmt.Sprintf("Command '%s' executed successfully", command)
						} else {
							// Truncate long output
							if len(output) > 500 {
								output = output[:500] + "\n... (output truncated)"
							}
							m.message = output
						}
						
						// Refresh directory listing for file operations
						if command == "touch" || command == "mkdir" || command == "rm" || 
						   command == "rmdir" || command == "cp" || command == "mv" {
							return m, changeDir(m.currentDir)
						}
					}
					return m, nil
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
				}
			}
		}

	case dirChangedMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.currentDir = msg.dir
			m.files = msg.files
			m.filteredFiles = m.files
			
			// Position cursor at the directory we came from
			if m.lastDir != "" {
				m.cursor = findCursorPosition(m.files, m.lastDir)
				m.lastDir = "" // Reset after use
			} else {
				m.cursor = 0
			}
			
			m.offset = 0
			m.err = nil
			m.updateScroll()
		}
	}

	return m, nil
}

// View function
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to continue...", m.err)
	}

	var s strings.Builder

	// Title
	if m.navMode {
		s.WriteString(navModeStyle.Render("Navigation Mode"))
	} else {
		s.WriteString(titleStyle.Render("File Manager"))
	}
	s.WriteString("\n\n")

	// Current directory
	s.WriteString(fmt.Sprintf("Directory: %s\n\n", m.currentDir))

	// Navigation mode info
	if m.navMode {
		s.WriteString(searchStyle.Render("Search: " + m.searchInput + "_"))
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("Sort: %s (1=name, 2=size, 3=time, 4=type)\n", m.sortBy))
		s.WriteString("\n")
	}

	// File list
	filesToShow := m.files
	if m.navMode {
		filesToShow = m.filteredFiles
	}

	if len(filesToShow) == 0 {
		if m.navMode && m.searchInput != "" {
			s.WriteString("No files match your search\n")
		} else {
			s.WriteString("No files in this directory\n")
		}
	} else {
		// Calculate visible range
		visibleItems := m.height - 8
		if m.navMode {
			visibleItems = m.height - 10
		}
		if visibleItems < 1 {
			visibleItems = 1
		}
		
		start := m.offset
		end := start + visibleItems
		if end > len(filesToShow) {
			end = len(filesToShow)
		}

		// Show only visible files
		for i := start; i < end; i++ {
			file := filesToShow[i]
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			style := itemStyle
			if m.cursor == i {
				style = selectedStyle
			}

			icon := "📄"
			if file.IsDir() {
				icon = "📁"
			}

			line := fmt.Sprintf("%s %s %s", cursor, icon, file.Name())
			s.WriteString(style.Render(line))
			s.WriteString("\n")
		}

		// Show scroll indicator if there are more files
		if len(filesToShow) > visibleItems {
			if m.offset > 0 {
				s.WriteString("... ↑ more files above ...\n")
			}
			if end < len(filesToShow) {
				s.WriteString("... ↓ more files below ...\n")
			}
		}
	}

	s.WriteString("\n")

	// Show message if any
	if m.message != "" {
		s.WriteString(messageStyle.Render(m.message))
		s.WriteString("\n\n")
	}

	// Command input or navigation mode info
	if m.navMode {
		s.WriteString("Navigation Mode: ESC to exit, type to search, 1-4 to sort, ↑↓ to navigate\n")
	} else {
		s.WriteString(inputStyle.Render("$ " + m.input + "_"))
		s.WriteString("\n")
		s.WriteString("Built-in: cd, ls, pwd, clear, quit | System commands: touch, mkdir, rm, cp, mv, cat, grep, find, etc. | Navigation: ↑↓, ←→, Enter | ESC: Nav Mode\n")
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}