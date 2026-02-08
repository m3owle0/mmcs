# Git Setup Instructions

## Step 1: Create a GitHub Repository

1. Go to https://github.com/new
2. Create a new repository (name it `mmcs` or whatever you prefer)
3. **Don't** initialize with README, .gitignore, or license (we already have these)
4. Click "Create repository"

## Step 2: Add Remote and Push

After creating the repository, GitHub will show you commands. Use these:

### Option A: If you haven't initialized git yet:

```powershell
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git
git push -u origin main
```

### Option B: If git is already initialized:

```powershell
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git
git push -u origin main
```

**Important:** Replace `YOUR_USERNAME` and `YOUR_REPO_NAME` with your actual GitHub username and repository name.

## Example:

If your GitHub username is `johndoe` and your repo is named `mmcs`, the command would be:

```powershell
git remote add origin https://github.com/johndoe/mmcs.git
```

## Troubleshooting

### If you get "remote origin already exists":
```powershell
git remote remove origin
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git
```

### If you get authentication errors:
- Use GitHub Personal Access Token instead of password
- Or use GitHub Desktop app for easier authentication
