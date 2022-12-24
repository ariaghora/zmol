# Quick introduction to Zmol

## Hello world
```
-- This is a comment.
-- Make sure to write effective comments to help others understand your code.
println("Hello world")
```

## Variable definition
```
an_integer = 10
a_float = 3.14
a_string = "hello"
a_boolean = true
a_list = [1, 2, 3, "Alice", "Bob"]

-- Function is a first-class citizen
a_function = @(x) { x + 1 }
println(a_function(10))

```

## Primitive types

| Type | Description |
| --- | --- |
| `Int` | Integer, a whole number value that can be positive, negative, or zero. In this programming language, an `Int` is a 32-bit signed integer. |
| `Float`| Floating point number, a numerical data type that can represent fractional values. In this programming language, a `Float` is a 64-bit double-precision floating point number. |
| `Bool`| Boolean, a data type that can have only two values: `true` or `false`. |
| `String`| String, a data type that represents a sequence of characters. Strings can be used to store and manipulate text data. |
| `List`| A container data type that can hold multiple values with any data type. Lists are ordered, mutable (can be modified), and can contain duplicates. They are often used to store and manipulate collections of data. A list can be accessed by an integer index. |
| `Table`| Key-value data structure. A table is a collection of key-value pairs, where the key is used to access the corresponding value. Tables are unordered, mutable, and do not allow duplicate keys. They are often used to store and manipulate data that needs to be quickly retrieved using a unique key. |
| `Function`| Functions in Zmol are first-class citizens, which means that they can be assigned to variables, passed as arguments to other functions, and returned as values from functions.|

## Functions

### Function definition

Following is how we define a function.
There is no return statement in Zmol.
The last expression in the block is the returned value.

```

is_even = @(x) {
    x % 2 == 0
}
```

Functions are first-class citizens. They can be passed as arguments to other functions.
```
add = @(x, y) { x + y }
sub = @(x, y) { x - y }

calculate = @(x, y, func) { func(x, y) }

print(calculate(10, 5, add))
print(calculate(10, 5, sub))
```

Anonymous functions are also supported.
```
println(calculate(10, 5, @(x, y) { x * y }))
```

## Built-in functions

| Function | Description |
| --- | --- |
| `print`, `println` | Prints the given value to the standard output. |
| `input` | Reads a line from the standard input. |
| `len` | Returns the length of the given list or string. |
| `type` | Returns the type of the given value. |

### Iterable-related functions
| Function | Description |
| --- | --- |
| `range_list` | Returns a list of integers in the given range. |
| `map` | Applies the given function to each element in the given list. |
| `filter` | Returns a list of elements that satisfy the given predicate. |
| `reduce` | Reduces the given list to a single value using the given function. |
| `zip` | Returns a list of lists, where the i-th list contains the i-th element from each of the argument lists. |
| `split` | Splits the given string into a list of strings using the given delimiter. |

## Conditional statements

If-else statement as you expect, parentheses are not required.

```
if n_wheels == 4 { 
    println("maybe a car")
    println("maybe a truck as well")
} else if n_wheels == 2 {
    println("maybe a motorcycle") 
} else if n_wheels == 3 { print("maybe a bubble car") 
} else {
    println("I don't know what it is")
    println("I need to learn more")
}
```

# Loops
Loops primarily iterate over iterables, such as lists, tables, and strings.
The loop statement uses `iter` keyword, followed by the iterable to be iterated over.
We can use `as` keyword to bind the current value to a variable.
```
iter [1, 2, 3] as i {
    println(i)    
}
```
(WIP) Sometimes we don't need the current value. We can skip the `as` keyword. 
```
-- This one is WIP
iter [0 ... 10] {
    println("blah")
}
```

## (WIP) Sentinel loop 
For infinite loops, we can use following syntax:

```
iter {
    -- do something

    if break_condition_is_true {
        break
    }
}
```

## Operators

| Operator | Description |
| --- | --- |
| `=` | Assignment operator |
| `+ - * / %` | Arithmetic operators |
| `== != < > <= >=` | Comparison operators |
| `&& ||` | Logical AND and OR operators |
| `!` | Logical negation operator |
| `[]` | Indexing operator |
| `@` | Function definition operator |

The `+` operator can be used to concatenate strings and lists

```
println("hello" + "world")
println([1, 2, 3] + [4, 5, 6])
```


## Special operators

Special operators will give us some functional programming taste, maybe? They enable function composition in an infix style.

Suppose we have following functions:
```
scale = @(x, factor) {
    x * factor
}

is_even = @(x) {
    x % 2 == 0
}

sum = @(list) {
    reduce(list, @(x, y) { x + y }, 0) 
}

```

We can use following operators:


| Operator | Description | Example |
| --- | --- | --- |
| `->` | Element-wise map: Applies a function to each element of a list or array and returns a new list or array with the modified elements. | `[1, 2, 3] -> scale{2}` returns `[2, 4, 6]` |
| `>-` | Element-wise filter: Filters a list or array based on a given predicate function, returning a new list or array with only the elements that satisfy the predicate. | `[1, 2, 3, 4] >- is_even{}` returns `[2, 4]` |
| `|>` | Pipe: Uses the evaluated expression on the left-hand side (LHS) as the input for the function on the right-hand side (RHS). | `[1, 2, 3] |> sum{}` returns `6` |

The special operators allow you to pass the result of the LHS expression as the input of the function on the RHS. For example, in the expression `[1, 2, 3] -> scale{2}`, the `->` operator passes the list `[1, 2, 3]` as the input to the `scale` function, which multiplies each element of the list by `2` and returns a new list with the modified elements.

If the function takes only one argument, you can write, for example, `[1, 2, 3] -> sum{}`, which is equivalent to `sum([1, 2, 3])`.

If the function on the RHS takes more than one argument, you can put the second and subsequent arguments in curly braces, for example, `[1, 2, 3] -> scale{2}`.

The nice part is that you can chain these operators together. Consider the following example:

```
-- Suppose we have a list of numbers
numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

-- Now we want to filter only even numbers, double them, and sum them up.
-- Instead of writing a long chain of function calls like this:
sum(map(filter(numbers, is_even), @(x) { scale(x, 2) }))

-- We can use the special operators in a more functional style:
result = numbers
    >- is_even{}
    -> scale{2}
    |> sum{}

```
