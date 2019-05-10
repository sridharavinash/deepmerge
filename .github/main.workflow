workflow "Run Tests" {
  resolves = ["Go Action"]
  on = "push"
}

action "Go Action" {
  uses = "actions-contrib/go@v0.1.0"
  runs = "go"
  args = "test -v ./..."
}
