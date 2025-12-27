#!/bin/bash
# Initialize versioning for automated releases
# This script creates the initial v0.1.0 tag required by the workflow

set -e

echo "🚀 Initializing automated versioning for goodbetterwurst..."
echo ""

# Check if we're on main branch
current_branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$current_branch" != "main" ]; then
    echo "⚠️  Warning: You're on branch '$current_branch', not 'main'"
    echo "   The workflow only runs on the main branch."
    read -p "   Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

# Check if tag already exists
if git rev-parse v0.1.0 >/dev/null 2>&1; then
    echo "✅ Tag v0.1.0 already exists!"
    echo ""
    echo "   Local:  $(git rev-parse v0.1.0)"
    if git ls-remote --tags origin | grep -q "v0.1.0"; then
        echo "   Remote: $(git ls-remote --tags origin | grep v0.1.0 | awk '{print $1}')"
        echo ""
        echo "✨ Your repository is already initialized for automated versioning."
        echo "   Next push to main will create v0.1.1"
        exit 0
    else
        echo "   Remote: Not found"
        echo ""
        read -p "📤 Push existing tag to remote? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git push origin v0.1.0
            echo "✅ Tag v0.1.0 pushed to remote!"
            echo "✨ Your repository is now ready for automated versioning."
            exit 0
        else
            echo "⚠️  Tag exists locally but not on remote. Workflow may not work correctly."
            exit 1
        fi
    fi
fi

# Make sure we're up to date
echo "📥 Fetching latest changes..."
git fetch --all

# Check if remote tag exists but not local
if git ls-remote --tags origin | grep -q "v0.1.0"; then
    echo "✅ Tag v0.1.0 exists on remote but not locally"
    echo "   Fetching tag..."
    git fetch --tags
    echo "✅ Tag v0.1.0 fetched!"
    echo "✨ Your repository is ready for automated versioning."
    exit 0
fi

echo "📝 Creating initial tag v0.1.0..."
echo ""
echo "   This tag will be the starting point for automated versioning."
echo "   Future pushes to main will auto-increment: v0.1.0 → v0.1.1 → v0.1.2 ..."
echo ""

# Create the tag
git tag -a v0.1.0 -m "Initial release - automated versioning enabled"

echo "✅ Tag v0.1.0 created locally!"
echo ""
read -p "📤 Push tag to remote now? (Y/n): " -n 1 -r
echo

if [[ $REPLY =~ ^[Nn]$ ]]; then
    echo ""
    echo "⚠️  Tag created locally but NOT pushed to remote."
    echo "   To push it later, run:"
    echo "   git push origin v0.1.0"
    exit 0
fi

# Push the tag
echo "📤 Pushing tag to remote..."
git push origin v0.1.0

echo ""
echo "✅ Success! Your repository is now ready for automated versioning."
echo ""
echo "📋 Next steps:"
echo "   1. Make changes to your code"
echo "   2. Commit and push to main (or merge a PR)"
echo "   3. The workflow will automatically:"
echo "      • Run tests"
echo "      • Create tag v0.1.1"
echo "      • Build and push Docker image"
echo "      • Create GitHub release"
echo ""
echo "📚 For more information, see: DEPLOYMENT.md"
