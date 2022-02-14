# Wordle

Wordle game unscrambler tool.

## Usage

```shell
wordle --database=/path/to/n5db/file --pattern=XXXXX --negative=Y...
```

Param `--database` points to text file, contains list of English nouns with length 5. Working [example](https://github.com/koykov/wordle/raw/master/nouns5.txt).

Param `--pattern` describes word template and may contains:
* `*` any char (gray)
* `[a-z]` known char (green)
* `^[a-z]` known char but placed in other position (yellow)

Param `--negative` describes list of chars that must be excluded from query. 
