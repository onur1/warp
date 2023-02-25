# data

Package data provides a set of useful data types for working with different types of data, such as data that contains an error, data that doesn't exist, data that doesn't exist _yet_, data that keeps changing, and etc.

This is a [functional programming](https://github.com/enricopolanski/functional-programming) library. The practical part of FP is all about higher level abstractions such as these and combining their behaviors to gain properties and functionality to deal with data in even more complex situations.

Data types obey certain laws to implement a typeclass instance (that contains implementations of functions defined in a typeclass, such as [functors, applicatives and monads](https://www.adit.io/posts/2013-04-17-functors,_applicatives,_and_monads_in_pictures.html)) specialized to a particular type. These laws are rooted in [category theory](https://www.youtube.com/watch?v=V10hzjgoklA) and together they form a type system that ensures safety and composability.
