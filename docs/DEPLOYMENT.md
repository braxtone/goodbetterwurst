# Deployment & Release Guide

This document explains how to use the automated versioning and release workflow for the GoodBetterWurst project.

## Overview

The GitHub workflow at [`.github/workflows/go.yml`](.github/workflows/go.yml) automatically:
- ✅ Builds and tests Go code
- ✅ Auto-increments patch version (v0.1.0 → v0.1.1 → v0.1.2...)
- ✅ Tags Docker images with semantic version AND `latest`
- ✅ Creates GitHub releases with auto-generated notes

## Prerequisites

### 1. Docker Hub Secrets
Ensure these secrets are configured in your GitHub repository settings:
- `DOCKER_HUB_USERNAME` - Your Docker Hub username
- `DOCKER_HUB_PASSWORD` - Your Docker Hub access token

**To check/add secrets:**
1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Verify `DOCKER_HUB_USERNAME` and `DOCKER_HUB_PASSWORD` are present

### 2. Initial Version Tag (REQUIRED BEFORE FIRST RUN)

Before the first workflow run, you must create the initial `v0.1.0` tag:

```bash
# From the repository root
git checkout main
git pull origin main
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

**Why?** The semantic versioning action needs an existing tag to calculate the next version. Without this initial tag, the workflow will fail.

## How It Works

### Workflow Trigger
- **Pull Requests to main**: Runs build and test only (no versioning/deployment)
- **Pushes to main**: Runs full pipeline (build → test → version → deploy → release)

### Job Flow
```
┌─────────────────┐
│ build-and-test  │ ← Runs on PRs and pushes
└────────┬────────┘
         │ (only on push to main)
         ↓
┌─────────────────┐
│ version-and-tag │ ← Creates v0.1.x tag
└────┬────────┬───┘
     │        │
     ↓        ↓
┌────────┐  ┌─────────────┐
│ docker │  │ create-     │
│ push   │  │ release     │
└────────┘  └─────────────┘
```

### Version Calculation
- **Base**: v0.1.x (pre-release/beta)
- **Increment**: Patch version auto-increments on each push to main
- **Format**: `v{major}.{minor}.{patch}`

Examples:
- First push: `v0.1.0` → `v0.1.1`
- Second push: `v0.1.1` → `v0.1.2`
- Third push: `v0.1.2` → `v0.1.3`

### Docker Tags
Each deployment creates two Docker tags:
- **Semantic version**: `braxtone/goodbetterwurst:v0.1.5`
- **Latest**: `braxtone/goodbetterwurst:latest`

### GitHub Releases
- **Title**: Version tag (e.g., `v0.1.5`)
- **Body**: Auto-generated from commit messages since last release
- **Type**: Full release (not draft or pre-release)

## Deployment Process

### Standard Deployment (Recommended)
1. Create a feature branch and make your changes
2. Open a pull request to `main`
3. Wait for tests to pass
4. Merge the PR to `main`
5. **Automatic actions happen:**
   - Tests run again
   - Version increments (e.g., v0.1.5 → v0.1.6)
   - Git tag created
   - Docker image built and pushed
   - GitHub release created

### Direct Push to Main (Not Recommended)
If you push directly to main (bypassing PRs):
```bash
git checkout main
git pull origin main
# Make changes
git add .
git commit -m "Your change description"
git push origin main
```
The same automatic versioning and deployment will occur.

## Monitoring Deployments

### View Workflow Status
1. Go to your GitHub repository
2. Click the "Actions" tab
3. Select the latest "Go CI/CD" workflow run
4. Check each job's status

### Check Docker Tags
Visit Docker Hub: https://hub.docker.com/r/braxtone/goodbetterwurst/tags

You should see:
- `latest` (most recent deployment)
- `v0.1.1`, `v0.1.2`, `v0.1.3`, etc. (all versions)

### View Releases
1. Go to your GitHub repository
2. Click "Releases" in the right sidebar
3. See all releases with auto-generated notes

## Using the Docker Image

### Pull Latest Version
```bash
docker pull braxtone/goodbetterwurst:latest
```

### Pull Specific Version
```bash
docker pull braxtone/goodbetterwurst:v0.1.5
```

### Run Container
```bash
docker run -p 8080:8080 braxtone/goodbetterwurst:latest
```

## Troubleshooting

### Workflow Fails with "No tags found"
**Problem**: Initial tag `v0.1.0` doesn't exist.

**Solution**: Create the initial tag:
```bash
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

### Docker Push Fails
**Problem**: Docker Hub credentials are missing or incorrect.

**Solution**: 
1. Verify secrets in Settings → Secrets and variables → Actions
2. Regenerate Docker Hub access token if needed
3. Update `DOCKER_HUB_PASSWORD` secret

### Release Creation Fails
**Problem**: Insufficient permissions.

**Solution**: The workflow uses `GITHUB_TOKEN` which is automatically provided. Ensure the repository settings allow GitHub Actions to create releases:
1. Go to Settings → Actions → General
2. Scroll to "Workflow permissions"
3. Select "Read and write permissions"

### Version Doesn't Increment
**Problem**: Multiple pushes create the same version.

**Solution**: The workflow includes concurrency controls to prevent this. If it still happens:
1. Check that tags are being pushed correctly
2. Ensure no manual tags were created with the same version
3. Review workflow logs for errors

### Tests Fail but Want to Deploy Anyway
**Problem**: Tests are failing but you need to deploy urgently.

**Solution**: Not recommended, but you can:
1. Fix the tests first (strongly recommended)
2. OR temporarily skip tests (modify workflow to not depend on build-and-test job)
3. OR manually create a tag and release

## Advanced Configuration

### Change Base Version (e.g., to v1.0.x)
When ready for production (v1.0.x):

1. Create a new tag manually:
```bash
git tag -a v1.0.0 -m "Production release"
git push origin v1.0.0
```

2. The workflow will automatically continue from v1.0.0 → v1.0.1 → v1.0.2...

### Manual Version Bump (Major or Minor)
To manually bump major or minor version:

```bash
# For minor version bump (v0.1.x → v0.2.0)
git tag -a v0.2.0 -m "Minor version bump"
git push origin v0.2.0

# For major version bump (v0.x.x → v1.0.0)
git tag -a v1.0.0 -m "Major version bump"
git push origin v1.0.0
```

The workflow will continue incrementing patches from the new base.

### Disable Deployments Temporarily
To disable deployments without removing the workflow:

1. Go to Settings → Environments
2. Add protection rules requiring approval
3. OR comment out the `docker-push` and `create-release` jobs in the workflow

## Rollback Procedures

### Rollback Docker Deployment
Pull and deploy a previous version:
```bash
docker pull braxtone/goodbetterwurst:v0.1.4
docker run -p 8080:8080 braxtone/goodbetterwurst:v0.1.4
```

### Delete Erroneous Tag
If a bad tag was created:
```bash
# Delete local tag
git tag -d v0.1.5

# Delete remote tag
git push --delete origin v0.1.5
```

### Delete Erroneous Release
1. Go to GitHub Releases
2. Click the release to delete
3. Click "Delete" button
4. Note: Docker images cannot be deleted, only overwritten

## Concurrency & Safety

The workflow includes concurrency controls:
```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: false
```

This ensures:
- Only one deployment per branch runs at a time
- Concurrent pushes wait for previous deployment to complete
- No version conflicts from simultaneous pushes

## Best Practices

1. **Always use Pull Requests** - Test changes before merging to main
2. **Write descriptive commit messages** - They appear in release notes
3. **Monitor workflow runs** - Check Actions tab after each merge
4. **Verify deployments** - Check Docker Hub and Releases after deployment
5. **Tag major releases manually** - Use manual tags for v1.0.0, v2.0.0, etc.
6. **Keep main branch protected** - Require PR reviews and status checks

## Need Help?

- Review workflow logs in Actions tab
- Check the planning document at [`plans/github-workflow-versioning-plan.md`](plans/github-workflow-versioning-plan.md)
- Verify Docker Hub credentials
- Ensure initial v0.1.0 tag exists
