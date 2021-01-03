/**
 * @Author Hatch
 * @Date 2021/01/03 11:29
**/
package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Test main function")
	retCode := m.Run()

	Run()
	os.Exit(retCode)
}


