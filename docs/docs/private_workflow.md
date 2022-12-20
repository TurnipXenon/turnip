# Private workflow

Personal notes about how to keep certain parts of the code and content hidden.

We want to be transparent as much as we want, but we also want to keep certain things private. Mostly due to personal
preferences. We have two ways to keep things from leaking:

1. Databases
2. Private branches

## Private branch

We work on four main branches:

1. origin/main
2. origin/wip
3. private/main
4. origin/wip

The branches interact like this:

```
time/steps ->

local/wip      *----------------*--
              1 \            7 / 
origin/main    --*-------*-*--*----
                2 \  ~4 / 5 \     
local/wip      ----*---*-----*-----
                  3 \ /     /6\
private/main   ------*-----*-------
```

**Steps:**

1. We make a PR to origin, and merge our changes from WIP to origin/main.
2. We pull changes from origin/main to local/wip.
3. We make sure we have the latest changes from private. We then make our private changes, make a PR to private and
   merge them to private/main.
4. Step 4 is special because we have to run a script which scrubs the private information from our templates. We squash
   our commits to hide these changes. We make a PR and merge to origin/main.
5. For step 5, we delete local/wip, and make a new local/wip. We merge origin/main to local/wip and handle merge
   conflicts.
6. We merge these resolved conflicts to private/main.
7. Merge the latest changes from origin/main to local/wip. We continue working on our local/wip.

## Future plan

We want this repository to be a plugin to the actual code we have.
