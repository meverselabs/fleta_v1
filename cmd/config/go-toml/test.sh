#!/bin/bash
# fail out of the script if anything here fails
set -e
set -o pipefail

# set the path to the present working directory
export GOPATH=`pwd`

function git_clone() {
  path=$1
  branch=$2
  version=$3
  if [ ! -d "src/$path" ]; then
    mkdir -p src/$path
    git clone https://$path.git src/$path
  fi
  pushd src/$path
  git checkout "$branch"
  git reset --hard "$version"
  popd
}

# Remove potential previous runs
rm -rf src test_program_bin toml-test

go get github.com/davecgh/go-spew/spew
go get gopkg.in/yaml.v2
go get github.com/BurntSushi/toml

# get code for BurntSushi TOML validation
git_clone github.com/BurntSushi/toml master a368813
git_clone github.com/BurntSushi/toml-test master 39e37e6

# build the BurntSushi test application
go build -o toml-test github.com/BurntSushi/toml-test

# vendorize the current lib for testing
# NOTE: this basically mocks an install without having to go back out to github for code
mkdir -p src/github.com/fletaio/fleta_v1/cmd/config/go-toml/cmd
mkdir -p src/github.com/fletaio/fleta_v1/cmd/config/go-toml/query
cp *.go *.toml src/github.com/fletaio/fleta_v1/cmd/config/go-toml
cp -R cmd/* src/github.com/fletaio/fleta_v1/cmd/config/go-toml/cmd
cp -R query/* src/github.com/fletaio/fleta_v1/cmd/config/go-toml/query
go build -o test_program_bin src/github.com/fletaio/fleta_v1/cmd/config/go-toml/cmd/test_program.go

# Run basic unit tests
go test github.com/fletaio/fleta_v1/cmd/config/go-toml -race -coverprofile=coverage.txt -covermode=atomic
go test github.com/fletaio/fleta_v1/cmd/config/go-toml/cmd/tomljson
go test github.com/fletaio/fleta_v1/cmd/config/go-toml/query

# run the entire BurntSushi test suite
if [[ $# -eq 0 ]] ; then
  echo "Running all BurntSushi tests"
  ./toml-test ./test_program_bin | tee test_out
else
  # run a specific test
  test=$1
  test_path='src/github.com/BurntSushi/toml-test/tests'
  valid_test="$test_path/valid/$test"
  invalid_test="$test_path/invalid/$test"

  if [ -e "$valid_test.toml" ]; then
    echo "Valid Test TOML for $test:"
    echo "===="
    cat "$valid_test.toml"

    echo "Valid Test JSON for $test:"
    echo "===="
    cat "$valid_test.json"

    echo "Go-TOML Output for $test:"
    echo "===="
    cat "$valid_test.toml" | ./test_program_bin
  fi

  if [ -e "$invalid_test.toml" ]; then
    echo "Invalid Test TOML for $test:"
    echo "===="
    cat "$invalid_test.toml"

    echo "Go-TOML Output for $test:"
    echo "===="
    echo "go-toml Output:"
    cat "$invalid_test.toml" | ./test_program_bin
  fi
fi
