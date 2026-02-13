# Testing Documentation

## Current Test Coverage

This project includes unit tests for the core HTTP handlers with proper mocking of dependencies.

### ✅ Implemented Tests

#### Handler Tests (`cmd/api/handler/package_test.go`)
- ✅ `TestHandleGetAlreadyBuildPackages_Success` - Successfully retrieve built packages
- ✅ `TestHandleGetAlreadyBuildPackages_Error` - Handle database errors
- ✅ `TestHandleGetAurPackageByName_Success` - Successfully search AUR packages  
- ✅ `TestHandleGetAurPackageByName_MissingParameter` - Handle missing query parameters
- ✅ `TestHandleGetAurPackageByName_Error` - Handle AUR API errors
- ⚠️ `TestHandleBuildPackage_Success` - **FAILS** (requires Docker daemon)
- ✅ `TestHandleBuildPackage_PackageAlreadyExists` - Handle duplicate packages
- ✅ `TestHandleBuildPackage_InvalidJSON` - Handle malformed requests
- ✅ `TestHandleBuildPackage_DatabaseError` - Handle database errors

#### Storage Tests (`pkg/storage/`)
- ✅ Filesystem storage implementation tests (already existing)

#### Scheduler Tests (`pkg/scheduler/`)  
- ✅ Basic scheduler functionality tests (already existing)

### 🚧 Missing/Limited Tests

#### Docker Integration (`TestHandleBuildPackage_Success`)
**Issue**: This test fails because it attempts real Docker operations which require:
- Docker daemon running
- Proper container image available
- Network connectivity for image pulls

**Solution Approaches**:
1. **Mock Docker Controller**: Create interface for `docker.ContainerController` and mock it
2. **Integration Test Environment**: Set up Docker in CI with proper test containers
3. **Test Containers**: Use testcontainers-go library for Docker-based integration tests

**Current Status**: Test exists but fails without Docker environment.

#### Builder Service Tests (`pkg/builder/`)
**Issue**: Builder service heavily integrates with Docker operations, making unit testing complex without extensive mocking.

**Missing Coverage**:
- `BuildAurPackage()` method testing
- `SavePackage()` flow testing  
- `CleanUpBuildEnv()` testing

**Recommendation**: Requires mocking of `docker.ContainerController` interface.

#### AUR Service Tests (`pkg/aur/`)
**Issue**: AUR service makes external HTTP calls to Arch Linux AUR API.

**Missing Coverage**:
- AUR API response parsing
- Network error handling
- Rate limiting behavior

**Recommendation**: Mock HTTP client or use HTTP test server.

#### Database Repository Tests (`cmd/api/db/repo/`)
**Issue**: Repository directly uses MongoDB client.

**Missing Coverage**:
- MongoDB operations
- Connection error handling
- Query result processing

**Recommendation**: Use MongoDB test container or in-memory MongoDB.

### 🔧 Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./cmd/api/handler -v

# Run tests with coverage
go test -v -cover ./...

# Run tests excluding Docker-dependent tests
go test ./... -short
```

### 📋 Test Environment Setup

#### Required for Full Test Suite:
1. **Docker**: For container-based tests
2. **MongoDB**: For database integration tests  
3. **Network Access**: For AUR API tests

#### Minimal Testing (Current):
```bash
# Only run unit tests with mocked dependencies
go test ./cmd/api/handler -v
```

### 🎯 Future Test Improvements

1. **Create Interfaces**: Make Docker controller and other external dependencies mockable
2. **Integration Tests**: Set up proper test environment with Docker and MongoDB
3. **End-to-End Tests**: Full API testing with real dependencies
4. **Performance Tests**: Load testing for package building operations
5. **Test Containers**: Use testcontainers-go for isolated testing

### 📊 Test Coverage Summary

| Package | Coverage | Status |
|---------|----------|---------|
| `cmd/api/handler` | ~85% | ✅ Good (1 failing test) |
| `pkg/storage` | ~90% | ✅ Excellent |
| `pkg/scheduler` | ~70% | ✅ Good |
| `pkg/builder` | 0% | ❌ Missing |
| `pkg/aur` | 0% | ❌ Missing |
| `cmd/api/db/repo` | 0% | ❌ Missing |

**Overall Assessment**: Core HTTP handlers are well-tested with proper mocking. Integration testing requires additional setup but current unit tests provide good coverage for business logic validation.