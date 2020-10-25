# Random Workout Generator

Simple program for generating a random workout based on weighted options.

## Syntax


### Random Selection

```
NAME { A : X | B : Y }
```

If `NAME` is evaluated, the weights `X` and `Y` are used for random selection. E.g. if `X` is 2 and
`Y` is 1, then `A` is selected 66% of the time, and `B` is selected `33%` of the time. If either `A`
or `B` are listed, then they are evaluated next, if selected.

### Sequence

```
NAME [ { A : X | B : Y }, (C) ]
```

Square brackets indicate expressions are evaluated in order (no random selection, unless a
sub-expression is a random selection, i.e. curly brackets.) For this example, `NAME` will resolve
to `A` then `C` OR `B` then `C`.

### Alias

```
NAME [ (A), (B) ]
```

If a name is surrounded by parentheses, then it is always resolved. For this example, the result
will always be `A` then `B`.

### Full Example

```
workout { run : 2 | swim : 1 }
run { tempo : 2 | LSD : 1 | interval: 5 }
interval [ (lap), (rest) ]
lap { 1/4 mile : 1 | 1/2 mile : 1 }
rest { 2 minutes : 2 | 3 minutes : 2 | 4 minutes : 1 }
```

Some possible results of this plan:

```
workout
  run
    LSD
```

```
workout
  swim
```

```
workout
  run
    interval
      lap
        1/4 mile
      rest
        3 minutes
```
