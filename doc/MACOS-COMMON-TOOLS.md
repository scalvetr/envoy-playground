# MacOS: Install common tools

Install Hombrew & kubectl

```shell
xcode-select --install
# Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

brew install kubectl@1.22.2
brew install helm@3.7.1

echo "source <(kubectl completion zsh)" >> ~/.zshrc
```