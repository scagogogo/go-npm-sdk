package npm

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestInitOptionsStruct(t *testing.T) {
	options := InitOptions{
		Name:        "test-package",
		Version:     "1.0.0",
		Description: "A test package",
		Author:      "Test Author",
		License:     "MIT",
		Private:     true,
		WorkingDir:  "/tmp/test",
		Force:       true,
	}

	// Test field values
	if options.Name != "test-package" {
		t.Errorf("Expected Name 'test-package', got '%s'", options.Name)
	}

	if options.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", options.Version)
	}

	if options.Description != "A test package" {
		t.Errorf("Expected Description 'A test package', got '%s'", options.Description)
	}

	if options.Author != "Test Author" {
		t.Errorf("Expected Author 'Test Author', got '%s'", options.Author)
	}

	if options.License != "MIT" {
		t.Errorf("Expected License 'MIT', got '%s'", options.License)
	}

	if !options.Private {
		t.Error("Expected Private to be true")
	}

	if !options.Force {
		t.Error("Expected Force to be true")
	}

	// Test JSON serialization (WorkingDir and Force should be excluded)
	data, err := json.Marshal(options)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled InitOptions
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// WorkingDir and Force should not be in JSON
	if unmarshaled.WorkingDir != "" {
		t.Error("WorkingDir should not be serialized to JSON")
	}

	if unmarshaled.Force {
		t.Error("Force should not be serialized to JSON")
	}

	// Other fields should be preserved
	if unmarshaled.Name != options.Name {
		t.Errorf("Name not preserved in JSON: expected '%s', got '%s'", options.Name, unmarshaled.Name)
	}
}

func TestInstallOptionsStruct(t *testing.T) {
	options := InstallOptions{
		SaveDev:       true,
		SaveOptional:  false,
		SaveExact:     true,
		Global:        false,
		Production:    true,
		WorkingDir:    "/tmp/project",
		Registry:      "https://registry.npmjs.org/",
		Force:         true,
		IgnoreScripts: false,
	}

	// Test all boolean fields
	if !options.SaveDev {
		t.Error("Expected SaveDev to be true")
	}

	if options.SaveOptional {
		t.Error("Expected SaveOptional to be false")
	}

	if !options.SaveExact {
		t.Error("Expected SaveExact to be true")
	}

	if options.Global {
		t.Error("Expected Global to be false")
	}

	if !options.Production {
		t.Error("Expected Production to be true")
	}

	if !options.Force {
		t.Error("Expected Force to be true")
	}

	if options.IgnoreScripts {
		t.Error("Expected IgnoreScripts to be false")
	}

	// Test string fields
	if options.WorkingDir != "/tmp/project" {
		t.Errorf("Expected WorkingDir '/tmp/project', got '%s'", options.WorkingDir)
	}

	if options.Registry != "https://registry.npmjs.org/" {
		t.Errorf("Expected Registry 'https://registry.npmjs.org/', got '%s'", options.Registry)
	}

	// Test JSON serialization
	data, err := json.Marshal(options)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled InstallOptions
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	// All fields should be preserved
	if unmarshaled.SaveDev != options.SaveDev {
		t.Error("SaveDev not preserved in JSON")
	}

	if unmarshaled.Registry != options.Registry {
		t.Error("Registry not preserved in JSON")
	}
}

func TestUninstallOptionsStruct(t *testing.T) {
	options := UninstallOptions{
		SaveDev:    true,
		Global:     false,
		WorkingDir: "/tmp/uninstall",
	}

	if !options.SaveDev {
		t.Error("Expected SaveDev to be true")
	}

	if options.Global {
		t.Error("Expected Global to be false")
	}

	if options.WorkingDir != "/tmp/uninstall" {
		t.Errorf("Expected WorkingDir '/tmp/uninstall', got '%s'", options.WorkingDir)
	}
}

func TestListOptionsStruct(t *testing.T) {
	options := ListOptions{
		Global:     true,
		Depth:      2,
		Production: false,
		WorkingDir: "/tmp/list",
		JSON:       true,
	}

	if !options.Global {
		t.Error("Expected Global to be true")
	}

	if options.Depth != 2 {
		t.Errorf("Expected Depth 2, got %d", options.Depth)
	}

	if options.Production {
		t.Error("Expected Production to be false")
	}

	if !options.JSON {
		t.Error("Expected JSON to be true")
	}

	if options.WorkingDir != "/tmp/list" {
		t.Errorf("Expected WorkingDir '/tmp/list', got '%s'", options.WorkingDir)
	}
}

func TestPublishOptionsStruct(t *testing.T) {
	options := PublishOptions{
		Tag:        "beta",
		Access:     "public",
		Registry:   "https://custom-registry.com/",
		WorkingDir: "/tmp/publish",
		DryRun:     true,
	}

	if options.Tag != "beta" {
		t.Errorf("Expected Tag 'beta', got '%s'", options.Tag)
	}

	if options.Access != "public" {
		t.Errorf("Expected Access 'public', got '%s'", options.Access)
	}

	if options.Registry != "https://custom-registry.com/" {
		t.Errorf("Expected Registry 'https://custom-registry.com/', got '%s'", options.Registry)
	}

	if options.WorkingDir != "/tmp/publish" {
		t.Errorf("Expected WorkingDir '/tmp/publish', got '%s'", options.WorkingDir)
	}

	if !options.DryRun {
		t.Error("Expected DryRun to be true")
	}
}

func TestPackage(t *testing.T) {
	repo := &Repository{
		Type: "git",
		URL:  "https://github.com/user/repo.git",
	}

	bugs := &Bugs{
		URL:   "https://github.com/user/repo/issues",
		Email: "bugs@example.com",
	}

	pkg := Package{
		Name:        "test-package",
		Version:     "1.2.3",
		Description: "A test package for testing",
		Dependencies: map[string]string{
			"lodash": "^4.17.21",
			"axios":  "^0.24.0",
		},
		DevDeps: map[string]string{
			"jest":   "^27.0.0",
			"eslint": "^8.0.0",
		},
		OptionalDeps: map[string]string{
			"fsevents": "^2.3.0",
		},
		PeerDeps: map[string]string{
			"react": ">=16.8.0",
		},
		Scripts: map[string]string{
			"test":  "jest",
			"build": "webpack",
		},
		Keywords:   []string{"test", "npm", "package"},
		Author:     "Test Author",
		License:    "MIT",
		Homepage:   "https://example.com",
		Repository: repo,
		Bugs:       bugs,
		Main:       "index.js",
		Private:    false,
	}

	// Test all fields
	if pkg.Name != "test-package" {
		t.Errorf("Expected Name 'test-package', got '%s'", pkg.Name)
	}

	if pkg.Version != "1.2.3" {
		t.Errorf("Expected Version '1.2.3', got '%s'", pkg.Version)
	}

	if len(pkg.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(pkg.Dependencies))
	}

	if pkg.Dependencies["lodash"] != "^4.17.21" {
		t.Errorf("Expected lodash version '^4.17.21', got '%s'", pkg.Dependencies["lodash"])
	}

	if len(pkg.DevDeps) != 2 {
		t.Errorf("Expected 2 dev dependencies, got %d", len(pkg.DevDeps))
	}

	if len(pkg.OptionalDeps) != 1 {
		t.Errorf("Expected 1 optional dependency, got %d", len(pkg.OptionalDeps))
	}

	if len(pkg.PeerDeps) != 1 {
		t.Errorf("Expected 1 peer dependency, got %d", len(pkg.PeerDeps))
	}

	if len(pkg.Scripts) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(pkg.Scripts))
	}

	if len(pkg.Keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(pkg.Keywords))
	}

	if pkg.Repository != repo {
		t.Error("Repository not set correctly")
	}

	if pkg.Bugs != bugs {
		t.Error("Bugs not set correctly")
	}

	if pkg.Private {
		t.Error("Expected Private to be false")
	}
}

func TestRepository(t *testing.T) {
	repo := Repository{
		Type: "git",
		URL:  "https://github.com/user/repo.git",
	}

	if repo.Type != "git" {
		t.Errorf("Expected Type 'git', got '%s'", repo.Type)
	}

	if repo.URL != "https://github.com/user/repo.git" {
		t.Errorf("Expected URL 'https://github.com/user/repo.git', got '%s'", repo.URL)
	}

	// Test JSON serialization
	data, err := json.Marshal(repo)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled Repository
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if unmarshaled.Type != repo.Type {
		t.Error("Type not preserved in JSON")
	}

	if unmarshaled.URL != repo.URL {
		t.Error("URL not preserved in JSON")
	}
}

func TestBugs(t *testing.T) {
	bugs := Bugs{
		URL:   "https://github.com/user/repo/issues",
		Email: "bugs@example.com",
	}

	if bugs.URL != "https://github.com/user/repo/issues" {
		t.Errorf("Expected URL 'https://github.com/user/repo/issues', got '%s'", bugs.URL)
	}

	if bugs.Email != "bugs@example.com" {
		t.Errorf("Expected Email 'bugs@example.com', got '%s'", bugs.Email)
	}

	// Test with only URL
	bugsURLOnly := Bugs{
		URL: "https://example.com/issues",
	}

	if bugsURLOnly.Email != "" {
		t.Error("Expected empty Email when not set")
	}

	// Test with only Email
	bugsEmailOnly := Bugs{
		Email: "support@example.com",
	}

	if bugsEmailOnly.URL != "" {
		t.Error("Expected empty URL when not set")
	}
}

func TestPackageInfo(t *testing.T) {
	now := time.Now()

	repo := &Repository{
		Type: "git",
		URL:  "https://github.com/lodash/lodash.git",
	}

	author := &Person{
		Name:  "John Doe",
		Email: "john@example.com",
		URL:   "https://johndoe.com",
	}

	packageInfo := PackageInfo{
		Name:        "lodash",
		Version:     "4.17.21",
		Description: "Lodash modular utilities.",
		Keywords:    []string{"modules", "stdlib", "util"},
		Homepage:    "https://lodash.com/",
		Repository:  repo,
		Author:      author,
		License:     "MIT",
		Dependencies: map[string]string{
			"core-js": "^3.0.0",
		},
		DevDeps: map[string]string{
			"webpack": "^5.0.0",
			"babel":   "^7.0.0",
		},
		Versions: map[string]interface{}{
			"4.17.21": map[string]interface{}{
				"name":    "lodash",
				"version": "4.17.21",
			},
			"4.17.20": map[string]interface{}{
				"name":    "lodash",
				"version": "4.17.20",
			},
		},
		Time: map[string]time.Time{
			"4.17.21": now,
			"4.17.20": now.Add(-24 * time.Hour),
		},
		DistTags: map[string]string{
			"latest": "4.17.21",
			"beta":   "4.18.0-beta.1",
		},
	}

	// Test basic fields
	if packageInfo.Name != "lodash" {
		t.Errorf("Expected Name 'lodash', got '%s'", packageInfo.Name)
	}

	if packageInfo.Version != "4.17.21" {
		t.Errorf("Expected Version '4.17.21', got '%s'", packageInfo.Version)
	}

	if len(packageInfo.Keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(packageInfo.Keywords))
	}

	if packageInfo.Repository != repo {
		t.Error("Repository not set correctly")
	}

	if packageInfo.Author != author {
		t.Error("Author not set correctly")
	}

	if len(packageInfo.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(packageInfo.Dependencies))
	}

	if len(packageInfo.DevDeps) != 2 {
		t.Errorf("Expected 2 dev dependencies, got %d", len(packageInfo.DevDeps))
	}

	if len(packageInfo.Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(packageInfo.Versions))
	}

	if len(packageInfo.Time) != 2 {
		t.Errorf("Expected 2 time entries, got %d", len(packageInfo.Time))
	}

	if len(packageInfo.DistTags) != 2 {
		t.Errorf("Expected 2 dist tags, got %d", len(packageInfo.DistTags))
	}

	if packageInfo.DistTags["latest"] != "4.17.21" {
		t.Errorf("Expected latest tag '4.17.21', got '%s'", packageInfo.DistTags["latest"])
	}
}

func TestPerson(t *testing.T) {
	// Test full person info
	person := Person{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		URL:   "https://janedoe.com",
	}

	if person.Name != "Jane Doe" {
		t.Errorf("Expected Name 'Jane Doe', got '%s'", person.Name)
	}

	if person.Email != "jane@example.com" {
		t.Errorf("Expected Email 'jane@example.com', got '%s'", person.Email)
	}

	if person.URL != "https://janedoe.com" {
		t.Errorf("Expected URL 'https://janedoe.com', got '%s'", person.URL)
	}

	// Test person with only name
	personNameOnly := Person{
		Name: "John Smith",
	}

	if personNameOnly.Name != "John Smith" {
		t.Errorf("Expected Name 'John Smith', got '%s'", personNameOnly.Name)
	}

	if personNameOnly.Email != "" {
		t.Error("Expected empty Email when not set")
	}

	if personNameOnly.URL != "" {
		t.Error("Expected empty URL when not set")
	}

	// Test JSON serialization
	data, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled Person
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if unmarshaled.Name != person.Name {
		t.Error("Name not preserved in JSON")
	}

	if unmarshaled.Email != person.Email {
		t.Error("Email not preserved in JSON")
	}

	if unmarshaled.URL != person.URL {
		t.Error("URL not preserved in JSON")
	}
}

func TestSearchResult(t *testing.T) {
	author := &Person{
		Name:  "Test Author",
		Email: "author@example.com",
	}

	publisher := &Person{
		Name:  "Test Publisher",
		Email: "publisher@example.com",
	}

	maintainer1 := &Person{
		Name:  "Maintainer 1",
		Email: "maintainer1@example.com",
	}

	maintainer2 := &Person{
		Name:  "Maintainer 2",
		Email: "maintainer2@example.com",
	}

	searchPackage := SearchPackage{
		Name:        "test-search-package",
		Version:     "2.1.0",
		Description: "A package for search testing",
		Keywords:    []string{"search", "test", "npm"},
		Date:        time.Now(),
		Links: map[string]string{
			"npm":        "https://www.npmjs.com/package/test-search-package",
			"homepage":   "https://example.com",
			"repository": "https://github.com/user/test-search-package",
		},
		Author:      author,
		Publisher:   publisher,
		Maintainers: []*Person{maintainer1, maintainer2},
	}

	scoreDetail := ScoreDetail{
		Quality:     0.85,
		Popularity:  0.72,
		Maintenance: 0.91,
	}

	searchScore := SearchScore{
		Final:  0.83,
		Detail: scoreDetail,
	}

	searchResult := SearchResult{
		Package:     searchPackage,
		Score:       searchScore,
		SearchScore: 0.95,
	}

	// Test SearchPackage fields
	if searchResult.Package.Name != "test-search-package" {
		t.Errorf("Expected package name 'test-search-package', got '%s'", searchResult.Package.Name)
	}

	if len(searchResult.Package.Keywords) != 3 {
		t.Errorf("Expected 3 keywords, got %d", len(searchResult.Package.Keywords))
	}

	if len(searchResult.Package.Links) != 3 {
		t.Errorf("Expected 3 links, got %d", len(searchResult.Package.Links))
	}

	if searchResult.Package.Author != author {
		t.Error("Author not set correctly")
	}

	if searchResult.Package.Publisher != publisher {
		t.Error("Publisher not set correctly")
	}

	if len(searchResult.Package.Maintainers) != 2 {
		t.Errorf("Expected 2 maintainers, got %d", len(searchResult.Package.Maintainers))
	}

	// Test SearchScore fields
	if searchResult.Score.Final != 0.83 {
		t.Errorf("Expected final score 0.83, got %f", searchResult.Score.Final)
	}

	if searchResult.Score.Detail.Quality != 0.85 {
		t.Errorf("Expected quality score 0.85, got %f", searchResult.Score.Detail.Quality)
	}

	if searchResult.Score.Detail.Popularity != 0.72 {
		t.Errorf("Expected popularity score 0.72, got %f", searchResult.Score.Detail.Popularity)
	}

	if searchResult.Score.Detail.Maintenance != 0.91 {
		t.Errorf("Expected maintenance score 0.91, got %f", searchResult.Score.Detail.Maintenance)
	}

	if searchResult.SearchScore != 0.95 {
		t.Errorf("Expected search score 0.95, got %f", searchResult.SearchScore)
	}
}

func TestCommandResult(t *testing.T) {
	// Test successful command result
	successResult := CommandResult{
		Success:  true,
		ExitCode: 0,
		Stdout:   "Command executed successfully",
		Stderr:   "",
		Duration: 2 * time.Second,
		Error:    nil,
	}

	if !successResult.Success {
		t.Error("Expected Success to be true")
	}

	if successResult.ExitCode != 0 {
		t.Errorf("Expected ExitCode 0, got %d", successResult.ExitCode)
	}

	if successResult.Stdout != "Command executed successfully" {
		t.Errorf("Expected Stdout 'Command executed successfully', got '%s'", successResult.Stdout)
	}

	if successResult.Stderr != "" {
		t.Errorf("Expected empty Stderr, got '%s'", successResult.Stderr)
	}

	if successResult.Duration != 2*time.Second {
		t.Errorf("Expected Duration 2s, got %v", successResult.Duration)
	}

	if successResult.Error != nil {
		t.Errorf("Expected no error, got %v", successResult.Error)
	}

	// Test failed command result
	failedResult := CommandResult{
		Success:  false,
		ExitCode: 1,
		Stdout:   "",
		Stderr:   "Command failed with error",
		Duration: 500 * time.Millisecond,
		Error:    errors.New("command failed"),
	}

	if failedResult.Success {
		t.Error("Expected Success to be false")
	}

	if failedResult.ExitCode != 1 {
		t.Errorf("Expected ExitCode 1, got %d", failedResult.ExitCode)
	}

	if failedResult.Stderr != "Command failed with error" {
		t.Errorf("Expected Stderr 'Command failed with error', got '%s'", failedResult.Stderr)
	}

	if failedResult.Error == nil {
		t.Error("Expected error to be set")
	}

	// Test JSON serialization (Error field should be omitted when nil)
	data, err := json.Marshal(successResult)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var unmarshaled CommandResult
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if unmarshaled.Success != successResult.Success {
		t.Error("Success not preserved in JSON")
	}

	if unmarshaled.Duration != successResult.Duration {
		t.Error("Duration not preserved in JSON")
	}
}
