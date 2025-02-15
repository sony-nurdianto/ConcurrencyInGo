package main

import "testing"

func BenchmarkUnblockingConnHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UnblockingConnHelloWorld()
	}
}

func BenchmarkConnHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConHelloWorld()
	}
}

func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HelloWorld()
	}
}
