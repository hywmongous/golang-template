#!/bin/bash

# If working on a "feature/x" the latest commit from
# "develop" is merged into the "feature/x" you are working on.
# However, if you are working on "develop" the latest commit
# from "release/x" is merged into the "develop" branch.
# Otherwise, if you are working on a "release/x" branch
# then the latest commit from "develop" is merged into it.
