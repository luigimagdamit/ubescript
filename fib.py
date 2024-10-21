import time

def fibonacci(n):
    n1, n2 = 0, 1
    i = 0
    while i < n:
        print(n1)
        tmp = n2 
        n2 = n1 + n2
        n1 = tmp 
        i+=1

def fibonacci_with_time(n):
    start_time = time.time()
    fib_number = fibonacci(n)
    end_time = time.time()
    completion_time = end_time - start_time
    return fib_number, completion_time

# Example usage
n = 10  # Change this to compute a different Fibonacci number
result, duration = fibonacci_with_time(n)
print(f"Fibonacci({n}) = {result}, computed in {duration:.6f} seconds.")