package main

import (
    "fmt"
)

func main() {
    var num1, num2 float64
    var op string

    fmt.Print("Enter first number: ")
    fmt.Scanln(&num1)

    fmt.Print("Enter operator (+, -, *, /): ")
    fmt.Scanln(&op)

    fmt.Print("Enter second number: ")
    fmt.Scanln(&num2)

    switch op {
    case "+":
        fmt.Println("Result:", num1+num2)
    case "-":
        fmt.Println("Result:", num1-num2)
    case "*":
        fmt.Println("Result:", num1*num2)
    case "/":
        if num2 != 0 {
            fmt.Println("Result:", num1/num2)
        } else {
            fmt.Println("Error: Division by zero")
        }
    default:
        fmt.Println("Invalid operator")
    }
}
