# 🛠️ Go Development Setup with mise

## 🚀 **Why mise for Go?**

Instead of using system-managed Go packages, we use [mise](https://mise.jdx.dev/) for several advantages:

- ✅ **Latest Go versions** - Always get the newest features and performance
- ✅ **No sudo required** - User-level installation
- ✅ **Version management** - Easy to switch between Go versions
- ✅ **Project isolation** - Different projects can use different Go versions
- ✅ **Cleaner system** - No conflicts with system packages

## 📦 **Installing mise**

### Quick Installation
```bash
# Install mise
curl https://mise.run | sh

# Add to your shell profile
echo 'eval "$(~/.local/bin/mise activate bash)"' >> ~/.bashrc  # for bash
echo 'eval "$(~/.local/bin/mise activate zsh)"' >> ~/.zshrc   # for zsh

# Reload your shell
source ~/.bashrc  # or source ~/.zshrc
```

### Alternative Installation Methods
```bash
# Via AUR (if you prefer)
paru -S mise-bin  # or yay -S mise-bin

# Via cargo (if you have Rust)
cargo install mise

# Via package manager
curl https://mise.jdx.dev/install.sh | sh
```

## 🔧 **Setting up Go**

### Global Go Installation
```bash
# Install latest Go globally
mise use -g go@latest

# Verify installation
go version

# Check mise status
mise list
```

### Project-specific Go
```bash
# In your project directory
cd /path/to/hyprland-monitor-tui

# Use specific Go version for this project
mise use go@1.22.0  # or go@latest

# This creates a .mise.toml file in your project
cat .mise.toml
```

### Example .mise.toml
```toml
[tools]
go = "latest"  # or "1.22.0" for specific version
```

## 🎯 **Benefits for This Project**

### Development Workflow
```bash
# Clone the project
git clone <repository>
cd hyprland-monitor-tui

# mise automatically activates the right Go version
mise use go@latest

# Build with the latest Go features
go build -o hyprland-monitor-tui .

# Install
./install.sh  # Will detect Go via mise
```

### Multiple Go Projects
```bash
# Project A uses Go 1.21
cd project-a
mise use go@1.21.0

# Project B uses Go 1.22
cd project-b  
mise use go@1.22.0

# Project C uses latest
cd hyprland-monitor-tui
mise use go@latest
```

## 🔍 **Troubleshooting**

### Go Not Found After Installation
```bash
# Check mise status
mise doctor

# List installed tools
mise list

# Manually activate mise
eval "$(mise activate bash)"

# Reinstall Go if needed
mise uninstall go
mise use -g go@latest
```

### Path Issues
```bash
# Check if mise is in PATH
which mise

# Check Go path
which go

# If issues, ensure mise is properly activated
echo 'eval "$(mise activate bash)"' >> ~/.bashrc
source ~/.bashrc
```

### Installation Script Fails
```bash
# If mise not found during install.sh
curl https://mise.run | sh
source ~/.bashrc
./install.sh

# Or install Go manually first
mise use -g go@latest
./install.sh
```

## 📋 **Commands Reference**

### Essential mise Commands
```bash
# Install latest Go
mise use -g go@latest

# Install specific version
mise use -g go@1.22.0

# List available Go versions
mise list-remote go

# List installed versions
mise list

# Update to latest
mise upgrade go

# Remove version
mise uninstall go@1.21.0

# Show current versions
mise current
```

### Integration with Installation
```bash
# Our install.sh will automatically:
1. Check if Go is available
2. If not, check for mise
3. Install Go via: mise use -g go@latest
4. Continue with build process
```

## 🎨 **Perfect for Development**

Using mise for this Tokyo Night TUI project gives you:

- 🚀 **Latest Go performance** for faster compilation
- 🔧 **Modern Go features** for cleaner code
- 📦 **No system conflicts** with other packages
- ⚡ **Quick setup** on any development machine
- 🎯 **Consistent environment** across team members

---

**Ready to build the most beautiful TUI with the latest Go! 🎉** 