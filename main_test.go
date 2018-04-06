package main

import (
	"testing"
	"os"
)

func TestMain(m *testing.M){
	println("Before all...")

	code := m.Run()

	println("After all...")

	os.Exit(code)
}

func TestChHomeDirSuccess(t *testing.T){
}

func TestChHomeDirFailed(t *testing.T){

}

func TestHandleCtrlCSuccess(t *testing.T){

}

func TestHandleCtrlCFailed(t *testing.T){

}

func TestEncodingJsonSuccess(t *testing.T){

}

func TestEncodingJsonFailed(t *testing.T){

}

func TestGetServicesInfoSuccess(t *testing.T){

}

func TestGetServicesInfoFailed(t *testing.T){

}
