# Git Setup Guide

Your local directory is now initialized as a Git repository. Here's how to connect it to your GitHub repo.

## Step 1: Connect to Your GitHub Repo

Your GitHub repo is: `https://github.com/m3owle0/mmcs`

Run these commands in your terminal:

```bash
cd "c:\Users\puppiesandkittens\Downloads\mmcs revamp"

# Add your GitHub repo as remote
git remote add origin https://github.com/m3owle0/mmcs.git

# Or if you prefer SSH (if you have SSH keys set up):
# git remote add origin git@github.com:m3owle0/mmcs.git
```

## Step 2: Check What's Already in GitHub

First, see what's in your GitHub repo:

```bash
# Fetch from GitHub
git fetch origin

# See what branch GitHub has (usually 'main' or 'master')
git branch -r
```

## Step 3: Add and Commit Your Files

```bash
# Add all files
git add .

# Or add specific files:
git add notifier/.github/workflows/notifier.yml
git add CNAME
git add index.html
git add *.svg

# Commit
git commit -m "Add GitHub Actions notifier workflow and deployment files"
```

## Step 4: Push to GitHub

**If your GitHub repo is empty or you want to replace everything:**

```bash
# Push to main branch (or 'master' if that's your default)
git push -u origin main

# If GitHub uses 'master' instead:
# git push -u origin master
```

**If your GitHub repo already has files:**

```bash
# Pull first to merge any existing files
git pull origin main --allow-unrelated-histories

# Then push
git push -u origin main
```

## Step 5: Verify

1. Go to: https://github.com/m3owle0/mmcs
2. Check that your files are there
3. Go to: https://github.com/m3owle0/mmcs/actions
4. You should see the workflow running automatically

## Troubleshooting

### "Remote origin already exists"
```bash
# Remove existing remote
git remote remove origin

# Add it again
git remote add origin https://github.com/m3owle0/mmcs.git
```

### "Authentication failed"
- Make sure you're logged into GitHub
- You may need to use a Personal Access Token instead of password
- Or set up SSH keys

### "Branch name mismatch"
- Check what branch GitHub uses: `git branch -r`
- Use that branch name instead of `main`

---

**Quick Command Summary:**

```bash
cd "c:\Users\puppiesandkittens\Downloads\mmcs revamp"
git remote add origin https://github.com/m3owle0/mmcs.git
git add .
git commit -m "Add GitHub Actions workflow"
git push -u origin main
```
