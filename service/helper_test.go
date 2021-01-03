/**
 * @Author Hatch
 * @Date 2021/01/03 11:25
**/
package service

import "testing"

func TestFileExist(t *testing.T) {
	filename := "./helper.go"
	if !FileExist(filename) {
		t.Fatalf("File doesn't exists:%v", filename)
	}
}

func TestIsDir(t *testing.T) {
	path := "../resource"
	if !IsDir(path) {
		t.Fatalf("Folder doesn't exists:%v", path)
	}
}

func TestIsFile(t *testing.T) {
	path := "./helper.go"
	if !IsFile(path) {
		t.Fatalf("The path %v isn't file", path)
	}
}

func TestIsAllowExtFile(t *testing.T) {
	filename := "./helper.go"
	allowExts := map[string]bool{
		"go": true,
		"php": false,
	}

	if !IsAllowExtFile(filename, allowExts) {
		t.Fatal("The function IsAllowExtFile has some error")
	}
}

// Benchmark
func BenchmarkFileExist(b *testing.B) {
	filename := "./helper.go"
	for i := 0; i < b.N; i++ {
		FileExist(filename)
	}
}

func BenchmarkIsDir(b *testing.B) {
	path := "../resource"
	for i := 0; i < b.N; i++ {
		IsDir(path)
	}
}

func BenchmarkIsFile(b *testing.B) {
	filename := "./helper.go"
	for i := 0; i < b.N; i++ {
		IsFile(filename)
	}
}

func BenchmarkIsAllowExtFile(b *testing.B) {
	filename := "./helper.go"
	allowExts := map[string]bool{
		"go": true,
		"php": false,
	}

	b.SetParallelism(2)
	b.RunParallel(func (pb *testing.PB) {
		for pb.Next() {
			IsAllowExtFile(filename, allowExts)
		}
	})

}
