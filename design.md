
```
-- This is a comment

-- Variable definition
let x = 10
print(x)

-- Just refer to the identifier to access a variable
let double = x * 2
print(double)

-- List
let a_list = [1, 2, 3]

-- Creating a list with range
let numbers = [1 ... high]
let numbers = [1 ... high, 2]

```

## Conditional
- Basic form: `<BOOL_EXP> ? <EXPR_WHEN_TRUE> | <EXPR_WHEN_FALSE>`
```
number = 5
result = if (number % 3 == 0) && (number % 5 == 0): "fizz" | "buzz" 

-- or in multi-line form

result = if (number % 3 == 0) && (number % 5 == 0):
    "fizz"
| 
    "buzz" 
end

-- or like this, when the blocks have one statement

n_wheels = input() -> int()

if n_wheels == 4: print("maybe a car")
| n_wheels == 2: print("maybe a motorcycle")
| n_wheels == 3: print("maybe a bubble car")
| print("I don't know what it is")
end

-- This is also fine:

if n_wheels == 4: 
    print("maybe a car")
    print("maybe a truck as well")
| n_wheels == 2: print("maybe a motorcycle")
| n_wheels == 3: print("maybe a bubble car")
|
    print("I don't know what it is")
    print("I need to learn more")
end
```

# Loops
```
for i in [0 ... 10]:
    print(i)    
end

while true:
    -- do something
end
```

## Functions
```
-- Define a function like this. No return statement.
-- Any last expression in the block is the returned value.
let is_even = @(x):
    x % 2 == 0
end

-- One-line function
let add = @(x, y): x + y

-- Recursion
let fib = @(n): n > 1 ? n + fib(n-1) | 1
```

## Special operators
Functional taste, maybe?
Special operators, operating on iterables, lowest precedence:
- `|>` Element-wise map
- `>-` Element-wise filter
- `->` Pipe: use the output of LHS as the input of function on RHS
```
let result = numbers
    |> scale(3)
    |> minus(1)
    >- is_even
    -> sum
    -> print
```

Parentheses are optional if the function definition has no parameter.

### Project Euler problem 1
```
let is_divisible = @(x): (x % 3 == 0) && (x % 5 == 0)
[1 ... 1000] >- is_divisible -> sum -> print
```

### Guessing game
```
let answer = 42
while (let in = input() -> int) != answer:
    let message = if in > answer: "too high"
    | in < answer: "too low"
    | "that's the correct answer!"
    end
    print(message)
end
```