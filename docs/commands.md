# Commands
This file is mostly for notes on the commands I wish to implement.

parago -j/--parallel N      # Obvious, the level of parallelism.
parago --dry-run            # Don't run any tasks, just pretend
parago --relax              # Need a better name. Keeps going with
                            # non-dependent tasks if one fails, instead
                            # of the default behavior of aborting the entire
                            # run.
