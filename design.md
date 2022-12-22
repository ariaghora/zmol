# WIP - Many things are missing and will change lol, as usual

```
-- This is a comment

-- Variable definition
x = 10
print(x)

-- Just refer to the identifier to access a variable
double = x * 2
print(double)

-- List
a_list = [1, 2, 3]

-- Creating a list with range
numbers = [1 ... high]
numbers = [1 ... high, 2]
```

## Conditional

If-elif-else statement as you expect

```
if n_wheels == 4 { 
    print("maybe a car")
    print("maybe a truck as well")
}
elif n_wheels == 2 { print("maybe a motorcycle") }
elif n_wheels == 3 { print("maybe a bubble car") }
else {
    print("I don't know what it is")
    print("I need to learn more")
}
```

# Loops
```
iter [0 ... 10] as i {
    print(i)    
}

iter [0 ... 10] {
    println("blah")
}
```

Sentinel loop plan

```
iter {
    -- do something
}
```

## Functions
```
-- Define a function like this. No return statement.
-- Any last expression in the block is the returned value.
is_even = @(x) {
    x % 2 == 0
}

-- Recursion
fib = @(n) {
    n <= 1 ? n : fib(n-1) + fib(n-2)
} 
```

## Special operators
Functional taste, maybe?
Special operators, operating on iterables, lowest precedence:
- `->` Element-wise map
- `>-` Element-wise filter
- `|>` Pipe: use the evaluated expression on LHS as the input of function on RHS
```
result = numbers
    -> scale(3)
    -> minus(1)
    >- is_even
    |> sum
result |> print
```

Parentheses are optional if the function definition has no parameter.

### Project Euler problem 1
```
is_divisible = @(x) { (x % 3 == 0) && (x % 5 == 0) }
[1 ... 1000] >- is_divisible |> sum |> print
```

### Guessing game
```
answer = 42
while (let in = input() -> int) != answer {
    let message = if in > answer { "too high" }
    elif in < answer { "too low" }
    else { "that's the correct answer!" }
    print(message)
}
```