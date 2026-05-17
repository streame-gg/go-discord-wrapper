# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added — Phase A1 (Collections package)

New `collection` package providing a generic `Collection[K, V]` type —
a map+slice hybrid with insertion-order iteration and 30+ utility methods
(`Filter`, `Find`, `Partition`, `Sort`, `Concat`, `Difference`, …), inspired by
[@discordjs/collection](https://github.com/discordjs/discord.js/tree/main/packages/collection).
This is the foundational data structure for the manager-pattern migration
planned in Phase A2.

The `collection` package is currently isolated; no existing API uses it
yet. Migration of cache and REST returns to Collection-based APIs comes
in Phase A2 (planned breaking change).
