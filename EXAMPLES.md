# Zylisp REPL Examples

## Basic Arithmetic

```zylisp
> (+ 1 2 3)
6

> (- 10 3)
7

> (* 2 3 4)
24

> (/ 24 3 2)
4

> (+ (* 2 3) 4)
10
```

## Variables

```zylisp
> (define x 10)
10

> x
10

> (define y 20)
20

> (+ x y)
30

> (* x y)
200
```

## Lambda Functions

```zylisp
> (define square (lambda (x) (* x x)))
<function>

> (square 5)
25

> (square 10)
100

> (define add (lambda (a b) (+ a b)))
<function>

> (add 3 7)
10
```

## Conditionals

```zylisp
> (if (> 5 3) "yes" "no")
"yes"

> (if (< 5 3) "yes" "no")
"no"

> (define abs (lambda (n) (if (< n 0) (- n) n)))
<function>

> (abs -42)
42

> (abs 17)
17
```

## Lists

```zylisp
> (list 1 2 3 4 5)
(1 2 3 4 5)

> (car (list 1 2 3))
1

> (cdr (list 1 2 3))
(2 3)

> (cons 0 (list 1 2 3))
(0 1 2 3)

> (null? (list))
true

> (null? (list 1))
false
```

## Recursive Functions

### Factorial

```zylisp
> (define factorial
    (lambda (n)
      (if (<= n 1)
          1
          (* n (factorial (- n 1))))))
<function>

> (factorial 5)
120

> (factorial 6)
720
```

### Fibonacci

```zylisp
> (define fib
    (lambda (n)
      (if (<= n 1)
          n
          (+ (fib (- n 1)) (fib (- n 2))))))
<function>

> (fib 0)
0

> (fib 1)
1

> (fib 10)
55
```

### List Sum

```zylisp
> (define sum
    (lambda (lst)
      (if (null? lst)
          0
          (+ (car lst) (sum (cdr lst))))))
<function>

> (sum (list 1 2 3 4 5))
15
```

### List Length

```zylisp
> (define length
    (lambda (lst)
      (if (null? lst)
          0
          (+ 1 (length (cdr lst))))))
<function>

> (length (list 1 2 3 4 5))
5

> (length (list))
0
```

## Comparison Operations

```zylisp
> (= 5 5)
true

> (= 5 3)
false

> (< 3 5)
true

> (> 5 3)
true

> (<= 3 3)
true

> (>= 5 3)
true
```

## Type Predicates

```zylisp
> (number? 42)
true

> (number? (quote x))
false

> (symbol? (quote x))
true

> (list? (list 1 2 3))
true

> (list? 42)
false
```

## Quoting

```zylisp
> (quote (+ 1 2))
(+ 1 2)

> (quote x)
x

> (list (quote a) (quote b) (quote c))
(a b c)
```

## Complex Examples

### Map Function (simple version)

```zylisp
> (define map1
    (lambda (f lst)
      (if (null? lst)
          (list)
          (cons (f (car lst))
                (map1 f (cdr lst))))))
<function>

> (map1 square (list 1 2 3 4))
(1 4 9 16)
```

### Filter Function

```zylisp
> (define is-positive (lambda (n) (> n 0)))
<function>

> (define filter
    (lambda (pred lst)
      (if (null? lst)
          (list)
          (if (pred (car lst))
              (cons (car lst) (filter pred (cdr lst)))
              (filter pred (cdr lst))))))
<function>

> (filter is-positive (list -2 -1 0 1 2 3))
(1 2 3)
```

### Range Function

```zylisp
> (define range
    (lambda (n)
      (if (<= n 0)
          (list)
          (cons n (range (- n 1))))))
<function>

> (range 5)
(5 4 3 2 1)
```

## REPL Commands

```
:help          - Show help message
:reset         - Reset the environment
exit or quit   - Exit the REPL
```

## Tips

- Use parentheses for all function calls: `(+ 1 2)` not `+ 1 2`
- Lambda creates anonymous functions: `(lambda (x) (* x x))`
- Define binds values to names: `(define x 42)`
- All Lisp code is expressions - everything returns a value
- Functions are first-class - they can be passed and returned
- Recursion is the primary iteration mechanism
