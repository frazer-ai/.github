run "greets_by_name" {
  command = apply

  variables {
    name = "test"
  }

  assert {
    condition     = output.greeting == "hello test"
    error_message = "unexpected greeting"
  }
}
