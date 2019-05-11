workflow "Run Tests" {
  resolves = ["Go Action"]
  on = "push"
}

action "Go Action" {
  uses = "actions-contrib/go@v0.1.0"
  runs = "go"
  args = "test -v ./..."
}

workflow "Lint" {
  on = "push"
  resolves = ["docker://golangci/golangci-lint:latest"]
}

workflow "New workflow 1" {
  on = "push"
}

action "docker://golangci/golangci-lint:latest" {
  uses = "docker://golangci/golangci-lint:latest"
  args = "run"
}
