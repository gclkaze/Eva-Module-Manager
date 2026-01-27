# Validation Tests Documentation

## Overview

Comprehensive test suite for all validation functions in `pkg/utils/validations.go`. The tests cover all validation functions with extensive permutations for different input types, edge cases, boundary conditions, and error scenarios.

## Test Coverage

### 1. **IsValidEmail** (15 test cases)
Tests email validation with different scenarios:
- ✅ Valid emails (basic, with subdomains, with dots/hyphens/numbers)
- ✅ Length boundaries (too short, too long, edge cases)
- ✅ Format validation (missing @, multiple @, invalid domain)
- ✅ Invalid characters and control characters

**Key Test Cases:**
```
valid basic email: "user@example.com"
valid subdomain: "user@mail.example.co.uk"
valid max length boundary
too short: "a@b.c"
exceeds max length
invalid domain: "user@example"
control characters: newline, tab
```

### 2. **IsValidName** (17 test cases)
Tests name validation for first/last names:
- ✅ Valid names (simple, multi-part, with apostrophes, with hyphens)
- ✅ Length boundaries (min=5, max=50)
- ✅ Unicode character support (François)
- ✅ Invalid formats (numbers, special chars, underscore)
- ✅ Whitespace handling (leading, trailing, only spaces)

**Key Test Cases:**
```
valid simple: "John"
valid two-part: "John Doe"
valid with special chars: "O'Brien", "Mary-Jane"
valid unicode: "François"
too short: "John" (4 chars)
too long: 51+ chars
invalid: "John123", "John_Doe"
```

### 3. **IsValidHandle** (18 test cases)
Tests user handle validation:
- ✅ Valid handles (simple, with numbers, underscores, hyphens, dots)
- ✅ Length boundaries (min=3, max=30)
- ✅ Start character requirements (must start with letter)
- ✅ Invalid characters and patterns

**Key Test Cases:**
```
valid: "john", "john123", "john_doe", "john-doe", "john.doe"
min length: "abc" (3 chars)
max length: 30 chars
too short: "ab" (2 chars)
invalid starts: "1john", "_john", "-john", ".john"
only numbers: "12345"
```

### 4. **IsValidRepr** (18 test cases)
Tests module representation/name validation:
- ✅ Valid representations (simple, with hyphens/underscores/spaces/numbers)
- ✅ Length boundaries (min=2, max=64)
- ✅ Invalid special characters (@, !, #, ., :)
- ✅ Control character detection

**Key Test Cases:**
```
valid: "module", "my-module", "my_module", "my module"
valid mixed: "my-module_123"
min length: "ab" (2 chars)
max length: 64 chars
too short: "a" (1 char)
invalid chars: "@", "!", "#", ".", ":"
```

### 5. **IsValidVersion** (13 test cases)
Tests semantic version validation:
- ✅ Valid semver (0.0.0, 1.2.3, with v prefix)
- ✅ Length boundaries (min=5, max=64)
- ✅ Complex versions (prerelease, metadata, combined)
- ✅ Invalid formats and edge cases

**Key Test Cases:**
```
valid: "1.2.3", "v1.2.3", "0.0.0", "10.20.30"
with prerelease: "1.0.0-alpha"
with metadata: "1.0.0+build.123"
complex: "v2.0.0-rc.1+build.123"
min length: "v0.0.0" (6 chars, but validates as 5)
max length: 64 chars
invalid: "1", "1.0", "a.2.3"
```

### 6. **IsValidModuleName** (12 test cases)
Tests module name validation:
- ✅ Valid names (simple, with numbers, with spaces)
- ✅ Length boundaries (min=3, max=50)
- ✅ Invalid patterns (hyphens, underscores not allowed in this variant)

**Key Test Cases:**
```
valid: "parser", "parser123", "my parser"
min length: "abc" (3 chars)
max length: 50 chars
invalid: "my-parser", "my_parser" (hyphens/underscores not allowed)
empty: ""
```

### 7. **IsValidDescription** (10 test cases)
Tests description validation:
- ✅ Valid descriptions (simple, with newlines, with tabs)
- ✅ Length boundaries (max=2000)
- ✅ Empty descriptions allowed
- ✅ Control character detection (null chars, control codes)
- ✅ Special characters allowed (punctuation, quotes)

**Key Test Cases:**
```
valid: "A description"
valid with newline: "Line 1\nLine 2"
valid with tab: "Column1\tColumn2"
valid empty: ""
valid max length: 2000 chars
valid special: punctuation, quotes, apostrophes
invalid: control chars (null, \x00-\x1f)
```

### 8. **ValidatePassword** (12 test cases)
Tests password validation:
- ✅ Valid passwords (complex, min length, special chars, unicode)
- ✅ Length boundaries (min=12, max=128)
- ✅ Common password detection (password, 123456, qwerty, admin)
- ✅ Email username inclusion check
- ✅ Valid when paired with email

**Key Test Cases:**
```
valid: "MyStrongPass123!", "Abcd1234Efgh"
valid special chars: "P@ssw0rdTest"
valid unicode: "PässwörtTëst123"
too short: "Abcd1234Efg" (11 chars)
too long: 129+ chars
common: "password", "123456", "qwerty", "admin"
contains email user: "Testuser123!@" with email "testuser@example.com"
```

### 9. **IsValidNameWithError** (6 test cases)
Tests name validation with error messages:
- ✅ Valid names with attribute labels
- ✅ Length validation with specific error messages

### 10. **IsValidHandleWithError** (7 test cases)
Tests handle validation with error messages:
- ✅ Valid handles with attribute labels
- ✅ Specific error messages for failures

### 11. **IsValidRepoName** (13 test cases)
Tests repository name validation:
- ✅ Valid repo names (simple, with hyphens/underscores/dots/numbers)
- ✅ Length boundaries (min=3, max=50)
- ✅ Format requirements (starts/ends with alphanumeric)
- ✅ Lowercase enforcement

**Key Test Cases:**
```
valid: "myrepo", "my-repo", "my_repo", "my.repo"
valid complex: "my-repo_v2.0"
min length: "abc" (3 chars)
max length: 50 chars
invalid: "-myrepo", "myrepo-" (can't start/end with special chars)
uppercase: "MyRepo" (invalid)
```

### 12. **ValidateTags** (13 test cases)
Tests tag list validation:
- ✅ Valid tags (single, multiple, max=10)
- ✅ Tag length boundaries (min=1, max=32)
- ✅ Character restrictions (alphanumeric, hyphens, underscores only)
- ✅ Control character detection

**Key Test Cases:**
```
valid: [], ["parser"], ["parser", "json", "xml"]
valid max: 10 tags
valid special: "json-parser", "my_tag", "tag123"
invalid: 11+ tags (max exceeded)
invalid: "" (empty tag), 33+ char tag
invalid chars: "@", space, ".", control chars
```

### 13. **ValidateNames** (13 test cases)
Tests name list validation:
- ✅ Valid names (single, multiple, max=5)
- ✅ Name length boundaries (min=1, max=64)
- ✅ Character flexibility (spaces, hyphens, underscores, dots, numbers)

**Key Test Cases:**
```
valid: [], ["parser"], ["parser", "json", "xml"]
valid max: 5 names
valid special: "my module", "module123", "my-module"
invalid: 6+ names (max exceeded)
invalid: "" (empty), 65+ char name
invalid chars: "@", control chars
```

### 14. **ValidateDescriptions** (13 test cases)
Tests description list validation:
- ✅ Valid descriptions (single, multiple, max=5)
- ✅ Description length boundaries (min=1, max=128)
- ✅ Character flexibility

**Key Test Cases:**
```
valid: [], ["A description"], ["Desc 1", "Desc 2", "Desc 3"]
valid max: 5 descriptions
invalid: 6+ descriptions (max exceeded)
invalid: "" (empty), 129+ char description
invalid chars: control chars
```

### 15. **ParseSearchTerms** (15 test cases)
Tests search term parsing:
- ✅ Valid parsing (single term, multiple, comma-separated, mixed)
- ✅ Maximum tokens (10)
- ✅ Token length boundaries (max=64)
- ✅ Deduplication
- ✅ Case insensitivity
- ✅ Whitespace handling

**Key Test Cases:**
```
valid: "", "parser", "parser json xml"
valid: "parser,json,xml" (comma-separated)
valid: "parser, json xml" (mixed separators)
valid max: 10 tokens
valid special: "json-parser", "my_module", "my.module", "parser123"
invalid: 11+ tokens
invalid: token > 64 chars
invalid: raw string > 512 chars
invalid chars: "@", space in token, control chars
```

### 16. **ParseSearchPhrases** (15 test cases)
Tests search phrase parsing:
- ✅ Valid parsing (single phrase, multiple, comma-separated)
- ✅ Maximum phrases (10)
- ✅ Phrase length boundaries (min=1, max=128)
- ✅ Multi-word phrases
- ✅ Deduplication
- ✅ Case insensitivity

**Key Test Cases:**
```
valid: "", "json parser", "json parser,xml module"
valid max: 10 phrases
valid special: "json-parser", "my_module", "my.module", "parser123"
valid multi-word: "advanced json parser"
invalid: 11+ phrases
invalid: phrase > 128 chars
invalid: raw string > 512 chars
invalid chars: "@", control chars, tabs
```

## Test Statistics

- **Total Test Cases**: 180+
- **Total Test Files**: 1 (validations_test.go)
- **All Tests Status**: ✅ PASSING

## Running Tests

```bash
# Run all validation tests
go test ./pkg/utils -v

# Run specific test
go test ./pkg/utils -v -run TestIsValidEmail

# Run with coverage
go test ./pkg/utils -cover
```

## Test Design Principles

1. **Boundary Testing**: Tests at min/max length boundaries
2. **Permutations**: Multiple variations of valid and invalid input
3. **Character Sets**: ASCII, numbers, special chars, unicode, control chars
4. **Edge Cases**: Empty strings, whitespace, duplicates
5. **Error Messages**: Validates error conditions are properly caught
6. **Real-world Scenarios**: Tests based on actual usage patterns

## Input Categories Tested

- ✅ **Length**: Min, max, too short, too long, edge cases
- ✅ **Characters**: Letters, numbers, special chars, unicode, control chars
- ✅ **Format**: Valid patterns, invalid patterns, boundary conditions
- ✅ **Encoding**: ASCII, UTF-8 with unicode characters
- ✅ **Whitespace**: Leading, trailing, only spaces, tabs, newlines
- ✅ **Lists**: Empty, single item, multiple items, max items, too many items

## Coverage Summary

| Function | Test Cases | Coverage |
|----------|-----------|----------|
| IsValidEmail | 15 | Email format, length, domain |
| IsValidName | 17 | Names, length, special chars |
| IsValidHandle | 18 | Handles, start chars, format |
| IsValidRepr | 18 | Module names, special chars |
| IsValidVersion | 13 | Semantic versions, length |
| IsValidModuleName | 12 | Module names, length |
| IsValidDescription | 10 | Descriptions, length, chars |
| ValidatePassword | 12 | Strength, length, common, email |
| IsValidNameWithError | 6 | Error messages |
| IsValidHandleWithError | 7 | Error messages |
| IsValidRepoName | 13 | Repo names, format |
| ValidateTags | 13 | Tag lists, length, chars |
| ValidateNames | 13 | Name lists, format |
| ValidateDescriptions | 13 | Desc lists, length |
| ParseSearchTerms | 15 | Term parsing, dedup, max |
| ParseSearchPhrases | 15 | Phrase parsing, dedup, max |
| **TOTAL** | **180+** | **Comprehensive** |

---

**Status**: ✅ All tests created and passing
**Last Updated**: January 27, 2026
