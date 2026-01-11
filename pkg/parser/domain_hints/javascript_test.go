package domain_hints

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/domain"
)

func TestJavaScriptExtractor_Extract_ES6Imports(t *testing.T) {
	source := []byte(`
import { test, expect } from '@playwright/test';
import axios from 'axios';
import * as lodash from 'lodash';
import '@testing-library/jest-dom';
import type { User } from './types';

test('should work', async () => {
  const mockUser = { name: 'test' };
  authService.validateToken();
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	t.Run("imports", func(t *testing.T) {
		expectedImports := map[string]bool{
			"@playwright/test":          true,
			"axios":                     true,
			"lodash":                    true,
			"@testing-library/jest-dom": true,
		}

		// type-only import should be excluded
		excludedImports := []string{"./types"}

		importSet := make(map[string]bool)
		for _, imp := range hints.Imports {
			importSet[imp] = true
		}

		for imp := range expectedImports {
			if !importSet[imp] {
				t.Errorf("expected import %q to be included", imp)
			}
		}

		for _, imp := range excludedImports {
			if importSet[imp] {
				t.Errorf("expected type-only import %q to be excluded", imp)
			}
		}
	})
}

func TestJavaScriptExtractor_Extract_CommonJS(t *testing.T) {
	source := []byte(`
const lodash = require('lodash');
const { get } = require('axios');
const path = require('path');

test('should work', async () => {
  const mockData = getData();
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageJavaScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"lodash": true,
		"axios":  true,
		"path":   true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestJavaScriptExtractor_Extract_Calls(t *testing.T) {
	source := []byte(`
import { test, expect } from '@playwright/test';

test('should work', async () => {
  authService.validateToken('token');
  userRepo.findById(1);
  const result = orderService.create(order);
  doSomething();
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"authService.validateToken": true,
		"userRepo.findById":         true,
		"orderService.create":       true,
		"doSomething":               true,
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"test", "expect"}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	for call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q to be included, got %v", call, hints.Calls)
		}
	}

	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected test framework call %q to be excluded", call)
		}
	}
}

func TestJavaScriptExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`// empty file`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageJavaScript}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestJavaScriptExtractor_Extract_MixedImports(t *testing.T) {
	source := []byte(`
import { test } from '@playwright/test';
const axios = require('axios');
import type { Response } from 'express';

test('mixed imports', async () => {
  const mockResponse = {};
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	// Both ES6 and CommonJS should be captured
	if !importSet["@playwright/test"] {
		t.Error("expected @playwright/test import")
	}
	if !importSet["axios"] {
		t.Error("expected axios import (CommonJS)")
	}

	// Type-only imports should be excluded
	if importSet["express"] {
		t.Error("expected type-only express import to be excluded")
	}
}

func TestJavaScriptExtractor_Extract_PlaywrightFile(t *testing.T) {
	source := []byte(`
import { test, expect } from '@playwright/test';
import { LoginPage } from './pages/login';

test.describe('authentication flow', () => {
  const mockCredentials = { email: 'test@example.com', password: 'secret' };

  test('should login successfully', async ({ page }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await authService.login(mockCredentials);
    await expect(page).toHaveURL('/dashboard');
  });
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Verify imports include both library and local imports
	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	if !importSet["@playwright/test"] {
		t.Error("expected @playwright/test import")
	}
	if !importSet["./pages/login"] {
		t.Error("expected ./pages/login import")
	}

	// Verify calls (excluding test framework)
	callSet := make(map[string]bool)
	for _, c := range hints.Calls {
		callSet[c] = true
	}

	if !callSet["authService.login"] {
		t.Errorf("expected authService.login call, got %v", hints.Calls)
	}
}

func TestJavaScriptExtractor_Extract_TSX(t *testing.T) {
	source := []byte(`
import React from 'react';
import { render, screen } from '@testing-library/react';
import { UserProfile } from './UserProfile';

test('should render user profile', () => {
  const mockUser = { id: 1, name: 'John' };
  render(<UserProfile user={mockUser} />);
  userService.getProfile(mockUser.id);
  expect(screen.getByText('John')).toBeInTheDocument();
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTSX}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Verify imports
	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	if !importSet["react"] {
		t.Error("expected react import")
	}
	if !importSet["@testing-library/react"] {
		t.Error("expected @testing-library/react import")
	}
	if !importSet["./UserProfile"] {
		t.Error("expected ./UserProfile import")
	}

	// Verify calls
	callSet := make(map[string]bool)
	for _, c := range hints.Calls {
		callSet[c] = true
	}

	if !callSet["userService.getProfile"] {
		t.Errorf("expected userService.getProfile call, got %v", hints.Calls)
	}
}

func TestJavaScriptExtractor_Extract_CallsNormalization(t *testing.T) {
	source := []byte(`
import { test } from '@playwright/test';

test('chained calls', async () => {
  // Long chains should be normalized to 2 segments
  e2eSelectors.queryEditor.resourcePicker.select.button().click();
  e2e.components.NavToolbar.editDashboard.editButton().should('be.visible');

  // Calls with newlines should be normalized
  e2eSelectors.configEditor
    .azureCloud
    .input()
    .find('input')
    .type('Azure');
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, c := range hints.Calls {
		callSet[c] = true
	}

	// Should be normalized to 2 segments
	expectedCalls := []string{
		"e2eSelectors.queryEditor",
		"e2e.components",
		"e2eSelectors.configEditor",
	}

	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected %q call (2-segment normalized), got %v", call, hints.Calls)
		}
	}

	// Full chains should NOT be present
	for call := range callSet {
		if len(call) > 50 {
			t.Errorf("call too long (should be normalized): %s", call)
		}
	}
}

func TestJavaScriptExtractor_Extract_ImportsNotOvertaken(t *testing.T) {
	// This test ensures describe/it strings are NOT captured as imports
	source := []byte(`
import { test } from '@playwright/test';

describe('Azure monitor datasource', () => {
  it('create dashboard with panels', () => {
    addVariable('subscription');
    addVariable('resourceGroups');
  });
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	// Only real imports should be captured
	if !importSet["@playwright/test"] {
		t.Error("expected @playwright/test import")
	}

	// describe/it arguments should NOT be in imports
	invalidImports := []string{
		"Azure monitor datasource",
		"create dashboard with panels",
		"subscription",
		"resourceGroups",
	}

	for _, inv := range invalidImports {
		if importSet[inv] {
			t.Errorf("describe/it argument %q should NOT be in imports", inv)
		}
	}

	// Verify only 1 import
	if len(hints.Imports) != 1 {
		t.Errorf("expected 1 import, got %d: %v", len(hints.Imports), hints.Imports)
	}
}

func TestJavaScriptExtractor_Extract_MockFunctionFiltering(t *testing.T) {
	source := []byte(`
import { test, expect, describe } from '@playwright/test';

describe('user service', () => {
	test('should create user', () => {
		const fn = jest.fn();
		fn(userData);
		userService.create(userData);

		const mockValidator = vi.fn();
		mockValidator();

		const mockAPI = jest.fn().mockResolvedValue({ id: 1 });
		mockAPI();
		apiClient.send(request);
	});
});
`)

	extractor := &JavaScriptExtractor{lang: domain.LanguageTypeScript}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	t.Run("mock functions filtered", func(t *testing.T) {
		callSet := make(map[string]bool)
		for _, call := range hints.Calls {
			callSet[call] = true
		}

		// Standalone fn() calls should be excluded
		if callSet["fn"] {
			t.Errorf("expected standalone fn() to be filtered, got %v", hints.Calls)
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		callSet := make(map[string]bool)
		for _, call := range hints.Calls {
			callSet[call] = true
		}

		// Domain calls should be included
		expectedCalls := []string{"userService.create", "apiClient.send"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("jest.fn and vi.fn filtered", func(t *testing.T) {
		callSet := make(map[string]bool)
		for _, call := range hints.Calls {
			callSet[call] = true
		}

		// jest and vi framework calls should already be filtered
		if callSet["jest"] || callSet["vi"] {
			t.Errorf("expected jest/vi framework calls to be filtered, got %v", hints.Calls)
		}
	})
}
