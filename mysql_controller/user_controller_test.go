package controller

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateAccount(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
