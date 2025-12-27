# Automated Versioning & Release Setup

This repository now has automated versioning and releases configured! 🎉

## What's Changed

### Files Modified
- **[`.github/workflows/go.yml`](.github/workflows/go.yml)** - Updated workflow with automated versioning

### Files Created
- **[`DEPLOYMENT.md`](DEPLOYMENT.md)** - Complete deployment and usage guide
- **[`plans/github-workflow-versioning-plan.md`](plans/github-workflow-versioning-plan.md)** - Detailed planning document
- **[`scripts/init-versioning.sh`](scripts/init-versioning.sh)** - Initialization script

## Quick Start

### Step 1: Initialize Versioning (ONE TIME ONLY)

Run the initialization script to create the initial `v0.1.0` tag:

```bash
./scripts/init-versioning.sh
```

Or manually:

```bash
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

### Step 2: Commit and Push Changes

Commit the workflow changes to your repository:

```bash
git add .
git commit -m "Add automated versioning and releases"
git push origin main
```

### Step 3: Verify Setup

After pushing, check:
1. **GitHub Actions** - Workflow should run automatically
2. **Docker Hub** - Images should be tagged with `v0.1.1` and `latest`
3. **GitHub Releases** - A new release should be created

## How It Works

```
Push to main → Tests Pass → Auto-increment version → Docker push + GitHub release
```

**Example flow:**
1. Merge PR to main
2. Tests run (Go build & test)
3. Version increments: `v0.1.0` → `v0.1.1`
4. Docker image pushed: `braxtone/goodbetterwurst:v0.1.1` and `latest`
5. GitHub release created with auto-generated notes

## What Happens on Each Push to Main

✅ **Build and Test** - Go code is built and tested  
✅ **Version Increment** - Patch version auto-increments  
✅ **Git Tag** - New tag created (e.g., `v0.1.5`)  
✅ **Docker Push** - Multi-platform image pushed to Docker Hub  
✅ **GitHub Release** - Release created with commit notes  

## Version Format

- **Base**: `v0.1.x` (pre-release/beta)
- **Auto-increment**: Patch version on each push
- **Examples**: `v0.1.0` → `v0.1.1` → `v0.1.2` → ...

## Docker Tags Generated

Each deployment creates TWO tags:
- `braxtone/goodbetterwurst:v0.1.5` (specific version)
- `braxtone/goodbetterwurst:latest` (always points to newest)

## Important Notes

⚠️ **Prerequisites:**
- Docker Hub credentials must be configured as GitHub secrets
- Initial `v0.1.0` tag must exist before first workflow run

⚠️ **Workflow Only Runs On:**
- Pushes to `main` branch (creates versions and deploys)
- Pull requests (runs tests only, no deployment)

⚠️ **Concurrency Control:**
- Only one deployment can run at a time
- Prevents version conflicts from simultaneous pushes

## Documentation

📚 **[DEPLOYMENT.md](DEPLOYMENT.md)** - Full deployment guide including:
- Detailed workflow explanation
- Troubleshooting steps
- Rollback procedures
- Advanced configuration options

📋 **[plans/github-workflow-versioning-plan.md](plans/github-workflow-versioning-plan.md)** - Planning document with:
- Architecture decisions
- Alternative approaches considered
- Risk analysis and mitigations

## Troubleshooting

### "No tags found" Error
**Solution:** Create initial tag with `./scripts/init-versioning.sh`

### Docker Push Fails
**Solution:** Verify `DOCKER_HUB_USERNAME` and `DOCKER_HUB_PASSWORD` secrets

### Release Creation Fails
**Solution:** Enable "Read and write permissions" in Settings → Actions → General

See [DEPLOYMENT.md](DEPLOYMENT.md#troubleshooting) for more troubleshooting help.

## Future Enhancements

Consider adding:
- Conventional commit parsing for smart version bumping
- Pre-release tags for feature branches
- Automated changelog generation
- Slack/Discord notifications

## Need Help?

1. Check [DEPLOYMENT.md](DEPLOYMENT.md) for detailed instructions
2. Review workflow logs in GitHub Actions tab
3. Check planning document at [plans/github-workflow-versioning-plan.md](plans/github-workflow-versioning-plan.md)
