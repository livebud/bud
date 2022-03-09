# `framework/`

Contains the compilers, drivers, parsers, generators and runtime packages for bud.

Packages are grouped by functionality and suffixed by type. For example, "appc" is short for "app compiler". Framework contains the following types

- **Compilers (c)**: A collection of parsers and generators used to build a binary.

  Currently there are two compilers: a cli compiler and an app compiler. The cli binary is a generated, project-specific task runner used to run generators and commands. The app compiler is used within the cli binary to build your application.

- **Generators (g)**: Generate code that glues your application into a runnable Go program. May use zero or more parsers.

- **Parsers (p)**: Parse the application code to create a structure that can be used by the generators. One parser can be used by multiple generators.

- **Drivers (d)**: Used by tests to programmatically drive the compilers and simplify testing.

- **Runtimes (r)**: Used by the generated code to simplify the generators. These runtimes are not used directly, but referred to by the public `runtime/` package.

# FAQ

## What's with the funny names?

It's the best I could come up with 😅

I wanted the following:

- Meaningful names.
- Similar files are nearby each other.
- Not too many package collisions.
