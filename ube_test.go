package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	// Run the command
	cmd := exec.Command("./ubescript", "./tests/hello.ube")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	// Write the command output to a file
	err = ioutil.WriteFile("output.txt", output, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Read the output file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := "hello world luigi\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}

	// Print the file content
	fmt.Println(string(data))
}
func TestFib(t *testing.T) {
	// Run the command
	cmd := exec.Command("./ubescript", "./tests/fib.ube")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	// Write the command output to a file
	err = ioutil.WriteFile("output.txt", output, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Read the output file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := "0\n1\n1\n2\n3\n5\n8\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}

	// Print the file content
	fmt.Println(string(data))
}
func TestAssign(t *testing.T) {
	// Run the command
	cmd := exec.Command("./ubescript", "./tests/assign.ube")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	// Write the command output to a file
	err = ioutil.WriteFile("output.txt", output, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Read the output file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := "6561\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}

	// Print the file content
	fmt.Println(string(data))
}

func TestAssignmentPrecedence(t *testing.T) {
	// Run the command
	cmd := exec.Command("./ubescript", "./tests/brekky.ube")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	// Write the command output to a file
	err = ioutil.WriteFile("output.txt", output, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Read the output file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := "beignets\nbeignets with cafe au lait\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}

	// Print the file content
	fmt.Println(string(data))
}

func TestLocalScopes(t *testing.T) {
	// Run the command
	cmd := exec.Command("./ubescript", "./tests/scopes.ube")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	// Write the command output to a file
	err = ioutil.WriteFile("output.txt", output, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Read the output file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := "12\n3000\n143\n24\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}

	// Print the file content
	fmt.Println(string(data))
}

func TestIf(t *testing.T) {
	// Run the command
	cmd := exec.Command("./ubescript", "./tests/if.ube")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	// Write the command output to a file
	err = ioutil.WriteFile("output.txt", output, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Read the output file
	data, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := "goodbye world!\nhello world!\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}

	// Print the file content
	fmt.Println(string(data))
}
