shortid
=======

> Amazingly short non-sequential url-friendly id generator

Go port of [dylang/shortid](https://github.com/dylang/shortid)


Usage
-----

```go
package main

import(
	"github.com/bradialabs/shortid"
)

func main() {
	s := shortid.New()
	id := s.Generate()
}
```

API
---

```go
import(
	"github.com/bradialabs/shortid"
)
```

---------------------------------------

### `Generate()`

__Returns__ `string` non-sequential unique id.

__Example__

```go
users.insert({
    _id: s.Generate()
    name: ...
    email: ...
    });
```

---------------------------------------

### `SetCharacters(string)`

__Default:__ `'0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_'`

__Returns:__ new alphabet as a `string` 

__Recommendation:__ If you don't like _ or -, you can to set new characters to use. 

__Optional__

Change the characters used.

You must provide a string of all 64 unique characters. Order is not important.

The default characters provided were selected because they are url safe.

__Example__

```go
// use $ and @ instead of - and _
s.SetCharacters('0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ$@');
```

```go
// any 64 unicode characters work, but I wouldn't recommend this.
s.SetCharacters('ⒶⒷⒸⒹⒺⒻⒼⒽⒾⒿⓀⓁⓂⓃⓄⓅⓆⓇⓈⓉⓊⓋⓌⓍⓎⓏⓐⓑⓒⓓⓔⓕⓖⓗⓘⓙⓚⓛⓜⓝⓞⓟⓠⓡⓢⓣⓤⓥⓦⓧⓨⓩ①②③④⑤⑥⑦⑧⑨⑩⑪⑫');
```

---------------------------------------

### `SetWorker(int64)`

__Default:__ `0`

__Recommendation:__ You typically won't want to change this.

__Optional__

If you are running multiple server processes then you should make sure every one has a unique `worker` id. Should be an integer between 0 and 16. 
If you do not do this there is very little chance of two servers generating the same id, but it is theatrically possible 
if both are generated in the exact same second and are generating the same number of ids that second and a half-dozen random numbers are all exactly the same. 

__Example__

```go
s.SetWorker(1);
```

---------------------------------------

### `SetSeed(int64)`

__Default:__ `1`

__Recommendation:__ You typically won't want to change this.

__Optional__

Choose a unique value that will seed the random number generator so users won't be able to figure out the pattern of the unique ids. Call it just once in your application before using `shortId` and always use the same value in your application.

Most developers won't need to use this, it's mainly for testing ShortId. 

If you are worried about users somehow decrypting the id then use it as a secret value for increased encryption.

__Example__

```go
s.SetSeed(1000);
```