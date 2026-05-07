package domain

import (
	"testing"

	"ez-admin-gin/server/internal/platform/model"
)

func TestBuildAncestors(t *testing.T) {
	tests := []struct {
		name   string
		parent model.Department
		want   string
	}{
		{"root", model.Department{}, "0"},
		{"child of root", model.Department{ID: 1, Ancestors: "0"}, "0,1"},
		{"deep node", model.Department{ID: 5, Ancestors: "0,1,3"}, "0,1,3,5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildAncestors(tt.parent); got != tt.want {
				t.Errorf("BuildAncestors() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFullPath(t *testing.T) {
	tests := []struct {
		name string
		item model.Department
		want string
	}{
		{"root node", model.Department{ID: 1, Ancestors: ""}, "1"},
		{"child", model.Department{ID: 3, Ancestors: "0,1"}, "0,1,3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FullPath(tt.item); got != tt.want {
				t.Errorf("FullPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsDescendantPath(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		target string
		want   bool
	}{
		{"exact match", "0,1,3", "0,1,3", true},
		{"prefix match", "0,1,3,5", "0,1,3", true},
		{"no relation", "0,2,4", "0,1,3", false},
		{"partial prefix not match", "0,1,30", "0,1,3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDescendantPath(tt.path, tt.target); got != tt.want {
				t.Errorf("IsDescendantPath(%q, %q) = %v, want %v", tt.path, tt.target, got, tt.want)
			}
		})
	}
}
