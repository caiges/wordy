# Overview

1. The program accepts as arguments a list of one or more file paths (e.g. ./solution.rb file1.txt file2.txt ...).
2. The program also accepts input on stdin (e.g. cat file1.txt | ./solution.rb).
3. The program outputs a list of the 100 most common three word sequences in the text, along with a count of how many times each occurred in the text. For example: 231 - i will not, 116 - i do not, 105 - there is no, 54 - i know not, 37 - i am not …
4. The program ignores punctuation, line endings, and is case insensitive (e.g. “I love\nsandwiches.” should be treated the same as "(I LOVE SANDWICHES!!)")
5. The program is capable of processing large files and runs as fast as possible.
6. The program should be tested. Provide a test file for your solution.

# Problems

- [x] Ignoring escape sequences, e.g. "\n".
- [x] Knowing the top N sequences implies we've processed all input data.
- [x] Ignoring punctuation. 
- [x] Sort collection when complete.

# Limitations

- More structured use of language, as with URLs, or compound formatting are collapsed and parsed as a single word.

# Developing

Running `make` with no specified target will print the available tasks and descriptions.

- `make` will clean and build the project.
- `make test` will run the tests.
- `go run main.go samples/<nameofsampletext.txt>

# Usage

```
wordy [-grouping=3] [-top=100] [-debug=false] <somefile>
```

or

```
cat <somefile> | wordy [-grouping=3] [-top=100] [-debug=false]
```

## Grouping

The default grouping is set to `3`. The `-grouping` flag can be provided:

`wordy -grouping 10 somefile` will return groupings of `10`.

## Number of results

The default is `100` of the most frequent groupings. The `-top` flag can be provided:

`wordy -top 30 somefile` will return the top `30` results.

## Debug

A `-debug` flag is available for displaying messages about the collection process.
