# Zmol

A toy language that aims to be:

- The abuser of Fira Code font ligatures
- Concise; you can learn over a weekend
- Able to solve non-trivial problems (like hello world)
- Zmol in size

The Euler Project problem 1 solution in Zmol:

```
sum = @(list): reduce(list, (@(x, y):x+y), 0) 
divisible = @(x): (x % 3 == 0) || (x % 5 == 0)

range_list(1, 1000) >- divisible |> sum |> print
```

output:
```
233168
```