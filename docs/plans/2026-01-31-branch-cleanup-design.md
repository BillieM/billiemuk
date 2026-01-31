# Branch Cleanup and Post Removal Design

**Goal:** Clean up git branches, update GitHub default branch, and remove example post.

## Summary

We're doing three things:
1. Change GitHub's default branch from `deploy-pages` to `main`
2. Delete the `deploy-pages` branch (local and remote) since it's no longer used
3. Remove the hello-world example post

The deployment workflow already runs on `main`, so this is safe. We're just updating GitHub's settings to match the current reality and cleaning up the unused branch.

## Current State

- GitHub default branch: `deploy-pages`
- Actual deployment: Runs from `main` via GitHub Actions
- Branches: `main`, `deploy-pages` (local and remote)
- Posts: `2026-01-27-hello-world.md` (to be removed), `2026-01-31-building-a-blog-that-just-works.md` (keep)

## Branch Cleanup Approach

The safest order of operations:

**1. Change GitHub default branch first**
Use `gh` CLI to change the default branch to `main`. This ensures that if someone visits the repo, they see the active branch.

```bash
gh repo edit --default-branch main
```

**2. Delete remote deploy-pages branch**
Once the default is switched, delete the remote branch. Since it's no longer the default, this is safe.

```bash
git push origin --delete deploy-pages
```

**3. Delete local deploy-pages branch**
Clean up the local branch.

```bash
git branch -D deploy-pages
```

This order prevents any issues where GitHub might reject deleting the default branch.

## Post Removal

Delete the hello-world post file:

```bash
rm content/posts/2026-01-27-hello-world.md
```

Then rebuild the site to verify it's gone:

```bash
go run . build
```

The build process cleans the `dist/` directory before building (`internal/builder/builder.go:31`), so the HTML file will be automatically removed.

## Verification

After making the changes, verify:

1. **GitHub default branch** - Check with `gh repo view --json defaultBranchRef`
2. **Branches deleted** - Run `git branch -a` to confirm no deploy-pages branches
3. **Post removed** - Check the built site in `dist/` and run the dev server to confirm hello-world is gone

Then commit the post deletion and push to main.

## Implementation Steps

1. Change GitHub default branch to main
2. Delete remote deploy-pages branch
3. Delete local deploy-pages branch
4. Delete hello-world markdown file
5. Rebuild site and verify
6. Commit and push changes
