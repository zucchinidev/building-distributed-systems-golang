package main

import "testing"

func TestNewPath(t *testing.T) {
	for _, testCase := range getTestCases() {
		newPath := NewPath(testCase.OriginalPath)
		if testCase.Path != newPath.Path {
			t.Errorf("Error parsing the path, we expected %s, got %s", testCase.Path, newPath.Path)
		}

		if testCase.HasId != newPath.HasID() {
			t.Errorf("Error checking the has id, we expected %v, got %v", testCase.HasId, newPath.HasID())
		}

		if testCase.ID != newPath.ID {
			t.Errorf("Error getting the identifier, we expected %s, got %s", testCase.ID, newPath.ID)
		}
	}
}

type newPathTestCase struct {
	OriginalPath string
	Path         string
	ID           string
	HasId        bool
}

func getTestCases() []newPathTestCase {
	return []newPathTestCase{
		{OriginalPath: "/", Path: "", ID: "", HasId: false},
		{OriginalPath: "/people/", Path: "people", ID: "", HasId: false},
		{OriginalPath: "/people/1", Path: "people", ID: "1", HasId: true},
	}
}
