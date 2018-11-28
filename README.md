# Speechmatics API Test

## Requirements (You can just run [./speechmatics](./speechmatics) if you are using Mac)
### Set up go working environment

> `bash` on Mac

```bash
brew install go
echo "export GOPATH=~/go
export PATH=$(go env GOPATH)/bin:$PATH" >> ~/.bash_profile
source ~./bash_profile
```

### Download source code

```bash
go get github.com/gilgameshskytrooper/speechmatics
cd ~/go/src/github.com/gilgameshskytrooper/speechmatics
```

### Define necessary environment variables

```bash
echo "export SPEECHMATICSUSERID=11111
export SPEECHMATICSAUTHTOKEN=abcdefg..." >> ~/.bash_profile
source ~/.bash_profile
```

## Run Program

```bash
go build
./speechmatics
```
