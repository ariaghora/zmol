
```
-- This is a comment

-- Variable definition
@x = 10

-- Just refer to the identifier to access a variable
@double = x * 2
print(double)

-- List
@a_list = [1, 2, 3]

-- Creating a list with range
@numbers = [1 ... high]

-- Define a function like this. No return statement.
-- Any last expression in the block is the returned value.
@scale(x, factor) = x * factor 
@minus(x, amount) = x - amount 

-- Multi-line definition
@is_even(x):
    x % 2 == 0
end
```


Special operators, operating on iterables, lowest precedence:
- `|>` Element-wise map
- `>-` Element-wise filter
- `->` Pipe: use the output of LHS as the input of function on RHS
```
let result = numbers
    |> scale(3)
    |> minus(1)
    >- is_even()
    -> sum()
    -> print()
```
