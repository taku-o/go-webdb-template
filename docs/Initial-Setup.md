# Initial Setup

## Install Package Application
* Docker
    * https://www.docker.com/ja-jp/
* Cursor
    * https://cursor.com/ja?from=home
* Go
    * https://go.dev/dl/

## Homebrew

### Homebrew
```bash
    /bin/bash -c"$(curl -fsSLhttps://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    eval "$(/opt/homebrew/bin/brew shellenv)"
    $(/opt/homebrew/bin/brew shellenv)
```

### GitHub CLI
```bash
    brew install gh
    gh auth login
    gh auth status
```

### Atlas
```bash
    brew install ariga/tap/atlas
```

## Node

### nvm
```bash
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.4/install.sh | bash
    nvm ls-remote
    nvm install v22.14.0
    nvm use v22.14.0
    nvm alias default v22.14.0
```

* add .bashrc
```bash
    if [ -f ~/.nvm/nvm.sh ]
    then
      source ~/.nvm/nvm.sh
    fi
```

### Claude Code
```bash
    npm install -g @anthropic-ai/claude-code
```

## uv

### uv
```bash
    brew install uv
```

### Serena
* use this command at project directory.
```bash
    claude mcp add serena -- uvx --from git+https://github.com/oraios/serena serena-mcp-server --context ide-assistant --enable-web-dashboard false --project $(pwd)
```

* use some times, to update serena index files.
```bash
    uvx --from git+https://github.com/oraios/serena index-project
```


