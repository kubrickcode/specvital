// Package all imports all parser strategies for side-effect registration.
// Usage: _ "github.com/kubrickcode/specvital/lib/parser/strategies/all"
package all

import (
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/cargotest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/cypress"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/gotesting"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/gtest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/jest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/junit4"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/junit5"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/kotest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/minitest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/mocha"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/mstest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/nunit"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/phpunit"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/playwright"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/pytest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/rspec"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/swift-testing"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/testng"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/unittest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/vitest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/xctest"
	_ "github.com/kubrickcode/specvital/lib/parser/strategies/xunit"
)
