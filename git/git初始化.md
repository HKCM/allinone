# git初始化

Git global setup

```bash
git config --global user.name "karl"
git config --global user.email "karl@mail.com"

# 或
git config user.name "karl"
git config user.email "karl@mail.com"
```

Create a new repository

```bash
git clone git@gitlab.okcoin.tokyo:devops/sysops/task/staging-cdk-setup.git
cd staging-cdk-setup
git switch --create main
touch README.md
git add README.md
git commit -m "add README"
git push --set-upstream origin main
```

Push an existing folder

```bash
cd existing_folder
git init --initial-branch=main
git remote add origin git@gitlab.okcoin.tokyo:devops/sysops/task/staging-cdk-setup.git
git add .
git commit -m "Initial commit"
git push --set-upstream origin main
```

Push an existing Git repository

```bash
cd existing_repo
git remote rename origin old-origin
git remote add origin git@gitlab.okcoin.tokyo:devops/sysops/task/staging-cdk-setup.git
git push --set-upstream origin --all
git push --set-upstream origin --tags
```

