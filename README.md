# FileRover - Terminal File Manager

A powerful, command-driven terminal file manager built with Go and the Bubble Tea TUI framework. FileRover provides an intuitive interface for file navigation, management, and system command execution.

## 🚀 Features

### 📁 **File Navigation**
- **Single-panel interface** with full-screen TUI
- **Arrow key navigation** (↑↓ for file list, ←→ for directory navigation)
- **Smart cursor positioning** - cursor stays at the directory you came from when going back
- **Scrollable file lists** for directories with many files
- **Visual file indicators** (📁 for directories, 📄 for files)

### 🔍 **Navigation Mode**
- **Press `ESC`** to enter/exit navigation mode
- **Real-time search** - type to filter files by name
- **Multiple sort options**:
  - `1` - Sort by name (directories first)
  - `2` - Sort by size (largest first)
  - `3` - Sort by modification time (newest first)
  - `4` - Sort by type (directories first)
- **Case-insensitive search** with instant filtering

### 💻 **Command System**
- **Always-on command prompt** - no need to press `:` to enter commands
- **Built-in commands**:
  - `cd [directory]` - Change directory
  - `ls` / `dir` - Refresh directory listing
  - `pwd` - Show current directory
  - `clear` - Clear messages
  - `quit` / `exit` - Exit application
- **System command execution** - Run any command available on your system:
  - `touch`, `mkdir`, `rm`, `cp`, `mv`, `cat`
  - `grep`, `find`, `chmod`, `chown`
  - `tar`, `zip`, `unzip`
  - And many more!

### 🎨 **User Interface**
- **Colorful TUI** with different styles for different elements
- **Full-screen mode** using alternate screen buffer
- **Command feedback** - shows success/error messages
- **Scroll indicators** - shows when there are more files above/below
- **Dynamic help text** - different instructions for normal and navigation modes

## 📦 Installation

### Prerequisites
- Go 1.19 or later
- Linux/macOS/Windows terminal

### Build from Source
```bash
git clone <repository-url>
cd filerover
go mod tidy
go build -o rover main.go
```

### Run
```bash
./rover
```

## 🎮 Controls

### **Normal Mode**
| Key | Action |
|-----|--------|
| `↑/↓` | Navigate file list |
| `←` | Go to parent directory |
| `→` | Enter selected directory |
| `Enter` | Execute command or enter directory |
| `ESC` | Enter navigation mode |
| `Ctrl+C` | Exit application |
| `Backspace` | Delete characters in command input |

### **Navigation Mode**
| Key | Action |
|-----|--------|
| `ESC` | Exit navigation mode |
| `↑/↓` or `j/k` | Navigate filtered results |
| `Enter` | Enter selected directory |
| `1-4` | Sort by name/size/time/type |
| `Backspace` | Delete search characters |
| Type | Add characters to search |

### **Command Input**
- Type any command and press `Enter` to execute
- Commands are executed in the current directory context
- Output is displayed in the message area
- Long output is truncated for readability

## 📋 Usage Examples

### **Basic Navigation**
```bash
# Navigate directories
cd Documents
cd /home/user/Downloads
cd ..  # Go to parent directory

# List files
ls
ls -la  # Detailed listing
```

### **File Operations**
```bash
# Create files and directories
touch newfile.txt
mkdir newfolder
mkdir -p path/to/nested/directory

# Copy and move files
cp file1.txt file2.txt
cp -r folder1 folder2
mv oldname.txt newname.txt

# Remove files
rm unwanted.txt
rm -rf directory
```

### **Search and Filter**
```bash
# Search for files
find . -name "*.txt"
grep "pattern" *.txt

# View file contents
cat readme.txt
head -n 20 largefile.txt
```

### **Navigation Mode**
1. Press `ESC` to enter navigation mode
2. Type "doc" to search for files containing "doc"
3. Press `2` to sort by size
4. Press `→` to enter a directory
5. Press `ESC` to exit navigation mode

## 🛠️ Technical Details

### **Architecture**
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Uses [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling
- Command execution via `os/exec` package
- Real-time file system monitoring

### **File Structure**
```
filerover/
├── main.go          # Main application code
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
└── README.md        # This file
```

### **Dependencies**
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling library

## 🎯 Design Philosophy

FileRover is designed to be:
- **Command-first**: Emphasizes command-line operations over GUI interactions
- **Keyboard-driven**: Optimized for keyboard navigation and shortcuts
- **System-integrated**: Uses native system commands for maximum compatibility
- **Efficient**: Minimal resource usage and fast performance
- **Intuitive**: Familiar interface patterns from traditional file managers

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

### **Development Setup**
```bash
git clone <repository-url>
cd filerover
go mod tidy
go run main.go
```

## 📄 License

This project is open source. Please check the license file for details.

