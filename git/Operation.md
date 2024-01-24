# 创建空分支

```bash
git checkout --orphan blankbranch
git rm -rf .
echo "# blank branch" > README.md
git add README.md 
git commit -m "new blank branch"
```